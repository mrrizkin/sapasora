# Message domain foundation (Track 7.1)

`internal/modules/message` is the additive, provider-neutral foundation for
logical messages and delivery history.

## Included

- Tenant/workspace-scoped `Message` records with inbound/outbound direction,
  provider message ID, idempotency key, request ID, and contact, conversation,
  and channel references.
- Normalized content retained for domain operations while safe JSON/String
  diagnostics expose only content metadata, never message text or opaque keys.
- `MessageAttachment` metadata with opaque storage/provider references; media
  bytes are not part of the domain model.
- Provider-neutral `DeliveryRecord` observations and immutable, sequence-numbered
  `MessageEvent` history.
- Validation of ownership, references, direction/status, timestamps, and
  metadata counts. Provider metadata is key-redacted before storage and is not
  emitted by diagnostics.
- Concurrency-safe in-memory repository with defensive copies, tenant-scoped
  lookups, idempotency-key uniqueness, association ownership checks, and
  append-only event storage.

## Deliberate boundary

This track does not implement provider send/API wiring, HTTP routes, UI, state
transition rules, queue/outbox behavior, webhook ingestion, migrations, or DI
registration. Message statuses and event names are stored vocabularies only;
transition validation belongs to Track 7.2. The in-memory repository is a
domain/test adapter, not production persistence.
