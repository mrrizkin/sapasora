# 02 — Fitur dan Use Case

## Capability map

### 1. Pengiriman pesan

- Personal: kirim ke satu nomor.
- Broadcast: kirim ke banyak nomor atau grup.
- Schedule: kirim pada waktu tertentu.
- Recurring: kirim berulang berdasarkan hari/minggu/bulan/tahun atau interval custom.
- Variable: personalisasi pesan, misalnya `{name}`, `{var1}`, dan seterusnya.
- Attachment: image, video, audio, dan file pada paket yang mendukung media.
- Location: kirim koordinat latitude/longitude.
- Poll: pilihan 2–12 item, dengan mode single atau multiple selection.
- Typing indicator dan custom duration.
- Random delay dan sequential sending.
- Preview link dapat diaktifkan/nonaktifkan.
- Rotator: beberapa token/device dalam satu akun untuk distribusi pengiriman.

### 2. Incoming message dan chatbot

- **Autoreply** berbasis keyword.
- Default reply untuk pesan yang tidak cocok dengan rule.
- Reply dengan attachment pada paket tertentu.
- **Webhook** untuk chatbot dinamis menggunakan data dari aplikasi sendiri.
- Webhook menyediakan metadata seperti device, sender, message, nama, anggota grup, lokasi, poll, timestamp, attachment, dan `inboxid`.
- Flow + AI untuk chatbot dengan data kustom; penggunaan AI mengurangi kuota AI.
- Default AI tidak memakai chat history; history dapat dibuat dengan flow loop dan menambah konsumsi quota.
- “Ask Fai” pada docs membuka contextual assistant Aksita.ai untuk pertanyaan tentang Fonnte, harga, fitur, dan testimonial.
- Fonnte menyebut autoreply tidak berjalan ketika webhook digunakan—pengguna perlu memilih pola kontrol yang diinginkan.

### 3. Dashboard dan operasional

- Management banyak device dalam satu akun.
- Connect dengan QR WhatsApp Web, reconnect, reset, edit, dan token.
- Pengaturan nama device, webhook, autoread.
- Phone book: simpan nomor, variable, dan group kontak.
- Message history dan status pengiriman.
- Template pesan.
- Recurring template.
- Inbox dan CS Template/multi-agent (terlihat di dokumentasi terkait).
- Invoice/order history.
- Device notification.

### 4. Integrasi dan developer surface

Integrasi/tutorial publik yang terlihat:

- PHP API.
- JavaScript API.
- Webhook PHP/Node.js.
- WordPress.
- Google Form.
- Contact Form 7.
- WooCommerce.
- Elementor Form.
- Aksita AI.
- Postman collection.

## Alur utama yang ditemukan

### Quick start pengguna baru

1. Register akun.
2. Tambahkan device.
3. Hubungkan device melalui QR WhatsApp Web.
4. Kirim pesan dari menu Send.
5. Jika membuat chatbot, aktifkan autoread, buat autoreply/flow, lalu tes keyword.

Dokumentasi menyebut device gratis memperoleh kuota 1.000 pesan/bulan; kuota berkurang berdasarkan jumlah target, bukan sekadar jumlah request.

### Alur developer

1. Tambah dan connect device.
2. Ambil token dari dashboard.
3. POST ke `https://api.fonnte.com/send` dengan header `Authorization: TOKEN`.
4. Gunakan `target`, `message`, dan parameter opsional seperti `url`, `file`, `schedule`, `delay`, `location`, `typing`, `data`, atau `sequence`.
5. Proses hasil berdasarkan `status`, `process`, `requestid`, dan detail error.

### Alur chatbot dinamis

1. Buat endpoint webhook.
2. Isi URL webhook pada Device → Edit.
3. Aktifkan autoread.
4. Terima payload incoming message.
5. Ambil `sender`/`inboxid` dan data lain.
6. Kirim balasan melalui endpoint send.

## Batasan dan risiko yang perlu dipahami

- Layanan unofficial, bukan WhatsApp Official API.
- Automasi berbasis WhatsApp Web membuat status device/maintenance menjadi dependency penting.
- Website menyatakan risiko banned tetap ada, terutama broadcast ke nomor yang belum pernah berinteraksi.
- Pengguna disarankan memakai nomor secondary.
- Tidak ada SLA fixed.
- Tidak ada refund.
- Fonnte menyatakan tidak bertanggung jawab atas banned WhatsApp.
- Attachment: format yang didokumentasikan `png/jpg/jpeg/webp`, `mp4`, `pdf/doc/docx/xls/xlsx/csv/txt`, `mp3`; batas maksimum dokumentasi file adalah 10 MB. Namun contoh error API juga menyebut file di bawah 4 MB—ini inkonsistensi yang perlu diklarifikasi sebelum dijadikan kontrak produk.
- Fitur attachment API URL/file hanya tersedia pada paket tertentu.

## Sumber

- <https://fonnte.com/>
- <https://docs.fonnte.com/its-my-first-time-using-fonnte-what-to-do/>
- <https://docs.fonnte.com/send-message/>
- <https://docs.fonnte.com/api-send-message/>
- <https://docs.fonnte.com/device/>
- <https://docs.fonnte.com/phone-book/>
- <https://docs.fonnte.com/autoreply/>
- <https://docs.fonnte.com/webhook-reply-message/>
- <https://docs.fonnte.com/ai-getting-started/>
- <https://docs.fonnte.com/file-limitation/>
