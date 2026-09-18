# Channel domain foundation

Track 4.1 provides the additive, provider-neutral channel domain in
`internal/modules/channel`.

## Included

- `ChannelAccount` with tenant scope, internal ID/public ID separation, channel
  type, provider, capabilities, lifecycle status, credential/session references,
  and quota/rate policy.
- Connection event value model with tenant/public-ID references.
- `ChannelAccountRepository` and a concurrency-safe in-memory implementation
  used by repository tests.
- Defensive repository snapshots, tenant-scoped public-ID lookups, soft-delete
  behavior, and append-only connection event storage.

Credential and session values are opaque references only. Secret material,
provider payloads, usage counters, and provider-specific state do not belong in
this model.

## Integration boundary

This slice does not add migrations, database wiring, API routes, DI bindings,
or changes to the existing WhatsApp/Telegram device and gateway providers.
The in-memory repository is a domain/test adapter, not a production persistence
implementation. A later persistence integration must map `TenantID` to the
workspace/tenant table, preserve the internal-ID/public-ID contract, and add
its own constraints and migration only when that storage boundary is approved.
Provider adapters remain behind the existing adapter contracts and must be
integrated provider-by-provider in a separate change.
