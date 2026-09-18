# 10 — Infrastructure, Deployment, dan Observability

## Deployment strategy

Mulai dari managed infrastructure untuk mengurangi operational load:

- Container image untuk Go API/worker.
- Managed PostgreSQL.
- Managed Valkey atau queue service.
- S3-compatible object storage.
- Managed DNS/TLS/WAF.
- Managed email provider.
- Error tracking dan metrics backend.

Kubernetes belum diperlukan pada fase pilot kecuali ada kebutuhan nyata.

## Runtime components

- `web`: Vue static assets melalui CDN.
- `api`: GoFiber + Fibertia HTTP/Inertia API.
- `worker-message`: outbound provider send.
- `worker-webhook`: customer webhook delivery.
- `worker-campaign`: schedule/fan-out.
- `worker-analytics`: aggregates.
- `scheduler`: recurring jobs.
- `migrator`: controlled database migrations.
- `official-adapter`: Cloud API/partner integration.
- `whatsmeow-adapter`: isolated experimental worker.

## Environments

### Local

Mock provider, seeded data, Mailpit, MinIO.

### Staging

Meta test assets, isolated credentials, no real production campaigns.

### Production

Real credentials via secret manager, strict deployment approval, audit, backup, and alerting.

## CI/CD pipeline

1. Format/lint.
2. Unit tests.
3. Integration tests.
4. Contract tests provider.
5. Frontend build/test.
6. OpenAPI compatibility check.
7. Secret/dependency/image scan.
8. Build immutable images.
9. Deploy staging.
10. Smoke/e2e test.
11. Manual approval production.
12. Migration and rollout.
13. Post-deploy health check.

## Observability pillars

### Logs

Structured JSON:

- timestamp;
- level;
- service;
- environment;
- request_id;
- workspace_id hashed/redacted;
- channel_id;
- message_id;
- provider;
- error_code;
- duration_ms.

Never log token, message body, raw PII, or full provider credential.

### Metrics

API:

- request count/error/latency;
- auth failures;
- rate limited requests.

Messaging:

- messages by status/provider/channel;
- provider latency;
- queue depth/age;
- retry count;
- dead-letter count;
- delivery lag;
- duplicate/idempotency conflicts.

Webhook:

- received/verified/rejected;
- dispatch success/failure;
- retry age;
- signature failure.

Product:

- onboarding step completion;
- first delivered message;
- active workspace;
- campaign completion;
- inbox response time.

### Traces

Trace boundaries:

- incoming API request;
- DB query group;
- queue publish/consume;
- provider call;
- webhook dispatch;
- automation run.

Propagate correlation ID across API → queue → provider → webhook.

## Dashboard and alerts

Alerts:

- API error rate high.
- Queue age above threshold.
- Provider error spike.
- Signature failures spike.
- Duplicate campaign signal.
- Database connections/locks high.
- Storage growth/anomalous attachment.
- Worker crash loop.
- Unofficial protocol/auth failure spike.

Customer-facing status page hanya menampilkan informasi yang aman dan relevan.

## Backup and disaster recovery

- Automated PostgreSQL backups.
- Point-in-time recovery.
- Encrypted backup.
- Restore drill berkala.
- Object storage versioning/lifecycle.
- Queue messages durable sesuai kebutuhan.
- Document RPO/RTO.
- Jangan menganggap provider sebagai backup data internal.

## Cost controls

- Per-tenant usage meter.
- Log/trace sampling.
- Storage lifecycle cleanup.
- Worker concurrency guard.
- Campaign max fan-out.
- Alert biaya provider/Meta.
- Budget limit untuk AI/media/queue.
