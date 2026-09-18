# 17 — Frontend Schema Validation dengan Zod

> Dokumen ini adalah kontrak wajib untuk semua data yang masuk/keluar dari layer Vue/Inertia. Tujuan utamanya adalah **crash early, fail loudly, dan tidak membiarkan silent error**.

## Status

- Status: **Accepted / Mandatory**
- Scope: Vue components, Inertia pages/props, shared props, form requests, public API requests/responses, webhook/test payload yang ditampilkan di UI, dan client-side persisted state
- Technology: Vue 3 + TypeScript + Inertia + Fibertia + Zod

## Tujuan

Tanpa runtime validation, TypeScript hanya melindungi saat compile-time. Payload dari server/provider/user tetap dapat berubah atau rusak saat runtime.

Dengan Zod:

- request divalidasi sebelum dikirim;
- response divalidasi setelah diterima;
- Inertia page props divalidasi saat page diterima;
- shape yang berubah gagal di lokasi terdekat dengan sumber masalah;
- error dapat dilaporkan dengan schema/page/request ID;
- UI tidak diam-diam bekerja dengan data invalid.

## Prinsip

1. **Validate at boundaries** — validasi ketika data masuk/keluar boundary, bukan hanya di component terdalam.
2. **Parse, jangan percaya type assertion** — hindari `as SomeType` untuk data runtime.
3. **No silent fallback** — fallback hanya boleh untuk field yang memang optional dan didesain demikian.
4. **Fail loudly, fail usefully** — tampilkan error state yang jelas dan kirim telemetry tanpa PII.
5. **Schema dekat feature** — schema berada di slice/feature yang memiliki contract tersebut.
6. **Schema versioned** — perubahan breaking pada props/API memiliki version/changelog.
7. **Public ID contract** — semua resource external memiliki `id: string`; `public_id` dan internal BIGSERIAL tidak boleh muncul.

## Boundary yang wajib divalidasi

### 1. Form/request input

- Login/register.
- Channel onboarding.
- Contact import.
- Send message.
- Campaign.
- Template.
- Automation.
- Billing/settings.

Validasi dilakukan sebelum action Inertia/API dipanggil.

### 2. API response

Setiap response dari public API harus diparse sebelum digunakan oleh UI:

```ts
const result = MessageResponseSchema.parse(await response.json())
```

Jangan langsung memakai `response.json()` sebagai object terpercaya.

### 3. Inertia page props

Setiap Page memiliki schema props sendiri. Shared props juga memiliki schema terpisah.

```ts
const props = MessagePagePropsSchema.parse(usePage().props)
```

Jangan menggunakan satu schema besar untuk seluruh aplikasi jika page tidak membutuhkan semua data.

### 4. Shared props

Minimal schema untuk:

- authenticated user;
- current workspace;
- permissions;
- flash messages;
- validation errors;
- feature flags;
- environment/build metadata.

### 5. Client-persisted state

Data dari localStorage, sessionStorage, IndexedDB, URL query, dan browser event harus divalidasi sebelum digunakan.

## Struktur schema frontend

```text
/resources/js
  /features
    /messages
      schemas.ts
      types.ts
      api.ts
      pages.ts
      components/
    /channels
      schemas.ts
      api.ts
    /campaigns
      schemas.ts
      api.ts
  /shared
    /schemas
      common.ts
      auth.ts
      inertia.ts
      errors.ts
    /lib
      parse.ts
      telemetry.ts
```

Jika feature berada di backend VSA, struktur frontend mengikuti feature tersebut. Jangan membuat satu folder `schemas` raksasa yang tidak jelas ownership-nya.

## Schema conventions

### Naming

- `MessageSchema` — domain/view object.
- `MessageResponseSchema` — public API response.
- `CreateMessageRequestSchema` — request payload.
- `MessagePagePropsSchema` — Inertia page props.
- `MessageListItemSchema` — list projection.
- `MessageEventSchema` — event payload.

### Object strictness

Gunakan object strict untuk boundary publik agar field yang tidak dikenal tidak silently stripped:

```ts
import { z } from 'zod'

export const MessageSchema = z.object({
  id: z.string().min(1),
  status: z.enum(['created', 'queued', 'sending', 'sent', 'delivered', 'read', 'failed', 'canceled']),
  created_at: z.string().datetime(),
}).strict()
```

Gunakan `.passthrough()` hanya jika menerima extension fields yang memang dirancang. Jangan menggunakan coercion tanpa alasan yang terdokumentasi.

### Shared primitives

```ts
export const PublicIdSchema = z.string().min(21).max(21)
export const PaginationSchema = z.object({
  page: z.number().int().positive(),
  per_page: z.number().int().positive().max(100),
  total: z.number().int().nonnegative(),
}).strict()
```

Jika format NanoID berubah, ubah schema/contract secara eksplisit melalui decision record; jangan diam-diam menerima internal ID numerik.

## Request validation

```ts
export const CreateMessageRequestSchema = z.object({
  channel_id: PublicIdSchema,
  to: z.string().min(1),
  type: z.enum(['text', 'template', 'media']),
  message: z.string().max(60_000).optional(),
  template: z.object({
    name: z.string().min(1),
    language: z.string().min(2),
    variables: z.record(z.string(), z.string()).default({}),
  }).strict().optional(),
}).strict()

export type CreateMessageRequest = z.infer<typeof CreateMessageRequestSchema>
```

Flow:

1. Form state dikumpulkan.
2. `safeParse` dijalankan untuk menampilkan validation error yang actionable.
3. Hanya data valid yang dikirim.
4. Server tetap melakukan validasi ulang—Zod bukan pengganti backend validation.

## Response validation

```ts
export const MessageResponseSchema = z.object({
  data: MessageSchema,
  request_id: z.string().min(1),
}).strict()

export async function getMessage(id: string) {
  const response = await fetch(`/v1/messages/${encodeURIComponent(id)}`)
  const json: unknown = await response.json()

  if (!response.ok) {
    throw ApiErrorResponseSchema.parse(json)
  }

  return MessageResponseSchema.parse(json)
}
```

Gunakan `unknown`, bukan `any`, pada data hasil network sebelum parse.

## Inertia props validation

Fibertia/Inertia page props memiliki kontrak per page:

```ts
export const MessagePagePropsSchema = z.object({
  message: MessageSchema,
  permissions: z.object({
    can_reply: z.boolean(),
    can_resolve: z.boolean(),
  }).strict(),
  flash: FlashSchema,
}).strict()

export type MessagePageProps = z.infer<typeof MessagePagePropsSchema>
```

Page component:

```ts
const page = usePage()
const props = MessagePagePropsSchema.parse(page.props)
```

Page props hanya berisi external DTO. Jangan memasukkan database model yang mengandung `id BIGSERIAL`, `public_id`, atau `workspace_id` internal.

## `parse` vs `safeParse`

### Gunakan `safeParse` ketika

- memvalidasi input user;
- memvalidasi query string;
- memvalidasi localStorage yang bisa rusak;
- memproses response yang error dan ingin menampilkan fallback terkontrol.

### Gunakan `parse` ketika

- Inertia page props adalah kontrak internal yang wajib;
- API response sukses harus memenuhi schema;
- developer invariant harus gagal segera;
- data invalid tidak boleh diproses lebih lanjut.

Error boundary global harus menangkap exception `parse`, menampilkan halaman error yang berguna, dan mengirim telemetry ter-redact. “Crash early” bukan berarti layar putih tanpa informasi.

## Error reporting

Validation error menyertakan:

- schema name/version;
- feature/page;
- source (`inertia_props`, `api_response`, `form`, `storage`);
- HTTP status;
- request ID jika tersedia;
- issue path/code.

Jangan mengirim:

- API key/token;
- message body;
- full phone number;
- raw PII;
- seluruh payload tanpa redaction.

Contoh:

```ts
try {
  return MessageResponseSchema.parse(json)
} catch (error) {
  reportSchemaFailure({
    schema: 'MessageResponseSchema@v1',
    source: 'api_response',
    requestId,
    error,
  })
  throw error
}
```

## Schema source of truth

- Public API contract tetap didokumentasikan di OpenAPI dan ditampilkan melalui Scalar.
- Zod adalah runtime guard di frontend.
- Inertia props memiliki contract schema sendiri karena bukan public REST response.
- Hindari dua schema yang berbeda untuk endpoint yang sama.
- Evaluasi generate Zod/types dari OpenAPI atau shared contract package bila drift mulai terjadi.
- Contract test memakai fixture yang sama untuk backend dan frontend bila memungkinkan.

## Testing requirements

### Unit

- valid payload diterima;
- missing field ditolak;
- wrong type ditolak;
- unknown field ditolak sesuai strictness;
- `public_id` tidak muncul sebagai external field;
- internal numeric ID ditolak pada external `id`.

### Contract

- Backend response cocok dengan Zod schema.
- OpenAPI example cocok dengan Zod schema.
- Inertia fixture cocok dengan Page Props schema.
- Breaking schema change memerlukan review.

### E2E

- Invalid form tidak mengirim request.
- Invalid API response menghasilkan error state/telemetry.
- Invalid Inertia props tidak dipakai diam-diam.
- Corrupt localStorage dihapus atau masuk recovery flow.
- Error state dapat dipahami user.

## Handoff checklist

- [ ] Setiap public request memiliki Zod schema.
- [ ] Setiap public response memiliki Zod schema.
- [ ] Setiap Inertia Page memiliki props schema.
- [ ] Shared props memiliki schema.
- [ ] Network JSON dimulai dari `unknown`.
- [ ] Tidak ada `as Type` untuk untrusted runtime data.
- [ ] Tidak ada `any` pada boundary.
- [ ] Error parse tidak disenyapkan.
- [ ] Schema failure memiliki telemetry ter-redact.
- [ ] DTO mengikuti ID convention: `id: string`, tanpa `public_id`/internal `id`.
- [ ] Schema test dan contract fixture tersedia.
