# 18 — Team Handoff Checklist

## Cara menggunakan

Dokumen ini digunakan saat kickoff, planning sprint, code review, staging review, dan release gate. Setiap checklist memiliki owner dan evidence.

## Product alignment

- [ ] Semua team membaca 01 Product Brief.
- [ ] Problem/primary segment sudah disepakati.
- [ ] Wedge MVP jelas: shared inbox + transactional notification.
- [ ] Official/unofficial positioning dipahami.
- [ ] Non-goals disepakati.
- [ ] Success metrics disepakati.

## Branding

- [ ] Working name dicatat tanpa menganggap final.
- [ ] Positioning option dipilih/test.
- [ ] Domain/social/trademark screening owner ditentukan.
- [ ] Voice/tone dan istilah official/experimental disetujui.

## UX/design

- [ ] IA/navigation disetujui.
- [ ] First-run checklist tersedia.
- [ ] Official channel flow tersedia.
- [ ] Unofficial flow memiliki risk acknowledgement.
- [ ] Send/status/failure recovery selesai.
- [ ] Inbox/campaign flow selesai.
- [ ] Nuxt UI tokens/components terdokumentasi.
- [ ] Accessibility baseline diuji.

## Architecture

- [ ] VSA adalah top-level organization.
- [ ] Feature slice boundaries ditentukan.
- [ ] Ports/adapters hanya di external boundary.
- [ ] Official/whatsmeow integration terisolasi.
- [ ] Async lifecycle/message events dirancang.
- [ ] Multi-tenancy dan workspace scope dijaga.
- [ ] Failure isolation/kill switch ada.

## Backend

- [ ] GoFiber/Fibertia conventions.
- [ ] Go error mapping/request ID.
- [ ] sqlx repository + prepared statement lifecycle.
- [ ] goose migrations.
- [ ] Valkey namespace/TTL/queue semantics.
- [ ] zerolog redaction.
- [ ] OpenAPI/Scalar updated.
- [ ] Zod schema contract tersedia untuk frontend.

## ID contract

- [ ] Internal DB: `id BIGSERIAL/BIGINT`.
- [ ] Internal DB: `public_id VARCHAR(21)` NanoID.
- [ ] External: `id: string = public_id`.
- [ ] External tidak memiliki `public_id`.
- [ ] BIGSERIAL tidak muncul pada API response.
- [ ] BIGSERIAL tidak muncul pada Inertia props/UI network.
- [ ] BIGSERIAL tidak muncul pada customer webhook/public analytics.
- [ ] API lookup memakai workspace + public ID + authorization.
- [ ] Database model tidak langsung diserialize.

## Frontend

- [ ] Vue 3 + TypeScript.
- [ ] Inertia navigation, bukan vue-router.
- [ ] Tidak ada Pinia tanpa ADR.
- [ ] Zod request/response/page props/shared props.
- [ ] Network data dimulai sebagai `unknown`.
- [ ] Tidak ada `any`/unsafe type assertion di boundary.
- [ ] Invalid schema fail early/loudly.
- [ ] Error boundary dan telemetry ter-redact.
- [ ] Public DTO `id` bertipe string.

## Security/privacy

- [ ] API keys scoped/rotatable/revocable.
- [ ] Provider credentials encrypted.
- [ ] HMAC webhook + replay protection.
- [ ] Consent/opt-out/suppression.
- [ ] Quiet hours/rate guard/approval.
- [ ] SSRF/media controls.
- [ ] PII/log redaction.
- [ ] Retention/delete/export.
- [ ] Official policy/legal review.
- [ ] Unofficial risk/legal/license review.

## Reliability/operations

- [ ] Timeout/circuit breaker.
- [ ] Retry classification.
- [ ] Queue backpressure/DLQ/replay.
- [ ] Provider health metrics.
- [ ] Status page.
- [ ] Alert/runbook.
- [ ] Backup/restore tested.
- [ ] RPO/RTO/SLO disepakati.
- [ ] Incident communication owner.

## QA

- [ ] Unit test domain/mapper/Zod.
- [ ] Repository integration test.
- [ ] Provider contract fixtures.
- [ ] Webhook signature/replay test.
- [ ] Tenant isolation test.
- [ ] ID exposure test.
- [ ] Inertia props schema test.
- [ ] E2E core user journey.
- [ ] Load/resilience test sesuai phase.

## Launch

- [ ] Pilot users dan use case dipilih.
- [ ] Support FAQ/runbook.
- [ ] Billing/usage tested.
- [ ] Docs/tutorial/changelog.
- [ ] Product analytics tanpa PII.
- [ ] Public/unofficial disclosure.
- [ ] Release/rollback plan.
- [ ] Post-launch review schedule.

## Handoff owner matrix

| Area | Owner | Evidence |
|---|---|---|
| Product scope | TBD | Approved brief/backlog |
| UX/design | TBD | Prototype/design spec |
| Backend/VSA | TBD | Architecture/ADR/code |
| Official provider | TBD | Meta/BSP setup |
| Whatsmeow | TBD | Experimental readiness review |
| Security/legal | TBD | Review record |
| Infrastructure | TBD | Runbook/alerts/DR test |
| QA | TBD | Test report |
| Support/ops | TBD | FAQ/incident process |
