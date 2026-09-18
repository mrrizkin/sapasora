# 07 — Blueprint SaaS WhatsApp Automation Versi Improve

## Executive direction

Jangan membangun sekadar “Fonnte clone”. Bangun produk **WhatsApp Engagement Platform** yang:

> **Catatan penting:** produk bisa mendukung dua jalur koneksi—official dan unofficial—tetapi keduanya harus diposisikan sebagai produk/capability yang berbeda, bukan disamarkan sebagai hal yang sama.

- **Official Cloud API**: jalur production/default untuk bisnis yang membutuhkan compliance, reliability, dan procurement.
- **Unofficial WhatsApp Web connector**: jalur experimental/advanced, misalnya connector berbasis `whatsmeow`, dengan disclosure risiko yang sangat jelas.

Jangan menjanjikan SLA yang sama untuk dua jalur tersebut.

1. aman dan transparan sejak awal;
2. bisa dipakai nonteknis tanpa bantuan developer;
3. tetap punya API/webhook yang serius untuk developer;
4. memiliki observability dan recovery yang jelas;
5. menggunakan WhatsApp Business Platform/Cloud API secara resmi sebagai jalur utama.

Working positioning:

> **Cara paling aman dan mudah untuk mengotomatisasi customer communication di WhatsApp.**

Target awal yang disarankan: bisnis Indonesia dengan 1–10 agent dan kebutuhan notifikasi, follow-up, inbox, automation, serta integrasi ringan.

## Strategi koneksi: official + unofficial

### Official Cloud API

Cocok untuk default SaaS dan pelanggan production. Kelebihan utama: jalur resmi, lifecycle credential lebih jelas, template dan policy lebih terdefinisi, serta lebih mudah diposisikan untuk bisnis/enterprise. Trade-off-nya adalah onboarding Meta/WABA, approval template, biaya provider, dan sebagian capability tidak sama dengan WhatsApp personal.

### Unofficial connector berbasis whatsmeow

`whatsmeow` adalah library Go untuk berkomunikasi dengan protokol WhatsApp multi-device/WhatsApp Web; ini bukan WhatsApp Business Cloud API resmi. Secara produk, connector ini dapat menarik use case seperti personal number, QR/pairing, dan capability yang tidak tersedia di Cloud API, tetapi memiliki risiko:

- akun dapat logout, terputus, atau terkena ban;
- protokol dapat berubah tanpa notice;
- maintenance dan kompatibilitas menjadi tanggung jawab platform;
- tidak ada jaminan dari Meta atau library;
- handling session credential sangat sensitif;
- kebijakan penggunaan dan legal review harus dilakukan sebelum komersialisasi.

Positioning yang aman: **experimental connector / bring-your-own-number**, opt-in, tanpa SLA official, dan tidak diaktifkan default untuk workspace baru.

### Arsitektur provider-agnostic

Pisahkan domain produk dari implementasi channel:

```text
Core domain
  contacts · inbox · campaigns · automations · billing · analytics
                         |
                 Channel Provider API
                 /                    \
       Official Cloud API       Unofficial whatsmeow
```

Gunakan capability matrix per provider, misalnya:

| Capability | Official | Unofficial |
|---|---|---|
| Connection | Meta/WABA onboarding | QR/pairing + session |
| Personal number | Tergantung eligibility/provider | Ya, dengan risiko |
| Template approval | Umumnya wajib untuk outbound tertentu | Behavior berbeda |
| Groups | Capability/policy terbatas | Bisa tersedia, tetap berisiko |
| Delivery status | Provider webhook | Event dari session/protocol |
| SLA | Bisa dikontrakkan sesuai provider | Jangan disamakan |
| Ban risk | Policy/compliance risk | Lebih tinggi dan langsung ke nomor |

Nama capability harus muncul di UI sebelum user memilih channel. Jangan membuat fitur yang terlihat universal jika behavior provider berbeda.

### Boundary teknis yang wajib

- Interface `ChannelProvider`: connect, disconnect, health, send, receive events, media, status, revoke.
- Normalized internal event model, dengan raw provider payload tetap disimpan untuk debugging terbatas.
- Queue dan worker terpisah per provider.
- Session store terenkripsi dan access audit untuk unofficial connector.
- Provider-specific rate guard dan circuit breaker.
- Feature flag, kill switch, dan compatibility version untuk unofficial connector.
- Data residency/retention policy yang jelas untuk session dan message payload.
- Terms, consent, risk acknowledgement, dan export/delete session sebelum aktivasi unofficial.

## Apakah SaaS perlu approval untuk Official API?

Jawabannya bergantung pada scope:

### Jika hanya untuk bisnis sendiri

Bisa mulai implementasi langsung dalam **Development Mode** menggunakan Meta App, test number, access token, Phone Number ID, dan webhook. Untuk production, tetap perlu menyiapkan Business Portfolio/WABA, business phone number, credential production, policy/opt-in, dan konfigurasi webhook. Tidak perlu menjadi Tech Provider jika aplikasinya hanya dipakai oleh bisnis sendiri.

### Jika SaaS meng-onboard banyak bisnis/customer

Tidak cukup hanya membuat satu token lalu dipakai semua customer. Jalur yang disarankan adalah menjadi **Meta Tech Provider** atau memakai infrastruktur partner/BSP yang sudah mendukung onboarding customer.

Kebutuhan umumnya:

- Meta App dan konfigurasi WhatsApp Embedded Signup.
- Business Verification untuk bisnis/platform Anda.
- App Review/Advanced Access untuk permission yang diperlukan.
- Facebook Login for Business, HTTPS domain, redirect URI, dan `config_id`.
- Backend untuk menerima hasil onboarding, mengelola WABA/Phone Number ID, template, credential, dan webhook subscription.
- Proses billing: customer membayar langsung provider/Meta, atau SaaS mengelola model billing sesuai eligibility partner.
- Terms, privacy policy, consent, opt-in, dan support process yang jelas.

Untuk validasi cepat, gunakan partner/BSP terlebih dahulu. Untuk kontrol biaya, onboarding, dan margin jangka panjang, pertimbangkan Tech Provider path. Status approval, permission, pricing, dan partner eligibility dapat berubah sehingga harus dicek kembali pada dokumentasi Meta terbaru.

**Rule arsitektur:** satu customer sebaiknya memiliki mapping yang jelas ke WABA, Phone Number ID, credential, billing owner, dan webhook context-nya sendiri. Jangan mencampur token antar-customer.

## Masalah yang harus diselesaikan

### Untuk business owner

- Tidak tahu harus mulai dari mana.
- Takut salah konfigurasi device, token, webhook, atau template.
- Sulit mengetahui pesan gagal karena apa.
- Broadcast berisiko spam dan merusak reputasi nomor.
- Tidak jelas paket mana yang sesuai.

### Untuk developer

- API cepat dibuat tetapi kurang production-grade.
- Webhook perlu signature, retry, idempotency, dan observability.
- Status pesan sering asynchronous dan sulit dilacak.
- Dokumentasi tidak selalu konsisten dengan behavior aktual.
- Perlu sandbox dan test number sebelum menghubungkan nomor produksi.

### Untuk customer support/team

- Memerlukan shared inbox, assignment, collision prevention, notes, SLA, dan audit trail.
- Perlu tahu agent mana yang membalas dan kapan.

## Prinsip produk

1. **Official-first** — gunakan API resmi; jelaskan limitasi WhatsApp dengan bahasa manusia.
2. **Safe by default** — consent, opt-out, suppression list, rate guard, dan approval campaign aktif secara default.
3. **Progressive disclosure** — user pemula melihat happy path; advanced settings tetap tersedia tanpa memenuhi layar awal.
4. **Observable by default** — setiap message punya lifecycle, reason, request ID, dan timeline.
5. **One source of truth** — pricing, limits, API docs, dan UI memakai schema/kontrak yang sama.
6. **Design for recovery** — disconnect, token revoke, template rejected, rate limit, dan webhook failure memiliki next action yang jelas.
7. **API parity** — fitur dashboard dan API menggunakan capability model yang sama.

## Information architecture yang disarankan

- **Home / Overview** — health device, delivery rate, quota, unresolved issues.
- **Inbox** — percakapan, assignment, tags, internal notes, canned replies.
- **Contacts** — contact, segments, consent, opt-out, import CSV.
- **Campaigns** — draft, approval, schedule, audience, preview, delivery report.
- **Automations** — trigger → condition → action flow builder.
- **Templates** — WhatsApp-approved templates, variables, version, approval status.
- **Integrations** — API keys, webhooks, SDK, n8n/Zapier, WordPress.
- **Channels** — nomor, WABA, device/channel health.
- **Analytics** — sent, delivered, read, failed, response time, conversion event.
- **Settings** — team, roles, billing, retention, audit log.
- **Developer Center** — sandbox, API docs, logs, webhook tester, code examples.

## Onboarding UX ideal

### Step 1 — Pilih tujuan

User memilih satu atau lebih:

- Notifikasi transaksi
- Customer support/inbox
- Broadcast/campaign
- Chatbot/automation
- Integrasi aplikasi

Pilihan ini menentukan checklist dan rekomendasi setup.

### Step 2 — Hubungkan channel

Wizard menjelaskan:

- official vs unsupported connection;
- prasyarat Meta/WABA;
- estimasi waktu;
- biaya dari Meta vs biaya platform;
- status approval dan template.

Jangan meminta token mentah tanpa penjelasan. Credential disimpan terenkripsi dan hanya ditampilkan sekali.

### Step 3 — Kirim test message

User memilih nomor test, menulis pesan, lalu melihat hasil dalam satu layar:

- accepted;
- queued;
- sent;
- delivered;
- read;
- failed + alasan + solusi.

### Step 4 — Pilih quickstart

- Kirim invoice setelah pembayaran.
- Reminder appointment.
- Lead dari form masuk ke inbox.
- Welcome message.
- FAQ chatbot.

### Step 5 — Health check

Checklist otomatis:

- channel connected;
- permission valid;
- webhook reachable;
- signature verified;
- template approved;
- test event received;
- opt-out handling active.

## Fitur MVP

### P0 — wajib untuk product-market learning

1. Multi-tenant account/workspace.
2. Official WhatsApp connection.
3. Contact + tags + consent status.
4. Shared inbox sederhana.
5. Send message dari dashboard.
6. Template management dan variable preview.
7. Delivery timeline per message.
8. Basic campaign dengan schedule dan approval.
9. Webhook dengan HMAC signature.
10. API send message.
11. API key create/revoke/rotate.
12. Usage/quota dashboard.
13. Billing plan dan invoice.
14. Audit log dasar.
15. Documentation + code examples cURL, Node.js, PHP, Python.

### P1 — pembeda kuat

- Visual automation/flow builder.
- Assignment dan round-robin untuk agent.
- SLA timer dan escalation.
- Broadcast safety score.
- Quiet hours dan per-country timezone.
- Retry policy dan dead-letter queue.
- Import CSV dengan validation preview.
- Conversation context dan AI assist dengan data boundaries jelas.
- Sandbox/test mode.
- n8n, Zapier, Make, WooCommerce, Shopify, Google Forms.

### P2 — scale

- Omnichannel.
- Advanced analytics dan conversion attribution.
- Multiple WABA/channel routing.
- Partner/reseller mode.
- SSO, SCIM, custom retention, DPA/SLA enterprise.

## API dan reliability contract

API minimal:

```http
POST /v1/messages
GET  /v1/messages/{id}
POST /v1/messages/{id}/cancel
POST /v1/contacts/import
GET  /v1/conversations
POST /v1/webhooks/test
POST /v1/api-keys/rotate
```

Message object sebaiknya memiliki:

```json
{
  "id": "msg_123",
  "status": "delivered",
  "recipient": "+628123456789",
  "provider_message_id": "wamid...",
  "request_id": "req_123",
  "created_at": "2026-09-04T10:00:00Z",
  "events": [
    {"status": "queued", "at": "..."},
    {"status": "sent", "at": "..."},
    {"status": "delivered", "at": "..."}
  ],
  "error": null
}
```

Reliability requirements:

- Idempotency key untuk request send.
- Exponential backoff dengan batas retry.
- Dead-letter queue untuk event gagal.
- Webhook signature + timestamp replay protection.
- Correlation/request ID di UI, API response, dan logs.
- Status provider dipisahkan dari status internal platform.
- Rate limit yang terlihat dan dapat dipahami.
- Status page dan incident history.

## Safety, trust, dan compliance

- Explicit opt-in dan sumber consent per contact.
- Unsubscribe/STOP otomatis.
- Suppression list global per workspace.
- Approval sebelum campaign besar dikirim.
- Preview audience dan estimasi biaya.
- Dry run/test recipients.
- Quiet hours.
- Rate guard berdasarkan reputasi, country, dan campaign type.
- PII masking di log.
- Encryption at rest/in transit.
- Role-based access: owner, admin, manager, agent, developer, billing.
- Audit trail untuk credential, campaign, contact export, dan perubahan automation.
- Retention settings yang mudah dipahami.
- Jangan menyimpan data lebih lama dari kebutuhan bisnis.

## Pricing yang lebih sehat

Gunakan pricing berbasis kombinasi **platform fee + usage transparan**, bukan quota yang membingungkan.

Contoh struktur:

- **Free/Sandbox** — test number, limit rendah, tanpa campaign produksi.
- **Starter** — 1 channel, 3 agents, basic inbox, usage-based.
- **Growth** — multiple agents, automation, campaigns, analytics.
- **Scale** — multiple channels, advanced permissions, SLA, support.

Di pricing page tampilkan:

- platform fee;
- biaya provider/Meta bila ada;
- biaya per conversation/message;
- jumlah channel dan agent;
- limit campaign;
- biaya overage;
- contoh simulasi untuk 1.000, 10.000, dan 100.000 pesan;
- kalkulator estimasi biaya.

## Design system dan content quality

- Gunakan satu bahasa utama; sediakan toggle bahasa yang konsisten.
- Hindari jargon tanpa tooltip.
- Label status: `Connected`, `Needs attention`, `Paused`, `Deprecated`.
- Semua fitur punya “apa manfaatnya”, “prasyarat”, dan “risikonya”.
- Empty state selalu memiliki CTA dan contoh data.
- Form menampilkan validation sebelum submit.
- Error harus menjawab: apa yang terjadi, penyebab, dan langkah berikutnya.
- Dokumen menampilkan last updated, API version, changelog, dan compatibility.
- Video tutorial memiliki caption, transcript, durasi, dan timestamp.

## North Star dan metrik

### Activation

- Time to first successful delivered message.
- % workspace yang menyelesaikan channel connection.
- % user yang menyelesaikan health check.

### Value

- Active workspaces per week.
- Messages delivered per active workspace.
- Inbox conversations resolved.
- Automation runs completed.
- Campaign repeat rate.

### Quality

- Delivery rate.
- Failure rate per reason.
- Webhook success/latency.
- Median first response time.
- Disconnect/incidents per channel.

### Business

- Free-to-paid conversion.
- Net revenue retention.
- Cost per delivered message.
- Support tickets per 1.000 messages.
- Churn reason by segment.

## Rencana 90 hari

### Hari 0–30: validate

- Interview 10 business owner dan 10 developer.
- Pilih satu wedge: transactional notification atau shared inbox.
- Validasi official API onboarding.
- Buat clickable prototype onboarding, send message, inbox, dan delivery timeline.
- Uji pricing dengan fake door/landing page.

### Hari 31–60: build MVP

- Multi-tenant, auth, channel connection, contacts, send, status, webhook signing, logs.
- Siapkan sandbox dan test number.
- Buat dokumentasi quickstart empat bahasa.

### Hari 61–90: pilot

- Onboard 5–10 bisnis dengan use case berbeda.
- Ukur time-to-value dan failure reasons.
- Tambahkan campaign approval, opt-out, dan dashboard usage.
- Jangan scale acquisition sebelum reliability dan compliance baseline terpenuhi.

## Keputusan strategis

**Rekomendasi:** mulai dari produk kecil dengan satu job-to-be-done yang sangat jelas, bukan seluruh daftar fitur Fonnte sekaligus. Diferensiasi paling kuat bukan “lebih banyak fitur”, tetapi:

- setup lebih cepat;
- status lebih jelas;
- risiko lebih aman;
- API lebih reliable;
- billing lebih transparan;
- support dan recovery lebih baik.
