# 08 — VSA System Architecture

## Primary decision

Gunakan **Vertical Slice Architecture (VSA)** sebagai primary organization. Sistem dimulai sebagai modular monolith; service extraction hanya setelah bottleneck/ownership jelas.

Ports/adapters tetap dipakai, tetapi hanya untuk boundary eksternal. Kita tidak mengorganisasi code sebagai global `/handlers`, `/services`, `/repositories`.

## High-level

```text
┌─────────────────────────┐
│ Inertia + Vue 3 + Nuxt UI│
└────────────┬────────────┘
             │ HTTP/Inertia
┌────────────▼────────────┐
│ GoFiber + Fibertia       │
│ Auth/RBAC/Tenant/VSA     │
└─────┬──────┬──────┬──────┘
      │      │      │
┌─────▼──┐ ┌─▼──────▼──┐ ┌───────────┐
│Postgres│ │Valkey/Queue│ │Object Store│
└────────┘ └─────┬──────┘ └───────────┘
                  │
          ┌───────▼────────┐
          │ Go workers       │
          │ send/status/jobs │
          └──────┬─────┬─────┘
                 │     │
      ┌──────────▼┐  ┌─▼────────────┐
      │Official API│  │Whatsmeow     │
      │adapter     │  │experimental  │
      └───────────┘  └──────────────┘
```

## Repository layout

```text
/internal
  /features
    /auth
    /workspace
    /onboarding
    /channel
    /contact
    /conversation
    /message
    /template
    /campaign
    /automation
    /billing
    /analytics
    /developer
  /integrations
    /whatsapp/official
    /whatsapp/whatsmeow
    /payment
    /email
  /platform
    /database
    /http
    /inertia
    /queue
    /storage
    /crypto
    /observability
/resources/js
  /Pages
  /Components
  /Composables
  /types
  /lib
/db/migrations
```

## Feature slice example

```text
/internal/features/message
  endpoint.go       # Fiber route/action
  page.go           # Inertia page props/actions
  service.go        # business/use case
  repository.go     # local port
  repository_sql.go # sqlx/prepared statement
  model.go          # internal/read model
  dto.go            # external request/response
  mapper.go         # internal → external
  policy.go
  event.go
  *_test.go
```

Tidak semua slice wajib memiliki semua file. Shared code dibuat hanya jika memiliki consumer nyata lebih dari satu.

## Request lifecycle outbound

1. Fiber endpoint menerima request + tenant/actor/idempotency.
2. Auth, capability, consent, quota, payload, template/policy validation.
3. Create message `accepted`/`queued`.
4. Publish `message.send.requested` ke Valkey.
5. Worker memilih channel/provider.
6. Provider integration mengirim.
7. Simpan provider ID/raw response ter-redact.
8. Update lifecycle event.
9. Provider webhook/status → sent/delivered/read/failed.
10. Broadcast event ke UI/customer webhook/analytics.

## Request lifecycle inbound

1. Provider webhook/session event diterima.
2. Signature/authenticity verified.
3. Raw event encrypted dengan retention terbatas.
4. Deduplicate provider event ID.
5. Normalize event.
6. Upsert contact/conversation/message.
7. Jalankan routing/automation.
8. Notify UI/agent.

## VSA + integration boundary

Feature memegang use case dan business rule. Integration memegang detail Meta/HTTP/Whatsmeow.

```go
type MessageSender interface {
    Send(ctx context.Context, req SendRequest) (SendResult, error)
}

type ChannelHealth interface {
    Health(ctx context.Context, channelID string) (Health, error)
}
```

`message` tidak perlu mengetahui seluruh capability provider. Gunakan interface kecil/capability-specific.

## Multi-tenancy

- Semua tenant-owned tables memiliki `workspace_id` internal.
- Repository selalu menerima workspace context.
- Query tidak hanya mengandalkan handler filter.
- Cache key: `env:workspace:feature:key`.
- Object path: `tenant/{workspace_id}/...`.
- Worker message selalu membawa workspace ID dan minimum authorization context.
- External response tidak pernah membawa internal workspace ID.

## Sync/async

Synchronous: auth, CRUD, validation, enqueue, read status.

Asynchronous: provider send, campaign fan-out, delivery event, webhook dispatch, media, analytics, notifications.

## Failure isolation

- Provider timeout.
- Circuit breaker per provider/channel.
- Queue backpressure.
- Retry hanya transient error.
- Policy/invalid recipient tidak di-retry.
- Dead-letter queue + controlled replay.
- Kill switch untuk whatsmeow.

## Scalability path

- Pilot: satu API + worker deployment.
- Growth: worker terpisah per messaging/campaign/webhook/analytics.
- Scale: partition event/message, durable stream, dedicated provider workers, per-tenant rate limiter, regional deployment.
