# Sapasora Full-Feature Product & Technical Specification

> **Status:** Target specification / blueprint
>
> **Dokumen utama:** [`PRODUCT_MASTER_PLAN.md`](./PRODUCT_MASTER_PLAN.md)
>
> **Dokumen terkait:** [`PRODUCTION_READINESS.md`](./PRODUCTION_READINESS.md)
>
> Dokumen ini mendefinisikan target Sapasora sebagai platform komunikasi omnichannel SaaS yang lengkap. Ini **bukan scope MVP** dan bukan klaim bahwa seluruh fitur sudah tersedia di source code saat ini. Status implementasi aktual tetap mengacu pada audit production readiness.

---

## 1. Visi Produk

Sapasora adalah platform untuk mengelola komunikasi bisnis melalui berbagai channel dari satu workspace:

- WhatsApp
- Telegram
- SMS
- Email
- Channel tambahan melalui adapter/provider di masa depan

Platform harus mendukung komunikasi manual, otomatis, terjadwal, berulang, berbasis event, dan berbasis campaign tanpa mengorbankan keamanan, auditability, reliability, maupun kontrol tenant.

### Prinsip produk

1. **Omnichannel by design** — fitur inti tidak mengunci bisnis pada satu provider.
2. **Tenant isolation first** — tidak ada object yang dapat diakses lintas workspace tanpa grant eksplisit.
3. **Provider independent** — provider dapat diganti tanpa mengubah kontrak domain.
4. **Observable by default** — setiap message, delivery, automation, webhook, dan perubahan akses dapat dilacak.
5. **Secure secret handling** — token, credential, session, dan media sensitif tidak bocor ke response atau log.
6. **Reliable delivery** — pengiriman bersifat asynchronous, idempotent, retryable, dan memiliki dead-letter path.
7. **Human-in-the-loop** — automation tidak menghilangkan kemampuan agent untuk mengambil alih percakapan.
8. **Accessible and localizable** — UI mendukung accessibility, timezone, locale, dan format regional.

---

## 2. Batasan dan Definisi

### 2.1 Workspace

Workspace adalah boundary tenant utama. Semua data bisnis harus memiliki `workspace_id` atau relasi owner yang tervalidasi ke workspace.

### 2.2 User

User adalah identitas manusia yang dapat menjadi anggota satu atau beberapa workspace dengan role berbeda.

### 2.3 Channel account

Channel account adalah koneksi provider milik workspace, misalnya nomor WhatsApp, bot Telegram, sender SMS, atau mailbox Email.

### 2.4 Contact

Contact adalah entitas pelanggan/prospect yang dimiliki workspace. Satu contact dapat memiliki beberapa address/identifier pada channel yang berbeda.

### 2.5 Conversation

Conversation adalah thread komunikasi dengan satu contact pada satu channel atau unified thread jika workspace mengaktifkan identity resolution lintas channel.

### 2.6 Message

Message adalah unit komunikasi yang dikirim atau diterima. Delivery, provider event, retry, dan automation run memiliki lifecycle terpisah dari message.

---

## 3. Target Pengguna dan Role

### 3.1 Persona

| Persona | Kebutuhan utama |
|---|---|
| Workspace owner | Mengelola subscription, anggota, credential, channel, dan seluruh data |
| Administrator | Mengatur anggota, role, policy, template, automation, dan audit |
| Supervisor | Memantau agent, inbox, campaign, dan analytics |
| Agent/operator | Menjawab conversation dan mengirim message sesuai assignment |
| Campaign manager | Mengelola audience, template, schedule, dan campaign |
| Developer/integrator | Menggunakan API, webhook, SDK, dan API key |
| Auditor | Read-only access ke audit log, delivery, dan perubahan konfigurasi |
| End customer/contact | Menerima atau mengirim komunikasi melalui provider |

### 3.2 Role bawaan

- `owner`
- `admin`
- `supervisor`
- `agent`
- `campaign_manager`
- `developer`
- `auditor`
- `viewer`

Role harus dapat dicustom melalui permission matrix tanpa mengandalkan nama role di business logic.

### 3.3 Permission domain

Permission minimal:

- `workspace.read`, `workspace.update`, `workspace.delete`
- `member.read`, `member.invite`, `member.update`, `member.remove`
- `role.read`, `role.create`, `role.update`, `role.delete`
- `channel.read`, `channel.create`, `channel.update`, `channel.connect`, `channel.disconnect`, `channel.delete`
- `credential.create`, `credential.read`, `credential.revoke`, `credential.rotate`
- `contact.read`, `contact.create`, `contact.update`, `contact.delete`, `contact.import`, `contact.export`
- `conversation.read`, `conversation.assign`, `conversation.update`, `conversation.archive`
- `message.read`, `message.send`, `message.delete`, `message.retry`
- `template.read`, `template.create`, `template.update`, `template.publish`, `template.delete`
- `campaign.read`, `campaign.create`, `campaign.launch`, `campaign.pause`, `campaign.cancel`
- `automation.read`, `automation.create`, `automation.publish`, `automation.disable`
- `webhook.read`, `webhook.create`, `webhook.rotate`, `webhook.disable`
- `analytics.read`, `audit.read`, `billing.read`, `billing.manage`

Setiap permission harus diperiksa bersama `workspace_id`, resource owner, dan scope API key.

---

# 4. Modul Fitur Lengkap

## 4.1 Identity, Authentication, dan Workspace

### Fitur

- Registrasi user dengan email dan password.
- Login, logout, session management, remember-me.
- Email verification.
- Forgot password dan reset password.
- Password change dan password history.
- MFA/2FA menggunakan TOTP dan recovery codes.
- SSO melalui OIDC/SAML sebagai enterprise feature.
- Login melalui magic link sebagai opsi.
- Device/session list dan remote logout.
- Deteksi login mencurigakan.
- Rate limit login dan protection terhadap credential stuffing.
- Workspace creation, rename, logo, timezone, locale, dan default settings.
- Workspace switcher untuk user multi-workspace.
- Invite, accept, expire, resend, dan revoke invitation.
- Suspend/unsuspend member.
- Ownership transfer dengan confirmation dan audit.
- Workspace deletion dengan grace period dan purge job.

### Acceptance criteria

- Login selalu meregenerasi session ID.
- Logout/revoke menghentikan session yang bersangkutan.
- Password reset token one-time, short-lived, hashed, dan tidak masuk log.
- Perubahan role/password dapat merevoke session sesuai policy.
- User tidak dapat membaca atau mengubah workspace lain melalui ID yang diketahui.
- Semua perubahan membership tercatat di audit log.

---

## 4.2 Authorization, Role, dan Policy

### Fitur

- RBAC berbasis workspace.
- Custom role.
- Permission inheritance yang eksplisit.
- Resource-level access untuk channel, contact list, campaign, dan conversation.
- API key scope.
- Service account untuk integrasi.
- Temporary elevated access dengan expiry.
- Approval workflow untuk tindakan berisiko tinggi.
- Policy untuk export data, bulk send, delete, dan credential rotation.

### Invariant

- Authorization harus diputuskan di server, bukan hanya di UI.
- Query repository wajib tenant-scoped.
- Object ID publik tidak boleh dianggap sebagai bukti ownership.
- API key tidak boleh memperoleh akses di luar scope yang diberikan.

---

## 4.3 Channel dan Provider Management

### Fitur umum

- Add/connect channel account.
- QR pairing atau OAuth/provider authorization.
- Credential/session encryption at rest.
- Connection status: `pending`, `connecting`, `connected`, `degraded`, `disconnected`, `expired`, `error`.
- Health check channel.
- Reconnect dengan backoff.
- Disconnect dan revoke.
- Delete channel dengan confirmation dan cleanup asynchronous.
- Channel aliases dan display name.
- Per-channel timezone dan sending policy.
- Per-channel rate limit dan daily quota.
- Provider capability discovery.
- Provider error mapping ke error code internal.
- Failover provider jika dikonfigurasi.
- Maintenance mode untuk channel.
- Connection history.

### WhatsApp

- QR login/pairing.
- Session lifecycle yang durable.
- Incoming text, image, audio, video, document, sticker, location, contact, reaction, reply, and button/list event bila provider mendukung.
- Outgoing text dan media.
- Group support bila provider mengizinkan.
- Read receipt, delivery receipt, typing/presence.
- Contact sync yang dapat diaktifkan atau dimatikan.
- Webhook/event mapping.
- Per-device session ownership dan distributed lock.

### Telegram

- Bot token flow.
- User/TDLib flow bila memang diperlukan oleh product policy.
- Contact/user lookup.
- Text, image, audio, video, document, sticker, location, contact, button, list, and presence sesuai capability.
- Update polling/webhook lifecycle.
- Flood wait handling.
- Rate-limit-aware queue.
- Correct connect/disconnect state.
- Error propagation dari goroutine/client.

### SMS

- Provider abstraction: Twilio-compatible dan provider lokal.
- Sender number/alphanumeric sender ID sesuai regulasi.
- Delivery receipt.
- Country routing.
- Opt-out STOP/START handling.
- Segment count dan cost estimation.
- Character encoding detection.
- Quiet hours.
- Per-country compliance rules.

### Email

- SMTP provider dan API provider.
- From identity, reply-to, CC/BCC policy.
- HTML/text multipart.
- Attachments.
- Bounce, complaint, delivery, open, click event bila provider mendukung.
- DKIM/SPF/DMARC guidance.
- Suppression list.
- Unsubscribe group.
- Template preview dan test send.

### Provider adapter contract

Setiap adapter harus menyediakan:

```text
Connect
Disconnect
Health
Capabilities
Send
GetContact/User
GetConversations (jika tersedia)
ReceiveEvent
NormalizeError
NormalizeInboundEvent
ValidateAddress
```

Adapter tidak boleh dipanggil langsung oleh controller. Domain service memilih adapter berdasarkan channel account dan capability.

---

## 4.4 Device / Channel Account Management

### Fitur

- List, search, filter, sort, pagination.
- Create/connect/edit/delete.
- QR display dan QR expiration.
- Connection logs.
- Manual reconnect.
- Auto-connect toggle.
- Default channel selection.
- Webhook configuration.
- API exposure settings.
- Inbound/outbound enablement.
- Quota and usage display.
- Last connected, last event, last error.
- Graceful disconnect before deletion.
- Credential/token revoke on deletion.
- Export-safe DTO tanpa secret.

### State transition

```text
created -> pending
pending -> connecting
connecting -> connected
connecting -> error
connected -> degraded
connected -> disconnected
error -> pending
any active state -> deleting -> deleted
```

Setiap transition harus idempotent dan dicatat.

---

## 4.5 Contact, Phonebook, dan Audience

### Fitur contact

- Create/update/delete/restore.
- Soft delete dan permanent purge sesuai retention policy.
- Multiple phone numbers, email addresses, usernames, dan channel identifiers.
- Custom fields dengan schema per workspace.
- Tags, segments, notes, owner, source, lifecycle stage.
- Consent status per channel.
- Opt-in/opt-out history.
- Suppression/global block list.
- Duplicate detection dan merge.
- Import CSV/XLSX dengan preview dan mapping.
- Import asynchronous dengan error report.
- Export permission-gated dan audited.
- Bulk update/tag/delete.
- Search by name, identifier, custom field, tag, and last interaction.
- Contact timeline.
- Identity linking lintas channel dengan merge confirmation.

### Audience dan segment

- Static list.
- Dynamic segment berdasarkan field, tag, event, delivery status, dan interaction time.
- Segment preview count.
- Segment snapshot saat campaign launch.
- Exclusion list dan suppression check sebelum send.

### Compliance

- Consent source, timestamp, actor, and evidence.
- Per-channel opt-out.
- Data export request.
- Data deletion request.
- Retention and anonymization.
- Do-not-contact list.

---

## 4.6 Unified Inbox dan Conversation

### Fitur

- Conversation list dengan search, filter, sort, pagination.
- Filter channel, status, assignee, team, tag, SLA, unread, and last activity.
- Conversation detail dengan message timeline.
- Reply composer multi-format.
- Attach media.
- Internal note yang tidak dikirim ke contact.
- Assignment ke agent/team.
- Claim/unclaim.
- Status: `open`, `pending`, `snoozed`, `resolved`, `closed`.
- Priority.
- Labels/tags.
- Customer profile sidebar.
- Conversation merge/split dengan audit.
- Read/unread.
- Typing indicator bila provider mendukung.
- Canned response.
- Quote/reply/thread.
- Search message history.
- SLA timer dan breach indicator.
- Collision detection saat dua agent membalas.
- Supervisor takeover.
- Export conversation sesuai permission.

### Message timeline

Timeline harus membedakan:

- Inbound message.
- Outbound message.
- Delivery status.
- Internal note.
- Assignment event.
- Automation event.
- Template event.
- Webhook/provider event.
- System event.

---

## 4.7 Messaging dan Delivery

### Message types

- Plain text.
- Markdown/rich text yang dinormalisasi per channel.
- Image.
- Audio.
- Video.
- Document.
- Sticker.
- Location.
- Contact card.
- Button/action.
- List/menu.
- Template message.
- Reaction.
- Reply/quoted message.
- Batch/bulk message.

### Lifecycle message

```text
accepted -> queued -> sending -> sent
                         |-> failed
sent -> delivered -> read
sent -> unknown
failed -> retrying -> sent|failed|dead_letter
```

### Fitur reliability

- Client-supplied idempotency key.
- Server-generated message ID.
- Outbox pattern.
- Queue partition per channel account/contact.
- Ordering guarantee yang terdokumentasi.
- Exponential backoff dengan jitter.
- Retry hanya untuk error transient.
- Dead-letter queue.
- Manual retry dengan permission.
- Delivery event deduplication.
- Provider message ID mapping.
- Rate-limit-aware scheduling.
- Per-contact quiet hours.
- Message cancellation sebelum provider accept.
- Cost/usage tracking.

### Safety

- Preview sebelum send.
- Confirmation untuk bulk send.
- Maximum recipient guard.
- Template/consent validation.
- Blocked contact check.
- Channel capability validation.
- Media malware/size/type validation.
- Secret redaction pada logs.

---

## 4.8 Template Management

### Fitur

- Template create/edit/duplicate/archive.
- Draft, review, approved, rejected, archived.
- Versioning dan rollback.
- Per-channel variant.
- Variables dengan type, required/optional, default, and example.
- Preview dengan sample data.
- Validation terhadap unsupported variable/type.
- Approval workflow.
- Template category dan tags.
- Localization per language.
- WhatsApp provider template mapping.
- Email HTML/text variant.
- SMS length/segment preview.
- Usage statistics.
- Template deprecation.

### Variable safety

- Escape sesuai channel.
- Tidak boleh memasukkan secret ke variable.
- PII masking pada preview/log.
- Missing variable harus menghasilkan validation error sebelum queue.

---

## 4.9 Campaign dan Broadcast

### Fitur campaign

- Campaign draft.
- Audience selection.
- Audience snapshot.
- Channel selection.
- Template selection.
- Personalization preview.
- Test send.
- Schedule dengan timezone.
- Throttling/rate limit.
- Quiet hours.
- Frequency cap.
- Exclusion and suppression.
- Approval workflow.
- Launch, pause, resume, cancel.
- Partial failure handling.
- A/B variant.
- Batch progress.
- Cost estimate.
- Delivery report.
- Export result.
- Campaign clone.

### Campaign states

```text
draft -> awaiting_approval -> approved -> scheduled
scheduled -> running -> paused
running -> completed
running -> failed
paused -> running|cancelled
any pre-completion state -> cancelled
```

### Guardrails

- Confirm recipient count.
- Confirm estimated cost.
- Require approval above configurable threshold.
- Prevent duplicate launch.
- Idempotent batch creation.
- Persist recipient-level result.
- Respect unsubscribe immediately.

---

## 4.10 Scheduled dan Recurring Message

### Fitur

- One-time scheduled message.
- Recurring schedule: daily, weekly, monthly, custom cron-like rule.
- Timezone per schedule.
- Start/end date.
- Pause/resume/cancel.
- Holiday/blackout calendar.
- Missed-run policy.
- Audience resolved at run time atau snapshot, configurable.
- Per-run delivery report.
- Retry policy.
- Schedule history.
- Duplicate execution prevention.

### Scheduler requirements

- Distributed lock.
- Durable job store.
- At-least-once execution dengan idempotency.
- Clock drift handling.
- Graceful shutdown.
- Retry/dead-letter.
- Admin replay dengan audit.

---

## 4.11 Autoresponder dan Automation Builder

### Trigger

- Inbound message.
- Keyword/regex.
- New contact.
- Contact tag added/removed.
- Conversation opened/resolved.
- Delivery failure.
- Scheduled time.
- Webhook event.
- Campaign event.
- Contact field changed.

### Condition

- Channel.
- Device/channel account.
- Contact field.
- Tag/segment.
- Message content.
- Business hours.
- Conversation state.
- Previous automation result.
- Rate/limit condition.

### Action

- Send message/template.
- Add/remove tag.
- Assign agent/team.
- Set conversation status/priority.
- Add internal note.
- Wait/delay.
- Call webhook.
- Create task.
- Update contact field.
- Stop other automation.
- Escalate to human.

### Fitur operasional

- Visual flow builder.
- Draft/publish/disable.
- Versioning.
- Test/simulation mode.
- Dry-run pada sample contact.
- Loop detection.
- Maximum execution depth.
- Per-contact concurrency lock.
- Execution log.
- Retry and dead-letter.
- Kill switch per automation/workspace.
- Permissioned publication.

---

## 4.12 Webhook dan Integrations

### Outgoing webhook

- Event subscription per workspace.
- Filter by channel, event type, resource.
- HTTPS enforcement.
- SSRF protection.
- Domain/IP allowlist.
- HMAC signature dengan key rotation.
- Timestamp dan replay protection.
- Idempotency event ID.
- Retry with exponential backoff.
- Delivery attempt log.
- Dead-letter.
- Pause/disable after repeated failure.
- Test delivery.
- Payload versioning.
- Secret redaction.

### Incoming webhook

- Provider signature verification.
- Timestamp/replay validation.
- Provider-specific parser.
- Duplicate event handling.
- Raw event quarantine jika parsing gagal.
- Fast acknowledgement dan asynchronous processing.
- Dead-letter/replay tooling.

### API integrations

- REST API versioning.
- OpenAPI specification yang sinkron dengan implementation.
- API keys/service accounts.
- OAuth application untuk third-party integration.
- Idempotency header.
- Pagination cursor.
- Filter/sort contract.
- Rate limit headers.
- Consistent error envelope.
- SDK generation untuk TypeScript, Go, dan PHP.

---

## 4.13 API Key, Device Token, dan Developer Platform

### Fitur

- Create credential dengan label dan owner.
- Scope permission.
- IP allowlist optional.
- Expiry optional/required sesuai policy.
- One-time secret display.
- Hash-at-rest.
- Rotate tanpa downtime.
- Revoke.
- Last used metadata.
- Usage analytics.
- Audit log.
- Environment separation: test/live.
- Per-workspace quota.

### Kontrak keamanan

- Credential tidak pernah dikirim melalui URL path.
- Credential tidak pernah dikembalikan pada list/get setelah creation.
- Revoked, expired, deleted, dan disabled credential selalu ditolak.
- Lookup harus tenant-aware dan mengikutkan semua soft-delete/status predicate.
- Error authentication tidak membedakan credential tidak ada vs credential revoked.

---

## 4.14 Media dan File Management

### Fitur

- Upload direct atau presigned upload.
- Resumable upload untuk file besar.
- MIME sniffing dan extension validation.
- Virus/malware scanning.
- Image/video metadata handling.
- Thumbnail/transcoding bila diperlukan.
- Encryption at rest.
- Signed download URL dengan expiry.
- Access check pada setiap download.
- Storage quota per workspace.
- Retention policy.
- Orphan cleanup.
- Deduplication optional.
- Provider-specific media conversion.

### Storage policy

- Media tidak boleh bergantung pada filesystem container lokal.
- Path harus canonicalized dan bebas path traversal.
- `os.ModePerm` tidak boleh digunakan untuk data sensitif.
- File download harus memiliki content disposition dan content type yang aman.

---

## 4.15 Analytics, Reporting, dan Billing

### Analytics

- Messages sent/received.
- Delivery/read/failure rate.
- Response time.
- First response time.
- Resolution time.
- Agent workload.
- Campaign performance.
- Template performance.
- Channel health.
- Automation success/failure.
- Webhook delivery rate.
- Cost by workspace/channel/campaign.

### Reporting

- Dashboard realtime near-real-time.
- Date/timezone filter.
- Comparison period.
- CSV export.
- Scheduled report email.
- Role-based visibility.
- Pre-aggregated metrics untuk skala besar.

### Billing dan usage

- Plan, feature entitlement, dan quota.
- Seat count.
- Message/channel/media usage.
- Overage policy.
- Usage alert.
- Invoice integration.
- Payment provider integration.
- Trial, pause, cancel, grace period.
- Billing audit.
- Data retention setelah cancellation.

Billing harus dipisahkan dari authorization core sehingga workspace tetap dapat dibaca untuk support dan recovery.

---

## 4.16 Admin, Audit, dan Support

### Admin console

- Workspace search.
- User/member search.
- Channel health.
- Job/queue health.
- Delivery lookup.
- Webhook replay.
- Dead-letter inspection.
- Feature flag.
- Rate limit override dengan expiry.
- Support impersonation dengan consent, banner, duration, dan audit.
- Maintenance mode.

### Audit log

Audit event minimal memiliki:

- `id`
- `workspace_id`
- `actor_type`
- `actor_id`
- `action`
- `resource_type`
- `resource_id`
- `before` dan `after` yang telah disamarkan
- `request_id`
- `ip`
- `user_agent`
- `created_at`

Audit log append-only, immutable bagi user biasa, searchable, exportable, dan memiliki retention policy.

---

# 5. Information Architecture Frontend

## Route target

### Public/auth

- `/login`
- `/register`
- `/verify-email`
- `/forgot-password`
- `/reset-password`
- `/mfa`

### Workspace

- `/workspaces`
- `/workspace/settings`
- `/workspace/members`
- `/workspace/roles`
- `/workspace/audit`
- `/workspace/billing`
- `/workspace/usage`

### Channels

- `/devices`
- `/devices/new`
- `/devices/:id`
- `/devices/:id/edit`
- `/devices/:id/connection`
- `/devices/:id/logs`

### CRM/inbox

- `/phonebook`
- `/phonebook/import`
- `/phonebook/segments`
- `/inbox`
- `/inbox/:conversationId`
- `/inbox/assigned`
- `/inbox/unassigned`

### Messaging

- `/messages`
- `/messages/compose`
- `/messages/scheduled`
- `/messages/failed`
- `/templates`
- `/templates/new`
- `/templates/:id`

### Automation/campaign

- `/campaigns`
- `/campaigns/new`
- `/campaigns/:id`
- `/recurring`
- `/autoresponder`
- `/autoresponder/:id/builder`

### Developer/integrations

- `/integrations`
- `/integrations/webhooks`
- `/integrations/api-keys`
- `/documentation`

### Operations

- `/analytics`
- `/reports`
- `/settings`
- `/notifications`
- `/support`

Setiap route membutuhkan loading state, empty state, error state, permission guard, responsive layout, keyboard navigation, dan confirmation untuk tindakan destruktif.

---

# 6. Domain Model Target

Entitas inti:

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
ContactTag
ContactSegment
ConsentRecord
Conversation
ConversationParticipant
ConversationAssignment
Message
MessageAttachment
MessageDelivery
MessageReaction
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

### Aturan data

- Semua tabel tenant-scoped memiliki `workspace_id` atau jalur relasi yang tervalidasi.
- Semua ID eksternal provider memiliki unique constraint bersama channel account.
- Soft delete tidak boleh menghapus object dari authorization predicate.
- Status transition divalidasi di domain service.
- Timestamp disimpan UTC; timezone hanya untuk presentation/scheduling.
- PII dan secret dipisahkan dari data operasional bila memungkinkan.
- Data purge harus asynchronous, resumable, dan audited.

---

# 7. API dan Event Contract

## Response envelope

```json
{
  "data": {},
  "meta": {
    "request_id": "req_..."
  }
}
```

Error:

```json
{
  "error": {
    "code": "MESSAGE_VALIDATION_FAILED",
    "message": "Message payload is invalid",
    "details": [],
    "request_id": "req_..."
  }
}
```

`message` tidak boleh berisi SQL error, URL credential, secret, stack trace, atau detail internal.

## API requirements

- Version prefix, misalnya `/api/v1`.
- Cursor pagination untuk collection besar.
- `Idempotency-Key` pada create/send/launch/retry.
- `ETag`/conditional request untuk resource yang diedit bersama.
- Rate-limit headers.
- Request ID diterima dari client hanya jika divalidasi; server tetap membuat ID aman.
- OpenAPI security scheme harus sesuai dengan authentication aktual.
- API docs private atau dilindungi authentication pada production.

## Event naming

```text
workspace.member.created
channel.connected
channel.disconnected
contact.created
contact.updated
conversation.created
conversation.assigned
message.accepted
message.sent
message.delivered
message.read
message.failed
campaign.started
campaign.completed
automation.started
automation.failed
webhook.delivery.failed
```

Event harus memiliki `event_id`, `event_type`, `version`, `occurred_at`, `workspace_id`, `actor`, dan resource reference.

---

# 8. Background Processing dan Reliability

## Komponen wajib

- Outbox processor.
- Message send queue.
- Inbound event queue.
- Webhook delivery queue.
- Campaign batch queue.
- Scheduler.
- Media processing queue.
- Cleanup/retention worker.
- Analytics aggregation worker.
- Dead-letter worker/tooling.

## Requirement worker

- Graceful shutdown.
- Context cancellation yang benar.
- Bounded concurrency.
- Retry policy per error class.
- Visibility timeout/lease.
- Idempotency.
- Metrics per queue/job.
- Backpressure.
- Poison-message handling.
- Replay dengan audit.
- Tidak kehilangan job ketika process restart.

## Provider connection manager

- Satu logical owner untuk satu channel session.
- Distributed lock bila multi-instance.
- Stop hook pada shutdown.
- Reconnect policy.
- State machine yang tidak dapat lompat sembarangan.
- Error goroutine dikirim ke supervisor.
- Connection state dipersistenkan secara aman.

---

# 9. Security dan Compliance Baseline

## Application security

- CSRF token yang benar-benar secret dan tidak otomatis dikirim sebagai cookie yang sama.
- CORS allowlist.
- Secure, HttpOnly, SameSite cookie.
- Session regeneration setelah login/elevated action.
- CSP aktif dan diuji.
- HSTS pada HTTPS.
- Security headers lengkap.
- Request/body/read/write timeout.
- Media upload limit.
- Rate limit per IP, user, workspace, API key, dan provider.
- Validation dan canonicalization.
- SSRF protection pada URL outbound.
- TLS certificate verification aktif.
- Secret redaction.
- Dependency scanning dan SBOM.

## Data security

- Argon2id atau password hashing modern.
- Token HMAC/hash at rest.
- Encryption at rest untuk provider session dan credential.
- KMS/secret manager pada production.
- PII minimization.
- Data retention dan deletion.
- Export/delete request.
- Backup encryption.

## Messaging compliance

- Consent per channel.
- Opt-out segera.
- Suppression list.
- Quiet hours.
- Country/provider policy.
- Audit campaign approval.
- Unsubscribe link untuk email.
- STOP/START untuk SMS.
- Template approval untuk channel yang mewajibkan.

---

# 10. Observability dan Operasional

## Metrics

### HTTP/API

- Request count.
- Latency p50/p95/p99.
- Error rate.
- Status code.
- Rate-limit rejection.

### Messaging

- Accepted/sent/delivered/read/failed.
- Provider latency.
- Queue depth/age.
- Retry count.
- Dead-letter count.
- Duplicate event count.

### Channel

- Connected/disconnected count.
- Reconnect count.
- Provider error rate.
- Last event age.
- Session expiry.

### Business

- Active workspace.
- Active contact.
- Campaign completion.
- Agent response time.
- Usage/cost.

## Logging

- Structured JSON.
- Request ID, workspace ID, user ID, resource ID.
- Secret/PII redaction.
- No full webhook URL with credentials.
- No raw SQL bind values jika mengandung PII/secret.
- Retention dan access policy.

## Alerting

- Error rate spike.
- Queue backlog.
- Dead-letter growth.
- Provider outage.
- Webhook failure threshold.
- Login attack.
- Database pool exhaustion.
- Storage quota.
- Backup failure.
- Certificate expiry.

## Health endpoints

- `/livez` — process hidup.
- `/readyz` — dependency siap.
- `/health/details` — hanya admin/internal, dengan detail dependency.

---

# 11. Non-Functional Requirements

Target awal untuk deployment production, dapat disesuaikan setelah load test:

| Area | Target |
|---|---|
| API availability | 99.9% per bulan |
| API read latency | p95 < 300 ms tanpa provider call |
| API write acceptance | p95 < 500 ms sampai masuk queue |
| Webhook acknowledgement | p95 < 1 detik |
| Message queue durability | Tidak kehilangan accepted message |
| Delivery deduplication | Tidak ada duplicate akibat retry internal |
| Recovery Point Objective | <= 15 menit |
| Recovery Time Objective | <= 60 menit |
| Audit retention | Minimal 1 tahun, configurable |
| Availability data retention | Configurable per plan/compliance |
| Accessibility | WCAG 2.1 AA target |
| Browser support | Dua versi mayor browser modern |
| Localization | ID dan EN sebagai baseline |
| Timezone | UTC storage, workspace/user timezone presentation |

Target angka harus divalidasi melalui load test, bukan hanya konfigurasi.

---

# 12. Deployment, Backup, dan Disaster Recovery

## Environment

- Local development.
- Shared development.
- Staging.
- Production.

Setiap environment memiliki credential, database, provider session, webhook secret, dan storage yang terpisah.

## Production topology target

```text
Load Balancer / TLS
        |
Web/API instances (stateless)
        |
Queue / Worker instances
        |
PostgreSQL -- Redis/lock/cache -- Object Storage
        |
Provider adapters / external providers
```

## Requirements

- Container non-root.
- Immutable image.
- Pinned dependency/build.
- Secret injection dari secret manager.
- Readiness-based deployment.
- Rolling/blue-green deployment.
- Migration lock dan migration observability.
- Rollback plan.
- Database backup otomatis.
- Point-in-time recovery bila tersedia.
- Restore drill berkala.
- Object storage versioning/lifecycle.
- Centralized logs and metrics.
- Capacity plan untuk database, queue, storage, provider quota.

---

# 13. Testing Strategy

## Unit test

- Domain state transition.
- Authorization policy.
- Token hashing/verification.
- Validation.
- Template rendering.
- Segmentation.
- Retry classification.
- Idempotency.
- Timezone/scheduler.
- Provider error normalization.

## Integration test

- Database repository dengan tenant isolation.
- Authentication/session.
- API key/device token revocation.
- Queue/outbox.
- Webhook signature and retry.
- Storage access.
- Migration.
- Provider fake adapter.

## Contract test

- OpenAPI vs route implementation.
- Provider adapter capability contract.
- Webhook event schema.
- SDK generated clients.

## E2E test

- Register/login/reset/MFA.
- Workspace/member/role.
- Connect/disconnect channel.
- Contact import.
- Inbox reply.
- Media send.
- Template approval.
- Campaign launch/pause/cancel.
- Schedule/recurring.
- Automation publish/run.
- Webhook failure/retry.
- Credential rotate/revoke.
- Billing/quota.

## Non-functional test

- Load test API and queue.
- Provider throttling.
- Large media.
- Failover/restart.
- Backup restore.
- Security scan/SAST/DAST.
- Dependency vulnerability scan.
- Accessibility audit.
- Browser compatibility.

---

# 14. Current Repository Gap Map

Bagian ini menghubungkan target full-feature dengan kondisi yang ditemukan pada audit.

| Area | Kondisi saat audit | Target full feature |
|---|---|---|
| Telegram | Banyak method belum diimplementasikan | Adapter lengkap dengan lifecycle dan contract test |
| WhatsApp | Provider state, panic vector, storage, dan lifecycle perlu diperbaiki | Durable provider manager dan capability-complete adapter |
| SMS | Ada referensi UI, backend belum menjadi adapter lengkap | Provider abstraction, compliance, delivery receipt |
| Email | Belum menjadi channel production lengkap | SMTP/API adapter, bounce, suppression, unsubscribe |
| Auth | Login/logout tersedia, hardening belum lengkap | Verification, reset, MFA, session/device management |
| Authorization | Permission ada, tenant/resource scope belum konsisten | Workspace RBAC + resource policy + API scope |
| Credential | Secret plaintext dan revocation issue | Hashed, one-time reveal, rotate/revoke/audit |
| CSRF | Konfigurasi cookie tidak efektif | Secret token dan pipeline verification yang benar |
| Devices | CRUD UI dan cleanup provider belum lengkap | Full channel account lifecycle |
| Phonebook | Link UI belum memiliki route lengkap | Contact, consent, segment, import/export |
| Inbox | Belum menjadi unified inbox lengkap | Conversation, assignment, SLA, notes, takeover |
| Messages | Sebagian provider dispatch tidak konsisten | Queue, idempotency, retry, delivery lifecycle |
| Templates | UI/route belum lengkap | Versioning, approval, variables, localization |
| Campaigns | Belum tersedia penuh | Audience snapshot, approval, throttling, reporting |
| Recurring | Placeholder/route belum lengkap | Durable scheduler dan run history |
| Autoresponder | Placeholder/route belum lengkap | Visual automation engine |
| Webhook | SSRF/signing/retry/idempotency belum lengkap | Secure event integration platform |
| Media | Local path dan limit belum production-safe | Object storage, scan, signed URL, retention |
| Analytics | Metrics/analytics belum lengkap | Operational and business reporting |
| Billing | Belum menjadi domain lengkap | Entitlement, quota, usage, invoice |
| Admin | Route/feature belum lengkap | Admin/support/audit console |
| API docs | Swagger security mismatch/public exposure | Versioned authenticated developer portal |
| Workers | Scheduler/mail/cache/storage integration belum lengkap | Durable queue and worker platform |
| Tests | Unit/E2E belum tersedia, Go belum lulus | Full test pyramid and CI gates |
| Deployment | Docker/config masih development | Multi-instance production topology and DR |

---

# 15. Definition of Full Feature Complete

Sapasora baru dapat disebut **full-feature complete** apabila seluruh kategori berikut memiliki implementation, test, documentation, monitoring, dan operational runbook:

- [ ] Identity, MFA, session, workspace, invitation, dan account recovery.
- [ ] RBAC, custom role, resource authorization, API scope, dan tenant isolation.
- [ ] WhatsApp adapter lengkap.
- [ ] Telegram adapter lengkap.
- [ ] SMS adapter lengkap.
- [ ] Email adapter lengkap.
- [ ] Channel lifecycle, credential lifecycle, dan provider health.
- [ ] Contact/phonebook, consent, import/export, segment, dan merge.
- [ ] Unified inbox, assignment, notes, SLA, search, dan conversation lifecycle.
- [ ] Semua message type yang didukung provider dengan capability negotiation.
- [ ] Queue, outbox, idempotency, retry, dead-letter, dan delivery tracking.
- [ ] Template versioning, approval, localization, dan variable safety.
- [ ] Campaign/broadcast, audience snapshot, approval, throttling, dan report.
- [ ] Scheduled/recurring message dan durable scheduler.
- [ ] Autoresponder/automation builder, versioning, execution log, dan kill switch.
- [ ] Incoming/outgoing webhook dengan signing, SSRF protection, retry, dan replay.
- [ ] API v1, OpenAPI, SDK, API key, device token, dan developer portal.
- [ ] Media storage, scan, quota, signed access, retention, dan cleanup.
- [ ] Analytics, report, usage, quota, billing, dan notification.
- [ ] Admin/support console serta immutable audit trail.
- [ ] Security baseline, compliance workflow, backup, restore, dan DR.
- [ ] Unit, integration, contract, E2E, load, security, accessibility, dan migration test.
- [ ] CI/CD, observability, alerting, runbook, and on-call process.

---

# 16. Urutan Implementasi Full Feature

Urutan berikut bukan pemangkasan menjadi MVP. Ini dependency order agar seluruh feature dapat dibangun tanpa fondasi yang rapuh.

## Fase A — Platform foundation

- Workspace/tenant model.
- Authentication hardening dan MFA.
- Authorization/resource policy.
- Database conventions dan migrations.
- Error contract, request ID, audit event.
- Queue, outbox, idempotency, and worker foundation.
- Secret manager, storage abstraction, and configuration validation.

## Fase B — Provider platform

- Adapter interface dan capability model.
- WhatsApp lifecycle.
- Telegram lifecycle.
- SMS adapter.
- Email adapter.
- Connection manager, locks, reconnect, provider event normalization.

## Fase C — CRM dan inbox

- Contact/consent/segment.
- Conversation/message model.
- Unified inbox.
- Assignment, SLA, internal notes, search.
- Message delivery and retry.

## Fase D — Productivity and automation

- Template engine.
- Scheduled and recurring messages.
- Campaign/broadcast.
- Autoresponder/automation builder.
- Webhook integration.

## Fase E — Platform completeness

- Media pipeline.
- Analytics/reporting.
- Billing/usage/entitlement.
- Developer portal and SDK.
- Admin/support tools.
- Data export/delete and compliance workflows.

## Fase F — Production excellence

- Full test suite.
- Load/security/accessibility testing.
- Multi-instance deployment.
- Backup/restore/DR exercise.
- SLO dashboards and alerting.
- Release, rollback, and incident runbooks.

Setiap fase harus ditutup dengan acceptance test dan dokumentasi operasional; tidak ada fase yang dianggap selesai hanya karena UI sudah tampil.

---

# 17. Dokumen Turunan yang Wajib Dibuat

Spesifikasi ini perlu dipecah menjadi dokumen implementasi berikut:

- `ARCHITECTURE.md` — boundary service, adapter, queue, storage, dan deployment.
- `DOMAIN_MODEL.md` — entity, relationship, invariant, dan migration plan.
- `API_SPECIFICATION.md` — endpoint, schema, auth, pagination, error, dan idempotency.
- `PROVIDER_ADAPTER_CONTRACT.md` — capability dan lifecycle adapter.
- `SECURITY_MODEL.md` — trust boundary, threat model, secret, CSRF, SSRF, and tenant isolation.
- `EVENT_CATALOG.md` — event type, schema, versioning, dan replay.
- `OPERATIONS_RUNBOOK.md` — deploy, rollback, queue, provider outage, backup, restore.
- `TEST_PLAN.md` — unit, integration, contract, E2E, load, security, and accessibility.
- `COMPLIANCE_AND_RETENTION.md` — consent, opt-out, export, deletion, retention, and audit.
- `CHANGELOG_AND_MIGRATIONS.md` — compatibility and rollout strategy.

---

## Catatan Penutup

Dokumen ini mendefinisikan **produk lengkap yang ingin dibangun**, sedangkan [`PRODUCTION_READINESS.md`](./PRODUCTION_READINESS.md) mendefinisikan risiko dan gap dari implementasi saat ini. Keduanya harus digunakan bersama:

- `FULL_FEATURE_SPECIFICATION.md` menjawab **apa yang harus tersedia**.
- `PRODUCTION_READINESS.md` menjawab **mengapa implementasi sekarang belum siap dan apa risikonya**.

Sebuah fitur hanya dianggap selesai jika tersedia di domain, API, UI, background processing, authorization, observability, test, migration, dokumentasi, dan runbook operasionalnya.
