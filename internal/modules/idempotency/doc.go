// Package idempotency provides a provider-neutral, transport-independent
// idempotency foundation for request processing.
//
// It validates bounded idempotency keys, scopes them to a workspace, actor,
// and endpoint, fingerprints payloads without retaining them, and coordinates
// an in-memory claim/result store. HTTP middleware, distributed storage, and
// durable database adapters remain integration concerns and are intentionally
// not part of this package.
package idempotency
