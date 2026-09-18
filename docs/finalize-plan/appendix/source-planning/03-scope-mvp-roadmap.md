# 03 — Scope MVP dan Roadmap

## Product wedge

**Shared inbox + transactional notifications untuk SMB**, dengan API yang sama untuk developer.

Alasan pemilihan:

- frequency tinggi;
- value mudah diukur;
- lebih aman daripada memulai dari unsolicited broadcast;
- mendorong retention melalui inbox dan automation;
- cocok untuk official channel.

## MVP P0

### Workspace dan identity

- Register/login.
- Workspace creation.
- Owner/admin/member roles.
- Invite member.
- Audit log dasar.

### Channel

- Connect official channel melalui flow resmi/partner.
- Channel health.
- Disconnect/reconnect sesuai capability provider.
- Provider capability display.
- Experimental unofficial channel hanya feature flag/closed beta.

### Contacts

- Create/edit/delete.
- Import CSV dengan preview/validation.
- Tags/segments.
- Consent status dan consent source.
- Opt-out/suppression list.

### Inbox

- Conversation list.
- Search/filter.
- Message timeline.
- Reply text/media yang supported.
- Assignment agent.
- Internal note.
- Tag.
- Open/pending/resolved.

### Messaging

- Send test message.
- Send individual message.
- Template selection.
- Variable preview.
- Delivery timeline.
- Error reason dan retry/copy request ID.

### Developer

- API key create/revoke/rotate.
- REST send message.
- Message status endpoint.
- Signed webhook.
- Webhook test/replay terbatas.
- cURL, Go, Node.js, PHP, Python examples.

### Billing

- Plan selection.
- Usage meter.
- Invoice/history.
- Overage behavior.
- Provider/Meta fees dipisahkan dari platform fee.

## Must not ship di MVP

- Arbitrary bulk broadcast tanpa consent/approval.
- Unofficial connector untuk semua user.
- AI autonomous reply tanpa human handoff.
- Multi-provider billing kompleks.
- Omnichannel.
- Visual flow builder yang terlalu kompleks.
- Enterprise SSO/SCIM.

## P1 — setelah activation terbukti

- Campaign dengan approval, audience preview, schedule, quiet hours.
- Basic automation builder.
- Round-robin assignment.
- Canned replies.
- SLA/escalation.
- Retry/dead-letter UI.
- Better analytics.
- n8n/Make/Zapier connector.
- WooCommerce/WordPress integration.
- Unofficial connector closed beta untuk segment terpilih.

## P2 — scale

- Multi-WABA/channel routing.
- Agency mode.
- Advanced role/permission.
- SSO/SCIM.
- Custom retention/DPA/SLA.
- Conversion attribution.
- AI assist berbasis approved knowledge.
- Omnichannel.

## Release gates

### Alpha

- Internal team only.
- Synthetic/test data.
- Core send/receive/status stable.
- No customer campaign.

### Private beta

- 5–10 workspace.
- Support langsung.
- Feature flags.
- Daily error review.
- Incident runbook siap.

### Public beta

- Onboarding self-serve.
- Pricing dan legal reviewed.
- Usage/billing reliable.
- Status page.
- Documented limitations.

### General availability

- Security review.
- Backup/restore tested.
- SLO defined.
- Support response policy.
- Provider failure drill.
- Data deletion verified.

## Definition of MVP success

- 70% pilot workspace menyelesaikan onboarding tanpa pendampingan penuh.
- Median time-to-first-delivered-message <15 menit setelah channel siap.
- Delivery/status pipeline dapat ditelusuri untuk 99% message events.
- Tidak ada credential exposure pada testing.
- Minimal 3 workspace aktif mingguan selama 4 minggu.
- Minimal 2 customer bersedia membayar.
