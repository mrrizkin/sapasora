# 12 — Delivery Process, Testing, dan QA

## Delivery model

Gunakan vertical slice: satu fitur selesai dari UI → API → domain → DB → worker → provider → status → analytics.

Jangan membangun seluruh frontend dulu lalu backend; validasi lifecycle end-to-end per slice.

## Definition of Ready

Sebuah story siap dikerjakan jika:

- user/problem jelas;
- acceptance criteria tersedia;
- provider capability diketahui;
- error/empty/loading state dirancang;
- security/privacy impact diidentifikasi;
- metric/event analytics ditentukan;
- dependencies dan migration plan jelas.

## Definition of Done

- Code reviewed.
- Unit/integration/e2e sesuai kebutuhan.
- OpenAPI/schema updated.
- Migration backward-compatible.
- Logs/metrics/traces tersedia.
- Permission/tenant isolation tested.
- Error dan recovery UX tersedia.
- Docs/changelog updated.
- Feature flag/rollback tersedia bila risky.
- Product owner acceptance selesai.

## Testing pyramid

### Unit

- Frontend Zod schema valid/invalid fixtures.
- Inertia props parse dan fail-early behavior.
- API response parsing dan error boundary.
- Domain rules.
- Consent/policy checks.
- Billing calculations.
- Status mapping.
- Retry classification.
- Idempotency.

### Integration

- PostgreSQL repositories.
- Queue publish/consume.
- Object storage lifecycle.
- Provider adapter with mock HTTP.
- Webhook signature.
- Tenant isolation.

### Contract

- OpenAPI request/response.
- Backend response ↔ frontend Zod schema.
- Inertia page props ↔ page-specific Zod schema.
- External ID contract: `id: string`, tanpa `public_id`/internal `id`.
- Provider response fixtures.
- Webhook event schema.
- Version compatibility.

### E2E

- Register/login.
- Create workspace.
- Connect test channel.
- Send test message.
- See delivery timeline.
- Receive inbound event.
- Reply/assign inbox.
- Create API key/webhook.
- Failed message recovery.
- Billing/usage display.

### Load/resilience

- API enqueue throughput.
- Queue backlog.
- Webhook fan-out.
- Campaign pause/cancel.
- Provider timeout/rate limit.
- DB failover/restore.
- Worker crash/restart.

## Test doubles

- Official provider mock dengan success/error/rate-limit fixtures.
- Webhook event replay fixtures.
- Fake clock untuk schedule/recurring.
- Fake object storage.
- Fake payment provider.
- Whatsmeow adapter hanya integration test terisolasi, bukan dalam default test suite.

## QA scenarios penting

### Channel

- Connect success.
- OAuth canceled.
- Permission denied.
- Phone number already connected.
- Disconnect during send.
- Provider unavailable.

### Messaging

- Invalid recipient.
- Missing consent.
- Template not approved.
- Duplicate idempotency key.
- Provider accepted but status webhook delayed.
- Permanent vs transient failure.
- User retries after timeout.

### Campaign

- Empty audience.
- Mixed consent.
- Variable missing/fallback.
- Quiet hours.
- Approval revoked.
- Pause while queue active.
- Cancel after partial delivery.

### Security

- Workspace A cannot access workspace B.
- Zod validation failures tidak membocorkan PII/credential.
- Revoked API key denied.
- Replayed webhook rejected.
- Token absent from logs.
- Attachment SSRF blocked.
- Unauthorized export denied.

## Release checklist

- [ ] Database migration tested on copy.
- [ ] Backward compatibility checked.
- [ ] Feature flag configured.
- [ ] Provider health checked.
- [ ] Alert/runbook updated.
- [ ] Rollback tested.
- [ ] Support FAQ prepared.
- [ ] Customer communication prepared.
- [ ] Analytics events verified.
- [ ] Post-release owner assigned.
