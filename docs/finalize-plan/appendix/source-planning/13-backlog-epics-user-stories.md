# 13 — Backlog Epics dan User Stories

## Epic A — Workspace dan authentication

- As a user, saya dapat register/login dengan aman.
- As an owner, saya dapat membuat workspace.
- As an owner, saya dapat mengundang team member.
- As an admin, saya dapat mengatur role.
- As an owner, saya dapat melihat audit log.

**Acceptance:** tenant isolation, session revoke, invite expiry, role enforcement.

## Epic B — Channel onboarding

- As a business owner, saya dapat melihat requirement channel sebelum mulai.
- As a user, saya dapat menghubungkan official channel.
- As a user, saya dapat melihat status channel dan next action.
- As a user, saya dapat menjalankan test message.
- As a user, saya dapat disconnect channel dan menghapus credential.
- As an advanced user, saya dapat mengaktifkan experimental whatsmeow channel setelah risk acknowledgement.

**Acceptance:** provider/capability visible, no raw secret leakage, health check, cancel/retry flow.

## Epic C — Contacts dan consent

- As a user, saya dapat menambah/edit contact.
- As a user, saya dapat import CSV dengan preview.
- As a user, saya dapat menambah tags/segments.
- As a user, saya dapat menyimpan sumber opt-in.
- As a recipient, saya dapat opt-out.
- As a user, saya dapat melihat suppression list.

**Acceptance:** normalized identifier, duplicate handling, consent audit, opt-out honored.

## Epic D — Inbox

- As an agent, saya dapat melihat conversation yang assigned.
- As a manager, saya dapat assign/reassign conversation.
- As an agent, saya dapat reply dari inbox.
- As an agent, saya dapat menambahkan internal note.
- As an agent, saya dapat tag dan resolve conversation.
- As a manager, saya dapat melihat conversation yang melewati SLA.

**Acceptance:** unread state, optimistic UI safe, collision handling, delivery result visible.

## Epic E — Messaging API

- As a developer, saya dapat membuat scoped API key.
- As a developer, saya dapat mengirim text/template message.
- As a developer, saya dapat menggunakan idempotency key.
- As a developer, saya dapat membaca status message.
- As a developer, saya dapat menerima signed webhook.
- As a developer, saya dapat menguji/replay webhook event.

**Acceptance:** OpenAPI, request ID, stable error envelope, signature verification, retry guidance.

## Epic F — Templates

- As a user, saya dapat melihat status approval template.
- As a user, saya dapat preview variables.
- As a user, saya dapat melihat alasan rejection.
- As a developer, saya dapat mengambil template lewat API.

## Epic G — Campaign safety

- As a user, saya dapat membuat draft campaign.
- As a user, saya dapat melihat audience preview.
- As a user, saya dapat melihat contact tanpa consent yang dikeluarkan.
- As a manager, saya dapat approve/pause/cancel campaign.
- As a user, saya dapat mengatur timezone/quiet hours/rate.
- As a user, saya dapat melihat delivery report.

## Epic H — Automation

- As a user, saya dapat memilih trigger.
- As a user, saya dapat membuat condition.
- As a user, saya dapat membuat action send/assign/tag.
- As a user, saya dapat melihat automation run.
- As a user, saya dapat pause automation.

## Epic I — Billing dan usage

- As an owner, saya dapat melihat usage real-time.
- As an owner, saya dapat melihat cost estimate.
- As an owner, saya dapat upgrade/downgrade.
- As an owner, saya dapat melihat invoice.
- As an owner, saya dapat mengatur overage cap.

## Epic J — Observability dan support

- As a developer, saya dapat mencari request berdasarkan request ID.
- As a support agent, saya dapat melihat diagnostic tanpa credential/message secret.
- As a user, saya dapat melihat status page.
- As a user, saya dapat mengirim support ticket dengan context otomatis.
- As an operator, saya dapat pause provider/campaign.

## Prioritization

| Priority | Scope |
|---|---|
| P0 | Workspace, official channel, test send, contacts/consent, inbox dasar, API key, send/status, signed webhook |
| P1 | Templates, campaign approval, automation dasar, billing, analytics, replay/runbook |
| P2 | Whatsmeow closed beta, agency mode, round-robin/SLA, AI assist, integrations |
| P3 | Omnichannel, enterprise SSO/SCIM, advanced attribution |

## Story writing rule

Setiap story wajib menjawab:

- siapa user-nya;
- job apa yang diselesaikan;
- provider/capability yang terlibat;
- acceptance criteria;
- error/recovery behavior;
- security/privacy impact;
- event analytics;
- docs/support impact.
