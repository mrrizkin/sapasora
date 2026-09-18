# Channel circuit breaker

The provider-neutral circuit breaker in `internal/modules/channel` protects a
provider/channel call boundary without importing provider SDKs. It has three
states: `closed`, `open`, and `half_open`. Configured consecutive failures open
the circuit; the cooldown admits exactly one half-open probe; a successful
operation resets the failure count and closes the circuit. Typed errors let
callers distinguish an open cooldown from a probe already in flight.

`Execute(ctx, operation)` checks cancellation before admission and passes the
same context to the operation. The breaker is safe for concurrent use and its
clock is injectable for deterministic tests. It does not retry, classify, or
normalize provider errors; the returned operation error remains the adapter's
error for the caller to normalize according to the existing channel error
catalog.

## Integration boundary

This slice is additive and is not wired into the existing WhatsApp, Telegram,
or other provider implementations. The future integration point is the
orchestration/service boundary that owns a provider/channel operation: create
one breaker with the intended provider/channel scope, then call the provider
through `Execute` and handle the typed open/half-open rejection there. Provider
adapters remain responsible for provider calls and error normalization, while
the lifecycle state machine remains responsible for connection state. A later
integration must choose breaker ownership/lifetime, persistence or process
scope, metrics, and which provider failures count as circuit failures.
