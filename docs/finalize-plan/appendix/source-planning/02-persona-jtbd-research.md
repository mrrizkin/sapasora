# 02 — Persona, JTBD, dan Research Plan

## Persona 1 — Business owner

**Profil:** owner/ops SMB, bukan developer, 1–5 nomor/channel.

**Jobs to be done:**

- Saat ada order atau appointment, saya ingin pelanggan menerima konfirmasi otomatis.
- Saat customer menghubungi bisnis, saya ingin agent bisa membalas dari satu inbox.
- Saat membuat campaign, saya ingin tahu siapa yang akan menerima dan berapa biayanya.
- Saat sesuatu gagal, saya ingin diberi alasan dan solusi tanpa membaca dokumentasi teknis.

**Pain points:** setup Meta, template rejected, nomor tidak connect, status tidak jelas, takut spam/ban.

**Success:** first value <15 menit; dapat menyelesaikan issue umum tanpa support.

## Persona 2 — Developer

**Profil:** Go/Node/PHP developer, membangun aplikasi internal atau client project.

**Jobs to be done:**

- Saya ingin mengirim event bisnis ke WhatsApp dengan API yang konsisten.
- Saya ingin menerima inbound/status event melalui webhook yang aman.
- Saya ingin replay/debug request tanpa meminta akses dashboard client.
- Saya ingin setiap request idempotent dan dapat dilacak.

**Pain points:** docs tidak sinkron, token manual, error generik, retry tidak jelas, webhook tanpa signature.

**Success:** integration proof-of-concept <1 jam; production checklist jelas.

## Persona 3 — Support lead

**Profil:** mengelola 3–20 agent/customer service.

**Jobs to be done:**

- Membagi conversation secara adil.
- Mengetahui conversation yang belum dijawab.
- Menyimpan internal notes dan canned replies.
- Melihat performance agent dan response time.

**Pain points:** tab WhatsApp personal, conversation hilang, double reply, tidak ada audit trail.

## Persona 4 — Agency/solution implementer

**Profil:** membangun automasi untuk beberapa bisnis.

**Jobs to be done:**

- Membuat banyak workspace secara terisolasi.
- Menghubungkan channel milik client tanpa berbagi credential.
- Memantau status semua deployment.
- Mengelola billing/support tiap client.

## JTBD prioritas MVP

1. Connect official channel.
2. Send test message dan melihat delivery result.
3. Receive inbound conversation.
4. Reply dan assign conversation.
5. Create template/automation sederhana.
6. Integrate via API/webhook.
7. Diagnose failed delivery.

## Research questions

### Market

- Use case mana yang terjadi setiap hari?
- Apakah customer sudah memiliki WABA/Meta Business?
- Apakah mereka lebih memilih official, BSP, atau personal QR?
- Biaya dan risiko apa yang dianggap dapat diterima?

### UX

- Bagian onboarding mana yang paling membingungkan?
- Status apa yang ingin dilihat owner vs developer?
- Apakah campaign approval dan consent dipahami?
- Berapa lama sampai first value?

### Technical

- Framework/backend yang sudah digunakan customer?
- Perlu SDK bahasa apa saja?
- Payload/event apa yang wajib untuk integrasi?
- Kebutuhan retry, latency, throughput, dan retention?

## Research plan

### Tahap 1 — discovery interview

- 5 business owner.
- 5 developer.
- 3 support lead.
- 2 agency.

Durasi 30–45 menit, tanpa menjual produk. Minta contoh workflow aktual, bukan hanya opini.

### Tahap 2 — prototype test

Prototype minimal:

- connect channel wizard;
- send test message;
- inbox;
- failed-message recovery;
- pricing calculator.

Ukuran: task completion, time on task, error count, confidence score.

### Tahap 3 — pilot

5–10 workspace dengan use case berbeda. Tetapkan success criteria sebelum onboarding.

## Interview guide

1. Ceritakan cara Anda mengirim/menjawab WhatsApp sekarang.
2. Apa yang terjadi ketika nomor atau pesan gagal?
3. Siapa yang mengelola nomor dan credential?
4. Apa yang Anda anggap sebagai pesan sukses?
5. Pernah terkena spam report/ban/disconnect?
6. Fitur apa yang paling sering dipakai? Mana yang tidak?
7. Berapa biaya operasional saat ini?
8. Apa yang harus tersedia agar Anda pindah ke produk baru?
9. Siapa yang menyetujui campaign dan template?
10. Apa yang tidak boleh terjadi?

## Research ethics

- Jangan meminta token, nomor pribadi, atau payload customer asli.
- Gunakan data sintetis.
- Jelaskan recording/consent.
- Redact PII dari screenshot dan logs.
- Jangan menguji broadcast ke nomor tanpa consent.
