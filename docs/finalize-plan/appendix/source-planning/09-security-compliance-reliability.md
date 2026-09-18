# 09 — Security, Compliance, dan Reliability

## Security objectives

1. Tidak ada cross-tenant data access.
2. Token/provider credential tidak bocor.
3. Webhook tidak mudah dipalsukan/replay.
4. Message/PII terlindungi.
5. Akses internal dapat diaudit.
6. User diberi kontrol revoke/delete/export.

## Threat model ringkas

| Threat | Dampak | Mitigasi |
|---|---|---|
| API key leaked | Pengiriman atas nama tenant | Scoped key, secret reveal once, rotate/revoke, redacted logs |
| Cross-tenant query | Kebocoran data | Tenant context, repository guard, tests, optional RLS |
| Fake webhook | Data/automation palsu | HMAC, timestamp, replay protection |
| SSRF dari media URL | Akses jaringan internal | Allowlist/URL validation, proxy isolation, no private IP |
| Malicious attachment | Malware/storage abuse | MIME sniffing, size limit, antivirus scan, signed URLs |
| Prompt/data leakage | PII terbocor ke AI | Data boundaries, redaction, retention, opt-in |
| Campaign abuse | Spam/report/ban | Consent, approval, rate guard, opt-out |
| Session theft | Unofficial account takeover | Encrypt session, restricted access, rotation/delete |
| Worker duplicate | Duplicate message | Idempotency, dedup keys, provider-aware retry |
| Supply-chain issue | Compromise runtime | Pin dependencies, scan, SBOM, minimal image |

## Credential management

- Browser hanya menerima secret saat create/rotate jika benar-benar perlu.
- Simpan credential terenkripsi dengan envelope encryption.
- KMS/secret manager untuk master key.
- Tidak menulis token di logs, analytics, traces, error message, atau support screenshot.
- Rotate/revoke tersedia di UI dan API.
- Internal access by break-glass/audit.
- Unofficial session database dipisahkan dan encrypted.

## RBAC baseline

| Role | Hak utama |
|---|---|
| Owner | Semua workspace, billing, delete |
| Admin | Team, channel, settings, audit |
| Manager | Campaign approval, inbox, analytics |
| Agent | Inbox dan reply sesuai assignment |
| Developer | API keys, webhooks, logs |
| Billing | Plans, invoice, payment |
| Viewer | Read-only |

Tambahkan permission granular sebelum enterprise.

## Messaging policy controls

- Contact consent state: unknown, opted_in, opted_out.
- Source dan timestamp consent.
- STOP/unsubscribe handler.
- Global suppression list.
- Quiet hours per timezone.
- Campaign approval threshold.
- Max recipients per campaign.
- Rate/concurrency guard.
- Blocklist untuk risky content/recipient behavior sesuai policy.
- Dry run/test recipient.

## Official channel compliance

- Ikuti WhatsApp Business Platform Terms, Messaging Policy, Commerce Policy, dan kebijakan Meta yang berlaku.
- Jelaskan bahwa approval template dan messaging window dapat memengaruhi kemampuan kirim.
- Sediakan policy link/version pada onboarding.
- Jangan menjanjikan “anti-ban” atau bypass policy.
- Pastikan consent/legal review untuk sektor sensitif seperti kesehatan/keuangan.

## Unofficial channel governance

- Tidak diaktifkan default.
- Hanya closed beta sampai legal/operational review selesai.
- Risk acknowledgement disimpan sebagai audit event.
- Label provider selalu terlihat.
- Tidak ada claim official support.
- Kill switch per tenant/provider.
- Monitoring disconnect, ban signal, auth failure, protocol error.
- Data/session deletion flow harus jelas.
- Review lisensi library dan kebijakan Meta sebelum commercial launch.

## Privacy

- Privacy policy mudah ditemukan.
- Data map: contacts, message body, media, webhook payload, credentials, analytics.
- Retention configurable.
- PII masking dan least privilege.
- Export/delete request.
- Subprocessor inventory.
- DPA untuk customer enterprise jika diperlukan.

## Reliability targets (proposal)

| Area | Pilot target | GA target |
|---|---:|---:|
| API availability | 99.5% | 99.9% |
| API p95 enqueue latency | <500 ms | <300 ms |
| Webhook dispatch success | 98% first attempt | 99.5% eventual |
| Status processing lag | <60 s | <30 s |
| RPO database | 24 h | 1 h |
| RTO | 8 h | 2 h |

Target ini adalah SLO internal, bukan janji provider message delivery.

## Incident response

Severity:

- **SEV1:** cross-tenant leak, credential leak, total send outage.
- **SEV2:** provider/channel outage besar, duplicate campaign, webhook outage.
- **SEV3:** feature degradation, isolated tenant issue.
- **SEV4:** minor UI/docs issue.

Runbook minimal:

1. Detect dan triage.
2. Contain: pause campaign/kill switch/revoke key.
3. Communicate status.
4. Recover dan replay hanya event safe.
5. Preserve evidence.
6. Postmortem tanpa blame.
7. Update test/runbook/docs.

## Security acceptance gates

- SAST/dependency scan.
- Secret scan.
- Tenant isolation tests.
- Auth/RBAC tests.
- Webhook signature/replay tests.
- Rate-limit tests.
- Backup restore test.
- Penetration test sebelum GA.
- External legal/privacy review.
