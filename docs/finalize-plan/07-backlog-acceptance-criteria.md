# 07 — Backlog dan Acceptance Criteria

## Epic A — Workspace/auth

- Register/login.
- Create workspace.
- Invite team.
- Roles: owner, admin, manager, agent, developer, billing, viewer.
- Session revoke.
- Audit log.

**Acceptance:** tenant isolation, invite expiry, permission enforcement, no secret exposure.

## Epic B — Channel onboarding

- Explain official requirements.
- Connect official channel.
- Channel health.
- Test message.
- Disconnect/delete credential.
- Experimental whatsmeow flow behind feature flag/risk acknowledgement.

**Acceptance:** provider/capability visible; cancel/retry; no raw credential leakage.

## Epic C — Contacts/consent

- CRUD contact.
- CSV preview/validation.
- Tags/segments/groups.
- Variables.
- Consent source/timestamp.
- Opt-in/opt-out.
- Suppression list.

**Acceptance:** normalized identifier, duplicate handling, consent audit, opt-out honored.

## Epic D — Inbox

- Conversation list/search/filter.
- Timeline.
- Reply text/media supported.
- Assignment/reassignment.
- Internal note/tag.
- Open/pending/resolved.
- SLA visibility.

**Acceptance:** unread state, collision handling, delivery result visible.

## Epic E — Messaging API

- Scoped API key.
- Send text/template/media sesuai capability.
- Idempotency key.
- Message/status endpoints.
- Signed webhook.
- Webhook test/replay.
- Go/Node/PHP/Python/cURL examples.

**Acceptance:** OpenAPI, request ID, stable errors, signature/retry guidance.

## Epic F — Templates

- List/create/update/delete sesuai provider capability.
- Approval status.
- Variable preview/fallback.
- Rejection reason.
- API access.

## Epic G — Campaign

- Draft/audience preview.
- Consent filtering.
- Approval.
- Schedule/timezone/quiet hours.
- Rate guard.
- Pause/cancel/report.

## Epic H — Automation

- Trigger.
- Condition.
- Action: send/assign/tag/delay/notify.
- Run history.
- Pause/enable.
- AI assist dengan boundary/handoff.

## Epic I — Billing

- Plan/usage.
- Cost estimate.
- Provider fee separation.
- Invoice/history.
- Overage cap.
- Upgrade/downgrade/proration explanation.

## Epic J — Operations

- Request search by request ID.
- Diagnostics ter-redact.
- Status page.
- Support ticket context.
- Provider/campaign pause.

## P0 acceptance checklist

- [ ] User dapat membuat workspace.
- [ ] User dapat connect official test channel.
- [ ] User dapat mengirim test message.
- [ ] User melihat accepted/queued/sent/delivered/failed.
- [ ] User dapat memperbaiki error umum.
- [ ] User dapat membuat contact dan consent.
- [ ] Agent dapat menerima/reply/assign conversation.
- [ ] Developer dapat membuat key dan send API.
- [ ] Customer webhook ter-signature dan retryable.
- [ ] Public response hanya memakai public `id: string`.
- [ ] Frontend mem-parse boundary dengan Zod.
- [ ] Logs/metrics/traces tidak mengandung credential/PII.

## Story template

```markdown
## [FEATURE] Judul

### User story
Sebagai ..., saya ingin ..., sehingga ...

### Scope
...

### Acceptance criteria
- [ ] ...

### Error/recovery
- ...

### Security/privacy
- ...

### Analytics
- ...

### Docs/support
- ...
```
