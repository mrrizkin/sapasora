# Contact domain foundation (Track 5.1)

`internal/modules/contact` is the additive, provider-neutral foundation for
contacts and contact identities.

## Boundary

- `Contact` is tenant-scoped by `TenantID`; `ContactAddress` is tenant-scoped
  and belongs to a tenant-owned contact.
- Phone, email, and username values are normalized before they enter the
  repository. Username identities may use an explicit provider namespace;
  phone and email default to the `global` namespace.
- `AddressIdentity` is the duplicate-detection key. A durable repository must
  enforce uniqueness for `(tenant_id, kind, namespace, normalized_value)` and
  tenant-scoped public IDs atomically.
- Consent is address-level metadata (`unknown`, `opted_in`, `opted_out`) with
  source, timestamp, opaque evidence reference, and opaque actor reference.
- String and JSON diagnostic projections intentionally omit names, notes,
  custom fields, raw values, normalized values, provider address IDs, and
  consent evidence. Explicit domain operations remain responsible for access
  control when a value is needed.

The in-memory repository is concurrency-safe and exists for domain tests/local
composition. It is not a production persistence implementation.

## Deliberately not wired yet

This track does **not** add database migrations/tables, controllers, routes,
provider adapters, provider sync, consent enforcement, CRUD APIs, or DI/module
registration. Future integration should map provider payloads at the provider
boundary, authorize the tenant before repository calls, and add a durable
unique constraint plus migration in a separate track. API DTOs should make an
explicit choice about which PII may be returned rather than serializing domain
models directly.
