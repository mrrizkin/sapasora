# 06 — Technical Stack: GoFiber + Fibertia + Vue

## Decision summary

Baseline teknis mengikuti pengalaman dan signature produk developer:

- **Backend:** Go
- **HTTP framework:** GoFiber
- **Full-stack rendering/navigation:** Inertia
- **Go ↔ Inertia integration:** library internal **Fibertia**
- **Frontend:** Vue 3 + TypeScript
- **Runtime schema validation:** Zod untuk request, response, Inertia props, shared props, storage, dan boundary runtime
- **UI components:** Nuxt UI
- **Frontend routing/state:** Inertia page navigation dan page props; tidak menggunakan `vue-router` atau Pinia sebagai default
- **Database access:** sqlx + prepared statements
- **Database migration:** goose
- **Cache/queue:** Valkey
- **Logging:** zerolog
- **API schema/documentation UI:** OpenAPI + Scalar
- **Container tooling:** Podman + Compose-compatible workflow
- **Unofficial provider:** `whatsmeow` di worker/slice terisolasi

Ini bukan “stack pilihan sementara”. Fibertia dapat menjadi bagian dari developer experience dan signature produk internal, karena pola server-driven page + Vue tetap konsisten dengan cara kerja Anda.

## VSA sebagai primary architecture

Vertical Slice Architecture menjadi cara utama menyusun code. Ports/adapters tetap dipakai pada dependency boundary, tetapi bukan top-level organization.

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
```

Contoh slice:

```text
/internal/features/message
  endpoint.go
  page.go
  service.go
  repository.go
  repository_sql.go
  model.go
  dto.go
  mapper.go
  policy.go
  event.go
  *_test.go
```

Aturan:

- Jangan membuat folder global `/handlers`, `/services`, `/repositories` sebagai boundary utama.
- Business rule diletakkan dekat dengan feature yang memilikinya.
- Shared code hanya dibuat bila memiliki lebih dari satu consumer yang nyata.
- Integrasi eksternal tetap berada di `/integrations` atau local adapter yang jelas.
- Satu feature boleh memiliki HTTP action, Inertia page, SQL repository, event, dan test dalam slice yang sama.

## GoFiber + Fibertia

### Request model

```text
Browser request
  → GoFiber route
  → feature action/use case
  → repository/provider
  → Inertia response/page props
  → Vue page/component
```

Gunakan Inertia untuk:

- server-driven navigation;
- page props;
- form submission/action;
- validation error;
- redirect/flash message;
- partial reload bila didukung Fibertia.

Gunakan JSON API endpoint hanya untuk:

- public developer API;
- provider webhooks;
- asynchronous polling bila diperlukan;
- integrasi eksternal;
- kebutuhan UI yang memang membutuhkan client fetch.

### Fibertia boundary

Buat convention Fibertia yang konsisten untuk:

- page name;
- shared props;
- auth/user/workspace context;
- validation errors;
- flash messages;
- pagination/filter state;
- authorization failure;
- partial reload/deferred props;
- file upload.

Fibertia sebaiknya memiliki integration tests untuk memastikan kontrak Go ↔ Vue stabil.

## Go backend

### Libraries/conventions

- GoFiber untuk HTTP app dan middleware.
- `context.Context` pada operasi I/O.
- `zerolog` untuk structured JSON logging.
- `sqlx` untuk query dan scanning.
- `goose` untuk migration.
- OpenAPI sebagai contract artifact.
- Scalar sebagai API reference UI.
- `go vet`, `staticcheck`, `govulncheck`, formatter/linter di CI.
- Standard `testing`, `httptest`, dan Testcontainers untuk integration test.

### Error handling

Gunakan typed application errors yang dapat dipetakan ke:

- HTTP status;
- public error code;
- user-facing message;
- technical detail;
- retryable flag;
- request ID.

Jangan expose SQL error, provider raw error, token, atau message PII ke response user.
- Jangan expose internal BIGSERIAL `id` atau field database `public_id` secara langsung; mapping wajib dilakukan di DTO/serializer.

### Goroutine/worker rules

- Setiap goroutine punya owner dan cancellation.
- Shutdown harus graceful.
- Worker memiliki concurrency, timeout, retry, dan metric sendiri.
- Tidak ada goroutine background tanpa lifecycle management.

## Repository dengan sqlx + prepared statements

### Repository rules

- Repository didefinisikan dekat dengan feature slice.
- Query menggunakan `sqlx` dan prepared statement untuk query yang dipanggil berulang.
- Transaction boundary berada di use case/repository sesuai kebutuhan atomicity.
- Semua query tenant-owned menerima `workspace_id` secara eksplisit.
- Internal database `id` tidak pernah keluar dari backend boundary.
- External DTO/page props memakai field `id: string` yang nilainya berasal dari internal `public_id`.
- Field `public_id` tidak pernah dikirim ke API response, Inertia props, UI network, atau webhook customer.
- Query tidak boleh mengambil kolom yang tidak diperlukan.
- `rows.Close()` dan error `rows.Err()` wajib ditangani.
- Prepared statement lifecycle harus mengikuti lifecycle DB/transaction dan tidak bocor.

Contoh pola prepared statement yang dibuat saat constructing repository:

```go
const findMessageByPublicID = `
    SELECT id, public_id, workspace_id, status, created_at
    FROM messages
    WHERE workspace_id = $1 AND public_id = $2
`

type Repository struct {
    db   *sqlx.DB
    find *sqlx.Stmt
}

func NewRepository(ctx context.Context, db *sqlx.DB) (*Repository, error) {
    find, err := db.PreparexContext(ctx, findMessageByPublicID)
    if err != nil {
        return nil, err
    }
    return &Repository{db: db, find: find}, nil
}

func (r *Repository) Close() error {
    return r.find.Close()
}

func (r *Repository) FindByPublicID(
    ctx context.Context,
    workspaceID int64,
    publicID string,
) (Message, error) {
    var out Message
    err := r.find.GetContext(ctx, &out, workspaceID, publicID)
    return out, err
}
```

Gunakan prepared statement untuk hot query yang benar-benar berulang dan ukur dengan benchmark/metrics. Untuk query yang jarang dipanggil, query biasa bisa lebih sederhana. Statement yang dibuat repository harus ditutup saat graceful shutdown.

Untuk hot path yang benar-benar berulang, gunakan `PreparexContext` dan ukur manfaatnya. Jangan mengorbankan kesederhanaan jika query hanya dijalankan sesekali.

## Database migration dengan goose

- Migration forward-only dan terurut.
- Satu migration fokus pada satu perubahan logis.
- Hindari lock panjang di production.
- Index besar dibuat dengan strategi online/concurrent sesuai database.
- Migration tested pada database kosong dan database berisi data realistis.
- Deploy app harus backward-compatible dengan schema satu versi sebelumnya.
- Seed/demo data dipisahkan dari production migration.

Contoh struktur:

```text
/db/migrations
  202609040001_create_workspaces.sql
  202609040002_create_channels.sql
  202609040003_create_messages.sql
```

## Frontend Vue + Nuxt UI

### Baseline

- Vue 3 + TypeScript.
- Vite sebagai bundler jika tidak ada kebutuhan Nuxt runtime.
- Nuxt UI sebagai component/design baseline; gunakan komponen yang kompatibel dengan setup Vue/Inertia.
- Vitest untuk unit/component test.
- Playwright untuk browser/e2e test.

### Routing dan state

Tidak memakai `vue-router` karena navigation berasal dari Inertia.

Tidak memakai Pinia sebagai default karena:

- page data berasal dari server melalui Inertia props;
- form state lokal berada di component/form helper;
- flash/auth/workspace context berada di shared page props;
- state yang benar-benar client-only dapat memakai composable lokal.

Jika nanti ada state global client-heavy, keputusan Pinia/state library dibuat melalui ADR, bukan otomatis.

### Frontend structure

```text
/resources/js
  /Pages
    /Auth
    /Onboarding
    /Overview
    /Inbox
    /Contacts
    /Campaigns
    /Automations
    /Templates
    /Integrations
    /Channels
    /Settings
  /Components
    /ui
    /inbox
    /messages
    /campaigns
  /Composables
  /types
  /lib
```

Page component mengikuti feature slice backend, sedangkan primitive UI berada di component layer bersama.

## OpenAPI + Scalar

- OpenAPI disimpan di repository dan versioned.
- Scalar menjadi UI dokumentasi utama untuk public API.
- API schema tidak digenerate terpisah dari implementasi tanpa review.
- Zod memvalidasi runtime request/response dan Inertia props; data network dimulai sebagai `unknown`.
- Invalid boundary data harus fail loudly/early dan tidak boleh silently diabaikan.
- CI memeriksa schema breaking change.
- Contoh request/response wajib mencakup cURL, Go, JavaScript, PHP, Python.
- Internal Inertia actions tidak harus dipaksa menjadi public OpenAPI API.

## Valkey

Valkey digunakan untuk:

- queue sederhana/worker coordination;
- cache ephemeral;
- rate limit;
- idempotency short-lived record;
- distributed lock dengan expiry;
- session/flash bila diperlukan.

Rules:

- Valkey bukan source of truth untuk message/billing.
- Semua key memiliki namespace: `env:workspace:feature:key`.
- TTL wajib untuk cache/lock/idempotency.
- Queue semantics dan failure/retry harus terdokumentasi.
- Jika durability/streaming semakin kompleks, evaluasi NATS JetStream/managed queue melalui ADR.

## Zerolog

Structured fields minimum:

- `request_id`;
- `workspace_id` yang sudah direpresentasikan aman;
- `channel_id`;
- `message_id`;
- `provider`;
- `operation`;
- `duration_ms`;
- `error_code`.

Tidak boleh logging:

- API key/token;
- password/session secret;
- full phone number bila tidak diperlukan;
- full message body;
- raw provider payload tanpa redaction.

## Podman workflow

Local development menggunakan Podman:

- `podman compose` atau compatibility layer yang disepakati;
- PostgreSQL;
- Valkey;
- MinIO/S3-compatible storage;
- Mailpit;
- mock provider/webhook fixture;
- optional whatsmeow worker disabled by default.

Contoh command convention:

```bash
make dev
make infra-up
make infra-down
make migrate-up
make migrate-down
make test
make test-integration
make lint
make generate
make openapi-check
```

## Technical decisions

| Area | Decision | Status |
|---|---|---|
| Architecture | VSA modular monolith | Accepted direction |
| HTTP | GoFiber | Accepted |
| Inertia bridge | Fibertia | Accepted/signature |
| Frontend | Vue 3 + TypeScript + Zod | Accepted |
| UI | Nuxt UI | Accepted direction |
| Navigation | Inertia, no vue-router | Accepted |
| Global state | No Pinia by default | Accepted |
| Database access | sqlx + prepared statements where useful | Accepted |
| IDs | Internal `BIGSERIAL`; external `id: string` = `public_id`; `public_id` absent | Accepted |
| Migration | goose | Accepted |
| Logging | zerolog | Accepted |
| Cache/queue | Valkey | Accepted direction |
| API docs | OpenAPI + Scalar | Accepted |
| Containers | Podman | Accepted |
| Provider | Official first; whatsmeow optional | Accepted direction |
