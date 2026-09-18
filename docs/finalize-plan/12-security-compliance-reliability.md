# 12 — Security, Compliance, dan Reliability

## Security objectives

- Tidak ada cross-tenant access.
- Credential/provider token aman.
- Webhook authentic/replay-safe.
- PII/message/media terlindungi.
- Internal access audited.
- User dapat revoke/delete/export.

## Threats and mitigation

| Threat | Mitigation |
|---|---|
| API key leak | Scoped key, reveal once, encryption, rotation/revoke, redacted logs |
| Cross-tenant leak | Workspace context, repository guard, RLS option, tests |
| Fake/replay webhook | HMAC, timestamp, event dedup |
| SSRF media URL | Validate public URL, block private IP, isolated fetcher |
| Malicious attachment | MIME/size validation, antivirus scan, signed URLs |
| Campaign abuse | Consent, opt-out, approval, quiet hours, rate guard |
| Duplicate send | Idempotency, dedup, provider-aware retry |
| Session theft | Encrypted session store, restricted access, delete flow |
| AI leakage | Data boundary, redaction, assist/handoff |

## Credential

Encrypt via envelope encryption/KMS/secret manager. Jangan log token, password, raw session, full phone/message body, atau raw payload tanpa redaction. Secret hanya tampil saat create/rotate.

## RBAC

Owner, Admin, Manager, Agent, Developer, Billing, Viewer dengan least privilege. Audit role/key/channel/campaign/export changes.

## Messaging safety

- `unknown/opted_in/opted_out` consent state.
- Source/timestamp consent.
- STOP handler.
- Global suppression.
- Quiet hours/timezone.
- Approval threshold.
- Recipient/rate limits.
- Dry run/test recipient.
- Jangan menjanjikan anti-ban atau 100% delivery.

## Official compliance

Ikuti Meta/WhatsApp Terms, Messaging Policy, Commerce Policy, template approval, messaging window, opt-in, dan policy yang berlaku. Legal review diperlukan untuk health/finance/sensitive sectors.

## Unofficial governance

`whatsmeow` hanya optional experimental: tidak default, closed beta, risk acknowledgement, label visible, no official SLA, protocol/ban monitoring, kill switch, session delete, legal/licence review.

## Privacy

Privacy policy, data map, retention, PII masking, subprocessor inventory, export/delete, DPA bila diperlukan.

Proposal retention:

- Message content: default 90 hari setelah resolved.
- Delivery events: 180 hari.
- Raw webhook: 7–30 hari encrypted.
- Media: lifecycle + cleanup.
- Audit: 1 tahun paid, 90 hari free.
- Credential: selama channel aktif.
- Unofficial session: hapus saat disconnect/delete.

## Reliability targets proposal

| Area | Pilot | GA |
|---|---:|---:|
| API availability | 99.5% | 99.9% |
| API p95 enqueue | <500 ms | <300 ms |
| Webhook eventual success | 98% | 99.5% |
| Status lag | <60 s | <30 s |
| RPO | 24 h | 1 h |
| RTO | 8 h | 2 h |

Target ini adalah internal SLO, bukan jaminan provider delivery.

## Incident response

SEV1: data/credential leak atau total outage. SEV2: provider outage/duplicate campaign/webhook outage. SEV3: isolated degradation. SEV4: minor issue.

Runbook: detect → contain/pause/revoke → communicate → recover/replay safely → preserve evidence → postmortem → update tests/docs.

## Acceptance gates

SAST, dependency/secret/image scan, tenant isolation, auth/RBAC, webhook signature/replay, rate-limit, backup restore, penetration test, legal/privacy review.
