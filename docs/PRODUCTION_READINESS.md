# Production Readiness Audit

> **Product blueprint:** [`PRODUCT_MASTER_PLAN.md`](./PRODUCT_MASTER_PLAN.md)
>
> **Full feature requirements:** [`FULL_FEATURE_SPECIFICATION.md`](./FULL_FEATURE_SPECIFICATION.md)

## Ringkasan

Dokumen ini mencatat hasil audit statis terhadap project Sapasora berdasarkan source code, konfigurasi, dan hasil build/test lokal.

**Status saat ini: belum siap untuk public production.** Project masih cocok untuk development atau internal testing terbatas.

Temuan paling kritis adalah:

1. Implementasi Telegram belum selesai dan beberapa endpoint masih melakukan `panic("unimplemented")`.
2. Proteksi CSRF efektifnya tidak bekerja dan revocation device token dapat dilewati.
3. API key dan device token masih diperlakukan sebagai secret plaintext.
4. Ownership/tenant isolation dan dispatch provider belum diterapkan secara konsisten.
5. Input validation, rate limiting, dan proteksi SSRF belum memadai.
6. UI, test suite, dan konfigurasi Docker masih berorientasi development.

---

## Status Verifikasi

| Pemeriksaan | Status | Catatan |
|---|---:|---|
| `pnpm type-check` | Lulus | TypeScript tidak memiliki error type-check |
| `pnpm build:assets` | Lulus | Asset frontend berhasil dibuild |
| `pnpm test:unit -- --run` | Gagal | Tidak ada file unit test frontend |
| `pnpm format:check` | Gagal | 3 file generated route belum terformat |
| `go test ./...` | Gagal | TDLib native header tidak tersedia, error vet, dan test config gagal |
| `go vet ./...` | Gagal | TDLib header tidak tersedia dan terdapat error format logger pada `support/sync`/`wayfinder` |
| `pnpm exec oxlint . -D correctness` | Lulus | 0 warning dan 0 error |
| `pnpm audit --prod` | Warning | 1 advisory low pada `esbuild` 0.27.7 melalui dependency frontend |
| CI/CD | Belum tersedia | Tidak ditemukan workflow CI di repository |
| E2E test | Gagal | Playwright tidak menemukan test dan exit code 1 |

---

# Temuan P0 — Blocker Production

## P0-01 — Telegram belum selesai

**Lokasi:** `internal/modules/telegram/service.go`

Method berikut masih menggunakan `panic("unimplemented")`:

- `GetContacts`
- `GetUser`
- `Logout`
- `SendAudio`
- `SendButton`
- `SendChatPresence`
- `SendContact`
- `SendDocument`
- `SendImage`
- `SendList`
- `SendLocation`
- `SendSticker`
- `SendVideo`

### Dampak

Request ke endpoint Telegram dapat menghasilkan HTTP 500 atau panic ketika fitur tersebut dipanggil.

### Rekomendasi

- Implementasikan seluruh method Telegram.
- Tambahkan integration test dengan TDLib.
- Jika Telegram belum menjadi bagian release pertama, disable endpoint Telegram secara eksplisit dan return `501 Not Implemented`, bukan panic.

---

## P0-02 — API key dan device token terekspos

**Lokasi:**

- `internal/modules/apikey/apikey.go`
- `internal/modules/devicetoken/devicetoken.go`
- `internal/app/http/middleware/authentication.go`

Field secret masih menggunakan JSON tag yang mengekspos nilai token:

```go
Key   string `json:"key"`
Token string `json:"token"`
```

Token juga dicari langsung berdasarkan nilai plaintext di database.

### Masalah

- Secret tersimpan plaintext.
- Endpoint list/get berpotensi mengembalikan secret.
- Status `active/inactive` belum divalidasi oleh middleware.
- Device token memiliki `expired_at`, tetapi field tersebut belum divalidasi oleh middleware; API key saat ini tidak memiliki field expiry.
- Revoke/soft-delete device token juga perlu memblokir lookup join secara eksplisit.
- Belum ada rotasi token dan revoke yang terstandarisasi.
- Belum ada audit penggunaan credential.

### Rekomendasi

- Hash token di database menggunakan HMAC atau hash kriptografis yang sesuai.
- Tampilkan token hanya satu kali saat dibuat.
- Gunakan prefix/key ID untuk lookup dan hash untuk verifikasi.
- Validasi status dan expiry setiap request.
- Tambahkan revoke, rotate, last-used-at, dan audit log.
- Pastikan response list/get memakai DTO yang tidak mengandung secret.

---

## P0-03 — Ownership dan tenant isolation belum diterapkan

**Lokasi contoh:**

- `internal/modules/device/repository.go`
- `internal/modules/account/repository.go`
- `internal/modules/apikey/repository.go`
- `internal/app/policies/`

Repository masih mengambil data secara global, sementara policy umumnya hanya mengecek permission boolean. Resource dan owner belum digunakan untuk membatasi akses.

### Dampak

User atau API key yang memiliki permission tertentu berpotensi membaca atau mengubah resource milik user lain.

### Rekomendasi

- Definisikan model tenancy secara eksplisit.
- Tambahkan filter owner/tenant pada semua query.
- Validasi bahwa API key/device token hanya mengakses resource yang menjadi scope-nya.
- Uji akses lintas user menggunakan authorization integration test.
- Jangan mengandalkan ID publik sebagai pengganti authorization.

---

## P0-04 — Input validation dan proteksi panic belum memadai

Sebagian besar controller menggunakan `ctx.BodyParser` langsung tanpa validasi schema.

**Lokasi contoh:** `internal/app/http/controllers/` dan `internal/modules/whatsapp/service.go`

Terdapat slicing string tanpa pemeriksaan panjang input:

- `internal/modules/whatsapp/service.go:381`
- `internal/modules/whatsapp/service.go:645`
- `internal/modules/whatsapp/service.go:739`
- `internal/modules/whatsapp/service.go:985`
- `internal/modules/whatsapp/service.go:1140`

### Dampak

Input pendek atau malformed dapat menyebabkan panic. Input media besar juga dapat menyebabkan konsumsi memory dan bandwidth berlebihan.

### Rekomendasi

- Gunakan `BodyParserValidate` secara konsisten.
- Tambahkan schema validation untuk semua request.
- Validasi ukuran media sebelum decode base64.
- Validasi nomor telepon, URL, MIME type, enum, dan array length.
- Tambahkan body limit dan timeout khusus untuk endpoint media.
- Return `400` atau `422`, bukan panic/`500`.

---

## P0-05 — Webhook berisiko SSRF

**Lokasi:** `internal/modules/whatsapp/whatsmeow.go`

Temuan:

- URL webhook berasal dari input konfigurasi device.
- Redirect diizinkan hingga 15 kali.
- TLS verification dimatikan melalui `InsecureSkipVerify: true`.
- Tidak ada allowlist hostname/IP.
- Pengiriman dilakukan asynchronous tanpa delivery state.

### Dampak

Webhook dapat disalahgunakan untuk mengakses endpoint internal, metadata cloud, atau host private network.

### Rekomendasi

- Hapus `InsecureSkipVerify`.
- Validasi URL hanya `https` pada production.
- Blokir localhost, private IP, link-local, loopback, dan metadata address.
- Gunakan allowlist domain jika memungkinkan.
- Batasi redirect atau disable redirect.
- Tambahkan signing HMAC pada payload webhook.
- Gunakan queue dengan retry, backoff, dead-letter, dan delivery log.

---

# Temuan P1 — Wajib untuk Production yang Stabil

## P1-01 — Authentication belum lengkap

**Lokasi:** `internal/app/http/controllers/auth/auth.go`

Saat ini hanya tersedia login dan logout.

Belum tersedia:

- Password reset/forgot password
- Email verification
- MFA/2FA
- Login rate limiting
- Account lockout
- Session revocation
- Audit login
- Remember-me, walaupun field tersedia di UI

### Rekomendasi

Tetapkan scope authentication release pertama. Minimal tambahkan brute-force protection, session regeneration setelah login, password policy, dan audit log.

---

## P1-02 — Potensi mass assignment pada API key update

**Lokasi:** `internal/app/http/controllers/apikey/apikey.go`

Pada proses update, object model diparse ulang dari request setelah field internal diset. Pola ini dapat membuka kemungkinan perubahan field sensitif seperti key, status, atau user ID.

### Rekomendasi

- Parse ke DTO khusus.
- Salin hanya field yang memang boleh diubah.
- Jangan pernah parse request langsung ke model persistence.
- Tambahkan test untuk memastikan key, owner, dan status tidak dapat diubah tanpa endpoint khusus.

---

## P1-03 — Pagination salah pada beberapa repository

**Lokasi:**

- `internal/modules/account/repository.go:48`
- `internal/modules/device/repository.go:51`
- `internal/modules/devicetoken/repository.go:36`
- `internal/modules/role/repository.go:48`

Implementasi menggunakan:

```go
Offset(page - 1)
```

Seharusnya:

```go
Offset((page - 1) * limit)
```

API key repository sudah menggunakan formula yang benar.

### Rekomendasi

- Samakan helper pagination di semua repository.
- Batasi nilai minimum dan maksimum `page` serta `limit`.
- Tambahkan test untuk page 1, 2, dan page terakhir.

---

## P1-04 — UI masih banyak placeholder

### Fitur yang belum memiliki route backend

Sidebar mereferensikan halaman berikut, tetapi route web belum tersedia:

- `/phonebook`
- `/messages`
- `/send`
- `/templates`
- `/recurring`
- `/autoresponder`
- `/admins`
- `/documentation`
- `/settings`

**Lokasi:** `resources/js/component/sidebar/links.ts`

### Fitur device yang masih dummy

- Delete device hanya menampilkan toast/log.
- Bulk delete hanya menunggu satu detik tanpa request API.
- Delete pada detail device belum memanggil API.
- Edit device mengirim POST ke route edit halaman GET, bukan endpoint update API.
- Pilihan `SMS` dan `Email` tersedia di UI, tetapi backend hanya mendukung `whatsapp` dan `telegram`.

**Lokasi:**

- `resources/js/component/devices/devices_delete_modal.vue`
- `resources/js/pages/device/index.vue`
- `resources/js/pages/device/show.vue`
- `resources/js/pages/device/edit.vue`
- `resources/js/pages/device/components/form.vue`

---

## P1-05 — Security middleware belum konsisten

**Lokasi:** `internal/app/routes/pipeline.go`

API dan WebSocket menggunakan konfigurasi CORS wildcard:

```go
AllowOrigins: "*"
```

API juga hanya mengizinkan header:

```go
Origin, Content-Type, Accept
```

Header `Authorization` belum dicantumkan.

Selain itu:

- Helmet hanya diterapkan pada web pipeline.
- Rate limiting belum ada.
- Timeout middleware belum ada.
- Body limit khusus media belum ada.
- Endpoint WebSocket ping belum memiliki authentication.

### Rekomendasi

Gunakan allowlist origin, tambahkan `Authorization` hanya jika memang diperlukan, implementasikan rate limiting berbasis IP/key, dan tetapkan timeout/body limit per endpoint.

---

## P1-06 — Storage media belum production-safe

**Lokasi:** `internal/modules/whatsapp/client.go`

File event ditulis ke path berdasarkan lokasi executable:

```text
<executable-directory>/files/user_<device-id>/
```

Path tersebut bukan volume Docker yang dikonfigurasi.

### Risiko

- File hilang saat container diganti.
- Tidak ada cleanup/retention policy.
- Tidak ada quota per device.
- Tidak ada integrasi storage terpusat untuk deployment multi-instance.

### Rekomendasi

Gunakan storage abstraction yang sudah tersedia, lalu implementasikan S3/MinIO dengan retention, quota, dan cleanup job.

---

## P1-07 — Horizontal scaling belum didukung

Client WhatsApp dan TDLib disimpan di memory process. State koneksi, kill channel, dan device info tidak terdistribusi.

### Dampak

Menjalankan lebih dari satu instance dapat menyebabkan:

- Device terkoneksi di beberapa instance.
- Event duplicate.
- Webhook duplicate.
- State koneksi tidak konsisten.

### Rekomendasi

Pilih salah satu strategi:

1. Batasi deployment menjadi single instance dengan dokumentasi jelas; atau
2. Implementasikan distributed lock, device ownership, queue, dan leader election.

---

# Temuan P2 — Hardening dan Operasional

## P2-01 — Health check belum memeriksa dependency

Endpoint `/api/v1/health` hanya mengembalikan:

```json
{"status":"ok"}
```

Tidak memeriksa database, session store, TDLib, WhatsApp, storage, atau queue.

### Rekomendasi

Tambahkan:

- `/livez` untuk liveness process.
- `/readyz` untuk readiness dependency.
- Status per dependency.
- Version/build information.

---

## P2-02 — Observability belum lengkap

Logging sudah tersedia, tetapi belum ditemukan implementasi yang memadai untuk:

- Metrics/Prometheus
- Distributed tracing
- Error tracking
- Alerting
- Webhook delivery metrics
- Message success/failure metrics
- Device connection metrics
- Audit log

Request ID sudah digunakan oleh server, tetapi perlu dipastikan masuk ke semua log dan response/error context.

---

## P2-03 — Error response perlu disanitasi

**Lokasi:** `platform/server/error.go`

Error dikembalikan menggunakan `err.Error()`. Error database atau internal dapat bocor ke client, terutama jika error tersebut berisi detail query atau dependency.

### Rekomendasi

- Gunakan error code publik yang stabil.
- Log detail internal di server.
- Kirim message generik ke client pada production.
- Sertakan request ID untuk troubleshooting.
- Bedakan response HTML, Inertia, dan API secara konsisten.

---

## P2-04 — Configuration dan deployment hardening belum tersedia

**Lokasi:**

- `.env.example`
- `docker-compose.yml`
- `Dockerfile`
- `internal/config/config.go`

Temuan:

- Default environment masih development.
- `APP_DEBUG=true` pada contoh konfigurasi.
- Credential PostgreSQL default.
- Container belum menggunakan user non-root.
- Tidak ada konfigurasi TLS/reverse proxy.
- Error konfigurasi di `NewAppConfig` diabaikan.
- Credential Telegram ditandai required tetapi kegagalan load tidak dipropagasi.
- Belum ada environment preset staging/production.

### Rekomendasi

- Fail fast jika secret/config wajib tidak valid.
- Pisahkan konfigurasi development, staging, dan production.
- Gunakan secret manager.
- Jalankan container sebagai non-root.
- Tambahkan healthcheck aplikasi.
- Dokumentasikan TLS, proxy, backup, restore, dan rollback.

---

## P2-05 — Database operation dan migration hardening

Migration runner sudah tersedia, tetapi deployment production masih memerlukan:

- Migration lock untuk mencegah dua instance menjalankan migration bersamaan.
- Backup sebelum migration destruktif.
- Prosedur rollback yang diuji.
- Seed role dan admin yang deterministic.
- Foreign key delete policy yang jelas.
- Index dan query performance review.
- Backup/restore test berkala.

Toolbox admin saat ini masih bersifat interaktif sehingga perlu mekanisme bootstrap admin yang aman untuk deployment otomatis.

---

# Temuan Tambahan — Audit Putaran Kedua (Cross-Cutting Review)

Audit kedua memeriksa ulang alur request, authentication, authorization, lifecycle device, konfigurasi, concurrency, database, frontend, dependency, dan deployment.

## Temuan P0 Tambahan

### P0-06 — Proteksi CSRF efektifnya tidak bekerja

**Lokasi:** `internal/app/routes/pipeline.go`

Status remediasi T0.2: middleware CSRF web sekarang memakai pola synchronizer-backed double-submit.

```go
KeyLookup:  fmt.Sprintf("header:%s", csrfKey)
Session:    sessionStore
CookieName: cookieName // readable same-origin source, not the request extractor
```

Browser memang otomatis mengirim cookie pada request cross-site, tetapi request
unsafe juga wajib membawa header CSRF yang hanya ditambahkan oleh JavaScript
same-origin. Token header harus sama dengan cookie dan token yang tersimpan di
server session.

API dan channel pipeline tidak menggunakan middleware CSRF web. Endpoint
webhook yang memakai signature provider harus tetap berada di pipeline tersebut
dan memvalidasi signature serta replay protection secara independen.

---

### P0-07 — Revocation device token dapat dilewati

**Lokasi:** `internal/modules/device/repository.go`

Authentication device memakai query join berikut secara konseptual:

```sql
JOIN m_device_tokens AS mdt ON mdt.device_id = m_devices.id
WHERE mdt.token = ?
```

Query tidak memvalidasi:

- `mdt.status`
- `mdt.expired_at`
- `mdt.deleted_at`

Karena `DeviceToken` menggunakan soft delete dan join ditulis manual, token yang dihapus secara soft-delete berpotensi tetap ditemukan dan diterima oleh middleware.

### Rekomendasi

- Validasi status dan expiry di query authentication.
- Tambahkan `mdt.deleted_at IS NULL` secara eksplisit.
- Simpan hash token, bukan token plaintext.
- Uji token aktif, inactive, expired, revoked, dan token dari device/user lain.

---

### P0-08 — Device token dapat diterbitkan lintas owner

**Lokasi:** `internal/app/http/controllers/devicetoken/devicetoken.go`

Request menerima `DeviceID` dan `UserID` terpisah. Controller hanya memastikan keduanya ada, tetapi tidak memastikan device memang dimiliki user tersebut. Hal serupa berlaku pada pembuatan API key dan device.

### Dampak

User dengan permission global dapat membuat credential untuk resource milik user lain. Jika ownership belum diterapkan, credential tersebut menjadi akses lintas tenant.

### Rekomendasi

- Derive owner dari subject yang terautentikasi, bukan dari request body.
- Validasi `device.UserID == user.ID` atau gunakan authorization service berbasis resource.
- Tambahkan unique/scope policy untuk relasi token-device-user.
- Tambahkan integration test lintas owner.

---

### P0-09 — Gateway tidak konsisten memilih provider berdasarkan device type

**Lokasi:** `internal/modules/gateway/service.go`

Hanya sebagian method melakukan switch `DeviceType`. `GetContacts`, `GetUser`, `Logout`, dan sebagian besar method send selalu memanggil `whatsappService`, termasuk ketika device bertipe Telegram.

### Dampak

Request menggunakan Telegram device dapat gagal, menggunakan client WhatsApp yang salah, atau menghasilkan perilaku tidak terdefinisi. Ini terjadi di luar daftar method Telegram yang masih `panic("unimplemented")`.

### Rekomendasi

- Buat dispatch matrix yang eksplisit untuk setiap capability/provider.
- Return `501 Not Implemented` untuk capability yang belum tersedia.
- Jangan silently fallback ke WhatsApp.
- Tambahkan contract test untuk setiap method dan setiap `DeviceType`.

---

## Temuan P1 Tambahan

### P1-07 — Bug pada hasil koneksi WhatsApp non-immediate

**Lokasi:** `internal/modules/whatsapp/service.go`

Kondisi berikut terbalik:

```go
if client != nil || !client.IsConnected() {
    return ErrFailedToConnect
}
```

Client yang valid justru selalu dianggap gagal. Jika `client == nil`, ekspresi berikutnya dapat melakukan dereference nil.

### Dampak

Endpoint connect dengan `immediate=false` gagal atau dapat panic walaupun koneksi berhasil.

---

### P1-08 — Credential device ditempatkan di URL path — resolved by T0.7

**Lokasi:**

- `internal/app/routes/api.go`
- `internal/app/http/controllers/device/device.go`
- `internal/app/http/middleware/authentication.go`

Route lama `GET /api/v1/device/{token}/token` telah dihapus. Endpoint migrasi
sekarang menggunakan header dan tidak menerima token dari URL:

```http
GET /api/v1/device/by-token
Authorization: sk-dat-…
```

Middleware memvalidasi token dan menaruh device hasil resolusi di
`c.Locals("device")`; controller hanya mempercayai local tersebut. Dengan
demikian credential tidak masuk ke access log, reverse proxy log, tracing,
browser history, bookmark, atau `Referer`. Kegagalan autentikasi dicatat
secara generik tanpa credential mentah.

---

### P1-09 — Soft-delete/delete device tidak menghentikan client aktif

**Lokasi:** `internal/app/http/controllers/device/device.go`

`Destroy` hanya menghapus record melalui repository. Tidak ada disconnect WhatsApp/TDLib, kill channel, revoke token, cleanup file, atau penghapusan client dari memory.

### Dampak

Device yang sudah dihapus dapat tetap menerima event, mengirim webhook, menyimpan session, atau mempertahankan koneksi ke provider.

### Rekomendasi

Buat lifecycle service transactional: revoke credential, disconnect provider, stop event handler, cleanup state, lalu soft-delete resource. Operasi harus idempotent.

---

### P1-10 — Device startup menghubungkan semua device, bukan hanya device yang eligible

**Lokasi:** `internal/modules/device/repository.go`

`GetAllWhatsappDevices` dan `GetAllTelegramDevices` hanya memfilter `type`. Tidak ada filter status aktif, expiry, atau flag auto-connect.

### Dampak

Device inactive, expired, atau yang seharusnya tidak auto-connect dapat dibuat aktif kembali saat aplikasi restart.

### Rekomendasi

Tambahkan `auto_connect`, validasi status/expiry, dan state machine koneksi yang jelas.

---

### P1-11 — Lifecycle Telegram belum menangani error dan state secara lengkap

**Lokasi:** `internal/modules/telegram/tdlib.go` dan `internal/modules/telegram/client.go`

- Error dari `client.Connect` pada goroutine tidak dipropagasi.
- `IsConnected` hanya memeriksa `c.client != nil`, bukan authorization/connection state aktual.
- Disconnect pada client nil dapat panic.
- Status/JID tidak diperbarui secara konsisten.
- `ConnectDevices` mengabaikan error koneksi.

### Rekomendasi

Gunakan state machine, error channel, reconnect/backoff, context cancellation, dan update status atomik.

---

### P1-12 — Ada beberapa jalur panic dari data provider atau input malformed

Contoh tambahan di luar slicing media yang sudah dicatat:

- `user.Usernames.ActiveUsernames[0]` pada Telegram tanpa memastikan list tidak kosong.
- `i.VerifiedName.Serial` pada gateway tanpa memastikan `VerifiedName` tidak nil.
- `mime.ExtensionsByType(...)[0]` pada event WhatsApp tanpa memastikan extension tersedia.
- Dereference `document.FileName` tanpa memastikan pointer tidak nil.
- Jalur parsing JID/provider response masih bergantung pada input yang sudah tervalidasi dan perlu fuzz test.
- `SendKillChannel` dapat mengirim ke channel tanpa receiver dan memblokir goroutine.

### Rekomendasi

Gunakan defensive parsing, typed errors, fuzz test untuk payload/provider response, dan jalankan race test untuk lifecycle connection.

---

### P1-13 — Parameter pagination dapat menjadi unbounded atau invalid

Controller menerima `page` dan `limit` langsung dari query tanpa batas minimum/maksimum. Nilai negatif dapat menghasilkan offset/limit tidak valid; pada GORM, `Limit(-1)` dapat berarti menghapus limit dan menghasilkan query tanpa batas.

### Rekomendasi

Validasi `page >= 1`, `1 <= limit <= maxLimit`, batasi panjang search, dan gunakan satu helper pagination untuk semua endpoint.

---

### P1-14 — Session fixation dan session invalidation belum ditangani

**Lokasi:** `internal/app/http/controllers/auth/auth.go`

Login memakai session yang sudah ada dan tidak melakukan session ID regeneration setelah autentikasi. Logout juga mengabaikan error dari `session.Destroy()`. Password/role change tidak mencabut session yang sudah terbit.

### Rekomendasi

Regenerate session setelah login, revoke semua session pada logout/password change, simpan session minimal (subject ID/version), dan load permission terbaru pada request.

---

### P1-15 — Konfigurasi wajib dan invalid tidak menyebabkan fail-fast

**Lokasi:** `internal/config/config.go`

Semua pemanggilan `cfg.LoadStruct(...)` mengabaikan error. Padahal beberapa environment variable ditandai `required` dan parser dapat gagal untuk nilai duration/int/bool yang invalid.

### Dampak

Aplikasi dapat berjalan dengan konfigurasi kosong/default yang tidak aman, atau baru gagal jauh setelah startup. Contoh `.env.example` memakai `CACHE_TTL=60`, sedangkan parser duration mengharapkan format seperti `60s`.

### Rekomendasi

Propagasikan error konfigurasi dari bootstrap, validasi berdasarkan environment, dan tolak default development pada production.

---

### P1-16 — Migration runner tidak transactional dan tidak memakai lock

**Lokasi:** `platform/database/migrations/runner.go`

Setiap migration dijalankan terpisah dari pencatatan migration. Tidak ada transaction boundary, advisory lock, atau distributed lock.

### Risiko

- Dua replica dapat menjalankan migration bersamaan.
- Schema berhasil berubah tetapi record migration gagal dibuat.
- Deployment dapat berhenti pada schema setengah jadi.

### Rekomendasi

Gunakan advisory lock, transaction per migration bila didukung, checksum/versioning, dan backup/rollback procedure yang diuji.

---

### P1-17 — QR code dan provider session adalah credential sensitif

QR code disimpan dalam database sebagai data URL dan tersedia melalui endpoint/device response. TDLib/WhatsApp session juga berada di local volume/filesystem.

### Rekomendasi

- Terapkan TTL dan one-time read untuk QR.
- Jangan sertakan QR pada list/get umum.
- Batasi akses dengan scope khusus.
- Enkripsi data at rest dan backup.
- Hapus QR setelah pairing/expiry.

---

### P1-18 — Tidak ada batas resource untuk request dan media

**Lokasi:** `platform/server/server.go` dan endpoint gateway.

Tidak ada global `ReadLimit`, request body limit, upload limit, concurrency limit, atau timeout server. Media dikirim sebagai base64 sehingga ukuran payload dan memory dapat membesar sekitar 33% sebelum diproses.

### Rekomendasi

Tambahkan limit per route, streaming upload, MIME/content sniffing, quota per tenant/device, server read/write/idle timeout, dan circuit breaker provider.

---

## Temuan P2 Tambahan

### P2-06 — Security headers belum lengkap

Helmet dipasang pada web pipeline, tetapi tidak ada konfigurasi CSP yang aktif walaupun `platform/server/security/csp.go` tersedia. HSTS, Permissions-Policy, dan kebijakan frame/object perlu ditetapkan pada reverse proxy atau middleware.

### P2-07 — Error response dan logging masih berisiko membocorkan data

`DefaultErrorHandler` mengembalikan `err.Error()` ke client. Webhook URL juga ditulis ke log penuh. Pada development, GORM dapat mencatat SQL dan nilai sensitif.

Gunakan public error code, redaction untuk token/URL/header, dan jangan log payload message/media.

### P2-08 — OpenAPI/Swagger public dan kontrak security tidak lengkap

`/docs/api` dan file OpenAPI static didaftarkan tanpa authentication. Annotation security mendefinisikan `X-API-KEY`, tetapi implementasi membaca header `Authorization`, dan route tidak secara konsisten memiliki security requirement.

Lindungi dokumentasi production atau sediakan spec terpisah dengan security scheme yang benar.

### P2-09 — Custom context cancellation salah

**Lokasi:** `platform/ui/view/context.go`

`CombinedContext.startWatching` menutup `done` pada goroutine induk setelah child watcher dibuat, bukan setelah salah satu context selesai. Akibatnya `Done()` dapat tertutup segera dan cancellation/deadline tidak direpresentasikan dengan benar.

### P2-10 — Komponen platform yang diklaim production-ready belum terintegrasi

- `platform/cache` tidak dimuat oleh `internal/app/module.go`; default config `file` juga tidak didukung oleh `NewCache`.
- `platform/storage` tersedia tetapi storage abstraction tidak dipakai oleh media WhatsApp.
- Mail service dan queue hanya berada sebagai library/example, bukan alur aplikasi.
- Scheduler advanced memiliki dashboard placeholder dan scheduler jobs tetap in-memory secara default.

Scope production perlu diputuskan: integrasikan dengan test/operational contract, atau keluarkan dari build/runtime.

### P2-11 — Queue dan scheduler memiliki risiko shutdown/concurrency

Mail retry membuat goroutine yang dapat mengirim ke channel setelah queue ditutup. `MemoryJobQueue` memegang mutex ketika operasi channel dapat blocking. `EnqueueIn` memakai `time.Sleep` dan mengabaikan cancellation/error.

### P2-12 — Local storage tidak aman untuk data sensitif

**Lokasi:** `platform/storage/local.go`

Path dibangun dari input relatif tanpa canonicalization/traversal guard dan file dibuat dengan `os.ModePerm`. Jika abstraction ini diekspos ke input user, terdapat risiko path traversal dan permission terlalu terbuka.

### P2-13 — Frontend test/e2e pipeline belum siap

- Tidak ada unit test maupun E2E test.
- Playwright memakai `npm run dev`, sedangkan repository mendokumentasikan `pnpm` dan backend Go tidak otomatis tersedia pada base URL test.
- `pnpm test:e2e` gagal karena tidak ada test.
- `pnpm test:unit -- --run` gagal karena tidak ada test.

### P2-14 — Dependency audit menemukan advisory

`pnpm audit --prod` menemukan satu advisory **low** pada `esbuild` versi `0.27.7` melalui dependency `@nuxt/ui > @nuxt/fonts > fontless`. Versi patched yang dilaporkan adalah `>=0.28.1`.

Go vulnerability scan/SBOM belum tersedia di repository dan tool `govulncheck`, `trivy`, serta `gitleaks` tidak tersedia pada environment audit.

### P2-15 — Database pool dan graceful shutdown belum ditetapkan

Tidak ditemukan konfigurasi `MaxOpenConns`, `MaxIdleConns`, `ConnMaxLifetime`, provider disconnect, atau stop hook khusus untuk client WhatsApp/TDLib. Shutdown server ada, tetapi lifecycle koneksi messaging belum dikelola secara eksplisit.

### P2-16 — SMTP library menyediakan TLS verification bypass dan raw header construction

**Lokasi:** `platform/mail/smtp.go`

`SkipSSL` dapat mengaktifkan `tls.Config.InsecureSkipVerify`. Selain itu, subject, sender, recipient, dan custom headers dirangkai menjadi header SMTP tanpa validasi CRLF.

Jika mail service kelak menerima input user, ini dapat menyebabkan koneksi SMTP rentan MITM dan header injection. `SkipSSL` harus dilarang pada production, dan semua header/email address harus divalidasi.

---

## Temuan Fungsional Tambahan

| Area | Temuan |
|---|---|
| WhatsApp avatar | `GetAvatar` menggunakan JID device, bukan target `payload.Phone`; hasil avatar dapat salah. |
| Device update | Mengubah type/owner device tidak menghentikan koneksi provider lama atau merekonsiliasi token. |
| Device status | Status dapat diubah manual tanpa verifikasi terhadap koneksi provider aktual. |
| Device token GET | Error lookup diperiksa setelah policy dibuat/dijalankan, sehingga alur error tidak konsisten. |
| Session permissions | Session menyimpan snapshot account+role; perubahan role tidak langsung berlaku pada session lama. |
| API errors | Record-not-found sering dipropagasi mentah, sehingga mapping 404/401/403/422 tidak konsisten. |
| Webhook delivery | Tidak ada idempotency key, event ID yang terdokumentasi, retry queue, delivery status, atau signature verification. |
| Data retention | Tidak ada kebijakan retensi untuk QR, provider session, history sync, media, log, dan webhook payload. |

---

## Hal yang Sudah Baik / Tidak Menjadi Temuan Kritis

- Password menggunakan Argon2id dan verifikasi constant-time.
- Query search menggunakan parameter binding, bukan interpolasi langsung.
- Server memiliki request ID dan recover middleware.
- Model memakai soft delete untuk resource utama.
- Build asset dan type-check frontend berhasil.
- Oxlint correctness selesai tanpa warning/error.

Catatan: poin positif tersebut tidak menggantikan kebutuhan integration/security testing.

---

# Definition of Done Production

Project dapat dianggap siap production jika minimal seluruh poin berikut terpenuhi:

- [ ] Tidak ada `panic("unimplemented")` pada endpoint yang aktif.
- [ ] Seluruh secret token di-hash dan tidak dikembalikan ulang pada list/get.
- [ ] Status API key/device token dan expiry device token divalidasi.
- [ ] Revoke/soft-delete credential selalu menghentikan authentication.
- [ ] CSRF menggunakan token header/form yang benar-benar tidak dapat dipalsukan cross-site.
- [ ] Ownership/tenant isolation diuji melalui integration test.
- [ ] Semua request memiliki validation dan batas ukuran.
- [ ] Webhook memiliki proteksi SSRF, HTTPS, signing, retry, dan audit delivery.
- [ ] Login memiliki rate limit dan session hardening, termasuk session regeneration.
- [ ] Delete/edit UI benar-benar terhubung ke API.
- [ ] Semua route pada sidebar tersedia atau dihapus.
- [ ] Semua gateway method memiliki dispatch provider yang benar.
- [ ] Delete device menghentikan client aktif dan mencabut credential.
- [ ] Pagination sudah benar dan memiliki limit maksimum.
- [ ] `go test ./...` dan seluruh frontend test lulus.
- [ ] CI menjalankan build, test, lint, format, security scan, dan migration check.
- [ ] Tersedia liveness/readiness endpoint.
- [ ] Server timeout, body limit, media quota, dan rate limit tersedia.
- [ ] Migration menggunakan lock/transaction strategy yang diuji.
- [ ] Metrics, error tracking, dan alerting tersedia.
- [ ] Docker production menggunakan non-root user dan secret eksternal.
- [ ] Backup dan restore database sudah diuji.
- [ ] Load test untuk endpoint message/media sudah dilakukan.

---

# Roadmap yang Disarankan

## Fase 1 — Security dan correctness

1. Perbaiki authentication token dan ownership authorization.
2. Tambahkan validation, body limit, timeout, dan rate limiting.
3. Hilangkan panic dan implementasikan/disable fitur Telegram yang belum selesai.
4. Perbaiki webhook SSRF dan TLS.
5. Perbaiki pagination serta mass assignment.

## Fase 2 — Product completeness

1. Selesaikan UI device CRUD.
2. Tambahkan halaman/admin feature yang memang masuk scope release.
3. Hapus link placeholder yang belum tersedia.
4. Definisikan kontrak API dan error response.

## Fase 3 — Quality dan operation

1. Tambahkan unit, integration, dan E2E test.
2. Tambahkan CI/CD.
3. Tambahkan metrics, tracing, error tracking, dan alerting.
4. Siapkan Docker production, TLS, backup, restore, dan rollback.
5. Jalankan load test dan security review sebelum release.
