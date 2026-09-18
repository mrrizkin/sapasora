// Package contact defines the provider-neutral, tenant-scoped contact domain.
//
// This package intentionally stops at the domain and repository boundary. It
// does not register HTTP controllers, provider adapters, or database
// migrations. Consent transitions and suppression checks are provider-neutral
// domain operations; provider/API wiring (including STOP and unsubscribe
// handlers) remains deferred. CSV import/export is exposed as a bounded,
// tenant-explicit stream foundation in import_export.go; HTTP upload/download,
// authorization, async jobs, and audit wiring remain integration concerns.
// Those integrations must map provider identities into ContactAddress through
// the normalization constructors and keep tenant authorization outside this
// package. Duplicate candidates and merge previews are safe, transport-neutral
// domain projections; merge mutation requires explicit confirmation and the
// repository's atomic merge capability.
package contact
