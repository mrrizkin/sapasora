# Sapasora Product Master Plan

> **Status:** Consolidated product, discovery, architecture, and delivery source of truth
>
> **Scope:** Full-feature product; bukan MVP-only plan
>
> **Last consolidation source:** `docs/finalize-plan/` dari `~/Downloads/finalize-plan.zip`
>
> Dokumen ini menggabungkan hasil discovery, planning, audit codebase, dan spesifikasi full-feature Sapasora. Dokumen discovery asli tetap dipertahankan sebagai audit trail di [`finalize-plan/`](./finalize-plan/).

---

## 1. Cara Membaca Dokumen

Dokumen ini membedakan tiga hal:

1. **Current** — hal yang sudah terlihat di repository saat audit.
2. **Target** — kapabilitas yang harus dimiliki Sapasora full-feature.
3. **Migration** — pekerjaan untuk membawa current menuju target.

Dokumen normatif untuk detail tertentu:

| Topik | Dokumen canonical |
|---|---|
| Product/architecture/API/security/delivery yang sudah dikonsolidasi | Dokumen ini |
| Full feature requirements dan acceptance checklist | [`FULL_FEATURE_SPECIFICATION.md`](./FULL_FEATURE_SPECIFICATION.md) |
| Risiko implementasi saat ini | [`PRODUCTION_READINESS.md`](./PRODUCTION_READINESS.md) |
| Discovery dan planning asli | [`finalize-plan/README.md`](./finalize-plan/README.md) |

Jika terjadi konflik:

1. Keamanan dan tenant isolation mengalahkan convenience.
2. Kontrak di dokumen ini mengalahkan contoh lama.
3. Provider capability mengalahkan asumsi universal.
4. Status implementasi harus diverifikasi ke source code, bukan dianggap selesai karena ada dokumentasi.

---

# 2. Product Definition

## 2.1 Produk

Sapasora adalah **customer communication operating system** untuk bisnis dan developer:

- shared inbox;
- transactional notification;
- customer conversation;
- templates;
- campaign/broadcast yang aman;
- scheduling dan recurring delivery;
- automation dan AI-assisted workflow;
- analytics, usage, dan billing;
- public API dan webhook;
- provider/channel management.

Sapasora bukan sekadar endpoint `send message`. Produk harus menjawab seluruh lifecycle:

```text
Connect → Configure → Consent → Send → Queue → Provider accept
→ Sent → Delivered/Read → Reply → Assign → Resolve → Analyze → Recover
```

## 2.2 Primary wedge

**Official WhatsApp shared inbox + transactional notification**.

Alasannya:

- nilai bisnis mudah diukur;
- frekuensi pemakaian tinggi;
- lebih aman daripada memulai dengan unsolicited broadcast;
- menjadi fondasi inbox, automation, API, billing, dan retention;
- cocok untuk owner nonteknis dan developer.

Campaign, automation, SMS, Email, Telegram, AI, dan channel eksperimental tetap bagian dari target full-feature, tetapi dibangun di atas fondasi delivery dan safety yang sama.

## 2.3 Target pengguna

- SMB dengan 1–10 agent: ecommerce, jasa, edukasi, klinik, appointment, property, komunitas.
- Developer dan agency yang membangun integrasi untuk banyak client.
- Supervisor/support lead yang membutuhkan assignment, SLA, dan audit.
- Enterprise sebagai expansion: SSO/SCIM, DPA, SLA, procurement, dedicated support.

## 2.4 Positioning

> Untuk bisnis dan developer yang membutuhkan komunikasi customer yang mudah digunakan, dapat diukur, dan dapat diintegrasikan, Sapasora menyediakan inbox, delivery visibility, safe automation, dan API production-grade. Berbeda dari gateway yang hanya fokus mengirim pesan, Sapasora mengelola consent, channel, delivery, reply, recovery, analytics, dan billing dalam satu workspace.

## 2.5 Non-goals

- Bulk spam atau arbitrary broadcast tanpa consent.
- Janji anti-ban atau 100% delivery.
- Menyamarkan konektor unofficial sebagai official.
- AI autonomous reply tanpa boundary dan human handoff.
- Menggantikan CRM/ERP sepenuhnya.
- Mengaktifkan semua provider sebelum lifecycle core stabil.

---

# 3. Canonical Product Decisions

## 3.1 Channel strategy

| Channel | Target position | Policy |
|---|---|---|
| WhatsApp Cloud API/official partner | Default production | Onboarding, template/policy, official capability, provider billing |
| WhatsApp Web/whatsmeow | Experimental/advanced | Feature flag, closed beta, risk acknowledgement, kill switch, no official SLA |
| Telegram | Full adapter target | Capability-based, lifecycle/error handling lengkap |
| SMS | Full adapter target | Provider abstraction, opt-out, delivery receipt, country/compliance rules |
| Email | Full adapter target | SMTP/API, bounce, suppression, unsubscribe, delivery events |
| Future channels | Adapter boundary | Tidak boleh mengubah domain message/inbox |

Current repository memiliki WhatsApp melalui `whatsmeow` dan Telegram melalui TDLib. Itu dianggap **legacy/provider implementation yang harus dimigrasikan ke channel abstraction**, bukan alasan untuk memposisikan unofficial channel sebagai default production.

## 3.2 Architecture

- Vertical Slice Architecture sebagai organisasi utama.
- Modular monolith terlebih dahulu.
- Ports/adapters hanya pada boundary eksternal.
- Worker dan queue untuk seluruh operasi provider/asynchronous.
- Service extraction hanya setelah bottleneck, ownership, dan operational need terbukti.

## 3.3 Stack target dan kondisi repository

| Area | Target | Kondisi saat audit | Migration decision |
|---|---|---|---|
| Backend | Go + GoFiber | Sudah ada | Pertahankan |
| Dashboard | Inertia + Vue 3 + TypeScript | Sudah ada | Pertahankan |
| UI | Nuxt UI | Sudah/direncanakan | Standardisasi |
| Navigation | Inertia, tanpa vue-router default | Inertia | Pertahankan |
| Client state | Server/page props, local state | Belum seragam | Hindari global state tanpa ADR |
| Validation | Zod di frontend + backend validation | Type-check ada, runtime boundary belum lengkap | Tambahkan bertahap |
| Database access | sqlx/prepared statements untuk target | Repository saat ini menggunakan pola GORM | Migrasi incremental, tidak big-bang |
| Migration | goose atau runner yang memiliki lock/transaction | Migration runner existing belum cukup hardened | Pilih satu owner dan tambahkan lock |
| Cache/queue | Valkey dengan durable job semantics | Belum lengkap | Implementasikan outbox/worker |
| Logging | zerolog/structured logs | Logging ada | Redaction dan field standard |
| API docs | OpenAPI + Scalar | Swagger ada dengan mismatch security | Regenerate/fix contract |
| Storage | S3-compatible object storage | Media masih local path | Migrasikan melalui storage abstraction |
| Local infra | Podman-compatible stack | Docker tersedia | Docker tetap didukung; make targets portable |
| IDs | Internal BIGINT + external NanoID string | Current code perlu migrasi | Terapkan pada public contract, internal join tetap numeric |

Tidak ada migrasi stack yang boleh dilakukan hanya demi nama library. Setiap perubahan database/query/migration harus melewati ADR, benchmark, compatibility plan, dan rollback.

## 3.4 External ID contract

```text
internal: id BIGINT/BIGSERIAL
external: id string = public_id (NanoID 21 karakter)
```

Aturan:

- `id` internal tidak pernah keluar.
- `public_id` juga tidak pernah keluar sebagai field terpisah.
- External request/response/Inertia/webhook memakai `id` string.
- Lookup selalu `workspace_id + public_id` dan tetap melalui authorization.
- NanoID bukan encryption dan bukan authorization.

## 3.5 Product safety defaults

- Consent sebelum campaign/marketing delivery.
- Opt-out per channel dan global suppression.
- Quiet hours dan timezone.
- Approval threshold untuk bulk/costly action.
- Hard recipient/rate limit.
- Test recipient dan dry-run.
- Idempotency pada send/launch/retry.
- HMAC webhook, replay protection, retry, DLQ.
- Credential reveal once, hash/encrypt at rest, rotate/revoke.
- Uofficial connector selalu terlihat sebagai experimental.

---

# 4. Full Feature Catalog

## 4.1 Identity, workspace, and team

- Register/login/logout.
- Email verification.
- Forgot/reset/change password.
- MFA/TOTP/recovery codes.
- Session/device list dan remote revoke.
- Login throttling dan suspicious login.
- Workspace create/switch/settings.
- Timezone, locale, branding, notification settings.
- Member invitation/expiry/resend/revoke.
- Suspend/remove member.
- Ownership transfer.
- Custom role dan permission.
- SSO/SAML/OIDC/SCIM untuk enterprise.
- Workspace delete, grace period, purge.

## 4.2 Channel onboarding

- Official WhatsApp/WABA onboarding.
- Provider/phone/channel identification.
- QR/pairing untuk experimental whatsmeow.
- Telegram bot/user/TDLib onboarding sesuai policy.
- SMS sender/provider setup.
- Email SMTP/API setup.
- Capability discovery.
- Health/status/last event/last error.
- Connect/reconnect/disconnect/delete.
- Credential/session rotation dan revoke.
- Auto-connect policy.
- Per-channel quota, rate limit, timezone, alias.
- Connection history dan diagnostics.

## 4.3 Messaging modes

- Transactional notification.
- Personal send.
- Reply dari inbox.
- Campaign/broadcast.
- One-time schedule.
- Recurring schedule: daily, weekly, monthly, annually, custom interval.
- Follow-up sequence.
- Automation-triggered message.
- API-triggered message.
- Provider event-driven reply.

## 4.4 Message types

Capability-dependent, tetapi target domain mendukung:

- Text.
- Image.
- Video.
- Audio.
- Document/file.
- Sticker.
- Location.
- Contact card.
- Poll.
- Button/action.
- List/menu.
- Quick reply.
- Template.
- Reaction.
- Quoted reply/thread.
- Form/submission flow jika provider mendukung.
- Link preview.
- Sender/display name sesuai provider.

Platform wajib melakukan capability negotiation. Fitur yang tidak didukung provider harus menghasilkan error actionable, bukan panic atau silent drop.

## 4.5 Unified inbox

- Conversation list/search/filter.
- Channel/team/agent/status/priority/tag/SLA filter.
- Timeline inbound/outbound/status/internal note/automation/system event.
- Reply text/media/template.
- Attachment download.
- Assignment/reassignment/claim/unclaim.
- Internal note.
- Tags dan priority.
- Open/pending/snoozed/resolved/closed.
- Unread/read.
- Quote/reply/thread.
- Collision detection.
- Human handoff.
- Supervisor takeover.
- SLA timer dan breach.
- Conversation merge/split.
- Export permission-gated.

## 4.6 Contacts and consent

- Contact CRUD/restore/purge.
- Multiple phone/email/channel identifiers.
- Tags, groups, segments, custom fields.
- Import CSV/XLSX dengan preview/mapping/error report.
- Export permission-gated dan audited.
- Duplicate detection/merge.
- Contact timeline.
- Consent source/timestamp/evidence/actor.
- Per-channel opt-in/opt-out.
- STOP/START SMS handling.
- Email unsubscribe/suppression group.
- Global do-not-contact list.
- Identity linking lintas channel.

## 4.7 Templates and quick replies

- Draft/review/approved/rejected/archived.
- Versioning/rollback/clone.
- Per-channel variant.
- Provider approval mapping.
- Localization.
- Variable type/required/default/example.
- Escape dan redaction per channel.
- Preview/sample data/test send.
- Rejection reason.
- Usage/performance.
- Deprecation.

## 4.8 Campaign and broadcast

- Draft/audience preview.
- Consent/suppression filtering.
- Audience snapshot atau resolve-at-run-time.
- Template/variable selection.
- Test send.
- Approval workflow.
- Schedule/timezone/quiet hours.
- Throttle/rate guard.
- Frequency cap.
- Pause/resume/cancel.
- A/B variant.
- Recipient-level result.
- Partial failure.
- Cost estimate.
- Delivery report/export.
- Campaign clone.

## 4.9 Automation, autoresponder, and AI

### Trigger

Inbound message, keyword/regex, new contact, tag change, conversation state, delivery failure, schedule, webhook, campaign event, contact field change.

### Condition

Channel, device, contact field, tag/segment, content, business hours, conversation status, previous run, quota/rate.

### Action

Send/template, assign/team, tag, delay, note, status/priority, webhook, task, contact update, stop automation, human escalation.

### AI assist

- Suggested reply.
- Summarize conversation.
- Classify intent/sentiment.
- Extract structured fields.
- Knowledge-base retrieval.
- Human approval/handoff.
- Data boundary and PII redaction.
- Per-workspace model/provider setting.
- Quota, cost, prompt/output audit.
- No autonomous send by default.

### Automation operations

- Visual builder.
- Draft/publish/disable.
- Versioning.
- Simulation/dry-run.
- Loop/depth protection.
- Per-contact concurrency lock.
- Execution history, retry, DLQ, replay.
- Kill switch.

## 4.10 Developer platform

- Scoped API key/service account.
- Test/live environment.
- One-time secret reveal.
- Hash-at-rest.
- Expiry/rotate/revoke.
- IP allowlist.
- Last-used/usage analytics.
- REST API v1.
- OpenAPI + Scalar.
- cURL, Go, JavaScript, PHP, Python examples.
- SDK generation.
- Sandbox/test recipients.
- Idempotency-Key.
- Cursor pagination.
- Stable error envelope.
- Rate-limit headers.
- Webhook tester/replay.

## 4.11 Webhook and integrations

- Incoming provider webhook verification.
- Outgoing subscription/filter.
- HTTPS/SSRF protection.
- HMAC signing/key rotation.
- Timestamp/replay protection.
- Event IDs and dedup.
- Retry/backoff/DLQ/replay.
- Pause after repeated failure.
- Payload versioning.
- n8n, Make, Zapier, WordPress, WooCommerce, Google Forms integration targets.
- Import/export and support ticket context.

## 4.12 Analytics and operations

- Message sent/delivered/read/failed.
- Delivery lag/provider error.
- First response/resolution time.
- Agent workload/SLA.
- Campaign/template/automation performance.
- Channel health.
- Webhook success.
- Usage/cost by workspace/channel/campaign.
- Dashboard, comparison, CSV export, scheduled report.
- Status page.
- Request ID diagnostics.
- Provider/campaign pause control.
- Dead-letter inspection and replay.

## 4.13 Billing and economics

- Sandbox/free plan.
- Starter/Growth/Scale plans.
- Platform subscription.
- Provider/Meta usage separated.
- Media/storage/AI overage.
- Usage meter/ledger.
- Quota/entitlement.
- Cost estimate before campaign.
- Hard overage cap.
- Invoice/history.
- Upgrade/downgrade/proration.
- Trial/grace/cancel.
- Tax/currency.
- Usage alerts.
- Contribution margin tracking.

## 4.14 Media and storage

- Direct/presigned/resumable upload.
- MIME sniffing, extension/size check.
- Virus/malware scan.
- Thumbnail/transcoding.
- Encryption at rest.
- Signed expiring downloads.
- Quota/retention/orphan cleanup.
- Object storage versioning/lifecycle.
- Provider-specific conversion.

## 4.15 Admin, audit, and compliance

- Admin workspace/user/channel/job lookup.
- Support impersonation with consent/banner/expiry/audit.
- Feature flags and kill switches.
- Immutable audit log.
- Data export/delete request.
- Retention/anonymization.
- Privacy policy/DPA/subprocessor inventory.
- Consent evidence.
- Campaign approval history.
- Security incident evidence.

---

# 5. Canonical User Flows

## 5.1 First-run activation

```text
Register → Verify email → Create workspace → Choose use case
→ Choose official/experimental channel with disclosure
→ Connect channel → Health check → Send test message
→ Configure team/consent → Open Inbox/API quickstart
```

Exit condition: user melihat accepted, provider status, dan delivery outcome; jika gagal, user mendapat next action yang jelas.

## 5.2 Outbound API

```text
API request + auth + idempotency
→ tenant/scope check
→ recipient/consent/capability/quota validation
→ create accepted message
→ outbox/queue
→ provider adapter
→ lifecycle events
→ customer webhook/UI/analytics
```

## 5.3 Inbound message

```text
Provider event
→ verify signature/authenticity
→ persist encrypted raw event with retention
→ deduplicate
→ normalize
→ upsert contact/conversation/message
→ routing/automation
→ notify inbox
→ optional reply/handoff
```

## 5.4 Campaign

```text
Draft → Select audience → Preview/exclusion/consent
→ Estimate cost → Test send → Approval
→ Schedule → Fan-out with throttle
→ Per-recipient delivery → Pause/retry/DLQ
→ Report and audit
```

## 5.5 Disconnect/delete channel

```text
Confirm → stop new sends → drain/cancel eligible jobs
→ disconnect provider session → revoke credential
→ cleanup media/session → preserve required audit
→ mark deleted → show recovery/retention status
```

---

# 6. Domain, Data, and Contract Baseline

## 6.1 Entities

```text
User
Workspace
WorkspaceMember
Role
Permission
Invitation
Session
MFAFactor
ApiCredential
ChannelAccount
ProviderSession
Contact
ContactAddress
ConsentEvent
ContactTag
ContactSegment
Conversation
ConversationParticipant
ConversationAssignment
Message
MessageAttachment
MessageDelivery
MessageEvent
Template
TemplateVersion
Campaign
CampaignRecipient
Schedule
Automation
AutomationVersion
AutomationRun
WebhookEndpoint
WebhookSubscription
WebhookDelivery
InboundEvent
OutboxEvent
Job
DeadLetter
AuditEvent
UsageRecord
Subscription
Invoice
Notification
```

## 6.2 Message states

```text
accepted → queued → sending → sent → delivered → read
                         └──────→ failed → retrying → sent|dead_letter
accepted/queued → canceled
```

Platform membedakan `accepted` oleh Sapasora dari `sent`/`delivered` oleh provider.

## 6.3 Required invariants

- Semua tenant-owned query memiliki `workspace_id`.
- Soft-deleted/revoked/expired credential selalu ikut predicate lookup.
- Public ID lookup tidak pernah melewati authorization.
- Provider event ID unique per channel account.
- Send/launch/retry idempotent.
- Status transition divalidasi domain service.
- UTC pada storage; timezone di scheduling/presentation.
- Secret/PII tidak masuk DTO publik.
- Purge asynchronous, resumable, idempotent, audited.

## 6.4 API envelope

```json
{
  "data": {},
  "meta": {"request_id": "req_..."}
}
```

```json
{
  "error": {
    "code": "MESSAGE_VALIDATION_FAILED",
    "message": "Message payload is invalid",
    "next_action": "Perbaiki payload lalu kirim ulang",
    "retryable": false,
    "request_id": "req_..."
  }
}
```

Public API wajib mendukung versioning, idempotency, cursor pagination, rate-limit headers, request ID, OpenAPI, dan changelog.

---

# 7. Frontend Information Architecture

## Public/auth

`/login`, `/register`, `/verify-email`, `/forgot-password`, `/reset-password`, `/mfa`

## Workspace/admin

`/workspaces`, `/workspace/settings`, `/workspace/members`, `/workspace/roles`, `/workspace/audit`, `/workspace/billing`, `/workspace/usage`, `/admins`

## Channels

`/devices`, `/devices/new`, `/devices/:id`, `/devices/:id/edit`, `/devices/:id/connection`, `/devices/:id/logs`

## CRM/inbox

`/phonebook`, `/phonebook/import`, `/phonebook/segments`, `/inbox`, `/inbox/:conversationId`, `/inbox/assigned`, `/inbox/unassigned`

## Messages/content

`/messages`, `/messages/compose`, `/messages/scheduled`, `/messages/failed`, `/templates`, `/templates/new`, `/templates/:id`

## Campaign/automation

`/campaigns`, `/campaigns/new`, `/campaigns/:id`, `/recurring`, `/autoresponder`, `/autoresponder/:id/builder`

## Developer/integrations

`/integrations`, `/integrations/webhooks`, `/integrations/api-keys`, `/documentation`, `/integrations/sandbox`

## Operations

`/analytics`, `/reports`, `/notifications`, `/settings`, `/support`, `/status`

Setiap route wajib memiliki loading, empty, error, permission, responsive, keyboard/accessibility, confirmation, dan telemetry state.

---

# 8. Architecture and Runtime

## Target topology

```text
TLS/load balancer
        |
Stateless GoFiber/Inertia API instances
        |
PostgreSQL — Valkey/lock/cache — S3-compatible object storage
        |
Workers: messaging | inbound | webhook | campaign | scheduler | media | analytics
        |
Official adapters | whatsmeow experimental | Telegram | SMS | Email
```

## Feature slice

```text
/internal/features/<feature>
  endpoint.go
  page.go
  service.go
  repository.go
  model.go
  dto.go
  mapper.go
  policy.go
  event.go
  *_test.go
/internal/integrations/<provider>
/internal/platform/<database|http|queue|storage|crypto|observability>
/resources/js/Pages
/resources/js/Components
/resources/js/composables
```

Current repository dapat memakai layout legacy selama boundary dan migration plan jelas. Refactor directory bukan prerequisite sebelum bug/security fixes.

## Sync vs async

Synchronous: auth, CRUD, validation, enqueue, read status.

Asynchronous: provider send, inbound normalization, campaign fan-out, delivery, webhook dispatch, media, analytics, email, notification.

## Scaling path

- Pilot: API + worker sederhana.
- Growth: worker terpisah per concern dan per-provider lock.
- Scale: partition event/message, durable stream, per-tenant limiter, dedicated provider workers, regional deployment bila perlu.

---

# 9. Security, Compliance, Reliability

## Mandatory controls

- CSRF yang benar: token tidak boleh memakai cookie yang sama sebagai source dan extractor.
- CORS allowlist.
- Secure/HttpOnly/SameSite cookies.
- Session regeneration setelah login/elevated action.
- CSP/security headers aktif.
- HTTPS/TLS verification.
- Rate limit, timeout, body/media/concurrency limit.
- Input validation dan canonicalization.
- SSRF block untuk private/loopback/link-local/metadata IP.
- Scoped, hashed/encrypted, one-time credential.
- Revoke/expiry/status/soft-delete check.
- Tenant/resource authorization di policy dan repository.
- HMAC webhook/replay protection.
- Attachment scan/size/type validation.
- Consent/opt-out/quiet hours.
- Structured logs dengan PII/secret redaction.
- SAST/dependency/image/vulnerability scan dan SBOM.

## Reliability targets

| Area | Pilot | GA target |
|---|---:|---:|
| API availability | 99.5% | 99.9% |
| API enqueue p95 | <500 ms | <300 ms |
| Webhook eventual success | 98% | 99.5% |
| Status lag | <60 s | <30 s |
| RPO | 24 h | 1 h atau lebih baik sesuai biaya |
| RTO | 8 h | 2 h atau lebih baik sesuai biaya |

Target ini tidak menjamin delivery provider.

## Retention proposal

- Message content: 90 hari setelah resolved.
- Delivery event: 180 hari.
- Raw webhook: 7–30 hari encrypted.
- Audit: minimal 90 hari free, 1 tahun paid atau sesuai kontrak.
- Media: lifecycle/quota/cleanup.
- Credential: selama channel aktif dan purge setelah disconnect/delete.
- Unofficial session: hapus pada disconnect/delete sesuai legal policy.

Final retention harus melalui legal/privacy review.

---

# 10. Migration Plan: Current → Full Feature

## Track 0 — Correctness/security blockers

- Implement/disable Telegram unimplemented methods tanpa panic.
- Perbaiki CSRF.
- Perbaiki device token revocation/expiry/deleted predicate.
- Validasi ownership device token/API key/device.
- Perbaiki provider dispatch berdasarkan channel type.
- Hilangkan panic vectors dan nil dereference.
- Perbaiki session fixation/logout errors.
- Fail-fast config.
- Hapus credential dari URL.
- Tambahkan validation, limits, timeout, SSRF/TLS protection.

## Track 1 — Platform foundation

- Workspace/tenant boundary.
- Public ID/DTO mapping.
- Auth/MFA/session.
- RBAC/resource policy.
- Audit event.
- Error envelope.
- OpenAPI contract.
- Outbox/idempotency/job model.
- Storage abstraction.

## Track 2 — Channel abstraction

- Define capability-specific adapter interfaces.
- Official WhatsApp adapter sebagai default.
- Pindahkan whatsmeow ke isolated experimental worker.
- Stabilkan Telegram adapter.
- Tambahkan SMS/Email adapters.
- Provider event normalization.
- Connection state machine, lock, reconnect, stop hooks.

## Track 3 — Core customer operations

- Contact/consent/segments/import.
- Conversation/message model.
- Unified inbox.
- Assignment/SLA/notes.
- Message lifecycle/delivery.
- Media pipeline.

## Track 4 — Productivity and growth

- Template approval/versioning.
- Campaign/broadcast guardrails.
- Scheduling/recurring/follow-up.
- Automation builder.
- Webhook/integration portal.
- API/SDK/sandbox.

## Track 5 — Business platform

- Analytics/reporting.
- Usage ledger/quota.
- Billing/subscription/invoice.
- AI assist with guardrails.
- Admin/support console.
- Data export/delete/compliance.

## Track 6 — Production excellence

- CI/CD/security gates.
- Unit/integration/contract/E2E/load tests.
- Observability/SLO/alerts.
- Backup/restore/DR drill.
- Multi-instance and provider failure exercise.
- Runbook, support process, status page.

Tidak ada track yang dianggap selesai hanya karena halaman UI sudah muncul. Minimal backend, authorization, worker, persistence, tests, observability, docs, dan recovery harus ada.

---

# 11. Consolidated Backlog Epics

| Epic | Outcome |
|---|---|
| A — Workspace/Auth | Identity, workspace, member, session, MFA |
| B — Authorization | RBAC, custom role, resource scope, tenant isolation |
| C — Channel Platform | Official/unofficial/provider lifecycle dan health |
| D — Contacts/Consent | Contact, import, tags, segments, opt-out |
| E — Inbox | Conversation, assignment, notes, SLA, human handoff |
| F — Messaging | All supported message types, queue, lifecycle, delivery |
| G — Templates | Version, approval, variable, localization |
| H — Campaign | Audience, approval, schedule, throttle, report |
| I — Automation/AI | Trigger/condition/action, run history, assist/handoff |
| J — Developer | API, key, SDK, OpenAPI, sandbox, webhook |
| K — Media | Upload, scan, storage, signed access, retention |
| L — Analytics | Operational/business metrics and reports |
| M — Billing | Plan, quota, usage, provider fee, invoice |
| N — Admin/Compliance | Audit, support, export/delete, policy |
| O — Platform Ops | Queue, observability, CI/CD, backup, DR |

Setiap story menggunakan template:

```markdown
## [FEATURE] Judul

### User story
Sebagai ..., saya ingin ..., sehingga ...

### Scope
...

### Acceptance criteria
- [ ] Happy path
- [ ] Validation
- [ ] Authorization/tenant isolation
- [ ] Retry/idempotency/recovery
- [ ] Observability
- [ ] Tests
- [ ] Documentation
```

---

# 12. Full-Feature Acceptance Gate

Sapasora disebut full-feature complete hanya jika:

- [ ] Semua epic A–O memiliki implementation dan owner.
- [ ] WhatsApp official menjadi production default.
- [ ] Whatsmeow tidak dapat aktif tanpa feature flag/risk acknowledgement/kill switch.
- [ ] Telegram tidak memiliki `panic("unimplemented")` pada route aktif.
- [ ] SMS dan Email memiliki adapter, capability, delivery event, dan compliance flow.
- [ ] Tidak ada cross-tenant access pada API, Inertia, repository, worker, cache, storage, atau webhook.
- [ ] Credential tidak plaintext/URL/list response dan dapat revoke/rotate.
- [ ] CSRF, SSRF, CORS, CSP, TLS, rate limit, timeout, body limit, validation sudah diuji.
- [ ] Contact/consent/suppression dihormati pada semua send path.
- [ ] Send, campaign, retry, webhook, dan automation idempotent.
- [ ] Queue, outbox, retry, DLQ, replay, dan graceful shutdown tersedia.
- [ ] UI routes tidak lagi placeholder dan memiliki state lengkap.
- [ ] Public API dan OpenAPI/Scalar tidak drift.
- [ ] Media tidak bergantung pada filesystem container lokal.
- [ ] Usage, quota, analytics, billing, audit, retention, export/delete tersedia.
- [ ] Unit/integration/contract/E2E/load/security/accessibility test tersedia.
- [ ] Backup restore, DR, incident, provider outage, dan rollback sudah didrill.
- [ ] SLO, dashboard, alerting, runbook, support, dan status page tersedia.

---

# 13. Open Decisions yang Harus Ditutup

Keputusan berikut berasal dari discovery/planning dan belum boleh diasumsikan selesai:

1. Direct Meta Tech Provider atau BSP sebagai jalur official pertama.
2. Nama/brand final dan domain.
3. Official WhatsApp onboarding yang tersedia untuk test asset.
4. Kapan whatsmeow masuk closed beta dan hasil legal review-nya.
5. Segment pilot dan minimum willingness-to-pay.
6. Billing ownership Meta/provider dan model pass-through.
7. Valkey queue semantics vs managed durable queue.
8. Auth/session provider vs internal implementation.
9. Payment/tax/invoice provider.
10. Retention final per plan/industry.
11. Migration strategy dari GORM ke sqlx, jika memang diperlukan.
12. Migration runner canonical: goose atau hardened existing runner.
13. Public ID migration strategy untuk tabel existing.
14. AI provider, data boundary, dan default assist-only policy.
15. Requirement multi-region.

Setiap keputusan harus dibuat sebagai ADR, dengan context, alternatives, consequences, owner, dan follow-up.

---

# 14. Discovery and Planning Archive

Seluruh hasil dari `finalize-plan.zip` telah dimasukkan ke repository sebagai referensi:

- `docs/finalize-plan/01-product-brief.md`
- `docs/finalize-plan/02-discovery-findings.md`
- `docs/finalize-plan/03-branding-positioning.md`
- `docs/finalize-plan/04-personas-research.md`
- `docs/finalize-plan/05-product-scope-roadmap.md`
- `docs/finalize-plan/06-ux-information-architecture.md`
- `docs/finalize-plan/07-backlog-acceptance-criteria.md`
- `docs/finalize-plan/08-vsa-system-architecture.md`
- `docs/finalize-plan/09-technical-stack.md`
- `docs/finalize-plan/10-domain-data-id-contract.md`
- `docs/finalize-plan/11-api-webhook-provider-contract.md`
- `docs/finalize-plan/12-security-compliance-reliability.md`
- `docs/finalize-plan/13-infrastructure-observability.md`
- `docs/finalize-plan/14-pricing-billing-unit-economics.md`
- `docs/finalize-plan/15-delivery-testing-qa.md`
- `docs/finalize-plan/16-metrics-launch-plan.md`
- `docs/finalize-plan/17-risk-decisions-open-questions.md`
- `docs/finalize-plan/18-team-handoff-checklist.md`
- `docs/finalize-plan/appendix/`

File-file tersebut tidak dihapus atau diringkas agar konteks discovery dan planning asli tetap dapat ditelusuri.

---

## Kesimpulan

Semua fitur penting dari full-feature specification dan finalize plan sekarang berada di bawah satu arah produk:

- **Sapasora** sebagai brand/product.
- **Official WhatsApp first** untuk production value.
- **Omnichannel adapter architecture** untuk Telegram, SMS, Email, dan future channel.
- **Inbox + notification + API + automation + campaign + analytics + billing** sebagai satu platform.
- **Security, consent, delivery visibility, recovery, dan tenant isolation** sebagai non-negotiable.
- **Current implementation audit** sebagai migration backlog, bukan diabaikan.

Implementasi sebaiknya dimulai dari Track 0 dan Track 1, lalu dibangun secara vertical slice sampai seluruh acceptance gate full-feature terpenuhi.
