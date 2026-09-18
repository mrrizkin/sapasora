package secret

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

type fixedClock struct{ now time.Time }

func (c *fixedClock) Now() time.Time                 { return c.now }
func (c *fixedClock) Advance(duration time.Duration) { c.now = c.now.Add(duration) }

func testEngine(t *testing.T, clock *fixedClock, current string, keys map[string][]byte) (*Engine, *StaticKeyring) {
	t.Helper()
	keyring, err := NewStaticKeyring(current, keys)
	if err != nil {
		t.Fatal(err)
	}
	engine, err := NewEngine(EngineConfig{
		Keyring: keyring,
		Clock:   clock.Now,
		Random:  bytes.NewReader(bytes.Repeat([]byte{0x42}, 1024)),
	})
	if err != nil {
		t.Fatal(err)
	}
	return engine, keyring
}

func testKey(value byte) []byte { return bytes.Repeat([]byte{value}, 32) }

func TestEnvelopeRoundTripAndSafeRepresentations(t *testing.T) {
	clock := &fixedClock{now: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)}
	engine, _ := testEngine(t, clock, "v1", map[string][]byte{"v1": testKey(1)})
	expiresAt := clock.now.Add(time.Hour)
	plaintext := []byte("provider credential must not be persisted")
	envelope, err := engine.Encrypt(context.Background(), plaintext, EncryptOptions{ExpiresAt: &expiresAt})
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err := engine.Decrypt(context.Background(), envelope)
	if err != nil {
		t.Fatal(err)
	}
	if string(decrypted) != string(plaintext) {
		t.Fatalf("round trip = %q, want %q", decrypted, plaintext)
	}
	wipe(decrypted)

	serialized, err := json.Marshal(envelope)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(serialized), string(plaintext)) {
		t.Fatalf("serialized envelope contains plaintext: %s", serialized)
	}
	if strings.Contains(fmt.Sprintf("%+v", envelope), string(plaintext)) {
		t.Fatalf("envelope diagnostic contains plaintext: %+v", envelope)
	}
}

func TestEnvelopeRejectsWrongKeyAndTampering(t *testing.T) {
	clock := &fixedClock{now: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)}
	engine, _ := testEngine(t, clock, "v1", map[string][]byte{"v1": testKey(1), "v2": testKey(2)})
	envelope, err := engine.Encrypt(context.Background(), []byte("secret"), EncryptOptions{})
	if err != nil {
		t.Fatal(err)
	}

	wrongEngine, _ := testEngine(t, clock, "v1", map[string][]byte{"v1": testKey(2)})
	if _, err := wrongEngine.Decrypt(context.Background(), envelope); !errors.Is(err, ErrAuthentication) {
		t.Fatalf("wrong key error = %v, want authentication failure", err)
	}

	tamperedCiphertext := envelope.Clone()
	tamperedCiphertext.Ciphertext[len(tamperedCiphertext.Ciphertext)-1] ^= 0x01
	if _, err := engine.Decrypt(context.Background(), tamperedCiphertext); !errors.Is(err, ErrAuthentication) {
		t.Fatalf("tampered ciphertext error = %v, want authentication failure", err)
	}

	tamperedMetadata := envelope.Clone()
	tamperedMetadata.KeyVersion = "v2"
	if _, err := engine.Decrypt(context.Background(), tamperedMetadata); !errors.Is(err, ErrAuthentication) {
		t.Fatalf("tampered key version error = %v, want authentication failure", err)
	}

	tamperedExpiry := envelope.Clone()
	expiresAt := clock.now.Add(time.Hour)
	tamperedExpiry.ExpiresAt = &expiresAt
	if _, err := engine.Decrypt(context.Background(), tamperedExpiry); !errors.Is(err, ErrAuthentication) {
		t.Fatalf("tampered expiry error = %v, want authentication failure", err)
	}
}

func TestServiceRotationPreservesExpiryAndRevealIsOneTime(t *testing.T) {
	clock := &fixedClock{now: time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)}
	engine, keyring := testEngine(t, clock, "v1", map[string][]byte{
		"v1": testKey(1),
		"v2": testKey(2),
	})
	store := NewMemoryStore()
	service, err := NewService(engine, store, clock.Now)
	if err != nil {
		t.Fatal(err)
	}
	expiresAt := clock.now.Add(2 * time.Hour)
	reference, err := service.Store(context.Background(), "credential-1", "oauth", []byte("refresh-token"), &expiresAt)
	if err != nil {
		t.Fatal(err)
	}
	if reference.KeyVersion != "v1" || reference.ExpiresAt == nil || !reference.ExpiresAt.Equal(expiresAt) {
		t.Fatalf("initial reference = %+v", reference)
	}

	if err := keyring.SetCurrent("v2"); err != nil {
		t.Fatal(err)
	}
	rotated, err := service.Rotate(context.Background(), "credential-1")
	if err != nil {
		t.Fatal(err)
	}
	if rotated.KeyVersion != "v2" {
		t.Fatalf("rotated key version = %q, want v2", rotated.KeyVersion)
	}
	if rotated.ExpiresAt == nil || !rotated.ExpiresAt.Equal(expiresAt) {
		t.Fatalf("rotation changed expiry metadata: %+v", rotated)
	}

	revealed, err := service.Reveal(context.Background(), "credential-1")
	if err != nil {
		t.Fatal(err)
	}
	if got := string(revealed.Bytes()); got != "refresh-token" {
		t.Fatalf("revealed value = %q", got)
	}
	if got := fmt.Sprintf("%+v", revealed); got != "[REDACTED]" {
		t.Fatalf("revealed diagnostic = %q", got)
	}
	revealed.Wipe()
	if _, err := service.Reveal(context.Background(), "credential-1"); !errors.Is(err, ErrAlreadyRevealed) {
		t.Fatalf("second reveal error = %v, want already revealed", err)
	}
}

func TestMemoryStoreRejectsStaleLifecycleReplacement(t *testing.T) {
	clock := &fixedClock{now: time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)}
	engine, _ := testEngine(t, clock, "v1", map[string][]byte{"v1": testKey(1)})
	service, err := NewService(engine, NewMemoryStore(), clock.Now)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Store(context.Background(), "credential-cas", "api_key", []byte("secret"), nil); err != nil {
		t.Fatal(err)
	}

	store := service.store.(*MemoryStore)
	stale, err := store.Get(context.Background(), "credential-cas")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.ClaimReveal(context.Background(), "credential-cas", clock.now); err != nil {
		t.Fatal(err)
	}
	if err := store.Replace(context.Background(), "credential-cas", stale); !errors.Is(err, ErrSecretConflict) {
		t.Fatalf("stale replacement error = %v, want conflict", err)
	}
}

func TestServiceExpiryIsEnforcedDeterministically(t *testing.T) {
	clock := &fixedClock{now: time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)}
	engine, _ := testEngine(t, clock, "v1", map[string][]byte{"v1": testKey(1)})
	service, err := NewService(engine, NewMemoryStore(), clock.Now)
	if err != nil {
		t.Fatal(err)
	}
	expiresAt := clock.now.Add(time.Minute)
	if _, err := service.Store(context.Background(), "credential-expiring", "api_key", []byte("secret"), &expiresAt); err != nil {
		t.Fatal(err)
	}
	clock.Advance(time.Minute)
	if _, err := service.Reveal(context.Background(), "credential-expiring"); !errors.Is(err, ErrExpired) {
		t.Fatalf("expired reveal error = %v, want expired", err)
	}
}

func TestRedactionHelpersNeverExposeSecretValues(t *testing.T) {
	revealed := newRevealedSecret([]byte("do-not-log"))
	for _, diagnostic := range []string{fmt.Sprint(revealed), fmt.Sprintf("%+v", revealed), fmt.Sprintf("%#v", revealed)} {
		if strings.Contains(diagnostic, "do-not-log") || diagnostic != "[REDACTED]" {
			t.Fatalf("unsafe diagnostic = %q", diagnostic)
		}
	}
	encoded, err := json.Marshal(revealed)
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) != `"[REDACTED]"` {
		t.Fatalf("JSON redaction = %s", encoded)
	}

	reference := SecretReference{ID: "ref-1", Kind: "session", KeyVersion: "v1", CreatedAt: time.Unix(1, 0).UTC()}
	if strings.Contains(fmt.Sprint(reference), "do-not-log") {
		t.Fatal("reference diagnostic leaked secret")
	}
}
