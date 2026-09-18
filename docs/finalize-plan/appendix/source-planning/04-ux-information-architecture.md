# 04 — UX, Information Architecture, dan User Flow

## UX goals

1. User baru tahu langkah berikutnya.
2. User dapat membedakan channel official/unofficial.
3. Status message mudah dipahami nonteknis.
4. Developer tetap dapat mengakses detail teknis.
5. Setiap error memiliki action, bukan hanya pesan error.

## Navigation

```text
Overview
├── Inbox
├── Contacts
├── Campaigns
├── Automations
├── Templates
├── Integrations
├── Channels
├── Analytics
└── Settings
    ├── Team & roles
    ├── Billing
    ├── Security
    ├── Data & retention
    └── Audit log
```

Developer center dapat muncul di Integrations atau mode navigasi advanced.

## First-run experience

### Empty state overview

Tampilkan progress:

```text
[ ] Create workspace
[ ] Connect WhatsApp channel
[ ] Verify channel health
[ ] Send test message
[ ] Invite teammate
[ ] Create first automation
```

Jangan menampilkan dashboard kosong tanpa guidance.

## Flow: connect official channel

1. Klik `Connect channel`.
2. Pilih `Official WhatsApp`.
3. Jelaskan requirement, estimasi waktu, ownership, dan fee.
4. Jalankan Embedded Signup/partner flow.
5. Kembali ke SaaS dengan channel context.
6. Verifikasi webhook dan permissions.
7. Tampilkan channel health.
8. Kirim test message.
9. Tandai onboarding selesai.

## Flow: connect unofficial channel

Hanya untuk closed beta/advanced mode:

1. User membuka `Experimental channels`.
2. Melihat disclosure: bukan official Meta API, risiko logout/ban, tidak ada SLA.
3. Confirm risk acknowledgement.
4. Pair via QR/device flow.
5. Tampilkan session health dan reconnect instructions.
6. Batasi feature/campaign sesuai policy.
7. Sediakan disconnect dan delete session yang mudah.

## Flow: send message

1. Pilih recipient/contact.
2. Pilih template atau free-form yang diperbolehkan channel.
3. Isi variable.
4. Preview actual message.
5. Lakukan consent/policy check.
6. Tampilkan estimated usage/cost.
7. Confirm/send.
8. Buka delivery timeline.

## Flow: failed message

Tampilan minimal:

```text
Failed
Reason: Template is not approved
Impact: Message was not sent
Next step: Open template manager
[Fix template] [Retry] [View request ID]
```

Tambahkan `technical details` accordion untuk developer.

## Flow: inbox agent

- Conversation list dengan unread, priority, assignee, SLA.
- Conversation detail dengan message timeline.
- Agent dapat reply, add note, tag, assign, resolve.
- Hindari double reply dengan presence/lock ringan.
- Tampilkan channel badge dan delivery status.

## Flow: campaign safety

1. Create draft.
2. Pilih audience/segment.
3. Validasi consent dan suppression.
4. Preview sample recipients.
5. Preview message + variable fallback.
6. Set schedule, timezone, quiet hours, rate.
7. Approval oleh manager.
8. Send/pause/cancel.
9. Monitor delivery/failure.

## Flow: developer integration

1. Open Developer Center.
2. Create sandbox/project.
3. Generate scoped API key.
4. Copy quickstart code.
5. Send test request.
6. Inspect request/response.
7. Configure webhook endpoint.
8. Verify signature.
9. Replay test event.
10. Promote key/environment to production.

## Content rules

- Bahasa default Indonesia; istilah API tetap English bila lazim.
- Gunakan “Pesan gagal karena…” bukan “Error 400”.
- Selalu tampilkan last updated untuk docs.
- Hindari janji “anti banned” atau “100% delivered”.
- Jelaskan batas provider tepat di tempat fitur dipakai.

## Accessibility requirements

- Keyboard navigation.
- Focus state terlihat.
- Form label dan error terhubung dengan input.
- Contrast memenuhi WCAG AA.
- Status tidak hanya dibedakan warna.
- Table pricing dapat dipakai screen reader.
- Loading state dan async result diumumkan.

## UX validation metrics

- Task completion rate.
- Time to first success.
- Number of support-assisted steps.
- Error recovery completion.
- Misunderstanding rate official vs unofficial.
- Confidence after setup.
