package secret

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Service is the provider-neutral credential/session secret boundary. Its
// persistence methods accept only SecretRecord values, which contain an
// authenticated ciphertext envelope and metadata but no plaintext.
type Service interface {
	Store(ctx context.Context, id, kind string, plaintext []byte, expiresAt *time.Time) (SecretReference, error)
	GetReference(ctx context.Context, id string) (SecretReference, error)
	Reveal(ctx context.Context, id string) (RevealedSecret, error)
	Rotate(ctx context.Context, id string) (SecretReference, error)
	Revoke(ctx context.Context, id string) error
}

// SecretService implements Service over an Engine and a persistence contract.
type SecretService struct {
	engine  *Engine
	store   SecretStore
	clock   Clock
	auditor AccessAuditor
}

// NewService creates the boundary service. It does not wire any existing
// provider session; callers must provide an explicit SecretStore contract.
func NewService(engine *Engine, store SecretStore, clock Clock) (*SecretService, error) {
	if engine == nil || store == nil {
		return nil, ErrInvalidSecret
	}
	if clock == nil {
		clock = time.Now
	}
	var auditor AccessAuditor
	if candidate, ok := store.(AccessAuditor); ok {
		auditor = candidate
	}
	return &SecretService{engine: engine, store: store, clock: clock, auditor: auditor}, nil
}

// NewServiceWithAudit is an additive constructor for integrations whose audit
// sink is separate from the SecretStore. Existing SecretStore implementations
// remain valid and NewService still discovers an optional AccessAuditor on the
// store itself.
func NewServiceWithAudit(engine *Engine, store SecretStore, clock Clock, auditor AccessAuditor) (*SecretService, error) {
	service, err := NewService(engine, store, clock)
	if err != nil {
		return nil, err
	}
	service.auditor = auditor
	return service, nil
}

// Store encrypts plaintext immediately and persists only a reference plus
// ciphertext. The caller owns plaintext and should wipe it after this call.
func (s *SecretService) Store(ctx context.Context, id, kind string, plaintext []byte, expiresAt *time.Time) (reference SecretReference, err error) {
	defer func() {
		err = s.finishAudit(ctx, id, kind, AuditActionStore, err)
		if err != nil {
			reference = SecretReference{}
		}
	}()
	if strings.TrimSpace(id) == "" || strings.TrimSpace(kind) == "" {
		return SecretReference{}, ErrInvalidSecret
	}
	envelope, err := s.engine.Encrypt(ctx, plaintext, EncryptOptions{ExpiresAt: expiresAt})
	if err != nil {
		return SecretReference{}, err
	}
	reference = SecretReference{
		ID:         id,
		Kind:       kind,
		KeyVersion: envelope.KeyVersion,
		CreatedAt:  envelope.CreatedAt,
		ExpiresAt:  cloneTime(envelope.ExpiresAt),
	}
	if err := s.store.Create(ctx, SecretRecord{Reference: reference, Envelope: envelope}); err != nil {
		return SecretReference{}, err
	}
	return reference, nil
}

// GetReference returns metadata safe for DTOs. It does not return the
// ciphertext envelope, and never returns plaintext.
func (s *SecretService) GetReference(ctx context.Context, id string) (SecretReference, error) {
	record, err := s.store.Get(ctx, id)
	if err != nil {
		return SecretReference{}, err
	}
	return record.Reference, nil
}

// Reveal performs an atomic one-time claim at the service boundary. The
// returned value is intentionally a redacting wrapper; use Bytes or Take only
// at the provider call site and wipe it when finished.
func (s *SecretService) Reveal(ctx context.Context, id string) (revealed RevealedSecret, err error) {
	kind := ""
	defer func() {
		err = s.finishAudit(ctx, id, kind, AuditActionReveal, err)
		if err != nil {
			revealed.Wipe()
			revealed = RevealedSecret{}
		}
	}()
	record, err := s.store.ClaimReveal(ctx, id, s.clock())
	if err != nil {
		if s.auditor != nil {
			if existing, getErr := s.store.Get(ctx, id); getErr == nil {
				kind = existing.Reference.Kind
			}
		}
		return RevealedSecret{}, err
	}
	kind = record.Reference.Kind
	plaintext, err := s.engine.Decrypt(ctx, record.Envelope)
	if err != nil {
		return RevealedSecret{}, err
	}
	return newRevealedSecret(plaintext), nil
}

// Rotate re-encrypts an active secret with the current key version while
// preserving its expiry and one-time reveal state. It returns only reference
// metadata; plaintext is never returned by rotation.
func (s *SecretService) Rotate(ctx context.Context, id string) (reference SecretReference, err error) {
	kind := ""
	defer func() {
		err = s.finishAudit(ctx, id, kind, AuditActionRotate, err)
		if err != nil {
			reference = SecretReference{}
		}
	}()
	record, err := s.store.Get(ctx, id)
	if err != nil {
		return SecretReference{}, err
	}
	kind = record.Reference.Kind
	if record.RevokedAt != nil {
		return SecretReference{}, ErrSecretRevoked
	}
	if record.isExpired(s.clock()) {
		return SecretReference{}, ErrExpired
	}
	reencrypted, err := s.engine.ReEncrypt(ctx, record.Envelope)
	if err != nil {
		return SecretReference{}, err
	}
	reference = record.Reference
	reference.KeyVersion = reencrypted.KeyVersion
	reference.CreatedAt = reencrypted.CreatedAt
	reference.ExpiresAt = cloneTime(reencrypted.ExpiresAt)
	record.Reference = reference
	record.Envelope = reencrypted
	if err := s.store.Replace(ctx, id, record); err != nil {
		return SecretReference{}, err
	}
	return reference, nil
}

// Revoke prevents future reveals. Revoke is idempotent at the service boundary.
func (s *SecretService) Revoke(ctx context.Context, id string) (err error) {
	kind := ""
	defer func() {
		err = s.finishAudit(ctx, id, kind, AuditActionRevoke, err)
	}()
	if s.auditor != nil {
		record, getErr := s.store.Get(ctx, id)
		if getErr != nil {
			return getErr
		}
		kind = record.Reference.Kind
	}
	return s.store.Revoke(ctx, id, s.clock())
}

// Purge removes revoked or expired records when the configured store exposes
// the optional SecretPurger capability. Legacy SecretStore implementations
// continue to work for all existing operations.
func (s *SecretService) Purge(ctx context.Context) (int, error) {
	purger, ok := s.store.(SecretPurger)
	if !ok {
		return 0, ErrSecretPurgeUnsupported
	}
	return purger.Purge(ctx, s.clock())
}

func (s *SecretService) finishAudit(ctx context.Context, id, kind, action string, operationErr error) error {
	if s.auditor == nil {
		return operationErr
	}
	outcome := AuditOutcomeSuccess
	if operationErr != nil {
		outcome = AuditOutcomeFailure
	}
	auditErr := s.auditor.RecordAccess(ctx, AccessAuditEvent{
		SecretID:  id,
		Kind:      kind,
		Action:    action,
		Outcome:   outcome,
		Timestamp: s.clock().UTC(),
	})
	if operationErr != nil {
		return operationErr
	}
	if auditErr != nil {
		return fmt.Errorf("record secret access audit: %w", auditErr)
	}
	return nil
}

// RevealedSecret is an ephemeral plaintext holder. It has no exported
// plaintext field and all diagnostic/serialization methods are redacted.
type RevealedSecret struct {
	value []byte
}

func newRevealedSecret(value []byte) RevealedSecret {
	return RevealedSecret{value: value}
}

// Bytes returns a defensive copy for the provider call site. The returned copy
// should be wiped by the caller as soon as the provider operation completes.
func (r RevealedSecret) Bytes() []byte { return append([]byte(nil), r.value...) }

// Take returns the value and wipes this wrapper. Prefer Take when the provider
// API accepts a byte slice that can be cleared after use.
func (r *RevealedSecret) Take() []byte {
	if r == nil {
		return nil
	}
	value := r.value
	r.value = nil
	return value
}

// Wipe clears the ephemeral value held by this wrapper.
func (r *RevealedSecret) Wipe() {
	if r == nil {
		return
	}
	wipe(r.value)
	r.value = nil
}

func (r RevealedSecret) String() string   { return "[REDACTED]" }
func (r RevealedSecret) GoString() string { return "[REDACTED]" }
func (r RevealedSecret) MarshalJSON() ([]byte, error) {
	return json.Marshal("[REDACTED]")
}

// FormatSecret is a safe helper for log fields and diagnostics. It avoids
// depending on fmt behavior for arbitrary caller-owned values.
func FormatSecret(_ any) string { return "[REDACTED]" }
