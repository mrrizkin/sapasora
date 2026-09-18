package secret

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// SecretReference is safe to place in provider-neutral DTOs. It identifies
// encrypted material without carrying plaintext or ciphertext.
type SecretReference struct {
	ID         string     `json:"id"`
	Kind       string     `json:"kind"`
	KeyVersion string     `json:"key_version"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

func (r SecretReference) validate() error {
	if r.ID == "" || r.Kind == "" || r.KeyVersion == "" || r.CreatedAt.IsZero() {
		return ErrInvalidSecret
	}
	return nil
}

// String intentionally exposes only the opaque ID, kind, key version, and
// expiry metadata; it never exposes ciphertext or plaintext.
func (r SecretReference) String() string {
	return fmt.Sprintf("secret.Reference{id:%q, kind:%q, key_version:%q, expires_at:%v}", r.ID, r.Kind, r.KeyVersion, r.ExpiresAt)
}

func (r SecretReference) GoString() string { return r.String() }

// SecretRecord is the persistence shape. A store implementation must persist
// this record (or an equivalent shape), never the plaintext passed to Service.
type SecretRecord struct {
	Reference  SecretReference `json:"reference"`
	Envelope   Envelope        `json:"envelope"`
	RevealedAt *time.Time      `json:"revealed_at,omitempty"`
	RevokedAt  *time.Time      `json:"revoked_at,omitempty"`
}

func (r SecretRecord) validate() error {
	if err := r.Reference.validate(); err != nil {
		return err
	}
	if err := r.Envelope.validate(); err != nil {
		return err
	}
	if r.Reference.KeyVersion != r.Envelope.KeyVersion || !sameTime(r.Reference.ExpiresAt, r.Envelope.ExpiresAt) {
		return ErrInvalidSecret
	}
	return nil
}

// String and GoString are safe for structured logging.
func (r SecretRecord) String() string {
	return fmt.Sprintf("secret.Record{%s, envelope:%s, revealed:%t, revoked:%t}", r.Reference.String(), r.Envelope.String(), r.RevealedAt != nil, r.RevokedAt != nil)
}

func (r SecretRecord) GoString() string { return r.String() }

func (r SecretRecord) clone() SecretRecord {
	clone := r
	clone.Envelope = r.Envelope.Clone()
	clone.Reference.ExpiresAt = cloneTime(r.Reference.ExpiresAt)
	clone.RevealedAt = cloneTime(r.RevealedAt)
	clone.RevokedAt = cloneTime(r.RevokedAt)
	return clone
}

// SecretStore is the persistence contract for encrypted secrets. ClaimReveal
// must be atomic in durable implementations so concurrent requests cannot both
// obtain a one-time reveal. Replace must compare-and-swap the lifecycle
// metadata (RevealedAt and RevokedAt), so a stale rotation cannot clear a
// completed reveal or revoke operation.
type SecretStore interface {
	Create(ctx context.Context, record SecretRecord) error
	Get(ctx context.Context, id string) (SecretRecord, error)
	Replace(ctx context.Context, id string, record SecretRecord) error
	ClaimReveal(ctx context.Context, id string, now time.Time) (SecretRecord, error)
	Revoke(ctx context.Context, id string, at time.Time) error
}

// MemoryStore is a concurrency-safe reference implementation for tests and
// local development. It is not a durable production persistence mechanism.
type MemoryStore struct {
	mu      sync.Mutex
	records map[string]SecretRecord
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{records: make(map[string]SecretRecord)}
}

func (s *MemoryStore) Create(_ context.Context, record SecretRecord) error {
	if err := record.validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.records[record.Reference.ID]; exists {
		return ErrSecretConflict
	}
	s.records[record.Reference.ID] = record.clone()
	return nil
}

func (s *MemoryStore) Get(_ context.Context, id string) (SecretRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.records[id]
	if !ok {
		return SecretRecord{}, ErrSecretNotFound
	}
	return record.clone(), nil
}

func (s *MemoryStore) Replace(_ context.Context, id string, record SecretRecord) error {
	if err := record.validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, exists := s.records[id]
	if !exists {
		return ErrSecretNotFound
	}
	if id != record.Reference.ID {
		return ErrInvalidSecret
	}
	if !sameTime(existing.RevealedAt, record.RevealedAt) || !sameTime(existing.RevokedAt, record.RevokedAt) {
		return ErrSecretConflict
	}
	s.records[id] = record.clone()
	return nil
}

func (s *MemoryStore) ClaimReveal(_ context.Context, id string, now time.Time) (SecretRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.records[id]
	if !ok {
		return SecretRecord{}, ErrSecretNotFound
	}
	if record.RevokedAt != nil {
		return SecretRecord{}, ErrSecretRevoked
	}
	if record.isExpired(now) {
		return SecretRecord{}, ErrExpired
	}
	if record.RevealedAt != nil {
		return SecretRecord{}, ErrAlreadyRevealed
	}
	revealedAt := now.UTC()
	record.RevealedAt = &revealedAt
	s.records[id] = record.clone()
	return record.clone(), nil
}

func (s *MemoryStore) Revoke(_ context.Context, id string, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.records[id]
	if !ok {
		return ErrSecretNotFound
	}
	if record.RevokedAt == nil {
		revokedAt := at.UTC()
		record.RevokedAt = &revokedAt
		s.records[id] = record.clone()
	}
	return nil
}

func (r SecretRecord) isExpired(now time.Time) bool {
	return r.Reference.ExpiresAt != nil && !now.Before(r.Reference.ExpiresAt.UTC())
}

func sameTime(left, right *time.Time) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return left.Equal(*right)
}
