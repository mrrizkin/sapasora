// Package contact defines the provider-neutral, tenant-scoped contact domain.
//
// This package intentionally stops at the domain and repository boundary. It
// does not register HTTP controllers, provider adapters, or database
// migrations. Consent transitions and suppression checks are provider-neutral
// domain operations; provider/API wiring (including STOP and unsubscribe
// handlers) remains deferred. Those integrations must map provider identities
// into ContactAddress through the normalization constructors and keep tenant
// authorization outside this package.
package contact
