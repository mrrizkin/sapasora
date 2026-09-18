# 05 — Evaluasi Product Design

Perspektif: product designer / product owner yang menilai produk publik, onboarding, trust, dan kesiapan penggunaan.

## Positioning

Fonnte menjual **otomatisasi WhatsApp berbasis WhatsApp Web** sebagai alternatif sederhana dan murah untuk official API. Narasi produknya kuat untuk developer/UMKM yang mengutamakan time-to-value, tetapi trade-off reliability dan compliance cukup besar.

## Yang sudah bagus

### 1. Jalur value cepat

Homepage langsung menyebut outcome (“kirim pesan WhatsApp secara otomatis”) dan CTA “Coba Gratis”. Dokumentasi first-time user merangkum happy path menjadi add device → send message → optional autoreply/chatbot.

### 2. Feature-to-use-case mapping

Fitur tidak hanya disebut sebagai jargon teknis; website memberi contoh notifikasi, chatbot, reminder, tagihan, informasi, OTP, broadcast, template, follow-up, dan integrasi form. Ini membantu calon user mengaitkan produk dengan pekerjaan nyata.

### 3. Developer enablement

API endpoint, parameter, contoh PHP, contoh response sukses/error, Postman, JavaScript/Node.js, serta webhook tersedia. Dokumentasi API cukup actionable untuk proof of concept.

### 4. Self-serve dan low-cost adoption

Free tier, pricing IDR/USD, monthly/annual toggle, serta registration flow sederhana mengurangi hambatan awal. Register meminta Name, WhatsApp, Email, Password, Password Confirm, dan Country.

### 5. Transparansi risiko

Fonnte cukup terbuka soal unofficial API, potensi banned, tidak ada fixed SLA, no refund, ketergantungan pada WhatsApp Web, serta status tutorial sebagai POC—not production-ready. Transparansi ini bagus untuk expectation setting walaupun memperbesar perceived risk.

### 6. Support dan bantuan kontekstual

Ada support ticket yang meminta nomor WhatsApp untuk follow-up, knowledge base yang terstruktur, dan “Ask Fai” di docs. Ask Fai membuka chat Aksita.ai dengan prompt starter “Tentang Fonnte”, “Harga”, “Fitur”, dan “Testimoni”. Ini adalah pola contextual help yang menjanjikan.

## Friction dan risiko UX

### 1. Trust/compliance adalah blocker utama

Produk mengotomatisasi akun WhatsApp pengguna melalui WhatsApp Web. Pengguna menanggung risiko banned; tidak ada SLA dan tidak ada refund. Untuk bisnis/instansi, ini dapat menghalangi procurement meskipun feature set menarik.

**Usulan:** buat risk center yang ringkas dan actionable: acceptable use, warming-up number, opt-in, rate limit, broadcast hygiene, banned recovery, escalation, dan compatibility status.

### 2. Homepage terlalu padat

Satu halaman menampung positioning, fungsi, kelebihan, puluhan fitur + demo, integrasi, testimonial, social proof, pricing, FAQ, disclaimer, dan CTA. Feature inventory sangat lengkap, tetapi hierarchy dan prioritas persona kurang terlihat.

**Usulan:** pecah menjadi jalur persona:

- “Saya developer” → API/webhook quickstart.
- “Saya pemilik bisnis” → dashboard, broadcast, template, autoreply.
- “Saya pengguna WordPress/form” → plugin/integrations.

### 3. Bahasa dan consistency

Ada campuran Indonesia/English, typo/capitalization seperti “Mudah DIgunakan”, “All Feature”, dan istilah yang kadang belum konsisten. Fitur deprecated masih muncul di feature list dan pricing.

**Usulan:** content design system: satu bahasa per mode, glossary, status badge (Stable/Beta/Deprecated), last updated, owner, dan compatibility notes.

### 4. Pricing kurang mudah dibandingkan

Dua dimensi toggle (IDR/USD dan bulanan/tahunan) ditambah dua kelompok paket dapat membuat user sulit menjawab “paket mana yang saya butuhkan?”. Entitlement tidak tersampaikan sejelas harga/kuota.

**Usulan:** comparison table dengan recommended plan, cost per 1.000 pesan, effective monthly annual price, included media limits, multi-agent details, dan calculator berdasarkan jumlah device/target.

### 5. Onboarding teknis berisiko salah konfigurasi

User harus add device, scan QR, copy token, connect, memahami country code, target format, quota per recipient, dan autoread. Untuk webhook, URL harus public dan autoread wajib on.

**Usulan:** onboarding checklist interaktif dengan health checks:

- device connected;
- test message delivered;
- token test berhasil;
- webhook reachable;
- autoread enabled;
- response payload sample;
- test chatbot pass.

### 6. Security surface belum matang untuk production

Dokumentasi mengingatkan token cukup untuk mengirim pesan atas nama WhatsApp sehingga harus disimpan aman. Tutorial webhook yang diamati juga menunjukkan implementasi endpoint publik; sebuah komentar Februari 2024 di tutorial menyebut webhook signature belum tersedia saat itu dan masih roadmap.

**Usulan prioritas tinggi:** webhook signing secret/HMAC, IP/rate controls, token scopes per device, rotate/revoke token, audit log, idempotency key, retry policy, request examples yang memvalidasi signature, dan redaction log.

### 7. Ambiguitas lifecycle dan data

Website menyatakan data tidak digunakan untuk kepentingan Fonnte dan memiliki lifecycle seperti webhook files 30 menit, temporary/scheduled files 7 hari, account 30 hari setelah tanpa device, tetapi user tetap perlu memahami data apa yang tersimpan dan kapan terhapus.

**Usulan:** privacy center yang menjelaskan data map, retention timer per jenis data, delete/export controls, dan status penghapusan.

### 8. Inconsistency dokumentasi

Batas attachment di halaman File Limitation tertulis maksimum 10 MB, sedangkan contoh error response API menyebut file harus di bawah 4 MB. Ini dapat menimbulkan kegagalan operasional dan tiket support.

**Usulan:** jadikan satu source of truth, tampilkan limit per MIME/type/package di form dan API docs, serta tambahkan automated documentation checks.

## Prioritas roadmap yang disarankan

### P0 — Trust dan safety

1. Webhook signing/authentication.
2. Token rotation/revoke dan scoped credentials.
3. Reconcile attachment limits.
4. Risk/acceptable-use onboarding sebelum user menghubungkan nomor.

### P1 — Activation

1. Guided setup + test message.
2. Persona-based quickstarts.
3. Delivery status, retry, failure explanation, dan monitoring.
4. Production-ready examples selain PHP.

### P2 — Monetization dan scale

1. Pricing calculator/comparison.
2. Multi-agent/CS workflow yang lebih jelas.
3. Better campaign controls: opt-in, suppression list, rate guardrails.
4. Audit log dan team roles.

### P3 — Content quality

1. Konsistensi bahasa.
2. Video captions/transcript dan timestamp.
3. Versioning artikel, changelog, deprecation migration guide.
4. Accessibility review untuk pricing toggle, form, dan embedded demo.

## Pertanyaan validasi berikutnya

- Berapa activation rate dari register ke device connected?
- Berapa waktu median sampai first successful message?
- Berapa % device disconnect/banned per cohort dan per use case?
- Apakah user memahami bahwa quota dihitung per target?
- Berapa banyak tiket support berasal dari token, QR, webhook, atau attachment?
- Apakah free-tier user berhasil naik paket karena quota, media, multics agent, atau faktor lain?
- Apakah bisnis/instansi membutuhkan SLA, DPA, invoice formal, atau team access?
- Apakah WhatsApp Web dependency dapat diterima untuk segmen target?
