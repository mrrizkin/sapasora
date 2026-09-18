# 16 — Metrics dan Launch Plan

## North Star

**Weekly Active Workspaces** dengan minimal satu delivered-and-responded conversation atau successful transactional notification.

## Funnel

```text
visit → signup → workspace → channel connected → test accepted
→ test delivered → inbound/reply → automation/campaign → weekly active → paid
```

## Activation

- Signup/workspace/channel completion.
- Time to first accepted/delivered message.
- Onboarding completion.
- Support-assisted onboarding.
- 7-day retention.

## Quality

- Delivery/failure rate by provider/channel/error.
- Status webhook lag.
- Webhook eventual success.
- Duplicate outbound.
- Queue age.
- Disconnect rate.
- Campaign pause/cancel.
- Consent violation blocked.

## Business

- MRR/ARR, ARPA, gross margin.
- Cost per delivered message.
- Support cost/workspace.
- Free-to-paid.
- Churn/reason.
- NRR.
- CAC payback/LTV:CAC.

## Events

`signup_completed`, `workspace_created`, `channel_connect_started`, `channel_connected`, `channel_health_failed`, `test_message_sent`, `test_message_delivered`, `message_send_failed`, `conversation_received`, `conversation_replied`, `campaign_created`, `campaign_approved`, `campaign_paused`, `api_key_created`, `webhook_verified`, `subscription_started`, `subscription_canceled`.

Event properties minimal/redacted: plan, provider, capability, country, error code, duration bucket. Jangan mengirim message body, phone, API key, atau PII ke analytics.

## Launch phases

### Phase 0 — Foundation/validation

Interview, provider feasibility, threat model, prototype, VSA/stack spike, policy review.

### Phase 1 — Internal alpha

Workspace/auth, mock + official test integration, contacts/consent, send/status, inbox, API/webhook, observability.

### Phase 2 — Private pilot

5–10 businesses: invoice, appointment, lead, support. Manual onboarding, strict campaign limit, weekly review.

### Phase 3 — Public beta

Self-serve official onboarding, pricing/billing, docs, campaign/automation basic, status page, changelog.

### Phase 4 — Expansion

Agency mode, integrations, optional whatsmeow closed beta, advanced analytics/AI.

## GTM experiments

- Landing page per use case.
- Developer sandbox/quickstart.
- WordPress/WooCommerce plugin.
- Agency partnership.
- Template library.
- WhatsApp automation readiness audit.
- Education on consent, status, reliability.

## Launch checklist

### Product

- [ ] Positioning/segment tested.
- [ ] P0 flows complete.
- [ ] Pricing/limits clear.

### Technical

- [ ] Official provider verified.
- [ ] HMAC/replay protection.
- [ ] Idempotency/retry.
- [ ] Backup/restore.
- [ ] Alerts/runbooks.

### Trust/legal

- [ ] Terms/privacy.
- [ ] Meta/WhatsApp policy.
- [ ] Consent/opt-out.
- [ ] Data retention.
- [ ] Unofficial disclosure separated.

### Operations

- [ ] Support/status page.
- [ ] Billing support.
- [ ] Incident communication.
- [ ] Customer onboarding material.
