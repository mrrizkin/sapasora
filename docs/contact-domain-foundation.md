# Contact domain foundation and CRUD boundary (Tracks 5.1–5.3)

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

## Track 5.2 CRUD boundary

`ContactService` in `internal/modules/contact/service.go` is the use-case
boundary for tenant-scoped create, get, list, update, and soft-delete
operations on contacts and addresses. It requires the tenant ID on every
operation, validates ownership assertions on write models, resolves contact
ownership before address operations, and exposes only the stable
`ErrNotFound`, `ErrConflict`, and `ErrInvalid` error families.

The service is intentionally transport-neutral. HTTP/controllers must select
and authorize the tenant before calling it, map domain models into explicit
API DTOs (including deliberate PII choices), and translate domain errors into
HTTP responses. This track does **not** add controllers, routes, provider
adapters, consent enforcement, database migrations/tables, or DI/module
registration. Durable adapters must preserve the tenant-scoped public-ID and
normalized identity uniqueness constraints atomically.

## Track 5.3 CSV import/export foundation

`import_export.go` provides bounded, transport-neutral CSV streams through
`ImportCSV` and `ExportCSV` (or their importer/exporter wrappers). Both APIs
require an explicit tenant ID. Import validates the exact `CSVHeader` schema,
limits rows and field bytes, normalizes addresses through `NormalizeAddress`,
reports duplicate identities by row/field, and supports dry-run validation
without repository writes. Export writes tenant-scoped contacts and normalized
address values directly to an `io.Writer`; contacts without addresses remain
representable as rows with empty address columns.

CSV issue/error diagnostics contain only row numbers, field names, and stable
codes. Raw or normalized contact values are never included in errors or logs.
The current repository interface lists contacts before streaming their rows, so
HTTP/storage integrations must still apply request/body limits and choose job
or transaction behavior. HTTP upload/download routes, authorization and
export permission checks, async jobs, XLSX/mapping support, audit events, and
expiring URLs remain outside this module.

## Track 5.6 duplicate and identity boundary

`duplicate.go` provides exact matching on canonical `(kind, namespace, value)`
address identities, but exposes only identity fingerprints and tenant-safe
contact identifiers in candidate and merge diagnostics. Fuzzy names, raw
values, cross-tenant lookup, and implicit identity linking are not supported.

Merge preview returns a deterministic confirmation token. The in-memory
repository applies a confirmed merge atomically only when that token is still
current, records immutable merge metadata, and supports append-only undo only
when the moved addresses and target have not changed in a conflicting way.
Durable adapters must provide the same `MergeRepository` capability
transactionally. HTTP/UI approval and integration remain deferred.
