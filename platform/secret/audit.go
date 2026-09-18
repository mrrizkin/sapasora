package secret

import (
	"context"
	"time"
)

// AccessAuditEvent is the complete provider-neutral audit payload for a secret
// operation. It intentionally contains no plaintext, ciphertext, envelope, or
// caller-controlled metadata beyond the opaque secret ID and kind.
type AccessAuditEvent struct {
	SecretID  string    `json:"secret_id"`
	Kind      string    `json:"kind"`
	Action    string    `json:"action"`
	Outcome   string    `json:"outcome"`
	Timestamp time.Time `json:"timestamp"`
}

// SecretAccessAuditEvent is an explicit alias for integrations that prefer the
// longer contract name.
type SecretAccessAuditEvent = AccessAuditEvent

const (
	AuditActionStore  = "store"
	AuditActionReveal = "reveal"
	AuditActionRotate = "rotate"
	AuditActionRevoke = "revoke"

	AuditOutcomeSuccess = "success"
	AuditOutcomeFailure = "failure"
)

// AccessAuditor records secret access events. Implementations must persist or
// forward only AccessAuditEvent fields; in particular, they must not require
// plaintext or ciphertext from the service.
type AccessAuditor interface {
	RecordAccess(ctx context.Context, event AccessAuditEvent) error
}

// SecretAccessAuditor is an explicit alias for the optional audit contract.
type SecretAccessAuditor = AccessAuditor

// SecretPurger is an optional SecretStore capability. Purge removes records
// that are revoked or expired at now and returns the number removed.
type SecretPurger interface {
	Purge(ctx context.Context, now time.Time) (removed int, err error)
}

// SecretStorePurger is an explicit alias for the optional purge capability.
type SecretStorePurger = SecretPurger
