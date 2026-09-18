# 05 — System Architecture

## Architecture principles

- Gunakan **Vertical Slice Architecture (VSA)** sebagai cara utama mengorganisasi code berdasarkan capability/feature.
- Modular monolith terlebih dahulu; extract service hanya saat ada bottleneck/ownership jelas.
- Ports/adapters tetap digunakan hanya pada boundary eksternal (database, queue, provider, storage), bukan sebagai struktur folder utama.
- Core feature tidak mengetahui detail Meta atau Whatsmeow secara langsung.
- Semua outbound/inbound message event dapat ditelusuri.
- Asynchronous processing untuk send, webhook, media, dan campaign.
- Tenant isolation pada query, cache key, object path, dan logs.
- Internal database ID tidak pernah keluar dari backend boundary: API response, Inertia props, UI network, logs publik, dan customer webhook hanya memakai public reference.
- Provider failure tidak boleh membuat dashboard ikut down.

## High-level architecture

```text
                         ┌────────────────────────┐
                         │ Inertia + Vue 3 + Nuxt UI │
                         └────────────┬───────────┘
                                      │ HTTP/Inertia
                         ┌────────────▼───────────┐
                         │ GoFiber + Fibertia       │
                         │ Auth/RBAC/Tenant/VSA     │
                         └────┬─────┬─────┬────────┘
                            │    │    │
                ┌───────────▼┐ ┌─▼────▼─────┐ ┌──────────────┐
                │ PostgreSQL │ │ Valkey/Queue │ │ Object Store │
                └────────────┘ └──────┬──────┘ └──────────────┘
                                      │
                         ┌────────────▼────────────┐
                         │ Go Workers               │
                         │ send/status/campaign    │
                         └───────┬──────────┬───────┘
                                 │          │
                    ┌────────────▼───┐  ┌───▼──────────────┐
                    │ Official Adapter│  │ Whatsmeow Adapter │
                    │ Cloud API       │  │ Experimental      │
                    └────────────┬────┘  └────────┬─────────┘
                                 │                │
                    ┌────────────▼────┐  ┌────────▼─────────┐
                    │ Meta Webhooks   │  │ WhatsApp session │
                    └─────────────────┘  └──────────────────┘
```

## Vertical Slice Architecture

Top-level code diorganisasi berdasarkan feature/capability, bukan berdasarkan folder global `handlers`, `services`, atau `repositories`.

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
    /whatsapp
      /official
      /whatsmeow
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
```

Setiap slice idealnya memiliki bagian yang dibutuhkan feature tersebut:

```text
/internal/features/message
  endpoint.go       # Fiber route/action
  page.go           # Inertia page props/actions
  service.go        # use case dan business rules
  repository.go     # port lokal feature
  repository_sql.go # implementasi sqlx/prepared statement
  model.go          # domain/read model feature
  dto.go            # request/response
  mapper.go          # map internal model → external DTO (id = public_id)
  policy.go
  event.go
  *_test.go
```

Tidak semua slice harus memiliki semua file. Ports/adapters boleh muncul lokal pada slice atau di `integrations/`, tetapi VSA tetap menjadi boundary utama.

## Request lifecycle outbound

1. API menerima request dengan tenant, actor, idempotency key.
2. Validate auth, capability, consent, quota, payload, template/policy.
3. Buat `message` status `accepted`/`queued`.
4. Publish `message.send.requested`.
5. Worker memilih provider/channel.
6. Provider adapter mengirim request.
7. Simpan provider ID/raw response ter-redact.
8. Update lifecycle event.
9. Provider webhook/status meng-update `sent`, `delivered`, `read`, `failed`.
10. Broadcast internal event ke UI/webhook customer.

## Request lifecycle inbound

1. Terima provider webhook/session event.
2. Verify signature atau provider-specific authenticity.
3. Simpan raw event terenkripsi/retention terbatas.
4. Deduplicate berdasarkan provider event ID.
5. Normalize ke internal event.
6. Upsert contact/conversation/message.
7. Jalankan automation/routing.
8. Notify UI/agent.

## Provider boundary di dalam VSA

Provider boundary tetap penting, tetapi dipanggil dari slice yang membutuhkan capability tersebut. Hindari satu interface besar yang memaksa semua provider memiliki kemampuan sama.

```go
type MessageSender interface {
    Send(ctx context.Context, req SendRequest) (SendResult, error)
}

type ChannelHealth interface {
    Health(ctx context.Context, channelID string) (Health, error)
}

type IncomingEventParser interface {
    VerifyAndParse(ctx context.Context, body []byte, headers http.Header) ([]ProviderEvent, error)
}
```

Contoh: `message` hanya membutuhkan `MessageSender`, sementara `channel` membutuhkan `ChannelHealth`. Implementasi official dan whatsmeow berada di integration boundary. VSA mengatur use case; ports/adapters mengisolasi dependency eksternal.

## Multi-tenancy

- Semua tabel tenant-owned memiliki `workspace_id`.
- Repository wajib menerima tenant context.
- Jangan mengandalkan filter di handler saja.
- PostgreSQL Row-Level Security dapat dipertimbangkan setelah model stabil.
- Cache key wajib menyertakan workspace/channel ID.
- Object path: `tenant/{workspace_id}/...`.
- Worker payload selalu membawa workspace ID dan authorization context minimal.

## Synchronous vs asynchronous

Synchronous:

- auth/session;
- CRUD contacts/templates;
- validate request;
- enqueue message;
- read current status.

Asynchronous:

- provider send;
- campaign fan-out;
- delivery status;
- webhook dispatch;
- media processing;
- analytics aggregation;
- notifications.

## Scalability path

### Stage 1 — pilot

Satu GoFiber/Fibertia API + worker deployment, PostgreSQL managed, Valkey, object storage.

### Stage 2 — growth

Pisahkan worker berdasarkan workload: messaging, campaigns, webhook delivery, analytics. Tambahkan read replicas dan queue metrics.

### Stage 3 — scale

Partition message/event tables, dedicated provider workers, durable event stream, per-tenant rate limiter, regional deployment.

## Failure isolation

- Timeout semua provider call.
- Circuit breaker per provider/channel.
- Queue backpressure.
- Retry hanya untuk error transient.
- Jangan retry error policy/invalid recipient.
- Dead-letter queue dengan replay permission.
- Kill switch untuk unofficial adapter.
