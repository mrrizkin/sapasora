// Package secret provides provider-neutral encrypted secret envelopes and the
// persistence contract used by credential and provider-session integrations.
//
// The package deliberately keeps plaintext out of envelopes, references,
// records, and JSON/log representations. Plaintext is accepted only at the
// service boundary and is returned only through RevealedSecret after an
// atomic, one-time reveal claim.
package secret

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	AlgorithmAES256GCM = "aes-256-gcm"
	envelopeVersion    = 1
)

var (
	ErrInvalidKey             = errors.New("invalid secret encryption key")
	ErrInvalidEnvelope        = errors.New("invalid secret envelope")
	ErrKeyNotFound            = errors.New("secret encryption key not found")
	ErrAuthentication         = errors.New("secret envelope authentication failed")
	ErrExpired                = errors.New("secret has expired")
	ErrInvalidSecret          = errors.New("invalid secret")
	ErrSecretNotFound         = errors.New("secret not found")
	ErrSecretRevoked          = errors.New("secret has been revoked")
	ErrAlreadyRevealed        = errors.New("secret has already been revealed")
	ErrSecretConflict         = errors.New("secret store conflict")
	ErrSecretPurgeUnsupported = errors.New("secret store does not support purge")
)

// Key is an encryption key held by a Keyring. Material is never serialized by
// this package; callers should treat it as process-local sensitive material.
type Key struct {
	Version  string
	Material []byte
}

func (k Key) valid() error {
	if strings.TrimSpace(k.Version) == "" || len(k.Material) != 32 {
		return ErrInvalidKey
	}
	return nil
}

// String intentionally does not include key material.
func (k Key) String() string {
	return fmt.Sprintf("secret.Key{version:%q, material:[REDACTED]}", k.Version)
}

// GoString prevents %#v diagnostics from exposing key material.
func (k Key) GoString() string { return k.String() }

// MarshalJSON omits key material even if a key is accidentally included in a
// diagnostic payload.
func (k Key) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Version string `json:"version"`
	}{Version: k.Version})
}

// Keyring resolves the active key and historical keys by version. A durable
// implementation should resolve keys from a KMS/secret manager rather than a
// database record containing plaintext key material.
type Keyring interface {
	Current(ctx context.Context) (Key, error)
	Lookup(ctx context.Context, version string) (Key, error)
}

// StaticKeyring is a small process-local keyring useful for tests and local
// development. Production integrations should implement Keyring against a
// managed key service.
type StaticKeyring struct {
	mu      sync.RWMutex
	current string
	keys    map[string][]byte
}

// NewStaticKeyring creates a keyring from versioned 32-byte AES-256 keys.
// Material is copied before it is retained.
func NewStaticKeyring(current string, keys map[string][]byte) (*StaticKeyring, error) {
	if strings.TrimSpace(current) == "" {
		return nil, ErrInvalidKey
	}
	keyring := &StaticKeyring{current: current, keys: make(map[string][]byte, len(keys))}
	for version, material := range keys {
		key := Key{Version: version, Material: material}
		if err := key.valid(); err != nil {
			return nil, err
		}
		keyring.keys[version] = append([]byte(nil), material...)
	}
	if _, ok := keyring.keys[current]; !ok {
		return nil, ErrKeyNotFound
	}
	return keyring, nil
}

// SetCurrent changes the active key version. It is intended for controlled key
// rotation; existing envelopes remain decryptable while their old key exists.
func (k *StaticKeyring) SetCurrent(version string) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	if _, ok := k.keys[version]; !ok {
		return ErrKeyNotFound
	}
	k.current = version
	return nil
}

func (k *StaticKeyring) Current(ctx context.Context) (Key, error) {
	k.mu.RLock()
	version := k.current
	k.mu.RUnlock()
	return k.Lookup(ctx, version)
}

func (k *StaticKeyring) Lookup(_ context.Context, version string) (Key, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	material, ok := k.keys[version]
	if !ok {
		return Key{}, ErrKeyNotFound
	}
	return Key{Version: version, Material: append([]byte(nil), material...)}, nil
}

// String and GoString expose only key versions, never key material.
func (k *StaticKeyring) String() string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	versions := make([]string, 0, len(k.keys))
	for version := range k.keys {
		versions = append(versions, version)
	}
	sort.Strings(versions)
	return fmt.Sprintf("secret.StaticKeyring{current:%q, versions:%v}", k.current, versions)
}

func (k *StaticKeyring) GoString() string { return k.String() }

func (k *StaticKeyring) MarshalJSON() ([]byte, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	versions := make([]string, 0, len(k.keys))
	for version := range k.keys {
		versions = append(versions, version)
	}
	sort.Strings(versions)
	return json.Marshal(struct {
		Current  string   `json:"current"`
		Versions []string `json:"versions"`
	}{Current: k.current, Versions: versions})
}

// Clock makes expiry behavior deterministic in tests and explicit in callers.
type Clock func() time.Time

// EngineConfig configures an AES-GCM envelope engine.
type EngineConfig struct {
	Keyring Keyring
	Clock   Clock
	// Random is injectable for deterministic tests. Leave nil to use
	// crypto/rand.Reader.
	Random io.Reader
}

// Engine encrypts and decrypts provider-neutral secret envelopes.
type Engine struct {
	keyring Keyring
	clock   Clock
	random  io.Reader
}

// NewEngine creates an AES-256-GCM envelope engine.
func NewEngine(config EngineConfig) (*Engine, error) {
	if config.Keyring == nil {
		return nil, ErrInvalidKey
	}
	if config.Clock == nil {
		config.Clock = time.Now
	}
	if config.Random == nil {
		config.Random = rand.Reader
	}
	return &Engine{keyring: config.Keyring, clock: config.Clock, random: config.Random}, nil
}

// NewEngineWithKeyring is a convenience constructor using the system clock and
// cryptographically secure randomness.
func NewEngineWithKeyring(keyring Keyring) (*Engine, error) {
	return NewEngine(EngineConfig{Keyring: keyring})
}

// EncryptOptions contains authenticated envelope metadata. ExpiresAt is
// authenticated as part of the envelope and is retained in the reference.
type EncryptOptions struct {
	ExpiresAt *time.Time
}

// Envelope is the only encrypted value persisted by a SecretStore. It contains
// no plaintext. Ciphertext and Nonce are safe to serialize, but should still
// be treated as confidential because they are protected secret material.
type Envelope struct {
	Version    int        `json:"version"`
	Algorithm  string     `json:"algorithm"`
	KeyVersion string     `json:"key_version"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	Nonce      []byte     `json:"nonce"`
	Ciphertext []byte     `json:"ciphertext"`
}

func (e Envelope) validate() error {
	if e.Version != envelopeVersion || e.Algorithm != AlgorithmAES256GCM || strings.TrimSpace(e.KeyVersion) == "" {
		return ErrInvalidEnvelope
	}
	if e.CreatedAt.IsZero() || len(e.Nonce) == 0 || len(e.Ciphertext) == 0 {
		return ErrInvalidEnvelope
	}
	return nil
}

// String includes only non-sensitive metadata and sizes.
func (e Envelope) String() string {
	return fmt.Sprintf("secret.Envelope{algorithm:%q, key_version:%q, created_at:%q, ciphertext_bytes:%d}", e.Algorithm, e.KeyVersion, e.CreatedAt.UTC().Format(time.RFC3339Nano), len(e.Ciphertext))
}

// GoString prevents verbose diagnostics from dumping ciphertext bytes.
func (e Envelope) GoString() string { return e.String() }

// Clone returns a defensive copy suitable for store implementations.
func (e Envelope) Clone() Envelope {
	clone := e
	clone.Nonce = append([]byte(nil), e.Nonce...)
	clone.Ciphertext = append([]byte(nil), e.Ciphertext...)
	if e.ExpiresAt != nil {
		expiresAt := *e.ExpiresAt
		clone.ExpiresAt = &expiresAt
	}
	return clone
}

// Encrypt seals plaintext with the current key. The plaintext is copied only
// for the duration of the cryptographic operation and is not retained.
func (e *Engine) Encrypt(ctx context.Context, plaintext []byte, options EncryptOptions) (Envelope, error) {
	key, err := e.keyring.Current(ctx)
	if err != nil {
		return Envelope{}, err
	}
	if err := key.valid(); err != nil {
		return Envelope{}, err
	}

	block, err := aes.NewCipher(key.Material)
	if err != nil {
		return Envelope{}, ErrInvalidKey
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return Envelope{}, ErrInvalidKey
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(e.random, nonce); err != nil {
		return Envelope{}, fmt.Errorf("generate secret nonce: %w", err)
	}
	envelope := Envelope{
		Version:    envelopeVersion,
		Algorithm:  AlgorithmAES256GCM,
		KeyVersion: key.Version,
		CreatedAt:  e.clock().UTC(),
		ExpiresAt:  cloneTime(options.ExpiresAt),
		Nonce:      nonce,
	}
	envelope.Ciphertext = gcm.Seal(nil, nonce, plaintext, envelope.associatedData())
	return envelope, nil
}

// Decrypt authenticates and decrypts an envelope. Authentication failures do
// not disclose whether ciphertext or authenticated metadata was modified.
func (e *Engine) Decrypt(ctx context.Context, envelope Envelope) ([]byte, error) {
	if err := envelope.validate(); err != nil {
		return nil, err
	}
	key, err := e.keyring.Lookup(ctx, envelope.KeyVersion)
	if err != nil {
		return nil, err
	}
	if err := key.valid(); err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key.Material)
	if err != nil {
		return nil, ErrInvalidKey
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil || len(envelope.Nonce) != gcm.NonceSize() {
		return nil, ErrInvalidEnvelope
	}
	plaintext, err := gcm.Open(nil, envelope.Nonce, envelope.Ciphertext, envelope.associatedData())
	if err != nil {
		return nil, ErrAuthentication
	}
	if envelope.ExpiresAt != nil && !e.clock().Before(envelope.ExpiresAt.UTC()) {
		wipe(plaintext)
		return nil, ErrExpired
	}
	return plaintext, nil
}

// ReEncrypt decrypts an envelope with its recorded key version and encrypts it
// with the current key. Expiry metadata is preserved and authenticated again.
func (e *Engine) ReEncrypt(ctx context.Context, envelope Envelope) (Envelope, error) {
	plaintext, err := e.Decrypt(ctx, envelope)
	if err != nil {
		return Envelope{}, err
	}
	defer wipe(plaintext)
	return e.Encrypt(ctx, plaintext, EncryptOptions{ExpiresAt: envelope.ExpiresAt})
}

func (e Envelope) associatedData() []byte {
	expiresAt := ""
	if e.ExpiresAt != nil {
		expiresAt = e.ExpiresAt.UTC().Format(time.RFC3339Nano)
	}
	return fmt.Appendf(nil, "sapasora/secret-envelope/%d\x00%s\x00%s\x00%s\x00%s", e.Version, e.Algorithm, e.KeyVersion, e.CreatedAt.UTC().Format(time.RFC3339Nano), expiresAt)
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	clone := value.UTC()
	return &clone
}

func wipe(value []byte) {
	for index := range value {
		value[index] = 0
	}
}
