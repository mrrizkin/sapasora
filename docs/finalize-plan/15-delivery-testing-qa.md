# 15 — Delivery, Testing, dan QA

## Delivery model

Gunakan vertical slice: UI → GoFiber/Fibertia action → feature service → sqlx/Valkey/provider → status → UI/analytics.

## Definition of Ready

- User/problem jelas.
- Acceptance criteria.
- Provider capability.
- Loading/empty/error/recovery UX.
- Security/privacy impact.
- Analytics event.
- Dependency/migration plan.

## Definition of Done

- Code review.
- Unit/integration/e2e/contract test.
- OpenAPI/schema updated.
- Goose migration tested/backward-compatible.
- Logs/metrics/traces.
- Tenant/security checks.
- Zod schemas/fixtures.
- Error/recovery UX.
- Docs/changelog/support FAQ.
- Feature flag/rollback bila risky.
- Product acceptance.

## Test pyramid

### Unit

- Domain/policy/consent.
- Billing.
- Status mapping/retry/idempotency.
- Zod valid/invalid fixtures.
- Inertia props parse/fail-early.
- API error parsing.
- DTO mapping: `id` berasal dari `public_id`.

### Integration

- sqlx repositories/PostgreSQL.
- Prepared statement lifecycle.
- Valkey queue/cache/lock.
- Object storage.
- Provider mock HTTP.
- Webhook signature.
- Tenant isolation.

### Contract

- OpenAPI request/response.
- Backend response ↔ Zod.
- Inertia props ↔ page schema.
- Webhook schema/version.
- External ID: `id: string`, tanpa internal `id`/`public_id`.

### E2E

- Register/login/workspace.
- Connect test channel.
- Test message dan delivery timeline.
- Inbound/reply/assignment.
- API key/webhook.
- Failed message recovery.
- Billing/usage.
- Invalid form/response/props menghasilkan error state, bukan silent failure.

### Resilience/load

- API enqueue throughput.
- Queue backlog.
- Webhook fan-out.
- Provider timeout/rate limit.
- Campaign pause/cancel.
- Worker crash/restart.
- DB restore/failover.

## QA scenarios

Channel: OAuth cancel, permission denied, duplicate number, disconnect during send, provider unavailable.

Messaging: invalid recipient, missing consent, template rejected, duplicate idempotency, delayed status, transient/permanent failure, retry after timeout.

Campaign: empty audience, mixed consent, missing variable, quiet hours, approval revoked, pause/cancel partial delivery.

Security: cross-tenant denied, revoked key denied, replay webhook rejected, token absent logs, SSRF blocked, unauthorized export denied, invalid Zod data reported.

## Release gates

- Alpha: internal/test data/mock provider, no real campaign.
- Private beta: 5–10 workspaces, feature flags, daily error review, runbook.
- Public beta: self-serve, billing/legal, status page, docs.
- GA: security/restore test, SLO, support policy, provider failure drill, deletion verification.

## Release checklist

- [ ] Migration tested empty + realistic DB.
- [ ] Backward compatibility.
- [ ] OpenAPI/Scalar updated.
- [ ] Zod schemas/contract fixtures updated.
- [ ] IDs checked at every boundary.
- [ ] Provider health.
- [ ] Alerts/runbooks.
- [ ] Rollback.
- [ ] Support/customer communication.
- [ ] Analytics events.
