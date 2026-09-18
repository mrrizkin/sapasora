# 07 — Domain Model dan Data

## Core entities

```text
User ──< WorkspaceMember >── Workspace
Workspace ──< Channel
Workspace ──< Contact ──< ConsentEvent
Workspace ──< Conversation ──< Message ──< MessageEvent
Workspace ──< Template
Workspace ──< Campaign ──< CampaignRecipient
Workspace ──< Automation ──< AutomationRun
Workspace ──< ApiKey
Workspace ──< WebhookEndpoint
Workspace ──< AuditLog
Workspace ──< Subscription ──< UsageLedger
Channel ──< ProviderCredential
```

## Entity responsibilities

### Workspace

Tenant boundary, billing owner, timezone, settings, retention policy.

### Channel

A logical WhatsApp connection. Menyimpan provider type, WABA/phone identifiers, status, capabilities, and health summary.

### ProviderCredential

Reference encrypted secret; tidak menyimpan raw token di application logs atau frontend. Mendukung rotation/revoke.

### Contact

Recipient profile dengan normalized phone/identifier, tags, locale, timezone, opt-in state, and source.

### ConsentEvent

Immutable event: opted_in, opted_out, source, timestamp, evidence reference, actor.

### Conversation

Thread per channel/contact/group context dengan status open/pending/resolved, assignee, priority, last message.

### Message

Logical outgoing/incoming message: sender/recipient, content type, status, provider ID, request ID, idempotency key, timestamps.

### MessageEvent

Append-only lifecycle events: accepted, queued, sent, delivered, read, failed, canceled.

### Template

Provider-scoped template name, language, category, version, approval status, components, variables.

### Campaign

Draft/approved/scheduled/running/paused/completed/canceled with audience snapshot and safety checks.

### Automation

Trigger, conditions, actions, version, enabled state, and run limits.

### AuditLog

Security-sensitive mutations: login, key create/revoke, channel connect/disconnect, export, campaign approval, role change.

## PostgreSQL guidelines

- Setiap entity memakai dual ID: internal `BIGSERIAL`/`BIGINT` untuk join/index/query cepat dan `public_id` `nanoid(21)` untuk URL/API/public reference.
- `public_id` wajib unik secara global per table/entity dan tidak boleh ditebak.
- `public_id` `nanoid(21)` dibuat di application layer sebelum insert; database unique constraint menjadi guard terakhir dan collision ditangani dengan retry.
- Internal `id` tidak pernah keluar dari backend boundary.
- External contract hanya memiliki `id: text/string` yang nilainya berasal dari `public_id`; field `public_id` tidak ikut dikirim.
- Aturan ini berlaku untuk API response, Inertia props, UI network payload, customer webhook, logs publik, dan analytics event.
- Mapping internal → external harus eksplisit melalui DTO/serializer, bukan `SELECT *` atau auto-serialization model database.
- `created_at`, `updated_at`; soft-delete hanya bila business need jelas.
- Unique constraint per tenant untuk natural keys.
- Foreign key dan index dirancang berdasarkan query aktual.
- JSONB hanya untuk provider metadata/raw payload terkontrol, bukan seluruh domain.
- Encrypt sensitive columns atau gunakan envelope encryption.
- Message/event tables dipertimbangkan partitioning setelah volume terbukti.

## Minimal schema sketch

```sql
create table workspaces (
  id bigserial primary key,
  public_id varchar(21) not null unique,
  name text not null,
  timezone text not null default 'Asia/Jakarta',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table channels (
  id bigserial primary key,
  public_id varchar(21) not null unique,
  workspace_id bigint not null references workspaces(id),
  provider text not null,
  status text not null,
  external_account_id text,
  external_phone_id text,
  capabilities jsonb not null default '{}',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index channels_workspace_idx on channels(workspace_id);
create unique index channels_workspace_external_idx
  on channels(workspace_id, provider, external_phone_id);
```

## Message status model

Internal status:

```text
created → queued → sending → sent → delivered → read
                         └──────→ failed
created/queued ────────────────→ canceled
```

Provider statuses disimpan sebagai metadata dan dipetakan ke status internal. Jangan mengasumsikan semua provider memberi status yang sama.

## External ID contract

```text
Database/backend model:
  id: BIGSERIAL          # internal only
  public_id: varchar(21) # internal storage/public identity source

API/Inertia/webhook DTO:
  id: string             # value copied from public_id
  # public_id: absent
```

Contoh serializer:

```go
type MessageDTO struct {
    ID     string `json:"id"`
    Status string `json:"status"`
}

func ToMessageDTO(m Message) MessageDTO {
    return MessageDTO{ID: m.PublicID, Status: m.Status}
}
```

`id` pada DTO bukan database `id`; `id` adalah alias eksternal untuk `public_id`.

## Data retention proposal

- Message content: configurable, default 90 hari setelah resolved.
- Delivery events: 180 hari untuk analytics dasar.
- Raw provider webhook: 7–30 hari, encrypted, akses terbatas.
- Media: lifecycle sesuai kebutuhan/provider; cleanup job wajib.
- Audit logs: 1 tahun untuk paid, minimum 90 hari untuk free.
- Provider credentials: selama channel aktif; hapus saat disconnect/delete.
- Unofficial sessions: hapus segera saat user meminta disconnect/delete.

Retention final harus melalui legal/privacy review dan komitmen provider.

## Data export/delete

Workspace owner dapat:

- export contacts/conversations/reports;
- delete contact dan consent evidence sesuai policy;
- disconnect channel;
- delete credential/session;
- request workspace deletion.

Deletion harus asynchronous, observable, idempotent, dan memiliki completion audit event.

## Analytics model

Jangan query raw message table untuk semua dashboard. Buat aggregate terjadwal:

- messages by day/status/channel;
- delivery rate;
- response time;
- campaign metrics;
- automation run/failure;
- usage/billing ledger.
