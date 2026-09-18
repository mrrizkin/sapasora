# 13 — Infrastructure dan Observability

## Runtime

- Vue/Inertia assets via CDN/static hosting.
- GoFiber + Fibertia API.
- Go workers: message sender, provider event, webhook dispatch, campaign scheduler/fan-out, media, analytics, notifications.
- PostgreSQL managed.
- Valkey/queue.
- S3-compatible object storage.
- Official adapter.
- Isolated whatsmeow adapter.
- Mail provider.

Mulai tanpa Kubernetes. Gunakan containerized managed infrastructure dengan Podman-compatible workflow.

## Environments

- Local: mock provider, synthetic data, Mailpit, MinIO.
- Development: Meta test assets, internal users.
- Staging: production-like, no real campaign.
- Production: real data, secret manager, strict access.

## CI/CD

1. Format/lint.
2. Unit/integration/contract tests.
3. Frontend build/typecheck/test.
4. OpenAPI compatibility.
5. Secret/dependency/image scan.
6. Build immutable image.
7. Staging deploy/smoke.
8. Manual production approval.
9. Backward-compatible migration.
10. Rollout/health check/rollback.

## Observability

### Logs

zerolog structured JSON: timestamp, level, service, environment, request ID, safe workspace ID, channel ID, message ID, provider, operation, duration, error code.

Dilarang log token, password, session secret, full phone, full message body, atau raw payload tanpa redaction.

### Metrics

- API request/error/latency/auth/rate limit.
- Message status/provider latency/queue depth/retry/DLQ/delivery lag/duplicate.
- Webhook receive/verify/reject/dispatch/retry/signature failure.
- Onboarding completion/first delivery/active workspace/campaign/inbox response.
- Unofficial auth/protocol/disconnect signal.

### Traces

Trace API → DB → Valkey → provider → webhook → automation. Propagate request/correlation ID.

## Alerts

API error spike, queue age, provider error, signature failure, duplicate signal, DB connection/lock, storage anomaly, worker crash loop, unofficial protocol/auth failure.

## Backups/DR

- Encrypted PostgreSQL backups.
- Point-in-time recovery.
- Restore drill.
- Object storage versioning/lifecycle.
- Durable queue semantics sesuai kebutuhan.
- RPO/RTO documented.

## Cost controls

Per-tenant usage meter, sampled logs/traces, storage lifecycle, worker/concurrency guard, campaign fan-out limit, provider cost alert, AI/media/queue budget.
