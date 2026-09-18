# 15 — Metrics, Analytics, dan Launch Plan

## North Star Metric

**Weekly Active Workspaces with at least one delivered-and-responded conversation or successful transactional notification.**

Metric harus mengukur value, bukan sekadar request API.

## Funnel

```text
Landing visit
  → signup
  → workspace created
  → channel connected
  → test message accepted
  → test message delivered
  → first inbound/reply
  → first automation/campaign
  → weekly active
  → paid
```

## Activation metrics

- Signup completion rate.
- Workspace creation rate.
- Channel connection success rate.
- Time to first accepted message.
- Time to first delivered message.
- Onboarding completion.
- Support-assisted onboarding rate.
- First 7-day retention.

## Product quality metrics

- Delivery rate by provider/channel.
- Failure rate by normalized error code.
- Status webhook lag.
- Webhook delivery eventual success.
- Duplicate message count.
- Queue age.
- Disconnect rate.
- Campaign pause/cancel rate.
- Consent violation attempts blocked.

## Business metrics

- MRR/ARR.
- Free-to-paid conversion.
- ARPA.
- Gross margin per message/workspace.
- Churn and reason.
- Net revenue retention.
- Support cost per workspace.
- CAC payback.

## Event taxonomy

```text
signup_completed
workspace_created
channel_connect_started
channel_connected
channel_health_failed
test_message_sent
test_message_delivered
message_send_failed
conversation_received
conversation_replied
campaign_created
campaign_approved
campaign_paused
api_key_created
webhook_verified
subscription_started
subscription_canceled
```

Event properties wajib minim dan redacted:

- workspace plan;
- provider type;
- channel capability;
- country;
- error code;
- duration bucket.

Jangan kirim message body, phone number, API key, atau PII ke product analytics.

## Launch plan

### Phase 0 — foundation

- Product brief dan research.
- Meta/provider feasibility.
- Threat model.
- Design system shell.
- Go/Vue repo setup.
- Mock provider.
- CI/CD.

### Phase 1 — internal alpha

- Workspace/auth.
- Channel mock + official test integration.
- Send/status lifecycle.
- Inbox basic.
- Contacts/consent.
- API/webhook.
- Observability.

### Phase 2 — private pilot

- 5–10 selected businesses.
- Use cases: invoice, appointment, lead follow-up, support.
- Manual onboarding/support.
- Weekly product review.
- Strict campaign limits.

### Phase 3 — public beta

- Self-serve official onboarding.
- Pricing/billing.
- Docs/tutorial.
- Status page.
- Support SLA internal.
- Public changelog.

### Phase 4 — expansion

- Automation builder.
- Campaign approval.
- Agency mode.
- Optional whatsmeow closed beta.
- Additional integrations.

## Go-to-market experiments

1. Landing page per use case.
2. Developer quickstart and sandbox.
3. WordPress/WooCommerce plugin.
4. Agency partnership.
5. Template library.
6. Free audit: “cek kesiapan WhatsApp automation”.
7. Educational content around consent, status, and reliability.

## Launch checklist

### Product

- [ ] Value proposition tested.
- [ ] Target segment selected.
- [ ] P0 flows complete.
- [ ] Pricing and limits clear.

### Technical

- [ ] Official provider integration verified.
- [ ] Webhook signature and replay protection.
- [ ] Idempotency and retry.
- [ ] Backup/restore.
- [ ] Alerts/runbooks.

### Trust/legal

- [ ] Terms/privacy reviewed.
- [ ] WhatsApp/Meta policy reviewed.
- [ ] Consent/opt-out implemented.
- [ ] Data retention explained.
- [ ] Unofficial channel separately disclosed, if enabled.

### Operations

- [ ] Support process.
- [ ] Incident communication.
- [ ] Status page.
- [ ] Billing support.
- [ ] Customer onboarding material.

## Pilot review cadence

- Daily: incidents, failed messages, onboarding blockers.
- Weekly: activation, delivery quality, support themes, qualitative feedback.
- Monthly: retention, unit economics, roadmap re-prioritization.
