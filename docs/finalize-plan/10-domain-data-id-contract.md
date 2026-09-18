# 10 — Domain Model, Data, dan ID Contract

## Domain entities

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

## Responsibilities

- Workspace: tenant, billing owner, timezone, settings, retention.
- Channel: provider, WABA/phone IDs, status, capabilities, health.
- ProviderCredential: encrypted secret reference, rotation/revoke.
- Contact: normalized identifier, tags, locale/timezone, consent.
- ConsentEvent: immutable opt-in/opt-out/source/time/evidence/actor.
- Conversation: thread, assignment, priority, SLA, status.
- Message: logical inbound/outbound content/status/provider/request/idempotency.
- MessageEvent: append-only lifecycle.
- Template: provider scope, language, category, components, approval/version.
- Campaign: audience snapshot, approval, schedule, status, safety checks.
- Automation: trigger/condition/action/version/run limits.
- AuditLog: security-sensitive mutation.

## Mandatory dual ID convention

Internal database/backend model:

```text
id:        BIGSERIAL/BIGINT       # internal only
public_id: TEXT/VARCHAR(21)       # nanoid(21)
```

External contract:

```text
id: TEXT = public_id
public_id: absent
```

Internal `id` tidak pernah keluar dari backend. Field `public_id` juga tidak pernah keluar. Berlaku untuk API response, API request reference, Inertia props, UI network payload, customer webhook, public logs, analytics, dan documentation examples.

`id` external bukan hasil cast BIGSERIAL ke string. Nilainya adalah NanoID dari `public_id`.

## Example

Database:

```text
id: 98231
public_id: V1StGXR8_Z5jdHi6B-myT
```

External:

```json
{
  "id": "V1StGXR8_Z5jdHi6B-myT",
  "status": "delivered"
}
```

Forbidden:

```json
{
  "id": 98231,
  "public_id": "V1StGXR8_Z5jdHi6B-myT"
}
```

## SQL schema baseline

```sql
CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    public_id VARCHAR(21) NOT NULL UNIQUE,
    workspace_id BIGINT NOT NULL REFERENCES workspaces(id),
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX messages_workspace_public_id_idx
    ON messages (workspace_id, public_id);
```

Generate `public_id` in application layer before insert. Unique constraint is final guard; collision retries. `public_id` immutable.

## Lookup/mapping

Request:

```http
GET /v1/messages/V1StGXR8_Z5jdHi6B-myT
```

Backend mencari dengan `workspace_id + public_id`, menggunakan internal `id` untuk join/query, lalu map ke DTO.

```go
type MessageDTO struct {
    ID     string `json:"id"`
    Status string `json:"status"`
}

func ToMessageDTO(m Message) MessageDTO {
    return MessageDTO{ID: m.PublicID, Status: m.Status}
}
```

Database model tidak boleh langsung diserialize.

## Message lifecycle

```text
created → queued → sending → sent → delivered → read
                         └──────→ failed
created/queued ────────────────→ canceled
```

Provider status dipetakan ke internal status; provider-specific metadata tetap terbatas/ter-redact.

## Data policies

- `workspace_id` ada pada tenant-owned tables.
- Public ID lookup selalu authorization/tenant-scoped.
- JSONB hanya untuk provider metadata/raw payload terkontrol.
- Message content default proposal: 90 hari setelah resolved.
- Delivery event: 180 hari.
- Raw webhook: 7–30 hari encrypted.
- Media mengikuti lifecycle/provider dan cleanup job.
- Audit log: 1 tahun paid, minimum 90 hari free.
- Credential selama channel aktif; hapus saat disconnect/delete.
- Unofficial session hapus saat disconnect/delete request.

Retention final memerlukan legal/privacy review.

## Export/delete

Owner dapat export contact/conversation/report dan request delete workspace. Deletion asynchronous, observable, idempotent, serta memiliki audit completion event.

## Security note

NanoID bukan encryption/authorization. Authorization tetap membutuhkan actor + workspace + permission + resource ownership.
