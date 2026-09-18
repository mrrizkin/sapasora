# 09 — Technical Stack dan Engineering Conventions

## Baseline

- Go stable + GoFiber.
- Inertia + Vue 3 + TypeScript.
- Fibertia sebagai GoFiber ↔ Inertia library/signature internal.
- Nuxt UI.
- sqlx + prepared statements.
- PostgreSQL.
- goose.
- Valkey.
- zerolog.
- OpenAPI + Scalar.
- Podman.
- Zod.
- `whatsmeow` hanya pada experimental worker/feature slice terisolasi.

## GoFiber + Fibertia

```text
Browser
  → GoFiber route
  → feature action/use case
  → repository/provider
  → Inertia response/page props
  → Vue page/component
```

Inertia dipakai untuk server-driven navigation, page props, forms, validation errors, redirect/flash, partial reload/deferred props, dan upload sesuai kemampuan Fibertia.

Public developer API tetap REST/JSON + OpenAPI; tidak semua internal Inertia action dipaksa menjadi public API.

Fibertia convention wajib mencakup:

- page name;
- shared props;
- auth/user/workspace context;
- validation errors;
- flash messages;
- pagination/filter;
- partial reload/deferred props;
- file upload;
- error boundary/exception handling.

## Vue + Nuxt UI

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

Tidak menggunakan `vue-router` karena navigation dikelola Inertia. Tidak menggunakan Pinia sebagai default karena state berasal dari page props/server dan local component state. State global client-heavy harus melalui ADR baru.

Nuxt UI menjadi baseline untuk AppShell, Sidebar, WorkspaceSwitcher, status, tables, forms, dialogs, timeline, empty/loading/error state, dan notifications.

## Zod runtime validation

Semua frontend runtime boundary wajib memakai Zod:

- form/request input;
- public API response;
- Inertia page props;
- shared props;
- localStorage/sessionStorage/URL query;
- external/provider payload yang ditampilkan ke UI.

Aturan:

- Network data dimulai sebagai `unknown`.
- Hindari `any` dan `as Type` untuk untrusted data.
- Gunakan strict object pada public boundary.
- `safeParse` untuk user input/recoverable data.
- `parse` untuk invariant/contract yang wajib fail early.
- Invalid schema tidak boleh di-silent-ignore.
- Global error boundary menampilkan actionable error dan mengirim telemetry ter-redact.
- Backend tetap memvalidasi ulang.

OpenAPI/Scalar tetap public contract; Zod adalah runtime guard frontend. Schema/fixture harus diuji agar tidak drift.

## Database: sqlx + prepared statements

- Repository berada di feature slice.
- `sqlx` untuk query/scanning.
- Prepared statement dibuat saat constructing repository untuk hot query yang dipanggil berulang.
- Statement ditutup saat graceful shutdown.
- Ukur manfaat dengan benchmark/metrics; query sesekali boleh tetap sederhana.
- Transaction boundary jelas.
- Semua query tenant-owned menerima `workspace_id`.
- External ID lookup memakai `workspace_id + public_id`.
- Jangan `SELECT *` untuk response model.

Contoh pola:

```go
const findMessageByPublicID = `
    SELECT id, public_id, workspace_id, status, created_at
    FROM messages
    WHERE workspace_id = $1 AND public_id = $2
`

type MessageRepository struct {
    find *sqlx.Stmt
}
```

## IDs

Internal:

```text
id: BIGSERIAL/BIGINT
public_id: TEXT/VARCHAR(21), nanoid(21)
```

External:

```text
id: string = public_id
public_id: tidak ada
```

Internal `id` dan `public_id` tidak boleh muncul di API, Inertia props, UI network, webhook, public logs, atau analytics.

## goose

- Forward-only ordered migration.
- Satu migration satu perubahan logis.
- Hindari long lock.
- Index besar memakai strategi aman.
- Test pada empty dan realistic DB.
- Deploy backward-compatible satu schema version.

## zerolog

Structured fields: request ID, safe workspace ID, channel ID, message ID, provider, operation, duration, error code.

Dilarang logging token, password, session secret, full phone, full message body, raw payload tanpa redaction.

## Valkey

Use cases: queue sederhana/worker coordination, ephemeral cache, rate limit, short-lived idempotency, distributed lock, session/flash jika diperlukan.

Valkey bukan source of truth untuk messages/billing. Semua cache/lock/idempotency memiliki TTL dan namespaced key.

Jika kebutuhan durability/stream semakin kompleks, evaluasi NATS JetStream/managed queue lewat ADR.

## Podman

Local stack:

- PostgreSQL;
- Valkey;
- MinIO/S3-compatible storage;
- Mailpit;
- mock provider/webhook fixture;
- optional whatsmeow worker disabled by default.

Commands:

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

## Quality tools

- `go vet`.
- `staticcheck`.
- `govulncheck`.
- frontend lint/typecheck.
- Vitest.
- Playwright.
- Testcontainers.
- secret/dependency/image scan.
- OpenAPI breaking-change check.
