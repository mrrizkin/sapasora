# 06 — UX dan Information Architecture

## UX goals

- User selalu tahu langkah berikutnya.
- Official/unofficial channel jelas.
- Status message mudah dipahami owner dan detail teknis tersedia untuk developer.
- Error memiliki next action.
- Dashboard tidak dimulai dari empty screen tanpa guidance.

## Navigation

```text
Overview
├── Inbox
├── Contacts
├── Campaigns
├── Automations
├── Templates
├── Integrations / Developer Center
├── Channels
├── Analytics
└── Settings
    ├── Team & roles
    ├── Billing
    ├── Security
    ├── Data & retention
    └── Audit log
```

## First-run checklist

```text
[ ] Create workspace
[ ] Connect WhatsApp channel
[ ] Verify channel health
[ ] Send test message
[ ] Invite teammate
[ ] Create first template/automation
```

## Official channel flow

1. Klik `Connect channel`.
2. Pilih Official WhatsApp.
3. Tampilkan requirement, estimasi waktu, ownership, fee, policy.
4. Jalankan Embedded Signup/partner flow.
5. Kembali ke SaaS dengan channel context.
6. Verifikasi permissions/webhook.
7. Tampilkan channel health.
8. Kirim test message.
9. Tampilkan onboarding complete.

## Unofficial channel flow

Hanya closed beta/advanced:

1. Buka Experimental channels.
2. Tampilkan bukan official Meta API, risiko logout/ban, protocol maintenance, no SLA.
3. Simpan risk acknowledgement.
4. QR/pairing.
5. Tampilkan session health/reconnect.
6. Batasi feature/campaign sesuai policy.
7. Sediakan disconnect/delete session.

## Send message flow

1. Pilih recipient.
2. Pilih template atau format yang didukung.
3. Isi variable.
4. Preview actual message.
5. Consent/policy check.
6. Estimasi usage/cost.
7. Confirm.
8. Delivery timeline.

## Failed message state

```text
Failed
Reason: Template is not approved
Impact: Message was not sent
Next step: Open template manager
[Fix template] [Retry] [View request ID]
```

Technical details berada dalam accordion, bukan menggantikan human-readable explanation.

## Inbox flow

- List: unread, priority, assignee, SLA.
- Detail: timeline, reply, internal note, tag, assign, resolve.
- Channel badge dan delivery status terlihat.
- Presence/lock ringan mencegah double reply.

## Campaign flow

1. Draft.
2. Audience/segment.
3. Consent/suppression validation.
4. Sample preview + variable fallback.
5. Schedule/timezone/quiet hours/rate.
6. Manager approval.
7. Send/pause/cancel.
8. Delivery report.

## Developer flow

1. Developer Center.
2. Create sandbox/project.
3. Create scoped API key.
4. Copy quickstart.
5. Send test request.
6. Inspect response/log.
7. Configure webhook.
8. Verify signature.
9. Replay test event.
10. Promote to live.

## Design system

- Nuxt UI sebagai baseline.
- AppShell, Sidebar, WorkspaceSwitcher.
- StatusBadge, HealthCard, DeliveryTimeline.
- EmptyState, ErrorState, LoadingSkeleton.
- DataTable, FilterBar, SearchInput.
- MessageComposer, TemplatePicker, VariablePreview.
- ContactImportWizard, CampaignReview.
- ConfirmRiskDialog, CredentialRevealOnce.
- Toast/NotificationCenter.

## Content/accessibility

- Default bahasa Indonesia; istilah API tetap English jika lazim.
- Error menjawab apa yang terjadi, penyebab, dan next step.
- Last updated, glossary, status badge Stable/Beta/Deprecated.
- Keyboard navigation, visible focus, WCAG AA contrast.
- Status tidak hanya dibedakan berdasarkan warna.
- Loading/async result diumumkan.
- Video memiliki caption, transcript, timestamp.
