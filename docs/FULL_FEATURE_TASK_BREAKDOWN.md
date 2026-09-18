# Sapasora Full-Feature Task Breakdown

> **Status:** Master implementation backlog
>
> **Format:** Semua item adalah task/ticket yang dapat dicentang. Item berindentasi adalah subtask dari item di atasnya.
>
> **Source of truth:** [`PRODUCT_MASTER_PLAN.md`](./PRODUCT_MASTER_PLAN.md), [`FULL_FEATURE_SPECIFICATION.md`](./FULL_FEATURE_SPECIFICATION.md), [`PRODUCTION_READINESS.md`](./PRODUCTION_READINESS.md), dan [`finalize-plan/`](./finalize-plan/).
>
> **Catatan:** Seluruh checkbox sengaja belum dicentang. Fitur dianggap selesai hanya jika implementation, test, observability, authorization, documentation, dan recovery path-nya selesai.

---

## Cara Menggunakan Backlog

- Satu checkbox leaf sebaiknya menjadi satu ticket yang dapat dikerjakan dalam satu PR kecil atau satu rangkaian PR yang jelas.
- Parent checkbox hanya boleh dicentang setelah seluruh child checkbox selesai.
- Setiap ticket wajib memiliki acceptance criteria, test, migration impact, dan rollback/recovery plan bila relevan.
- Jangan mengerjakan UI seolah fitur selesai sebelum domain, API, authorization, worker, dan error state tersedia.
- Gunakan dependency order: **Track 0 → Foundation → Provider → Core Operations → Growth → Business Platform → Production Excellence**.
- Label yang disarankan: `security`, `backend`, `frontend`, `provider`, `data`, `worker`, `billing`, `docs`, `qa`, `ops`.

---

# Track 0 — Remediasi Blocker dari Audit Saat Ini

## T0.1 Telegram tidak boleh panic

- [x] Inventarisasi seluruh method Telegram yang masih `panic("unimplemented")`.
  - [x] `GetContacts`.
  - [x] `GetUser`.
  - [x] `Logout`.
  - [x] `SendAudio`.
  - [x] `SendButton`.
  - [x] `SendChatPresence`.
  - [x] `SendContact`.
  - [x] `SendDocument`.
  - [x] `SendImage`.
  - [x] `SendList`.
  - [x] `SendLocation`.
  - [x] `SendSticker`.
  - [x] `SendVideo`.
- [x] Definisikan behavior untuk capability Telegram yang tidak didukung.
- [ ] Implementasikan method Telegram yang masuk scope release.
- [x] Untuk method yang belum siap, return `501 Not Implemented` dengan error code stabil.
- [x] Tambahkan fake Telegram adapter untuk unit test.
- [ ] Tambahkan integration test dengan TDLib fixture.
- [ ] Pastikan error dari goroutine TDLib sampai ke supervisor.
- [x] Pastikan shutdown TDLib menutup client dan goroutine dengan aman.
- [x] Tambahkan test untuk `ActiveUsernames` kosong.
- [x] Tambahkan test untuk disconnect ketika client nil.
- [x] Tambahkan test untuk state `connecting`, `connected`, `disconnected`, dan `error`.
- [x] Hapus seluruh panic provider dari route yang aktif.

## T0.2 Perbaiki CSRF

- [x] Petakan seluruh route web yang menggunakan session/cookie.
- [x] Pisahkan tempat token CSRF disimpan dari tempat token dibaca oleh extractor.
- [x] Jangan gunakan cookie yang sama sebagai `KeyLookup` dan `CsrfFromCookie`.
- [x] Pilih pola synchronizer token atau double-submit cookie yang benar.
- [x] Tandai cookie CSRF sesuai model keamanan yang dipilih.
- [x] Pastikan token tidak otomatis tersedia sebagai header pada cross-site request.
- [x] Tambahkan middleware test untuk valid request.
- [x] Tambahkan test request tanpa token.
- [x] Tambahkan test token salah.
- [x] Tambahkan test token lama/expired.
- [ ] Tambahkan browser/E2E test terhadap form Inertia.
- [x] Dokumentasikan pengecualian hanya untuk webhook dengan signature provider.

## T0.3 Perbaiki device token authentication

- [x] Audit query authentication device token di `internal/modules/device/repository.go`.
- [x] Hentikan raw join yang mengabaikan status credential.
- [x] Tambahkan predicate `status = active`.
- [x] Tambahkan predicate `expired_at IS NULL OR expired_at > now()`.
- [x] Tambahkan predicate `deleted_at IS NULL` pada token.
- [x] Pastikan device/channel parent juga belum dihapus.
- [x] Pastikan workspace/owner scope masuk query.
- [x] Tambahkan test token aktif.
- [x] Tambahkan test token revoked.
- [x] Tambahkan test token expired.
- [x] Tambahkan test token soft-deleted.
- [x] Tambahkan test device soft-deleted.
- [x] Tambahkan test token dari workspace lain.
- [x] Samakan behavior dengan API key authentication.
- [x] Tambahkan audit event untuk authentication failure yang aman dari enumeration.

## T0.4 Perbaiki ownership credential dan device

- [x] Hilangkan kepercayaan pada `DeviceID` dan `UserID` yang dikirim client.
- [x] Turunkan owner dari authenticated workspace/actor.
- [x] Validasi device benar-benar dimiliki workspace actor.
- [x] Validasi API key owner dan scope.
- [x] Validasi device token hanya dapat dibuat untuk device yang authorized.
- [x] Tambahkan test create token lintas owner.
- [x] Tambahkan test update device lintas owner.
- [x] Tambahkan test delete device lintas owner.
- [x] Tambahkan test API key lintas owner.
- [x] Tambahkan repository helper untuk tenant-scoped lookup.
- [x] Larang controller melakukan lookup global berdasarkan ID saja.

## T0.5 Perbaiki provider dispatch gateway

- [x] Buat satu resolver adapter berdasarkan `DeviceType`/channel type.
- [x] Pastikan `GetContacts` dispatch Telegram ke Telegram adapter.
- [x] Pastikan `GetUser` dispatch Telegram ke Telegram adapter.
- [x] Pastikan `Logout` dispatch ke adapter yang benar.
- [x] Pastikan seluruh method send dispatch berdasarkan provider.
- [x] Return error capability jika provider tidak mendukung message type.
- [x] Hapus pemanggilan WhatsApp hard-coded dari gateway.
- [x] Tambahkan table-driven test untuk setiap device type dan method.
- [x] Tambahkan test unsupported device type.
- [x] Tambahkan test nil adapter/configuration.

## T0.6 Perbaiki bug dan panic WhatsApp

- [x] Perbaiki kondisi koneksi `client != nil || !client.IsConnected()`.
- [ ] Tambahkan test immediate connect.
- [ ] Tambahkan test non-immediate connect.
- [ ] Tambahkan test nil client.
- [ ] Tambahkan test already-connected client.
- [x] Validasi MIME extension sebelum indexing hasil split.
- [x] Validasi filename document sebelum dereference.
- [x] Validasi seluruh payload media sebelum decode.
- [x] Amankan `SendKillChannel` dari deadlock/unbuffered send.
- [x] Tambahkan timeout pada operasi provider.
- [ ] Tambahkan test reconnect race.
- [ ] Tambahkan test disconnect saat send berjalan.
- [x] Pastikan GetAvatar menggunakan target contact/JID, bukan device JID yang salah.

## T0.7 Credential tidak boleh bocor melalui URL

- [x] Hapus device token dari route URL `/api/v1/device/{token}/token`.
- [x] Ganti authentication dengan `Authorization: sk-dat-*` header.
- [x] Sediakan migration path melalui `GET /api/v1/device/by-token`; client lama
  harus bermigrasi karena route lama tidak lagi diterima.
- [x] Redact token dari access log: token tidak pernah diterima sebagai route
  parameter dan middleware tidak mencatat credential.
- [x] Redact token dari error log: authentication failure hanya mencatat alasan
  generik tanpa raw credential/error mentah.
- [x] Dokumentasikan penghapusan/deprecation endpoint lama di
  `docs/security/device-token-authentication.md`.
- [x] Tambahkan test bahwa endpoint hanya memakai `c.Locals("device")`, route
  lama tidak aktif, dan credential tidak muncul pada URL/history/referrer.

## T0.8 Delete/disconnect device harus lengkap

- [x] Definisikan state machine deletion device.
- [ ] Stop accepting new outbound jobs untuk device yang akan dihapus.
- [ ] Drain atau cancel job yang masih eligible.
- [x] Disconnect provider client aktif.
- [x] Revoke semua credential device.
- [ ] Hapus/rotate provider session sesuai retention policy.
- [ ] Hapus atau pindahkan media sesuai retention policy.
- [ ] Emit audit event.
- [ ] Emit device deleted event.
- [x] Pastikan restart tidak menghidupkan device yang sudah dihapus.
- [x] Tambahkan integration test delete dengan client aktif.

## T0.9 Startup dan lifecycle provider

- [x] Hubungkan hanya device yang `auto_connect = true`.
- [x] Filter status aktif.
- [x] Filter credential belum expired.
- [x] Skip soft-deleted device.
- [x] Tambahkan startup concurrency limit.
- [x] Tambahkan bounded, context-aware per-device reconnect backoff untuk percobaan koneksi provider yang gagal.

Catatan: backoff saat ini mencakup percobaan koneksi startup/manual yang gagal; deteksi dan pemulihan otomatis setelah koneksi yang sudah established terputus masih belum aman/tersedia di lifecycle provider saat ini.

- [x] Jangan block seluruh startup karena satu provider gagal.
- [x] Simpan last startup error.
- [x] Tambahkan graceful stop hook untuk semua provider.
- [x] Pastikan `IsConnected` benar-benar merepresentasikan state provider.
- [x] Selaraskan status dan JID/device identity.
- [x] Tambahkan provider lifecycle metrics.

## T0.10 Hardening config, migration, server, dan error

- [x] Propagate seluruh error `LoadStruct` configuration.
- [x] Fail fast untuk secret/provider config wajib.
- [x] Validasi duration seperti `CACHE_TTL` dengan unit yang eksplisit.
- [x] Tambahkan config schema test untuk valid dan invalid environment.
- [x] Tambahkan migration transaction policy.
- [x] Tambahkan PostgreSQL advisory lock/migration lock.
- [x] Dokumentasikan migration rollback.
- [x] Tambahkan server read timeout.
- [x] Tambahkan write timeout.
- [x] Tambahkan idle timeout.
- [x] Tambahkan request body limit.
- [x] Tambahkan media upload limit.
- [x] Tambahkan bounded concurrency provider.
- [x] Sanitasi error response production.
- [x] Jangan mengembalikan `err.Error()` internal secara langsung.
- [x] Redact webhook URL, SQL values, credential, PII, dan message body dari log.
- [x] Aktifkan CSP dan security headers lengkap.
- [x] Amankan Swagger/Scalar di production.
- [x] Selaraskan security annotation dengan `Authorization` yang sebenarnya.
- [x] Perbaiki `CombinedContext` agar `Done()` tidak tertutup terlalu awal.

---

# Track 1 — Product, Research, dan Decisions

## P1.1 Product alignment

- [ ] Tetapkan nama produk canonical: Sapasora.
- [ ] Tetapkan positioning final.
- [ ] Tetapkan primary wedge: official WhatsApp shared inbox + transactional notification.
- [ ] Tetapkan primary customer segment.
- [ ] Tetapkan secondary segment developer/agency.
- [ ] Tetapkan enterprise expansion scope.
- [ ] Dokumentasikan non-goals.
- [ ] Dokumentasikan risiko unofficial channel.
- [ ] Tetapkan North Star metric.
- [ ] Tetapkan activation metric.
- [ ] Tetapkan delivery quality metric.
- [ ] Tetapkan retention/revenue metric.

## P1.2 Discovery validation

- [ ] Susun interview screener untuk owner/ops.
- [ ] Susun interview screener untuk developer.
- [ ] Susun interview screener untuk support lead.
- [ ] Susun interview screener untuk agency.
- [ ] Jalankan minimal 15 interview target user.
- [ ] Validasi pain onboarding channel.
- [ ] Validasi pain delivery visibility.
- [ ] Validasi kebutuhan shared inbox.
- [ ] Validasi kebutuhan transactional notification.
- [ ] Validasi willingness-to-pay.
- [ ] Validasi official WhatsApp onboarding path.
- [ ] Validasi BSP vs direct Meta path.
- [ ] Validasi kebutuhan Telegram/SMS/Email.
- [ ] Dokumentasikan hasil riset dan keputusan.

## P1.3 ADR dan open decisions

- [ ] Buat ADR official Meta Tech Provider vs BSP.
- [ ] Buat ADR whatsmeow closed beta policy.
- [ ] Buat ADR GORM ke sqlx atau tetap GORM.
- [ ] Buat ADR migration runner canonical.
- [ ] Buat ADR queue semantics Valkey vs managed queue.
- [ ] Buat ADR auth/session provider.
- [ ] Buat ADR payment/tax/invoice provider.
- [ ] Buat ADR retention per plan/industry.
- [ ] Buat ADR public ID migration.
- [ ] Buat ADR AI provider dan data boundary.
- [ ] Buat ADR multi-region requirement.
- [ ] Assign owner dan deadline untuk setiap open decision.

## P1.4 Brand dan content

- [ ] Finalisasi naming/domain.
- [ ] Finalisasi tagline.
- [ ] Finalisasi brand voice.
- [ ] Finalisasi visual direction.
- [ ] Buat logo/favicon assets.
- [ ] Buat product terminology glossary.
- [ ] Buat legal disclaimer official/unofficial channel.
- [ ] Buat copy onboarding channel.
- [ ] Buat copy failure/recovery.
- [ ] Buat copy pricing/provider fee.

---

# Track 2 — Engineering Foundation

## P2.1 Repository conventions

- [ ] Dokumentasikan feature-slice directory convention.
- [ ] Dokumentasikan naming convention backend.
- [ ] Dokumentasikan naming convention frontend.
- [ ] Dokumentasikan error code convention.
- [ ] Dokumentasikan event naming convention.
- [ ] Dokumentasikan migration naming convention.
- [ ] Dokumentasikan DTO/model separation.
- [ ] Tambahkan lint/format command yang reproducible.
- [ ] Tambahkan make/task command untuk local development.
- [ ] Tambahkan pre-commit/changed-files checks bila disepakati.

## P2.2 Configuration

- [ ] Definisikan typed application config.
- [ ] Validasi required config saat startup.
- [ ] Pisahkan local/development/staging/production config.
- [ ] Hilangkan credential default production.
- [ ] Definisikan secret provider interface.
- [ ] Integrasikan secret manager untuk production.
- [ ] Tambahkan config redaction pada diagnostics.
- [ ] Tambahkan config documentation.
- [ ] Tambahkan config compatibility test.
- [ ] Tambahkan safe config dump tanpa secret.

## P2.3 Database foundation

- [ ] Inventarisasi seluruh tabel existing.
- [ ] Inventarisasi seluruh foreign key dan index.
- [ ] Inventarisasi soft-delete semantics.
- [ ] Definisikan `workspace_id` untuk setiap tenant-owned table.
- [ ] Definisikan timestamp convention UTC.
- [ ] Definisikan transaction boundary convention.
- [ ] Definisikan connection pool min/max.
- [ ] Definisikan connection lifetime/idle timeout.
- [ ] Tambahkan DB ping readiness.
- [ ] Tambahkan slow query logging ter-redact.
- [ ] Tambahkan migration lock.
- [ ] Tambahkan migration smoke test pada database kosong.
- [ ] Tambahkan migration smoke test pada realistic database.
- [ ] Tambahkan rollback/recovery procedure.

## P2.4 Public ID contract

- [ ] Buat public ID generator NanoID.
- [ ] Tambahkan unique constraint untuk setiap `public_id`.
- [ ] Buat mapper internal model ke external DTO.
- [ ] Larang serialization model database langsung.
- [ ] Migrasikan resource utama ke public ID.
- [ ] Migrasikan message ke public ID.
- [ ] Migrasikan channel/device ke public ID.
- [ ] Migrasikan contact/conversation ke public ID.
- [ ] Migrasikan campaign/template ke public ID.
- [ ] Migrasikan webhook/credential reference ke public ID.
- [ ] Tambahkan collision retry test.
- [ ] Tambahkan test internal numeric ID tidak muncul di response.
- [ ] Tambahkan test `public_id` tidak muncul sebagai field terpisah.

## P2.5 Error and request foundation

- [x] Definisikan public error code catalog.
- [x] Definisikan error category retryable/non-retryable.
- [x] Definisikan error envelope API.
- [x] Sertakan request ID pada semua response.
- [ ] Propagate request/correlation ID ke worker.
- [ ] Map validation error ke `400/422`.
- [ ] Map authentication error ke `401`.
- [ ] Map authorization error ke `403`.
- [ ] Map not found tanpa resource enumeration.
- [ ] Map conflict/idempotency error ke `409`.
- [ ] Map provider unavailable ke `503`.
- [ ] Buat production exception sanitizer.
- [ ] Buat developer diagnostics yang hanya bisa diakses authorized admin.

## P2.6 Validation foundation

- [x] Definisikan validation library/pattern backend.
- [x] Validasi semua request body.

> Scope note (P2.6 slice): the existing `platform/validator` is now used for request-body parsing and validation in the device, account, API-key, device-token, role, auth, and gateway controllers. List endpoints now share bounded pagination/filter parsing (`page >= 1`, `1 <= limit <= 100`, malformed numeric values rejected, and bounded search); resource routes validate URL-safe public IDs before service lookup. Query/path validation is evidenced; media/frontend/provider validation remains open.
- [x] Validasi query pagination/filter.
- [x] Validasi path parameter.
- [ ] Validasi URL/hostname.
- [ ] Validasi phone/email/identifier.
- [ ] Validasi enum/capability.
- [ ] Validasi media MIME/size.
- [ ] Validasi array length dan nested payload.
- [ ] Validasi template variable.
- [ ] Tambahkan frontend Zod schemas.
- [ ] Parse Inertia props dengan Zod.
- [ ] Parse API response dengan Zod.
- [ ] Parse local/session storage dengan Zod.
- [ ] Parse external/provider payload dengan Zod.
- [ ] Tambahkan schema drift tests.

---

# Track 3 — Identity, Workspace, dan Authorization

## P3.1 Authentication

- [ ] Implementasikan register.
- [ ] Implementasikan login.
- [ ] Regenerate session ID setelah login.
- [ ] Implementasikan logout.
- [ ] Handle logout destroy error.
- [ ] Implementasikan email verification.
- [ ] Implementasikan resend verification.
- [ ] Implementasikan forgot password.
- [ ] Implementasikan reset password one-time.
- [ ] Hash reset token.
- [ ] Set reset token expiry.
- [ ] Implementasikan password change.
- [ ] Implementasikan password history/policy.
- [ ] Implementasikan MFA enrollment.
- [ ] Implementasikan TOTP verification.
- [ ] Implementasikan recovery codes.
- [ ] Implementasikan MFA recovery/revoke.
- [ ] Implementasikan suspicious login notification.
- [ ] Tambahkan login brute-force rate limit.
- [ ] Tambahkan account lockout/backoff policy.
- [ ] Tambahkan session list.
- [ ] Tambahkan remote session revoke.
- [ ] Revoke session saat perubahan role/password sesuai policy.
- [ ] Tambahkan auth audit log.
- [ ] Tambahkan unit/integration/E2E auth tests.

## P3.2 Workspace

- [ ] Buat workspace.
- [ ] Edit workspace name.
- [ ] Edit workspace logo.
- [ ] Edit timezone.
- [ ] Edit locale/language.
- [ ] Edit notification settings.
- [ ] Implementasikan workspace switcher.
- [ ] Pastikan semua request membawa active workspace context.
- [ ] Validasi user membership saat switch.
- [ ] Implementasikan ownership transfer.
- [ ] Implementasikan workspace suspend.
- [ ] Implementasikan workspace delete request.
- [ ] Implementasikan grace period delete.
- [ ] Implementasikan asynchronous purge.
- [ ] Audit seluruh workspace mutation.

## P3.3 Member dan invitation

- [ ] Buat invitation token yang hashed.
- [ ] Set invitation expiry.
- [ ] Kirim invitation email.
- [ ] Accept invitation.
- [ ] Reject invitation.
- [ ] Resend invitation.
- [ ] Revoke invitation.
- [ ] List workspace members.
- [ ] Suspend member.
- [ ] Reactivate member.
- [ ] Remove member.
- [ ] Transfer ownership.
- [ ] Cegah invite lintas workspace tanpa permission.
- [ ] Audit semua member mutation.

## P3.4 RBAC dan policy

- [ ] Definisikan permission catalog.
- [ ] Seed built-in roles.
- [ ] Implementasikan custom role.
- [ ] Implementasikan role create/update/delete.
- [ ] Validasi role tidak dapat menghapus owner terakhir.
- [ ] Implementasikan resource-level policy.
- [ ] Implementasikan API key scope.
- [ ] Implementasikan developer/service account scope.
- [ ] Implementasikan temporary elevated permission dengan expiry.
- [ ] Tambahkan policy tests untuk setiap resource.
- [ ] Tambahkan cross-tenant authorization tests.
- [ ] Tambahkan API/UI permission guard.
- [ ] Pastikan UI guard tidak menggantikan server authorization.

---

# Track 4 — Channel and Provider Platform

## P4.1 Channel domain

- [x] Definisikan `ChannelAccount` model.
- [x] Definisikan channel type enum.
- [x] Definisikan provider enum.
- [x] Definisikan capability model.
- [x] Definisikan channel status state machine.
- [x] Definisikan connection event model.
- [x] Definisikan provider credential reference.
- [x] Definisikan provider session reference.
- [x] Definisikan channel quota/rate policy.
- [x] Tambahkan channel tenant scope.
- [x] Tambahkan channel public ID.
- [x] Tambahkan channel repository tests.

## P4.2 Adapter contract

- [x] Definisikan `Connect` contract.
- [x] Definisikan `Disconnect` contract.
- [x] Definisikan `Health` contract.
- [x] Definisikan `Capabilities` contract.
- [x] Definisikan `Send` contract.
- [x] Definisikan `GetContact/User` contract.
- [x] Definisikan `ReceiveEvent` contract.
- [x] Definisikan `ValidateAddress` contract.
- [x] Definisikan `NormalizeError` contract.
- [x] Definisikan `NormalizeInboundEvent` contract.
- [x] Pisahkan interface berdasarkan capability.
- [x] Buat fake adapter untuk test.
- [x] Buat contract test yang wajib lulus semua adapter.
- [x] Buat provider error mapping catalog.

## P4.3 Channel lifecycle

- [x] Implementasikan `pending` state.
- [x] Implementasikan `connecting` state.
- [x] Implementasikan `connected` state.
- [x] Implementasikan `degraded` state.
- [x] Implementasikan `disconnected` state.
- [x] Implementasikan `expired` state.
- [x] Implementasikan `error` state.
- [x] Validasi semua state transition.
- [x] Buat reconnect backoff.
- [x] Buat provider circuit breaker.
- [ ] Buat per-channel distributed lock.
- [ ] Buat connection history.
- [ ] Buat manual reconnect action.
- [ ] Buat graceful disconnect.
- [ ] Buat delete cleanup workflow.
- [ ] Buat channel health metrics.

## P4.4 Credential/session handling

- [x] Encrypt credential at rest.
- [ ] Encrypt provider session at rest.
- [x] Simpan secret reference, bukan secret di DTO.
- [x] Implementasikan credential rotation.
- [x] Implementasikan credential revoke.
- [x] Implementasikan credential expiry.
- [x] Redact secret pada log.
- [x] Redact secret pada diagnostics.
- [x] One-time reveal saat creation/rotation.
- [x] Tambahkan access audit untuk credential.
- [x] Tambahkan purge setelah disconnect/delete.

## P4.5 WhatsApp official

- [ ] Tentukan Meta/BSP onboarding path.
- [ ] Implementasikan provider credential setup.
- [ ] Implementasikan WABA identity mapping.
- [ ] Implementasikan phone number mapping.
- [ ] Implementasikan official webhook verification.
- [ ] Implementasikan webhook event ingestion.
- [ ] Implementasikan template synchronization.
- [ ] Implementasikan template approval status.
- [ ] Implementasikan media upload/download mapping.
- [ ] Implementasikan messaging window validation.
- [ ] Implementasikan provider error mapping.
- [ ] Implementasikan delivery/read status mapping.
- [ ] Implementasikan official health check.
- [ ] Implementasikan rate/quota handling.
- [ ] Tambahkan Meta/provider contract tests.
- [ ] Tambahkan onboarding E2E dengan test asset.
- [ ] Dokumentasikan policy/cost/limitation.

## P4.6 WhatsApp whatsmeow experimental

- [ ] Pisahkan whatsmeow dari official adapter.
- [ ] Jalankan dalam worker/process boundary yang jelas.
- [ ] Tambahkan feature flag.
- [ ] Tambahkan closed-beta allowlist.
- [ ] Tambahkan risk acknowledgement.
- [ ] Tambahkan label experimental pada UI/API/docs.
- [ ] Implementasikan QR/pairing.
- [ ] Implementasikan encrypted session persistence.
- [ ] Implementasikan reconnect/backoff.
- [ ] Implementasikan event normalization.
- [ ] Implementasikan session stop hook.
- [ ] Implementasikan kill switch workspace/channel.
- [ ] Implementasikan protocol/disconnect metrics.
- [ ] Implementasikan ban/logout risk monitoring.
- [ ] Implementasikan safe delete session.
- [ ] Legal review terms/licensing.
- [ ] Provider failure drill.

## P4.7 Telegram

- [ ] Pilih bot vs TDLib support scope.
- [ ] Implementasikan authentication/onboarding.
- [ ] Implementasikan connection lifecycle.
- [ ] Implementasikan `GetContacts`.
- [ ] Implementasikan `GetUser`.
- [ ] Implementasikan `Logout`.
- [ ] Implementasikan text send.
- [ ] Implementasikan audio send.
- [ ] Implementasikan button send.
- [ ] Implementasikan chat presence.
- [ ] Implementasikan contact send.
- [ ] Implementasikan document send.
- [ ] Implementasikan image send.
- [ ] Implementasikan list send.
- [ ] Implementasikan location send.
- [ ] Implementasikan sticker send.
- [ ] Implementasikan video send.
- [ ] Implementasikan inbound update normalization.
- [ ] Implementasikan flood-wait handling.
- [ ] Implementasikan rate limit.
- [ ] Implementasikan TDLib error supervisor.
- [ ] Tambahkan all-capability contract tests.

## P4.8 SMS

- [ ] Definisikan SMS adapter interface.
- [ ] Implementasikan provider pertama.
- [ ] Implementasikan sender identity.
- [ ] Implementasikan country routing.
- [ ] Implementasikan character encoding detection.
- [ ] Implementasikan segment count.
- [ ] Implementasikan cost estimate.
- [ ] Implementasikan delivery receipt.
- [ ] Implementasikan STOP/START handling.
- [ ] Implementasikan suppression list.
- [ ] Implementasikan quiet hours.
- [ ] Implementasikan country compliance rules.
- [ ] Tambahkan provider fake integration test.

## P4.9 Email

- [ ] Definisikan Email adapter interface.
- [ ] Implementasikan SMTP provider.
- [ ] Implementasikan API email provider.
- [ ] Implementasikan sender/from identity.
- [ ] Implementasikan reply-to.
- [ ] Implementasikan HTML/text multipart.
- [ ] Implementasikan attachment.
- [ ] Implementasikan bounce event.
- [ ] Implementasikan complaint event.
- [ ] Implementasikan delivery event.
- [ ] Implementasikan open/click event bila tersedia.
- [ ] Implementasikan suppression list.
- [ ] Implementasikan unsubscribe group.
- [ ] Validasi TLS SMTP.
- [ ] Hindari raw header injection.
- [ ] Tambahkan email test send.

---

# Track 5 — Contacts, Consent, dan Audience

## P5.1 Contact model

- [x] Buat contact model (table/migration deferred).
- [x] Buat contact address/identity model (table/migration deferred).
- [x] Normalisasi phone number.
- [x] Normalisasi email.
- [x] Simpan channel identifier provider.
- [x] Tambahkan timezone/locale.
- [x] Tambahkan owner/assignment.
- [x] Tambahkan lifecycle stage.
- [ ] Tambahkan custom fields schema.
- [x] Tambahkan tags.
- [x] Tambahkan notes.
- [x] Tambahkan source metadata.
- [x] Tambahkan workspace scope.
- [x] Tambahkan duplicate detection key.

## P5.2 Contact CRUD

- [x] Create contact.
- [x] Read contact.
- [x] Update contact.
- [x] Soft delete contact.
- [ ] Restore contact.
- [ ] Permanent purge dengan permission.
- [ ] Search contact.
- [x] Filter contact.
- [ ] Sort contact.
- [ ] Cursor pagination contact.
- [ ] Bulk update contact.
- [ ] Bulk tag contact.
- [ ] Bulk delete contact.
- [ ] Audit contact mutation.

## P5.3 Import/export

- [ ] Upload CSV.
- [ ] Upload XLSX bila diperlukan.
- [ ] Preview import.
- [ ] Mapping columns.
- [x] Validate rows.
- [x] Detect duplicate.
- [x] Show row-level errors.
- [ ] Async import job.
- [ ] Import progress.
- [ ] Download error report.
- [ ] Retry failed rows.
- [ ] Export permission check.
- [ ] Async export job.
- [ ] Export audit event.
- [ ] Expiring download URL.

## P5.4 Consent and suppression

- [x] Definisikan consent state per channel.
- [x] Simpan consent source.
- [x] Simpan consent timestamp.
- [x] Simpan evidence/reference.
- [x] Simpan actor.
- [x] Buat immutable consent event.
- [x] Implementasikan opt-in.
- [x] Implementasikan opt-out.
- [ ] Implementasikan global suppression.
- [x] Implementasikan per-channel suppression.
- [ ] Implementasikan STOP handler SMS.
- [ ] Implementasikan unsubscribe email.
- [x] Block send ke opted-out recipient.
- [x] Tambahkan consent audit.
- [ ] Tambahkan consent export/delete handling.

## P5.5 Segment dan audience

- [x] Buat static list.
- [x] Buat dynamic segment.
- [x] Filter berdasarkan field.
- [x] Filter berdasarkan tag.
- [ ] Filter berdasarkan event.
- [ ] Filter berdasarkan delivery status.
- [ ] Filter berdasarkan interaction time.
- [x] Preview segment count.
- [x] Simpan audience snapshot.
- [x] Simpan exclusion list.
- [x] Revalidate suppression sebelum send.
- [x] Test segment query tenant isolation.

## P5.6 Duplicate and identity

- [x] Definisikan duplicate matching rules.
- [x] Tampilkan duplicate candidates.
- [x] Implementasikan merge preview.
- [x] Implementasikan merge contact.
- [x] Simpan merge audit.
- [x] Implementasikan undo bila policy memungkinkan.
- [ ] Link identity lintas channel dengan confirmation.
- [x] Cegah accidental cross-contact merge.

---

# Track 6 — Conversation dan Unified Inbox

## P6.1 Conversation model

- [x] Buat conversation table/model.
- [x] Buat participant model.
- [x] Buat assignment model.
- [x] Tambahkan channel association.
- [x] Tambahkan contact association.
- [x] Tambahkan status.
- [x] Tambahkan priority.
- [x] Tambahkan tags.
- [x] Tambahkan SLA fields.
- [x] Tambahkan last activity.
- [x] Tambahkan unread count.
- [x] Tambahkan workspace scope.
- [x] Tambahkan conversation public ID.

## P6.2 Inbox list

- [x] Tampilkan conversation list.
- [x] Search conversation.
- [x] Filter channel.
- [x] Filter status.
- [x] Filter assignee.
- [x] Filter team.
- [x] Filter priority.
- [x] Filter tag.
- [x] Filter unread.
- [x] Filter SLA breach.
- [x] Sort last activity.
- [x] Cursor pagination.
- [ ] Preserve filter in URL query.
- [ ] Parse URL query dengan Zod.
- [x] Empty state.
- [x] Loading state.
- [x] Error state.

## P6.3 Conversation detail

- [ ] Tampilkan message timeline.
- [ ] Tampilkan delivery events.
- [ ] Tampilkan internal notes.
- [ ] Tampilkan automation events.
- [ ] Tampilkan assignment history.
- [ ] Tampilkan contact profile.
- [ ] Tampilkan consent status.
- [ ] Tampilkan channel health warning.
- [ ] Tampilkan SLA timer.
- [ ] Tampilkan unread/read state.
- [ ] Search dalam conversation.
- [ ] Load older messages.
- [ ] Handle message pagination.

## P6.4 Assignment and teamwork

- [ ] Assign ke agent.
- [ ] Assign ke team.
- [ ] Claim conversation.
- [ ] Unclaim conversation.
- [ ] Round-robin assignment.
- [ ] Manual reassignment.
- [ ] Supervisor takeover.
- [ ] Collision detection.
- [ ] Presence/typing indicator bila tersedia.
- [ ] Assignment audit.
- [ ] Assignment notification.

## P6.5 Conversation actions

- [ ] Set open.
- [ ] Set pending.
- [ ] Snooze.
- [ ] Resolve.
- [ ] Reopen.
- [ ] Close.
- [ ] Set priority.
- [ ] Add tag.
- [ ] Remove tag.
- [ ] Add internal note.
- [ ] Merge conversation.
- [ ] Split conversation.
- [ ] Export conversation.
- [ ] Permission-check setiap action.

## P6.6 SLA

- [ ] Definisikan business hours.
- [ ] Definisikan response SLA.
- [ ] Definisikan resolution SLA.
- [ ] Hitung timezone workspace.
- [ ] Tampilkan countdown.
- [ ] Tampilkan breach.
- [ ] Notify supervisor.
- [ ] Simpan SLA events.
- [ ] Tampilkan SLA analytics.

---

# Track 7 — Message Domain, Composer, dan Delivery

## P7.1 Message model

- [x] Buat message model.
- [x] Buat message attachment model.
- [x] Buat message delivery model.
- [x] Buat message event append-only.
- [x] Simpan inbound/outbound direction.
- [x] Simpan provider message ID.
- [x] Simpan idempotency key.
- [x] Simpan request ID.
- [x] Simpan normalized content.
- [x] Simpan redacted provider metadata.
- [x] Tambahkan workspace scope.
- [x] Tambahkan contact/conversation/channel association.

## P7.2 Message state machine

- [x] Implementasikan `accepted`.
- [x] Implementasikan `queued`.
- [x] Implementasikan `sending`.
- [x] Implementasikan `sent`.
- [x] Implementasikan `delivered`.
- [x] Implementasikan `read`.
- [x] Implementasikan `failed`.
- [x] Implementasikan `retrying`.
- [x] Implementasikan `dead_letter`.
- [x] Implementasikan `canceled`.
- [x] Validasi transition.
- [x] Simpan event setiap transition.
- [x] Deduplicate provider status event.
- [x] Bedakan platform accepted dan provider accepted.

## P7.3 Composer

- [ ] Compose text.
- [ ] Compose template.
- [ ] Compose image.
- [ ] Compose audio.
- [ ] Compose video.
- [ ] Compose document.
- [ ] Compose sticker.
- [ ] Compose location.
- [ ] Compose contact card.
- [ ] Compose poll bila capability tersedia.
- [ ] Compose button/list bila capability tersedia.
- [ ] Compose quick reply.
- [ ] Compose quoted reply.
- [ ] Preview variable.
- [ ] Validate capability sebelum send.
- [ ] Validate consent sebelum send.
- [ ] Validate quota sebelum send.
- [ ] Show cost estimate.
- [ ] Show quiet-hours warning.
- [ ] Confirmation untuk risky send.

## P7.4 Queue and outbox

- [ ] Buat outbox event table.
- [x] Buat message send job.
- [x] Publish outbox safely.
- [x] Implementasikan worker claim/lease.
- [x] Implementasikan bounded concurrency.
- [ ] Implementasikan queue per provider/channel.
- [ ] Implementasikan ordering policy.
- [x] Implementasikan retryable classification.
- [x] Implementasikan exponential backoff.
- [x] Implementasikan jitter.
- [x] Implementasikan dead-letter.
- [x] Implementasikan manual replay.
- [x] Implementasikan cancellation sebelum provider accept.
- [x] Implementasikan graceful shutdown.
- [ ] Tambahkan queue metrics.
- [ ] Tambahkan queue age alert.

## P7.5 Idempotency

- [x] Parse `Idempotency-Key`.
- [x] Scope key by workspace/actor/endpoint.
- [x] Store request fingerprint.
- [x] Return original result untuk duplicate request.
- [x] Reject same key dengan payload berbeda.
- [x] Set TTL sesuai operation.
- [x] Test concurrent duplicate requests.
- [ ] Test worker retry duplicate.
- [ ] Test campaign launch duplicate.
- [ ] Test webhook replay duplicate.

## P7.6 Delivery and failure UI

- [ ] Tampilkan accepted.
- [ ] Tampilkan queued.
- [ ] Tampilkan sending.
- [ ] Tampilkan sent.
- [ ] Tampilkan delivered.
- [ ] Tampilkan read.
- [ ] Tampilkan failed.
- [ ] Tampilkan retrying.
- [ ] Tampilkan canceled.
- [ ] Tampilkan actionable failure reason.
- [ ] Tampilkan `next_action`.
- [ ] Sediakan retry jika aman.
- [ ] Sediakan copy request ID.
- [ ] Hindari menampilkan provider secret/internal error.

---

# Track 8 — Media dan File

## P8.1 Upload

- [ ] Definisikan upload contract.
- [ ] Implementasikan direct upload.
- [ ] Implementasikan presigned upload.
- [ ] Implementasikan resumable upload untuk file besar.
- [ ] Validasi content length.
- [ ] Validasi MIME dengan sniffing.
- [ ] Validasi extension.
- [ ] Validasi filename.
- [ ] Normalize filename.
- [ ] Cegah path traversal.
- [ ] Set storage ownership workspace.
- [ ] Tambahkan upload progress.

## P8.2 Media security

- [ ] Scan malware/virus.
- [ ] Quarantine file gagal scan.
- [ ] Block executable content.
- [ ] Limit image dimensions.
- [ ] Limit video duration/size.
- [ ] Strip sensitive metadata bila diperlukan.
- [ ] Encrypt at rest.
- [ ] Redact media URL dari log.
- [ ] Audit media access.

## P8.3 Storage lifecycle

- [ ] Migrasikan WhatsApp media dari local path ke storage abstraction.
- [ ] Konfigurasi S3-compatible storage.
- [ ] Buat workspace prefix.
- [ ] Buat signed download URL.
- [ ] Expire signed URL.
- [ ] Cek authorization saat download.
- [ ] Implementasikan quota.
- [ ] Implementasikan retention.
- [ ] Implementasikan orphan cleanup.
- [ ] Implementasikan object versioning/lifecycle.
- [ ] Test container restart tidak menghilangkan media.

---

# Track 9 — Templates dan Content

## P9.1 Template domain

- [ ] Buat template model.
- [ ] Buat template version model.
- [ ] Tambahkan channel/provider scope.
- [ ] Tambahkan language.
- [ ] Tambahkan category.
- [ ] Tambahkan approval status.
- [ ] Tambahkan rejection reason.
- [ ] Tambahkan variables schema.
- [ ] Tambahkan usage counters.
- [ ] Tambahkan workspace scope.

## P9.2 Template lifecycle

- [ ] Create draft.
- [ ] Edit draft.
- [ ] Duplicate template.
- [ ] Submit review.
- [ ] Sync provider approval.
- [ ] Show approved.
- [ ] Show rejected.
- [ ] Show rejection reason.
- [ ] Publish version.
- [ ] Rollback version.
- [ ] Archive template.
- [ ] Deprecate template.
- [ ] Prevent use archived/rejected template.

## P9.3 Variables and preview

- [ ] Define variable name/type.
- [ ] Define required/optional.
- [ ] Define default.
- [ ] Define example.
- [ ] Validate missing variable.
- [ ] Escape per channel.
- [ ] Redact PII preview/log.
- [ ] Preview sample message.
- [ ] Preview SMS segment count.
- [ ] Preview email HTML/text.
- [ ] Preview WhatsApp provider component.
- [ ] Test send ke internal recipient.

## P9.4 Quick replies and content helpers

- [ ] Buat quick reply CRUD.
- [ ] Scope quick reply per workspace/team.
- [ ] Permission quick reply.
- [ ] Search quick reply.
- [ ] Insert quick reply ke composer.
- [ ] Buat canned response usage tracking.
- [ ] Tambahkan link preview policy.
- [ ] Tambahkan sender name policy.

---

# Track 10 — Scheduling, Recurring, dan Follow-up

## P10.1 Schedule domain

- [ ] Buat schedule model.
- [ ] Buat run history model.
- [ ] Simpan timezone.
- [ ] Simpan start/end date.
- [ ] Simpan audience mode snapshot/live.
- [ ] Simpan missed-run policy.
- [ ] Simpan blackout calendar.
- [ ] Tambahkan schedule public ID.

## P10.2 One-time schedule

- [ ] Schedule message.
- [ ] Validate future timestamp.
- [ ] Show timezone conversion.
- [ ] Edit pending schedule.
- [ ] Cancel pending schedule.
- [ ] Show run history.
- [ ] Handle missed execution.
- [ ] Prevent duplicate execution.
- [ ] Audit schedule mutation.

## P10.3 Recurring schedule

- [ ] Daily recurrence.
- [ ] Weekly recurrence.
- [ ] Monthly recurrence.
- [ ] Annual recurrence.
- [ ] Custom interval.
- [ ] Business-day option.
- [ ] End condition.
- [ ] Pause/resume.
- [ ] Cancel.
- [ ] Preview next runs.
- [ ] Handle DST/timezone.
- [ ] Handle provider quiet hours.
- [ ] Per-run delivery report.

## P10.4 Scheduler reliability

- [ ] Durable scheduler job store.
- [ ] Distributed scheduler lock.
- [ ] Lease/claim job.
- [ ] Clock drift handling.
- [ ] Retry policy.
- [ ] Dead-letter schedule run.
- [ ] Replay schedule run.
- [ ] Graceful shutdown.
- [ ] Scheduler metrics.
- [ ] Scheduler alerting.

---

# Track 11 — Campaign dan Broadcast

## P11.1 Campaign domain

- [ ] Buat campaign model.
- [ ] Buat campaign recipient model.
- [ ] Buat audience snapshot model.
- [ ] Tambahkan campaign status state machine.
- [ ] Tambahkan approval record.
- [ ] Tambahkan schedule.
- [ ] Tambahkan channel/template reference.
- [ ] Tambahkan cost estimate.
- [ ] Tambahkan workspace scope.

## P11.2 Campaign creation

- [ ] Create draft.
- [ ] Select channel.
- [ ] Select template.
- [ ] Select audience.
- [ ] Preview recipient count.
- [ ] Apply suppression.
- [ ] Apply consent filter.
- [ ] Apply exclusion list.
- [ ] Preview personalization.
- [ ] Validate variables.
- [ ] Estimate cost.
- [ ] Set quiet hours.
- [ ] Set throttle.
- [ ] Set frequency cap.
- [ ] Configure schedule.
- [ ] Send test.

## P11.3 Campaign approval and execution

- [ ] Implement approval threshold.
- [ ] Implement approval/rejection.
- [ ] Prevent self-approval bila policy melarang.
- [ ] Snapshot recipients saat launch.
- [ ] Recheck suppression saat execution.
- [ ] Prevent duplicate launch.
- [ ] Start campaign job.
- [ ] Pause campaign.
- [ ] Resume campaign.
- [ ] Cancel campaign.
- [ ] Handle partial failure.
- [ ] Handle provider outage.
- [ ] Store recipient-level status.
- [ ] Emit campaign events.
- [ ] Generate delivery report.
- [ ] Export campaign report.

## P11.4 Campaign safety

- [ ] Recipient maximum per workspace.
- [ ] Daily quota guard.
- [ ] Provider rate guard.
- [ ] Cost hard cap.
- [ ] Require confirmation above threshold.
- [ ] Block no-consent audience.
- [ ] Honor opted-out contact immediately.
- [ ] Add emergency pause.
- [ ] Add workspace kill switch.
- [ ] Add abuse monitoring.
- [ ] Audit campaign lifecycle.

## P11.5 Campaign analytics

- [ ] Sent count.
- [ ] Delivered count.
- [ ] Read count.
- [ ] Failed count.
- [ ] Opt-out count.
- [ ] Cost.
- [ ] Delivery rate.
- [ ] Failure reasons.
- [ ] Template performance.
- [ ] Variant performance.
- [ ] Agent/support follow-up.

---

# Track 12 — Automation, Autoresponder, dan AI

## P12.1 Automation model

- [ ] Buat automation model.
- [ ] Buat automation version model.
- [ ] Buat automation run model.
- [ ] Tambahkan draft/published/disabled state.
- [ ] Tambahkan workspace scope.
- [ ] Tambahkan execution limits.
- [ ] Tambahkan kill switch.

## P12.2 Triggers

- [ ] Inbound message trigger.
- [ ] Keyword trigger.
- [ ] Regex trigger.
- [ ] New contact trigger.
- [ ] Tag added trigger.
- [ ] Tag removed trigger.
- [ ] Conversation opened trigger.
- [ ] Conversation resolved trigger.
- [ ] Delivery failed trigger.
- [ ] Schedule trigger.
- [ ] Webhook event trigger.
- [ ] Campaign event trigger.
- [ ] Contact field changed trigger.

## P12.3 Conditions

- [ ] Channel condition.
- [ ] Device condition.
- [ ] Contact field condition.
- [ ] Tag/segment condition.
- [ ] Message content condition.
- [ ] Business hours condition.
- [ ] Conversation status condition.
- [ ] Previous automation result condition.
- [ ] Quota/rate condition.
- [ ] Consent condition.

## P12.4 Actions

- [ ] Send text.
- [ ] Send template.
- [ ] Send media.
- [ ] Add tag.
- [ ] Remove tag.
- [ ] Assign agent.
- [ ] Assign team.
- [ ] Set status.
- [ ] Set priority.
- [ ] Add internal note.
- [ ] Delay/wait.
- [ ] Call webhook.
- [ ] Create task.
- [ ] Update contact field.
- [ ] Stop automation.
- [ ] Escalate to human.
- [ ] Notify supervisor.

## P12.5 Builder and execution

- [ ] Buat visual flow builder.
- [ ] Buat node validation.
- [ ] Buat draft save.
- [ ] Buat versioning.
- [ ] Buat publish flow.
- [ ] Buat disable flow.
- [ ] Buat dry-run sample contact.
- [ ] Buat simulation output.
- [ ] Cegah loop tanpa batas.
- [ ] Batasi execution depth.
- [ ] Lock per contact.
- [ ] Retry transient action.
- [ ] Dead-letter failed run.
- [ ] Replay run dengan audit.
- [ ] Tampilkan execution timeline.
- [ ] Tampilkan failure reason.
- [ ] Tambahkan automation metrics.

## P12.6 AI assist

- [ ] Pilih AI provider.
- [ ] Definisikan data boundary.
- [ ] Redact PII sebelum prompt.
- [ ] Definisikan allowed context.
- [ ] Implementasikan suggested reply.
- [ ] Implementasikan conversation summary.
- [ ] Implementasikan intent classification.
- [ ] Implementasikan sentiment classification.
- [ ] Implementasikan field extraction.
- [ ] Implementasikan knowledge-base retrieval.
- [ ] Wajibkan human approval untuk send default.
- [ ] Implementasikan handoff ke agent.
- [ ] Simpan AI run audit.
- [ ] Tambahkan model/provider cost.
- [ ] Tambahkan workspace AI quota.
- [ ] Tambahkan prompt/output redaction.
- [ ] Tambahkan AI kill switch.
- [ ] Tambahkan hallucination/failure UX.

---

# Track 13 — Webhook, API, dan Developer Platform

## P13.1 Public API foundation

- [ ] Tetapkan `/api/v1` contract.
- [ ] Implementasikan authentication header.
- [ ] Implementasikan workspace scope.
- [ ] Implementasikan API error envelope.
- [ ] Implementasikan request ID.
- [ ] Implementasikan cursor pagination.
- [ ] Implementasikan filtering/sorting contract.
- [ ] Implementasikan rate-limit headers.
- [ ] Implementasikan idempotency header.
- [ ] Implementasikan ETag/conditional update bila perlu.
- [ ] Tambahkan deprecation policy.
- [ ] Tambahkan API changelog.

## P13.2 API credential

- [ ] Buat API credential model.
- [ ] Buat scoped permission.
- [ ] Buat test/live environment.
- [ ] Generate prefix/key ID.
- [ ] Hash secret at rest.
- [ ] One-time secret reveal.
- [ ] Implementasikan expiry.
- [ ] Implementasikan rotate.
- [ ] Implementasikan revoke.
- [ ] Implementasikan IP allowlist.
- [ ] Simpan last-used metadata.
- [ ] Simpan usage metrics.
- [ ] Jangan return secret di list/get.
- [ ] Jangan log secret.
- [ ] Tambahkan credential audit.

## P13.3 Message API

- [ ] `POST /v1/messages`.
- [ ] Text payload.
- [ ] Template payload.
- [ ] Media payload.
- [ ] Location payload.
- [ ] Button/list payload.
- [ ] Recipient validation.
- [ ] Consent validation.
- [ ] Capability validation.
- [ ] Quota validation.
- [ ] Idempotent response.
- [ ] `GET /v1/messages/{id}`.
- [ ] `GET /v1/messages/{id}/events`.
- [ ] `POST /v1/messages/{id}/cancel`.
- [ ] Error/retry guidance.

## P13.4 Contact/conversation API

- [ ] Contact CRUD endpoints.
- [ ] Contact import endpoint.
- [ ] Contact export endpoint.
- [ ] Segment endpoint.
- [ ] Conversation list endpoint.
- [ ] Conversation detail endpoint.
- [ ] Assignment endpoint.
- [ ] Message timeline endpoint.
- [ ] Permission/scoping test setiap endpoint.

## P13.5 Webhook outgoing

- [ ] Buat webhook endpoint model.
- [ ] Buat subscription model.
- [ ] Buat delivery attempt model.
- [ ] Subscribe event types.
- [ ] Filter resource/channel.
- [ ] Enforce HTTPS.
- [ ] Implementasikan SSRF protection.
- [ ] Block localhost/private/link-local/metadata IP.
- [ ] Implementasikan hostname allowlist.
- [ ] Hapus/limit redirect.
- [ ] Generate signing secret.
- [ ] HMAC raw body.
- [ ] Timestamp tolerance.
- [ ] Replay protection.
- [ ] Secret rotation overlap.
- [ ] Retry exponential backoff.
- [ ] Dead-letter delivery.
- [ ] Pause failing endpoint.
- [ ] Test delivery.
- [ ] Manual replay.
- [ ] Delivery log.
- [ ] Payload versioning.

## P13.6 Webhook incoming/provider

- [ ] Verify provider signature.
- [ ] Verify timestamp.
- [ ] Reject replay.
- [ ] Acknowledge quickly.
- [ ] Persist raw event encrypted.
- [ ] Deduplicate event ID.
- [ ] Normalize event.
- [ ] Quarantine parse failures.
- [ ] Queue processing.
- [ ] Dead-letter event.
- [ ] Admin replay.
- [ ] Retention cleanup.

## P13.7 OpenAPI/developer portal

- [ ] Generate/maintain OpenAPI.
- [ ] Fix security scheme `Authorization`.
- [ ] Tambahkan route security requirements.
- [ ] Protect production API docs.
- [ ] Integrasikan Scalar.
- [ ] Add request examples.
- [ ] Add error examples.
- [ ] Add idempotency examples.
- [ ] Add webhook signature examples.
- [ ] Add cURL example.
- [ ] Add Go example.
- [ ] Add JavaScript example.
- [ ] Add PHP example.
- [ ] Add Python example.
- [ ] Add sandbox credentials.
- [ ] Add webhook tester.
- [ ] Add changelog.
- [ ] Add contract CI.

## P13.8 SDK dan integrations

- [ ] Generate TypeScript client.
- [ ] Generate Go client.
- [ ] Generate PHP client.
- [ ] Generate Python client.
- [ ] Add SDK retry guidance.
- [ ] Add SDK idempotency helper.
- [ ] Add n8n integration.
- [ ] Add Make integration.
- [ ] Add Zapier integration.
- [ ] Add WordPress integration.
- [ ] Add WooCommerce integration.
- [ ] Add Google Forms integration.
- [ ] Add integration credential revoke.
- [ ] Add integration health page.

---

# Track 14 — Analytics, Usage, dan Reporting

## P14.1 Event instrumentation

- [ ] Define product event catalog.
- [ ] Define technical event catalog.
- [ ] Define event schema/version.
- [ ] Track registration.
- [ ] Track workspace activation.
- [ ] Track channel connected.
- [ ] Track first test message.
- [ ] Track first delivered message.
- [ ] Track inbox activity.
- [ ] Track campaign launch.
- [ ] Track automation publish.
- [ ] Track API key use.
- [ ] Track webhook delivery.
- [ ] Track upgrade/cancel.
- [ ] Redact PII from analytics.

## P14.2 Operational metrics

- [ ] API request count.
- [ ] API latency p50/p95/p99.
- [ ] API error rate.
- [ ] Queue depth.
- [ ] Queue age.
- [ ] Retry count.
- [ ] Dead-letter count.
- [ ] Provider latency.
- [ ] Provider error rate.
- [ ] Message accepted/sent/delivered/read/failed.
- [ ] Delivery lag.
- [ ] Duplicate event count.
- [ ] Webhook success/failure.
- [ ] Automation success/failure.
- [ ] Channel connected/disconnected.
- [ ] Media storage usage.

## P14.3 Business analytics

- [ ] Active workspace.
- [ ] Active channel.
- [ ] Active contact.
- [ ] First response time.
- [ ] Resolution time.
- [ ] Agent workload.
- [ ] SLA breach.
- [ ] Campaign performance.
- [ ] Template performance.
- [ ] Automation conversion.
- [ ] Opt-out rate.
- [ ] Cost by channel.
- [ ] Cost by campaign.
- [ ] Cost by workspace.

## P14.4 Reports

- [ ] Dashboard realtime/near-realtime.
- [ ] Date range filter.
- [ ] Timezone-aware report.
- [ ] Comparison period.
- [ ] CSV export.
- [ ] Scheduled report.
- [ ] Email report delivery.
- [ ] Role-based visibility.
- [ ] Report retention.
- [ ] Report generation async.

---

# Track 15 — Billing, Plan, dan Usage Ledger

## P15.1 Plan and entitlement

- [ ] Define Sandbox plan.
- [ ] Define Starter plan.
- [ ] Define Growth plan.
- [ ] Define Scale plan.
- [ ] Define feature entitlements.
- [ ] Define seat limits.
- [ ] Define channel limits.
- [ ] Define message quota.
- [ ] Define media/storage quota.
- [ ] Define AI quota.
- [ ] Define retention per plan.
- [ ] Define support level.
- [ ] Store entitlement version.

## P15.2 Usage ledger

- [ ] Create usage record model.
- [ ] Record accepted message.
- [ ] Record delivered message.
- [ ] Define billing event semantics.
- [ ] Record provider fee separately.
- [ ] Record media storage.
- [ ] Record AI usage.
- [ ] Prevent duplicate usage count.
- [ ] Support correction/credit event.
- [ ] Build usage aggregation.
- [ ] Build quota check.
- [ ] Build usage alert threshold.
- [ ] Build hard cap.

## P15.3 Subscription and invoice

- [ ] Integrate payment provider.
- [ ] Create subscription.
- [ ] Upgrade plan.
- [ ] Downgrade plan.
- [ ] Handle proration.
- [ ] Trial period.
- [ ] Grace period.
- [ ] Cancel subscription.
- [ ] Pause send after quota policy.
- [ ] Invoice preview.
- [ ] Invoice history.
- [ ] Tax/currency handling.
- [ ] Refund/credit policy.
- [ ] Billing audit.

## P15.4 Pricing UX

- [ ] Pricing page.
- [ ] Plan comparison.
- [ ] 1K/10K/100K calculator.
- [ ] Transactional vs marketing cost explanation.
- [ ] Provider fee disclaimer.
- [ ] Overage estimate.
- [ ] Sandbox CTA.
- [ ] No hidden feature toggle.
- [ ] Usage dashboard.
- [ ] Upgrade CTA.
- [ ] Quota warning.

---

# Track 16 — Admin, Audit, Support, dan Compliance

## P16.1 Audit log

- [ ] Define audit event schema.
- [ ] Store actor type/id.
- [ ] Store workspace.
- [ ] Store action.
- [ ] Store resource type/id.
- [ ] Store redacted before/after.
- [ ] Store request ID.
- [ ] Store IP/user agent policy.
- [ ] Make audit append-only.
- [ ] Prevent user mutation.
- [ ] Search audit.
- [ ] Filter audit.
- [ ] Export audit.
- [ ] Apply retention.
- [ ] Alert sensitive actions.

## P16.2 Admin console

- [ ] Search workspace.
- [ ] Search user/member.
- [ ] Search channel.
- [ ] Search message by request ID.
- [ ] Inspect queue/job.
- [ ] Inspect dead-letter.
- [ ] Replay webhook.
- [ ] Replay job.
- [ ] Pause campaign.
- [ ] Pause channel.
- [ ] Apply feature flag.
- [ ] Apply temporary rate override.
- [ ] Show provider health.
- [ ] Show incident state.

## P16.3 Support

- [ ] Buat support ticket context.
- [ ] Attach request ID.
- [ ] Attach redacted diagnostics.
- [ ] Support impersonation consent flow.
- [ ] Show impersonation banner.
- [ ] Set impersonation expiry.
- [ ] Audit impersonation.
- [ ] Add customer-visible incident status.
- [ ] Add status page.
- [ ] Add support macros.
- [ ] Add recovery guide links.

## P16.4 Data rights and compliance

- [ ] Publish privacy policy.
- [ ] Publish terms.
- [ ] Publish provider/unofficial disclosure.
- [ ] Map PII data inventory.
- [ ] Map subprocessors.
- [ ] Implement data export request.
- [ ] Implement data deletion request.
- [ ] Implement anonymization.
- [ ] Implement retention jobs.
- [ ] Implement consent evidence export.
- [ ] Implement DPA workflow if required.
- [ ] Implement sensitive industry restrictions.
- [ ] Legal review health/finance use cases.

---

# Track 17 — Frontend, UX, dan Design System

## P17.1 App shell

- [ ] Implement workspace switcher.
- [ ] Implement responsive sidebar.
- [ ] Implement mobile navigation.
- [ ] Implement topbar notifications.
- [ ] Implement user/session menu.
- [ ] Implement breadcrumbs.
- [ ] Implement permission-aware navigation.
- [ ] Remove dead sidebar links.
- [ ] Add global loading indicator.
- [ ] Add global error boundary.
- [ ] Add flash/toast conventions.

## P17.2 Design system

- [ ] Standardize Nuxt UI components.
- [ ] Define color tokens.
- [ ] Define spacing tokens.
- [ ] Define typography.
- [ ] Define status colors.
- [ ] Define destructive action pattern.
- [ ] Define empty state.
- [ ] Define loading/skeleton state.
- [ ] Define error state.
- [ ] Define confirmation dialog.
- [ ] Define table/pagination.
- [ ] Define timeline.
- [ ] Define form validation.
- [ ] Define keyboard navigation.
- [ ] Add dark mode policy if required.

## P17.3 Required product routes

- [ ] `/login`.
- [ ] `/register`.
- [ ] `/verify-email`.
- [ ] `/forgot-password`.
- [ ] `/reset-password`.
- [ ] `/mfa`.
- [ ] `/workspaces`.
- [ ] `/workspace/settings`.
- [ ] `/workspace/members`.
- [ ] `/workspace/roles`.
- [ ] `/workspace/audit`.
- [ ] `/workspace/billing`.
- [ ] `/workspace/usage`.
- [ ] `/devices`.
- [ ] `/devices/new`.
- [ ] `/devices/:id`.
- [ ] `/devices/:id/edit`.
- [ ] `/devices/:id/connection`.
- [ ] `/devices/:id/logs`.
- [ ] `/phonebook`.
- [ ] `/phonebook/import`.
- [ ] `/phonebook/segments`.
- [ ] `/inbox`.
- [ ] `/inbox/:conversationId`.
- [ ] `/messages`.
- [ ] `/messages/compose`.
- [ ] `/messages/scheduled`.
- [ ] `/messages/failed`.
- [ ] `/templates`.
- [ ] `/templates/new`.
- [ ] `/templates/:id`.
- [ ] `/campaigns`.
- [ ] `/campaigns/new`.
- [ ] `/campaigns/:id`.
- [ ] `/recurring`.
- [ ] `/autoresponder`.
- [ ] `/autoresponder/:id/builder`.
- [ ] `/integrations`.
- [ ] `/integrations/webhooks`.
- [ ] `/integrations/api-keys`.
- [ ] `/documentation`.
- [ ] `/analytics`.
- [ ] `/reports`.
- [ ] `/notifications`.
- [ ] `/settings`.
- [ ] `/support`.

## P17.4 Existing placeholder cleanup

- [ ] Hubungkan delete device ke API.
- [ ] Hubungkan bulk delete ke API.
- [ ] Hubungkan delete detail device ke API.
- [ ] Hubungkan edit device ke endpoint update yang benar.
- [ ] Hapus `console.log` sebagai mekanisme bisnis.
- [ ] Hapus fake delay pada operation.
- [ ] Hapus toast sukses palsu.
- [ ] Tambahkan optimistic/pessimistic state yang benar.
- [ ] Sinkronkan device type UI dengan backend.
- [ ] Jika SMS/Email belum aktif, tampilkan disabled/coming soon yang jujur.
- [ ] Tampilkan experimental warning untuk whatsmeow.

## P17.5 Accessibility and localization

- [ ] Semantic headings.
- [ ] Keyboard navigation.
- [ ] Focus management dialog.
- [ ] Screen reader labels.
- [ ] Color contrast.
- [ ] Form error association.
- [ ] Motion/reduced-motion policy.
- [ ] ID locale.
- [ ] EN locale.
- [ ] Date/timezone formatting.
- [ ] Number/currency formatting.
- [ ] Accessibility audit.

---

# Track 18 — Security Engineering

## P18.1 Application security

- [ ] Enforce HTTPS production.
- [ ] Set secure cookie.
- [ ] Set HttpOnly cookie.
- [ ] Set SameSite policy.
- [ ] Enable CSP.
- [ ] Enable HSTS.
- [ ] Set X-Content-Type-Options.
- [ ] Set Referrer-Policy.
- [ ] Set Permissions-Policy.
- [ ] Configure CORS allowlist.
- [ ] Add Authorization allowed header where needed.
- [ ] Add API rate limit.
- [ ] Add login rate limit.
- [ ] Add body limit.
- [ ] Add read/write timeout.
- [ ] Add media timeout.
- [ ] Add concurrency limit.
- [ ] Add provider timeout.

## P18.2 SSRF and outbound request security

- [ ] Centralize outbound HTTP client.
- [ ] Enforce TLS verification.
- [ ] Block localhost.
- [ ] Block loopback.
- [ ] Block private IPv4.
- [ ] Block private IPv6.
- [ ] Block link-local.
- [ ] Block cloud metadata address.
- [ ] Resolve DNS and revalidate IP.
- [ ] Limit redirects.
- [ ] Restrict schemes to HTTPS where required.
- [ ] Apply domain allowlist where possible.
- [ ] Add outbound request audit.
- [ ] Add SSRF test suite.

## P18.3 Secrets and credentials

- [ ] Inventory all secrets.
- [ ] Remove plaintext token responses.
- [ ] Hash API/device tokens.
- [ ] Encrypt provider sessions.
- [ ] Use secret manager.
- [ ] Implement one-time reveal.
- [ ] Implement rotation.
- [ ] Implement revoke.
- [ ] Implement expiry.
- [ ] Redact logs.
- [ ] Scan repository for secrets.
- [ ] Scan built images for secrets.
- [ ] Rotate existing exposed credentials.

## P18.4 Tenant isolation

- [ ] Add workspace context middleware.
- [ ] Add workspace repository scope.
- [ ] Review every repository query.
- [ ] Review every service lookup.
- [ ] Review every worker job.
- [ ] Review cache keys.
- [ ] Review storage paths.
- [ ] Review webhook subscriptions.
- [ ] Review analytics aggregation.
- [ ] Review exports.
- [ ] Add cross-tenant automated test matrix.
- [ ] Evaluate PostgreSQL RLS as defense in depth.

## P18.5 Messaging abuse prevention

- [ ] Enforce consent.
- [ ] Enforce opt-out.
- [ ] Enforce suppression.
- [ ] Enforce quiet hours.
- [ ] Enforce daily limit.
- [ ] Enforce rate limit.
- [ ] Enforce campaign approval.
- [ ] Add recipient confirmation.
- [ ] Add test send.
- [ ] Add emergency pause.
- [ ] Add abuse detection.
- [ ] Add support escalation.

## P18.6 Dependency/security scanning

- [ ] Run frontend dependency audit.
- [ ] Upgrade vulnerable esbuild chain.
- [ ] Add Go `govulncheck`.
- [ ] Add SAST.
- [ ] Add secret scan.
- [ ] Add container image scan.
- [ ] Generate SBOM.
- [ ] Add license scan.
- [ ] Define severity policy.
- [ ] Define remediation SLA.
- [ ] Run periodic penetration test.

---

# Track 19 — Infrastructure, Workers, dan Observability

## P19.1 Runtime topology

- [ ] Define API deployment.
- [ ] Define worker deployment.
- [ ] Define scheduler deployment.
- [ ] Define provider worker deployment.
- [ ] Define PostgreSQL deployment.
- [ ] Define Valkey deployment.
- [ ] Define object storage.
- [ ] Define mail provider.
- [ ] Define TLS/load balancer.
- [ ] Define scaling limits.
- [ ] Define environment separation.

## P19.2 Queue and workers

- [ ] Implement message sender worker.
- [ ] Implement provider event worker.
- [ ] Implement webhook worker.
- [ ] Implement campaign worker.
- [ ] Implement scheduler worker.
- [ ] Implement media worker.
- [ ] Implement analytics worker.
- [ ] Implement notification worker.
- [ ] Implement retention worker.
- [ ] Implement dead-letter worker.
- [ ] Implement worker heartbeat.
- [ ] Implement worker graceful shutdown.
- [ ] Implement worker crash recovery.
- [ ] Implement worker concurrency config.
- [ ] Implement worker retry config.
- [ ] Implement worker metrics.

## P19.3 Cache and lock

- [ ] Namespace cache keys by environment/workspace.
- [ ] Set cache TTL.
- [ ] Prevent secret caching.
- [ ] Implement distributed lock.
- [ ] Implement lock expiry.
- [ ] Implement lock owner token.
- [ ] Handle lock loss safely.
- [ ] Implement rate limiter.
- [ ] Implement short-lived idempotency storage.
- [ ] Decide whether Valkey durability is sufficient.
- [ ] Document fallback behavior saat Valkey unavailable.

## P19.4 Logs

- [ ] Standardize structured JSON logs.
- [ ] Add request ID.
- [ ] Add workspace context safely.
- [ ] Add user/actor context safely.
- [ ] Add provider/channel/message context.
- [ ] Add operation/duration.
- [ ] Add error code.
- [ ] Redact token/password/session.
- [ ] Redact full phone/email.
- [ ] Redact message body.
- [ ] Redact webhook URL/query.
- [ ] Redact SQL bind values.
- [ ] Define retention/access.

## P19.5 Metrics and tracing

- [ ] Add HTTP metrics.
- [ ] Add auth metrics.
- [ ] Add rate-limit metrics.
- [ ] Add queue metrics.
- [ ] Add message lifecycle metrics.
- [ ] Add provider metrics.
- [ ] Add webhook metrics.
- [ ] Add campaign metrics.
- [ ] Add automation metrics.
- [ ] Add media metrics.
- [ ] Add database pool metrics.
- [ ] Add cache metrics.
- [ ] Add worker metrics.
- [ ] Add distributed traces.
- [ ] Propagate trace/correlation ID.
- [ ] Redact trace attributes.

## P19.6 Health and alerts

- [ ] Implement `/livez`.
- [ ] Implement `/readyz`.
- [ ] Implement internal detailed health endpoint.
- [ ] Check database readiness.
- [ ] Check Valkey readiness.
- [ ] Check storage readiness.
- [ ] Check worker readiness.
- [ ] Check provider health.
- [ ] Alert API error spike.
- [ ] Alert queue backlog.
- [ ] Alert dead-letter growth.
- [ ] Alert provider outage.
- [ ] Alert webhook failure.
- [ ] Alert database pool exhaustion.
- [ ] Alert storage quota.
- [ ] Alert backup failure.
- [ ] Alert certificate expiry.

## P19.7 Docker/Podman/deployment

- [ ] Add production Dockerfile stage.
- [ ] Run container as non-root.
- [ ] Remove development debug defaults.
- [ ] Remove default database password.
- [ ] Add image pinning.
- [ ] Add image SBOM.
- [ ] Add healthcheck.
- [ ] Add resource limits.
- [ ] Add graceful signal handling.
- [ ] Add staging compose/Podman stack.
- [ ] Add local PostgreSQL.
- [ ] Add local Valkey.
- [ ] Add local MinIO.
- [ ] Add local Mailpit.
- [ ] Add mock provider.
- [ ] Document TLS/reverse proxy.

## P19.8 Backup and disaster recovery

- [ ] Configure encrypted PostgreSQL backup.
- [ ] Configure point-in-time recovery.
- [ ] Configure object storage versioning.
- [ ] Configure object lifecycle.
- [ ] Document RPO.
- [ ] Document RTO.
- [ ] Test database restore.
- [ ] Test object restore.
- [ ] Test queue recovery.
- [ ] Test provider session recovery policy.
- [ ] Run periodic restore drill.
- [ ] Document incident evidence preservation.

---

# Track 20 — Testing dan Quality Assurance

## P20.1 Unit tests

- [ ] Auth service.
- [ ] Password/token hashing.
- [ ] Session policy.
- [ ] RBAC/policy.
- [ ] Tenant scope.
- [ ] Validation.
- [ ] Public ID mapper.
- [ ] Message state machine.
- [ ] Idempotency.
- [ ] Retry classification.
- [ ] Template renderer.
- [ ] Variable validation.
- [ ] Segment evaluator.
- [ ] Consent evaluator.
- [ ] Campaign guard.
- [ ] Scheduler timezone/DST.
- [ ] Automation evaluator.
- [ ] Webhook signature.
- [ ] Error mapping.
- [ ] Provider capability.

## P20.2 Repository/integration tests

- [ ] Auth/session database flow.
- [ ] Workspace isolation.
- [ ] Member invitation.
- [ ] API key revoke/rotate.
- [ ] Device token status/expiry/delete.
- [ ] Device ownership.
- [ ] Pagination page/limit.
- [ ] Negative limit rejection.
- [ ] Contact CRUD.
- [ ] Contact import.
- [ ] Consent/opt-out.
- [ ] Conversation/message persistence.
- [ ] Outbox transaction.
- [ ] Queue claim/retry.
- [ ] Webhook delivery.
- [ ] Campaign recipient snapshot.
- [ ] Scheduler lock.
- [ ] Billing usage ledger.
- [ ] Migration up/down/smoke.

## P20.3 Provider contract tests

- [ ] Official WhatsApp fake adapter.
- [ ] Whatsmeow fake/fixture.
- [ ] Telegram fake adapter.
- [ ] SMS fake adapter.
- [ ] Email fake adapter.
- [ ] Send text.
- [ ] Send media.
- [ ] Send template.
- [ ] Receive event.
- [ ] Delivery event.
- [ ] Rate limit error.
- [ ] Authentication error.
- [ ] Provider timeout.
- [ ] Provider malformed payload.
- [ ] Unsupported capability.
- [ ] Disconnect/reconnect.

## P20.4 API/contract tests

- [ ] Validate OpenAPI route coverage.
- [ ] Validate response schema.
- [ ] Validate request schema.
- [ ] Validate error envelope.
- [ ] Validate public ID contract.
- [ ] Validate auth security scheme.
- [ ] Validate pagination.
- [ ] Validate idempotency.
- [ ] Validate rate-limit headers.
- [ ] Validate webhook event schema.
- [ ] Run breaking-change check.
- [ ] Generate SDK and compile examples.

## P20.5 Frontend unit/component tests

- [ ] Auth forms.
- [ ] Workspace switcher.
- [ ] Permission guards.
- [ ] Device forms.
- [ ] Device delete modal.
- [ ] Contact import mapping.
- [ ] Consent UI.
- [ ] Inbox list/filter.
- [ ] Conversation timeline.
- [ ] Composer.
- [ ] Delivery status.
- [ ] Template editor.
- [ ] Campaign wizard.
- [ ] Automation builder nodes.
- [ ] Webhook form.
- [ ] Billing/usage display.
- [ ] Zod invalid payload behavior.
- [ ] Empty/loading/error states.

## P20.6 E2E tests

- [ ] Register/login/logout.
- [ ] Verify email.
- [ ] Reset password.
- [ ] MFA.
- [ ] Create workspace.
- [ ] Invite member.
- [ ] Change role.
- [ ] Connect official test channel.
- [ ] Connect experimental channel with warning.
- [ ] Send test message.
- [ ] Receive inbound message.
- [ ] Reply from inbox.
- [ ] Assign conversation.
- [ ] Import contact.
- [ ] Respect opt-out.
- [ ] Create template.
- [ ] Create campaign.
- [ ] Approve/launch campaign.
- [ ] Pause campaign.
- [ ] Schedule recurring message.
- [ ] Publish automation.
- [ ] Test webhook.
- [ ] Rotate API key.
- [ ] Revoke device token.
- [ ] View analytics.
- [ ] Upgrade plan.
- [ ] Request export/delete.

## P20.7 Resilience/security/accessibility tests

- [ ] Provider timeout.
- [ ] Provider outage.
- [ ] Queue restart.
- [ ] Worker crash.
- [ ] Database failover behavior.
- [ ] Valkey unavailable behavior.
- [ ] Duplicate send race.
- [ ] Webhook replay.
- [ ] Webhook endpoint down.
- [ ] Large media.
- [ ] Malicious media.
- [ ] SSRF payload.
- [ ] Cross-tenant access.
- [ ] CSRF bypass attempt.
- [ ] Credential leakage scan.
- [ ] Rate-limit bypass attempt.
- [ ] Load test.
- [ ] Accessibility audit WCAG 2.1 AA.
- [ ] Browser compatibility test.

## P20.8 Quality gates

- [ ] `pnpm type-check` lulus.
- [ ] Frontend lint/oxlint lulus.
- [ ] Format check lulus tanpa generated-file exception.
- [ ] Frontend unit test tersedia dan lulus.
- [ ] Playwright E2E tersedia dan lulus.
- [ ] `go test ./...` lulus.
- [ ] `go vet ./...` lulus.
- [ ] Staticcheck lulus.
- [ ] Govulncheck lulus.
- [ ] Dependency audit policy lulus.
- [ ] Secret scan lulus.
- [ ] Container scan lulus.
- [ ] OpenAPI compatibility lulus.
- [ ] Migration test lulus.
- [ ] Backup restore test lulus.

---

# Track 21 — Documentation, Support, dan Developer Experience

## P21.1 Product documentation

- [ ] Update product brief.
- [ ] Update full feature specification.
- [ ] Update master plan.
- [ ] Update release scope.
- [ ] Update known limitations.
- [ ] Update official/unofficial channel policy.
- [ ] Update pricing rules.
- [ ] Update retention policy.

## P21.2 User documentation

- [ ] Getting started.
- [ ] Workspace setup.
- [ ] Invite team.
- [ ] Connect official WhatsApp.
- [ ] Connect Telegram.
- [ ] Connect SMS.
- [ ] Connect Email.
- [ ] Experimental whatsmeow warning.
- [ ] Send first message.
- [ ] Read delivery status.
- [ ] Resolve failed message.
- [ ] Use inbox.
- [ ] Manage contacts.
- [ ] Import contacts.
- [ ] Manage consent.
- [ ] Create template.
- [ ] Create campaign.
- [ ] Schedule recurring message.
- [ ] Configure automation.
- [ ] Configure webhook.
- [ ] Read analytics.
- [ ] Manage billing.
- [ ] Export/delete data.

## P21.3 Developer documentation

- [ ] API authentication.
- [ ] API quickstart.
- [ ] Send text.
- [ ] Send template.
- [ ] Send media.
- [ ] Get status.
- [ ] Handle idempotency.
- [ ] Handle errors.
- [ ] Handle rate limits.
- [ ] Verify webhook signature.
- [ ] Replay webhook safely.
- [ ] API key rotate/revoke.
- [ ] SDK examples.
- [ ] Sandbox usage.
- [ ] API versioning/deprecation.

## P21.4 Operations runbooks

- [ ] Deploy.
- [ ] Rollback.
- [ ] Database migration.
- [ ] Migration failure.
- [ ] Backup.
- [ ] Restore.
- [ ] Provider outage.
- [ ] Duplicate campaign.
- [ ] Queue backlog.
- [ ] Dead-letter replay.
- [ ] Credential leak.
- [ ] Cross-tenant incident.
- [ ] Webhook outage.
- [ ] Storage outage.
- [ ] Certificate expiry.
- [ ] Security incident.
- [ ] Customer data deletion.

## P21.5 Customer support readiness

- [ ] Define support channels.
- [ ] Define SLA per plan.
- [ ] Create troubleshooting macros.
- [ ] Create provider error guide.
- [ ] Create escalation matrix.
- [ ] Link request ID to ticket.
- [ ] Define incident communication.
- [ ] Define refund/credit policy.
- [ ] Train support on unofficial risk.
- [ ] Train support on consent/compliance.

---

# Track 22 — CI/CD dan Release Management

## P22.1 CI pipeline

- [ ] Run changed-file detection.
- [ ] Run Go format.
- [ ] Run frontend format.
- [ ] Run Go lint.
- [ ] Run frontend lint.
- [ ] Run type-check.
- [ ] Run unit tests.
- [ ] Run integration tests.
- [ ] Run contract tests.
- [ ] Run frontend build.
- [ ] Run E2E.
- [ ] Run OpenAPI check.
- [ ] Run migration check.
- [ ] Run secret scan.
- [ ] Run dependency scan.
- [ ] Run container scan.
- [ ] Generate SBOM.
- [ ] Upload test/coverage artifacts.

## P22.2 Deployment

- [ ] Build immutable image.
- [ ] Push signed image.
- [ ] Deploy staging.
- [ ] Run staging smoke test.
- [ ] Run migration compatibility check.
- [ ] Require production approval.
- [ ] Run backward-compatible migration.
- [ ] Deploy API.
- [ ] Deploy workers.
- [ ] Check readiness.
- [ ] Run production smoke test.
- [ ] Monitor rollout.
- [ ] Support rollback.
- [ ] Document release notes.

## P22.3 Feature flags

- [ ] Add channel feature flag.
- [ ] Add campaign feature flag.
- [ ] Add automation feature flag.
- [ ] Add AI feature flag.
- [ ] Add whatsmeow allowlist flag.
- [ ] Add billing rollout flag.
- [ ] Define flag owner.
- [ ] Define flag expiry.
- [ ] Audit flag changes.
- [ ] Add emergency kill switch.

## P22.4 Release gates

- [ ] Alpha gate.
  - [ ] Internal users only.
  - [ ] Mock/provider fixture available.
  - [ ] No real bulk campaign.
  - [ ] Critical security findings closed.
- [ ] Private beta gate.
  - [ ] 5–10 pilot workspaces.
  - [ ] Manual support.
  - [ ] Daily error review.
  - [ ] Incident runbook.
  - [ ] Strict send limits.
- [ ] Public beta gate.
  - [ ] Self-serve onboarding.
  - [ ] Billing/usage visible.
  - [ ] Status page.
  - [ ] Public docs.
  - [ ] Security review.
- [ ] GA gate.
  - [ ] SLO dashboard.
  - [ ] Restore drill.
  - [ ] Provider failure drill.
  - [ ] Rollback drill.
  - [ ] Support SLA.
  - [ ] No unresolved critical security issue.

---

# Track 23 — Launch, Metrics, dan Commercial Readiness

## P23.1 Activation funnel

- [ ] Define visitor-to-register event.
- [ ] Define register-to-workspace event.
- [ ] Define workspace-to-channel event.
- [ ] Define channel-to-test-message event.
- [ ] Define test-to-delivered event.
- [ ] Define delivered-to-inbox event.
- [ ] Define inbox-to-paid event.
- [ ] Dashboard activation funnel.
- [ ] Dashboard drop-off reasons.

## P23.2 Launch targets

- [ ] Set onboarding completion target.
- [ ] Set first-delivery time target.
- [ ] Set delivery success target.
- [ ] Set response SLA target.
- [ ] Set support response target.
- [ ] Set pilot workspace target.
- [ ] Set paying customer target.
- [ ] Set churn target.
- [ ] Set gross margin target.
- [ ] Review targets per phase.

## P23.3 Go-to-market readiness

- [ ] Create landing page.
- [ ] Create pricing page.
- [ ] Create demo flow.
- [ ] Create sample use cases.
- [ ] Create onboarding emails.
- [ ] Create developer quickstart.
- [ ] Create comparison positioning without unsupported claims.
- [ ] Explain provider fees.
- [ ] Explain unofficial risk.
- [ ] Prepare pilot agreement.
- [ ] Prepare feedback loop.
- [ ] Prepare changelog.

---

# Track 24 — Final Full-Feature Definition of Done

- [ ] Semua Track 0–23 sudah memiliki owner dan status.
- [ ] Tidak ada route aktif yang memanggil `panic("unimplemented")`.
- [ ] Tidak ada placeholder UI yang terlihat sebagai fitur selesai.
- [ ] Official WhatsApp production path tersedia.
- [ ] Telegram adapter lengkap atau seluruh capability yang tidak tersedia ditangani sebagai error terkontrak.
- [ ] SMS adapter tersedia.
- [ ] Email adapter tersedia.
- [ ] Whatsmeow terisolasi, diberi label experimental, dan memiliki kill switch.
- [ ] Workspace dan tenant isolation diuji.
- [ ] RBAC/resource authorization diuji.
- [ ] Auth/MFA/session/recovery selesai.
- [ ] Credential tidak plaintext, tidak bocor ke URL/log/response, dan dapat rotate/revoke.
- [ ] CSRF/SSRF/CORS/CSP/TLS/rate-limit/timeout/body-limit selesai.
- [ ] Contact, consent, opt-out, suppression selesai.
- [ ] Inbox, assignment, SLA, notes, handoff selesai.
- [ ] Semua message type yang dijanjikan memiliki capability check.
- [ ] Message lifecycle, queue, outbox, idempotency, retry, DLQ, replay selesai.
- [ ] Media storage aman dan tidak bergantung pada filesystem container.
- [ ] Template versioning/approval/variable safety selesai.
- [ ] Campaign approval/audience/throttle/cost/report selesai.
- [ ] Schedule/recurring/follow-up durable dan timezone-safe.
- [ ] Automation builder/run history/limits/kill switch selesai.
- [ ] AI assist memiliki data boundary, human approval, quota, dan audit.
- [ ] API v1/OpenAPI/Scalar/SDK/sandbox selesai.
- [ ] Webhook signing/replay/retry/DLQ/replay UI selesai.
- [ ] Analytics/reporting/usage selesai.
- [ ] Billing/plan/quota/invoice/overage selesai.
- [ ] Admin/audit/support/compliance selesai.
- [ ] CI/CD/security scan/SBOM selesai.
- [ ] Unit/integration/contract/E2E/load/security/accessibility test selesai.
- [ ] Health/metrics/tracing/logging/alerting selesai.
- [ ] Backup/restore/DR/rollback drill selesai.
- [ ] Runbook dan user/developer documentation selesai.
- [ ] Pilot/public-beta/GA release gate selesai.
- [ ] Product owner, engineering, security, QA, support, dan operations menyetujui release.
