# Channel lifecycle state machine

Track 4.3 currently provides the provider-neutral lifecycle state machine in
`internal/modules/channel`. It validates transitions without coupling the
channel domain to WhatsApp, Telegram, or another provider adapter.

## States

The supported states are:

- `pending`: waiting for a connection attempt or retry;
- `connecting`: a connection attempt is in progress;
- `connected`: the provider connection is healthy;
- `degraded`: the connection exists but health is impaired;
- `disconnected`: no active connection is available;
- `expired`: a credential, session, or pairing window is no longer valid;
- `error`: the latest lifecycle operation failed.

A state transition to the current state is accepted as an idempotent no-op.
Rejected transitions return `*channel.InvalidTransitionError`, which wraps
`channel.ErrInvalidTransition`; the current state is unchanged.

## Valid transitions

```text
pending     -> connecting, expired, error
connecting  -> connected, degraded, disconnected, expired, error
connected   -> degraded, disconnected, expired, error
degraded    -> connecting, connected, disconnected, expired, error
disconnected -> pending
expired     -> pending
error       -> pending
```

`disconnected`, `expired`, and `error` re-enter through `pending` so a future
connection attempt is explicit. The provider-neutral package also provides a
bounded, context-aware `ReconnectBackoffPolicy` with deterministic exponential
delays by default, optional bounded jitter, and no wait after the final
attempt. Circuit breaking, locks, history, metrics, and provider integration
remain outside this focused slice.

The policy is an integration boundary, not a provider migration: existing
WhatsApp and Telegram startup paths continue using their `providerstartup`
retry helper. Provider adapters can adopt `channel.ReconnectBackoffPolicy` in a
separate integration change once their lifecycle ownership and observability
contracts are ready.

The existing `FakeAdapter` remains an additive adapter-contract fake. It is not
wired to the lifecycle state machine yet: its contract tests intentionally
exercise direct `Connect`/`Disconnect` behavior, and changing that behavior
would silently expand this Track 4.3 slice into adapter migration.
