# 05 — Product Scope dan Roadmap

## Product scope map

### Core messaging

- Personal send.
- Group send bila capability provider mengizinkan.
- Broadcast/campaign dengan consent, approval, rate guard.
- Schedule.
- Recurring (daily, weekly, monthly, annually, custom interval).
- Variable/personalization.
- Text, image, video, audio, file, location, poll.
- Template dan quick reply.
- Follow-up.
- Random/sequential delay sesuai provider/policy.
- Link preview.
- CSV import.

### Inbox dan reply

- Shared inbox.
- Default reply.
- Keyword reply.
- Sender name.
- Submission/form flow.
- Webhook dynamic reply.
- Attachment handling.
- Download file.
- Human handoff.
- Assignment, notes, tags, resolve.

### Channel/device

- Official Cloud API channel.
- WABA/phone onboarding.
- Channel health.
- Credential rotation/revoke.
- Provider capability matrix.
- Optional QR/pairing whatsmeow channel dalam closed beta.
- Device/session lifecycle jika unofficial provider diaktifkan.

### Contacts

- Contact CRUD.
- Tags/segments/groups.
- Variable fields.
- CSV import preview/validation.
- Consent source/timestamp.
- Opt-in/opt-out.
- Suppression list.

### Developer

- API key scoped.
- Send API.
- Get message/status.
- Cancel/reschedule bila didukung.
- Webhook HMAC.
- Webhook tester/replay.
- SDK/examples cURL, Go, JavaScript, PHP, Python.
- Scalar API reference.
- Sandbox/test mode.

### Automation/AI

- Trigger → condition → action.
- Send, assign, tag, delay, notify.
- Run history/pause.
- AI assist dengan data boundary dan human handoff.
- AI history/quota transparan.

### Operations/billing

- Usage/quota.
- Delivery analytics.
- Campaign report.
- Agent performance/response time.
- Subscription/invoice.
- Platform/provider fee separation.
- Audit log.
- Status page/support ticket context.

## MVP P0

1. Workspace/auth/RBAC dasar.
2. Official channel onboarding atau partner integration.
3. Channel health + test message.
4. Contacts, tags, consent, opt-out.
5. Inbox basic: list, timeline, reply, assign, note, resolve.
6. Text/template message.
7. Delivery lifecycle dan actionable failure.
8. API key create/revoke/rotate.
9. Send API + message status.
10. Signed webhook + retry.
11. Usage meter.
12. OpenAPI + Scalar + quickstart.
13. Audit/security baseline.

## P1

- Campaign approval/audience preview/schedule/quiet hours.
- Media/poll/location sesuai official capability.
- Recurring/follow-up.
- Automation builder basic.
- Round-robin/SLA.
- Template management/approval status.
- CSV import.
- Retry/DLQ/replay UI.
- Analytics.
- n8n, Make, Zapier, WordPress, WooCommerce, Google Forms.
- AI assist.

## P2

- Whatsmeow closed beta dengan label experimental.
- Agency mode/multi-workspace.
- Multiple WABA/channel routing.
- Advanced RBAC/audit/retention.
- Enterprise SLA/DPA/SSO/SCIM.
- AI knowledge base.
- Conversion attribution.
- Omnichannel.

## Out of scope awal

- Arbitrary bulk send tanpa consent.
- Claim anti-ban.
- AI autonomous reply tanpa review/handoff.
- Omnichannel penuh.
- Enterprise procurement sebelum reliability terbukti.

## Roadmap

### Phase 0 — Validate, minggu 1–3

**Output:** problem/segment/channel validation.

- Interview 15 target users.
- Validasi official onboarding vs BSP.
- Prototype onboarding/inbox/delivery/pricing.
- Feasibility spike GoFiber/Fibertia/VSA.
- Threat model dan legal/policy review awal.

**Exit:** wedge dipilih, success metrics disepakati, provider path jelas.

### Phase 1 — Foundation, minggu 4–6

- Repo GoFiber + Fibertia + Vue/Nuxt UI.
- VSA conventions.
- Podman local stack.
- PostgreSQL/sqlx/goose.
- Valkey.
- zerolog/observability.
- Auth/workspace/RBAC.
- OpenAPI/Scalar.
- Mock provider.

**Exit:** vertical slice end-to-end dengan mock dan CI.

### Phase 2 — Internal alpha, minggu 7–10

- Official test channel.
- Channel onboarding/health.
- Contacts/consent.
- Send test/text/template.
- Status timeline.
- Inbox basic.
- API/webhook/idempotency.
- Zod frontend contracts.

**Exit:** internal team dapat menyelesaikan core flow tanpa manual DB intervention.

### Phase 3 — Private pilot, minggu 11–14

- 5–10 workspace.
- Invoice/appointment/lead/support use case.
- Support-assisted onboarding.
- Strict campaign limits.
- Usage/billing beta.
- Incident/runbook review.

**Exit:** 70% onboarding completion, first delivered <15 menit setelah channel siap, dua customer bersedia membayar.

### Phase 4 — Public beta, minggu 15–20

- Self-serve onboarding.
- Pricing/overage cap.
- Campaign P1.
- Automation basic.
- Status page.
- Public docs/changelog.
- Security review.

**Exit:** operational metrics stabil, support process siap, no critical security gaps.

### Phase 5 — GA dan scale, setelah minggu 20

- Reliability hardening.
- Agency mode.
- More integrations.
- Optional whatsmeow closed beta.
- Enterprise controls.
- Advanced analytics/AI.

## Release gates

### Alpha

Internal/test data, mock/provider fixtures, no real campaign.

### Private beta

Feature flags, daily error review, manual support, incident runbook.

### Public beta

Self-serve, billing, legal reviewed, status page, documented limitations.

### GA

Security/restore test, SLO, support policy, provider failure drill, data deletion verification.
