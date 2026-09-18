# 16 — ID Convention dan Public Contract

> Dokumen ini adalah **kontrak wajib** antara database, backend, API, Inertia, frontend, webhook, dan tim engineering. Jangan menginterpretasikan field `id` external sebagai database `id`.

## Status

- Status: **Accepted / Mandatory**
- Scope: seluruh entity/resource yang memiliki identity
- Berlaku untuk: API response, API request path/body, Inertia props, UI network payload, customer webhook, public logs, analytics event, dan dokumentasi developer

## Keputusan utama

Setiap resource memiliki dua identifier di internal backend/database:

```text
id:        BIGSERIAL / BIGINT       # internal database identity
public_id: TEXT / VARCHAR(21)       # public identity source, nanoid(21)
```

Namun external contract hanya memiliki satu identifier:

```text
id: TEXT = internal public_id
```

`public_id` **tidak pernah ikut dikirim keluar**.

## Contoh konkret

### Database row

```text
messages
────────────────────────────────────────
id          = 98231
public_id   = V1StGXR8_Z5jdHi6B-myT
status      = delivered
workspace_id = 17
```

### Internal Go model

Model database/internal boleh merepresentasikan dua field:

```go
type Message struct {
    ID          int64  `db:"id"`
    PublicID    string `db:"public_id"`
    WorkspaceID int64  `db:"workspace_id"`
    Status      string `db:"status"`
}
```

### External response

```json
{
  "id": "V1StGXR8_Z5jdHi6B-myT",
  "status": "delivered"
}
```

Response berikut **dilarang**:

```json
{
  "id": 98231,
  "public_id": "V1StGXR8_Z5jdHi6B-myT",
  "status": "delivered"
}
```

Response berikut juga **dilarang**:

```json
{
  "id": "98231",
  "status": "delivered"
}
```

`id` external bukan hasil cast `BIGSERIAL` menjadi string. Nilainya adalah `public_id`.

## Mapping rule

```text
                    INTERNAL                 EXTERNAL
Database             id: BIGINT              tidak keluar
Database             public_id: TEXT        sumber value

API response                                  id: TEXT = public_id
Inertia props                                  id: TEXT = public_id
UI network payload                             id: TEXT = public_id
Customer webhook                              id: TEXT = public_id
```

## Request flow

### URL/path parameter

```http
GET /v1/messages/V1StGXR8_Z5jdHi6B-myT
```

Backend:

1. Menerima `id` external sebagai `string`.
2. Memvalidasi format/panjang identifier.
3. Mengambil row berdasarkan `workspace_id + public_id`.
4. Menggunakan internal `id BIGINT` untuk join/query lanjutan.
5. Mengubah hasil ke external DTO.
6. Mengembalikan `id = public_id`.

### Request body

```json
{
  "message_id": "V1StGXR8_Z5jdHi6B-myT"
}
```

Nama field request boleh mengikuti konteks (`message_id`, `channel_id`, dan lain-lain), tetapi nilainya selalu public ID berbentuk string.

## DTO wajib

Database model tidak boleh dikembalikan langsung dari handler atau Inertia action.

```go
type MessageDTO struct {
    ID     string `json:"id"`
    Status string `json:"status"`
}

func ToMessageDTO(m Message) MessageDTO {
    return MessageDTO{
        ID:     m.PublicID,
        Status: m.Status,
    }
}
```

Untuk Inertia:

```go
type MessagePageProps struct {
    Message MessageDTO `json:"message"`
}
```

Jangan memasukkan struct database yang mengandung `ID`, `PublicID`, atau `WorkspaceID` ke page props secara langsung.

## Repository/query rule

Internal query boleh menggunakan `id` untuk performa:

```sql
SELECT id, public_id, workspace_id, status
FROM messages
WHERE workspace_id = $1
  AND public_id = $2;
```

API/resource lookup wajib selalu tenant-scoped:

```text
WHERE workspace_id = current_workspace
  AND public_id = external_id
```

Jangan lookup public ID tanpa workspace/authorization check.

## Generation rule

- Generate `public_id` menggunakan NanoID 21 karakter di application layer sebelum insert.
- `public_id` disimpan sebagai `TEXT` atau `VARCHAR(21)` sesuai convention schema.
- Database memiliki unique constraint pada `public_id`.
- Collision ditangani dengan retry insert.
- `public_id` immutable; tidak boleh berubah sepanjang lifecycle resource.
- Internal `id BIGSERIAL` juga immutable.
- Tidak memakai internal ID sebagai fallback public ID.
- Tidak mengungkap sequence, row count, atau mapping internal ID.

Contoh migration:

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

## Serialization policy

### Allowed

- Explicit DTO.
- Explicit mapper.
- Query scan ke internal model lalu map ke DTO.
- API documentation yang hanya memperlihatkan public contract.

### Forbidden

- `SELECT *` lalu auto-serialize.
- Mengembalikan database model dari Fiber handler.
- Mengirim `public_id` untuk “membantu debugging”.
- Mengirim `id` BIGSERIAL ke Vue/Inertia.
- Menaruh internal ID di HTML/data attribute yang dapat dibaca browser.
- Menaruh internal ID di response header atau client-visible error.
- Menggunakan public ID sebagai satu-satunya authorization check.

## Boundary checklist

### Backend

- [ ] Semua response melewati DTO/serializer.
- [ ] Semua list item menggunakan `id` dari `PublicID`.
- [ ] Tidak ada field `public_id` pada response struct JSON.
- [ ] Internal `ID` tidak memiliki JSON tag yang dapat keluar.
- [ ] Query lookup memakai tenant scope.

### Fiber/API

- [ ] Path parameter bertipe string.
- [ ] Error response tidak membocorkan internal ID.
- [ ] OpenAPI hanya mendokumentasikan `id: string`.
- [ ] Scalar tidak menampilkan database schema internal.

### Fibertia/Inertia

- [ ] Page props hanya berisi DTO/view model.
- [ ] Shared props tidak mengandung database model.
- [ ] Partial reload tetap menggunakan external DTO.
- [ ] Form action mengirim public ID string.

### Vue/UI

- [ ] TypeScript resource menggunakan `id: string`.
- [ ] Tidak ada `public_id` di type/props public.
- [ ] Tidak ada internal numeric ID pada URL, DOM, request, atau analytics.
- [ ] Component tidak mengandalkan urutan/sequential ID.

### Webhook/analytics

- [ ] Customer webhook hanya menggunakan public `id` string.
- [ ] Event payload tidak mengandung internal `id`.
- [ ] Analytics properties tidak menyimpan internal database ID.
- [ ] Logs yang dapat diakses customer tidak menyimpan internal ID.

## Testing requirements

### Unit test mapper

```go
func TestToMessageDTOUsesPublicID(t *testing.T) {
    got := ToMessageDTO(Message{
        ID:       98231,
        PublicID: "V1StGXR8_Z5jdHi6B-myT",
        Status:   "delivered",
    })

    if got.ID != "V1StGXR8_Z5jdHi6B-myT" {
        t.Fatalf("unexpected id: %s", got.ID)
    }
}
```

### Contract test

Pastikan serialized JSON:

- memiliki `id` bertipe string;
- nilai `id` sama dengan internal `public_id`;
- tidak memiliki key `public_id`;
- tidak memiliki nilai internal BIGSERIAL.

### Integration/e2e test

- Request menggunakan public ID berhasil.
- Request menggunakan internal numeric ID ditolak/tidak ditemukan.
- Workspace A tidak dapat membaca resource Workspace B walaupun mengetahui public ID.
- Inertia response tidak mengandung `public_id`, `workspace_id`, atau internal `id`.
- Webhook customer tidak mengandung internal identifier.

## Security clarification

NanoID bukan encryption dan bukan authorization. Public ID hanya berfungsi sebagai opaque identifier untuk menghindari sequential enumeration dan menjaga contract publik.

Authorization tetap wajib dilakukan berdasarkan:

```text
authenticated actor
+ current workspace/tenant
+ resource ownership/membership
+ role/permission
+ public_id lookup
```

## Handoff statement

> Di database, resource memiliki `id BIGSERIAL` dan `public_id nanoid(21)`. Internal `id` tidak pernah keluar dari backend. Semua external response, API request reference, Inertia props, UI network payload, customer webhook, dan public analytics menggunakan field `id` bertipe text yang nilainya berasal dari `public_id`. Field `public_id` tidak pernah dikirim keluar.
