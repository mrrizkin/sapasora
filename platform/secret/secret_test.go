package secret

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

type fixedClock struct{ now time.Time }

func (c *fixedClock) Now() time.Time                 { return c.now }
func (c *fixedClock) Advance(duration time.Duration) { c.now = c.now.Add(duration) }

type recordingAuditor struct {
	mu     sync.Mutex
	events []AccessAuditEvent
}

func (a *recordingAuditor) RecordAccess(_ context.Context, event AccessAuditEvent) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = append(a.events, event)
	return nil
}

func (a *recordingAuditor) Events() []AccessAuditEvent {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]AccessAuditEvent(nil), a.events...)
}

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

func TestServiceAuditsSecretLifecycleWithoutSecretMaterial(t *testing.T) {
	clock := &fixedClock{now: time.Date(2026, 4, 5, 6, 7, 8, 0, time.UTC)}
	engine, keyring := testEngine(t, clock, "v1", map[string][]byte{
		"v1": testKey(1),
		"v2": testKey(2),
	})
	auditor := &recordingAuditor{}
	service, err := NewServiceWithAudit(engine, NewMemoryStore(), clock.Now, auditor)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Store(context.Background(), "audit-1", "oauth", []byte("audit-secret"), nil); err != nil {
		t.Fatal(err)
	}
	if err := keyring.SetCurrent("v2"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Rotate(context.Background(), "audit-1"); err != nil {
		t.Fatal(err)
	}
	revealed, err := service.Reveal(context.Background(), "audit-1")
	if err != nil {
		t.Fatal(err)
	}
	revealed.Wipe()
	if err := service.Revoke(context.Background(), "audit-1"); err != nil {
		t.Fatal(err)
	}

	events := auditor.Events()
	if len(events) != 4 {
		t.Fatalf("audit event count = %d, want 4", len(events))
	}
	wantActions := []string{AuditActionStore, AuditActionRotate, AuditActionReveal, AuditActionRevoke}
	for index, event := range events {
		if event.SecretID != "audit-1" || event.Kind != "oauth" || event.Action != wantActions[index] || event.Outcome != AuditOutcomeSuccess || !event.Timestamp.Equal(clock.now) {
			t.Fatalf("audit event %d = %+v", index, event)
		}
		encoded, err := json.Marshal(event)
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(encoded, &fields); err != nil {
			t.Fatal(err)
		}
		if len(fields) != 5 {
			t.Fatalf("audit event fields = %v, want only contract fields", fields)
		}
		for _, field := range []string{"secret_id", "kind", "action", "outcome", "timestamp"} {
			if _, ok := fields[field]; !ok {
				t.Fatalf("audit event missing field %q: %s", field, encoded)
			}
		}
		if strings.Contains(string(encoded), "audit-secret") || strings.Contains(string(encoded), "ciphertext") {
			t.Fatalf("audit event contains secret material: %s", encoded)
		}
	}
}

func TestServicePurgeRemovesOnlyRevokedOrExpiredRecords(t *testing.T) {
	clock := &fixedClock{now: time.Date(2026, 5, 6, 7, 8, 9, 0, time.UTC)}
	engine, _ := testEngine(t, clock, "v1", map[string][]byte{"v1": testKey(1)})
	store := NewMemoryStore()
	service, err := NewService(engine, store, clock.Now)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Store(context.Background(), "active", "credential", []byte("active-secret"), nil); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Store(context.Background(), "revoked", "credential", []byte("revoked-secret"), nil); err != nil {
		t.Fatal(err)
	}
	expiredAt := clock.now
	if _, err := service.Store(context.Background(), "expired", "credential", []byte("expired-secret"), &expiredAt); err != nil {
		t.Fatal(err)
	}
	futureAt := clock.now.Add(time.Hour)
	if _, err := service.Store(context.Background(), "future", "credential", []byte("future-secret"), &futureAt); err != nil {
		t.Fatal(err)
	}
	if err := service.Revoke(context.Background(), "revoked"); err != nil {
		t.Fatal(err)
	}

	removed, err := service.Purge(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if removed != 2 {
		t.Fatalf("purged records = %d, want 2", removed)
	}
	for _, id := range []string{"revoked", "expired"} {
		if _, err := store.Get(context.Background(), id); !errors.Is(err, ErrSecretNotFound) {
			t.Fatalf("purged %q lookup error = %v, want not found", id, err)
		}
	}
	for _, id := range []string{"active", "future"} {
		if _, err := store.Get(context.Background(), id); err != nil {
			t.Fatalf("eligible %q was removed: %v", id, err)
		}
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
