# 02 — Discovery Findings: Fonnte

## Apa yang ditemukan

Fonnte adalah SaaS/API gateway WhatsApp Indonesia dengan positioning **“Fonnte WHATSAPP API”** dan janji **“Kirim pesan WhatsApp secara otomatis”**. Fonnte menyatakan layanannya unofficial dan menggunakan WhatsApp Web untuk otomatisasi melalui API/webhook.

Audience yang terlihat: developer, UMKM/business owner, digital marketer, WordPress/form user, dan organisasi yang membutuhkan automation sederhana.

## Use case

- Notifikasi aktivitas.
- Chatbot/autoreply.
- Reminder.
- Tagihan.
- Informasi.
- OTP.
- Personal send dan group send.
- Broadcast.
- Schedule dan recurring.
- Variable/personalization.
- Attachment.
- CSV.
- Follow-up.
- Lokasi.
- Random delay/random device.
- Simple API.

## Feature inventory

### Messaging

- Personal dan WhatsApp Group.
- Single, broadcast, schedule.
- Variable.
- Attachment.
- CSV.
- Follow-up.
- Location.
- Random delay.
- Random device/rotator.
- Simple API.
- Poll.
- Typing indicator/custom duration.
- Link preview.
- Sequential sending.
- `data` parameter untuk menggabungkan request.

### Reply/chatbot

- Default reply.
- Keyword-based reply.
- Sender name.
- Submission/form flow.
- Quick Reply.
- Download file.
- Webhook dynamic reply.
- AI/flow dengan custom data dan optional history.

### Dashboard

- Management device.
- QR connect/reconnect/reset/edit/token.
- Phone book, variable, grouping.
- Message history.
- Send berbagai tipe message.
- Template.
- Button (deprecated).
- Recurring.
- Autoreply.
- Status pengiriman real-time.
- Invoice/order history.
- Device notification.
- Inbox dan CS Template/multi-agent terlihat di knowledge base.

### Integrasi

- PHP API.
- JavaScript/Node.js.
- PHP/Node.js webhook.
- WordPress.
- Google Form.
- Contact Form 7.
- WooCommerce.
- Elementor Form.
- Caldera/Formidable Form.
- n8n.
- Aksita AI.
- Postman.

## Quick-start Fonnte yang ditemukan

1. Register.
2. Add device.
3. Connect melalui QR WhatsApp Web.
4. Copy token.
5. Send message.
6. Untuk chatbot: buat autoreply/flow, aktifkan autoread, dan tes keyword.

Dokumentasi menyebut free device mendapat 1.000 pesan/bulan dan quota berkurang berdasarkan jumlah target.

## Pricing yang diamati

### Text Only — bulanan IDR

| Paket | Kuota | Multics | Harga |
|---|---:|---:|---:|
| Free | 1.000/bulan | 0 | Rp0 |
| Lite | 1.000/bulan | 0 | Rp25.000 |
| Regular | 10.000/bulan | 2 | Rp66.000 |
| Regular Pro | 25.000/bulan | 2 | Rp110.000 |
| Master | Unlimited | 4 | Rp175.000 |

### All Feature — bulanan IDR

| Paket | Kuota | Multics | Harga |
|---|---:|---:|---:|
| Super | 10.000/bulan | 2 | Rp165.000 |
| Advanced | 25.000/bulan | 2 | Rp255.000 |
| Ultra | Unlimited | 4 | Rp355.000 |

### Tahunan IDR

- Lite Rp250.000.
- Regular Rp660.000.
- Regular Pro Rp1.100.000.
- Master Rp1.750.000.
- Super Rp1.650.000.
- Advanced Rp2.550.000.
- Ultra Rp3.550.000.
- Free Rp0.

Harga tahunan yang diamati setara sekitar 16,7% lebih murah dibanding 12 bulan harga bulanan. USD yang tampil: bulanan Lite `$2.5`, Regular `$6`, Regular Pro `$10`, Master `$14`, Super `$13`, Advanced `$22`, Ultra `$29`; tahunan Lite `$25`, Regular `$60`, Regular Pro `$100`, Master `$140`, Super `$130`, Advanced `$220`, Ultra `$290`.

Daftar label feature pada pricing: personal/group, text, schedule, recurring, template, button (deprecated), attachment, autoreply, autoreply spreadsheet, webhook, API, remove watermark, device notification, Multics Agent.

## Trust, limitations, dan legal positioning

- Bukan official WhatsApp API.
- Nomor pengirim adalah nomor user; secondary number disarankan.
- Risiko banned tetap ada, terutama broadcast ke nomor tanpa interaksi.
- Kecepatan yang disebut FAQ: 10 pesan/detik, sampai 600/menit selama tidak terkena banned.
- Tidak ada fixed SLA.
- Tidak ada refund.
- Fonnte tidak bertanggung jawab atas banned WhatsApp.
- Tutorial/docs dinyatakan sebagai POC, bukan production-ready.
- Dependensi pada WhatsApp Web berarti perubahan/penutupan WhatsApp Web dapat menghentikan layanan.

## Data/lifecycle yang didokumentasikan Fonnte

- Device: lifetime.
- WhatsApp session: lifetime selama tidak terputus.
- Temporary files: 7 hari.
- Webhook files: 30 menit.
- Scheduled files: 7 hari.
- Account: 30 hari setelah tidak memiliki device, dengan login mereset timer.
- Fonnte menyatakan data tidak digunakan untuk kepentingan personal/Fonnte dan message/contact dapat dikelola/dihapus user sesuai fitur yang tersedia.

## Discovery terhadap docs/tutorial/video

Knowledge base memiliki kategori Getting Started, Dashboard, How To, API, Device API, Webhook, Webhook JS, AI, dan Misc. Blog memiliki setidaknya 4 halaman arsip dengan topik PHP, webhook, n8n, WordPress, form, OTP, banned, affiliate, dan integrasi.

Video yang ditinjau melalui `playwright-cli`:

1. [Sending WhatsApp Messages Using PHP API](https://www.youtube.com/watch?v=A9klOsh1Uj0) — sekitar 1:17, Fonnte, snapshot sekitar 29K views, Auto-dubbed.
2. [How to send WhatsApp blasts/broadcasts using PHP](https://www.youtube.com/watch?v=2YdwSRWgzOI) — sekitar 1:08, Fonnte, snapshot sekitar 12K views, Auto-dubbed.

Player menampilkan captions unavailable/transcript panel tanpa segments, tetapi automatic Indonesian caption track dapat diambil. Ringkasan caption:

> Video pertama: buka dokumentasi, buat `api.php`, copy kode, isi `target`/`message`, jalankan dari browser, lalu cek pesan masuk.

> Video kedua: pisahkan banyak nomor dengan koma, tambahkan delay 5–10 detik, jalankan file PHP, dan cek hasil pengiriman.

## Pelajaran product design

Yang bagus:

- CTA free tier.
- Happy path sederhana.
- Use case konkret.
- API/docs actionable.
- Banyak tutorial/integrasi.
- Transparansi risiko.
- Contextual help Ask Fai melalui Aksita.ai.

Yang harus kita improve:

- Persona-based onboarding.
- Homepage dan pricing hierarchy.
- Official/unofficial disclosure.
- Delivery timeline dan recovery.
- Webhook authentication/signature.
- Token rotation/scoping.
- Consent/opt-out/rate guard.
- Docs versioning, captions/timestamp, dan consistency.
- Attachment limit inconsistency: satu docs menyebut 10 MB, contoh API error menyebut <4 MB.

## Source URLs

- <https://fonnte.com/>
- <https://docs.fonnte.com/>
- <https://fonnte.com/tutorial/>
- <https://docs.fonnte.com/api-send-message/>
- <https://docs.fonnte.com/its-my-first-time-using-fonnte-what-to-do/>
- <https://docs.fonnte.com/terms-conditions/> (melalui artikel Terms & Conditions Fonnte)
- <https://support.fonnte.com/>
