# Device token authentication

`DeviceRepository.GetDeviceByToken` only authenticates a token when all of the
following are true:

- the token status is `active`;
- `expired_at` is `NULL` or later than the database current time;
- the token row is not soft-deleted;
- the device row is not soft-deleted and has the `active` device status; and
- the token and device have the same `user_id` owner relation.

Callers that already have an authenticated API-key or session owner can use the
optional `GetDeviceByTokenForUser` extension. It adds that owner to the query
rather than trusting a token or path value.

The current schema has no `workspace_id` and no workspace relation. The
`user_id` relationship is therefore the strongest owner scope available today;
it is not a substitute for workspace isolation. A future workspace-aware
schema must add an explicit workspace predicate and migration before claiming
workspace-level isolation.

Authentication failures are logged as a generic device-token failure without
the raw credential or a reason that distinguishes missing, revoked, or expired
credentials.
