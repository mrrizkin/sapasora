# 04 — Tutorial dan Video YouTube

## Discovery tutorial

Website memiliki knowledge base terpisah di <https://docs.fonnte.com/> dan blog tutorial di <https://fonnte.com/tutorial/>.

### Struktur knowledge base

- Getting Started: register, login, forgot password.
- Dashboard: device, phone book, history, send message, autoreply, template, recurring, profile.
- How To: connect device, kirim pesan, PHP API, chatbot, group ID, Aksita AI.
- API: token, send, QR, validate number, group list, rotator, GET send, disconnect device.
- Webhook: URL, reply message, attachment, submission, message status, device status.
- API/Webhook JavaScript: contoh JavaScript dan Node.js.
- AI: getting started, quota, data, cost, best practice, history.
- Misc: variable, file limitation, package, status, rotation, lifecycle.

### Index artikel tutorial publik

| Tema | Artikel |
|---|---|
| PHP API | [Cara mengirim pesan WhatsApp dengan PHP](https://fonnte.com/tutorial/mengirim-pesan-whatsapp-php-api/) |
| PHP webhook/chatbot | [Membuat WhatsApp bot dengan PHP](https://fonnte.com/tutorial/membuat-whatsapp-bot-dengan-php/) |
| n8n | [Integrasi WhatsApp ke n8n](https://fonnte.com/tutorial/whatsapp-n8n-fonnte/) |
| Form/order | [Membuat form order di WhatsApp](https://fonnte.com/tutorial/membuat-form-order-di-whatsapp/) |
| Google Form | [Menghubungkan Google Form](https://fonnte.com/tutorial/google-form-whatsapp/) |
| WordPress | [Daftar plugin Fonnte](https://fonnte.com/tutorial/daftar-plugin-fonnte/) |
| WooCommerce | [Plugin WooCommerce WhatsApp](https://fonnte.com/tutorial/plugin-woocommerce-whatsapp/) |
| Elementor | [Elementor Form](https://fonnte.com/tutorial/elementor-form-whatsapp/) |
| Contact Form 7 | [Autoresponder Form](https://fonnte.com/tutorial/plugin-autoresponder-form-whatsapp/) |
| OTP | [Membuat OTP dengan WhatsApp](https://fonnte.com/tutorial/membuat-otp-whatsapp/) |
| Risiko | [Informasi terkait banned](https://fonnte.com/tutorial/informasi-terkait-banned/) |

Blog menampilkan setidaknya 4 halaman arsip tutorial, dengan topik tambahan Caldera Form, Formidable Form, affiliate agency, plugin broadcast, forgot password, dan pendaftaran vaksin.

## Video demo yang tertaut dari homepage

- [Kirim pesan single](https://www.youtube.com/watch?v=tnBASuROWaE)
- [Broadcast](https://www.youtube.com/watch?v=NKBJxauketA)
- [Schedule](https://www.youtube.com/watch?v=EnJ4YOqPQ0A)
- [Button — deprecated](https://www.youtube.com/watch?v=wqL770wLOKU)
- [Variable](https://www.youtube.com/watch?v=JKjeBvUJ5qI)
- [Attachment](https://www.youtube.com/watch?v=UhLi3JCG0Ps)
- [Follow up](https://youtu.be/OranA2itGL8)
- [Random delay](https://youtu.be/2YdwSRWgzOI)
- [Random device](https://youtu.be/g7eWFeisDk8)
- [Simple API playlist](https://youtube.com/playlist?list=PL7yPMr9J3bAnZlhcOUx-NvI5UnDZO9rKA)
- [Default reply](https://www.youtube.com/watch?v=9zNYAHilh8k)
- [Keyword based reply](https://www.youtube.com/watch?v=2B1IN1YzUIw)
- [Nama pengirim](https://www.youtube.com/watch?v=-7S7SGNRbkM)
- [Submission](https://youtu.be/rys6sS0DdQY)
- [Quick reply](https://youtu.be/3mbo-3lfQGA)
- [Download file](https://youtu.be/C4z4-tdSVxQ)
- [Webhook](https://www.youtube.com/watch?v=huglxjOeu3s)
- [Management device](https://www.youtube.com/watch?v=LH776-NCAio)
- [Simpan/grouping kontak](https://www.youtube.com/watch?v=4qIILzYKHwU)
- [Riwayat pesan](https://www.youtube.com/watch?v=JsXEu5D9XoM)
- [Kirim dari dashboard](https://www.youtube.com/watch?v=KyWnmI3XSaQ)
- [Template pesan](https://www.youtube.com/watch?v=ofTd_-nJQSE)
- [Pesan berulang](https://www.youtube.com/watch?v=kRa9oI8akdM)
- [Autoreply](https://www.youtube.com/watch?v=HS5N6JBGv60)
- [Invoice](https://www.youtube.com/watch?v=0v117Bv3_Rk)

## Video yang dibuka dan ditinjau

### 1. Sending WhatsApp Messages Using PHP API

- URL: <https://www.youtube.com/watch?v=A9klOsh1Uj0>
- Channel: Fonnte
- Durasi yang tampil di YouTube: sekitar 1:17.
- Snapshot saat discovery: sekitar 29K views, 3 tahun lalu, label **Auto-dubbed**.
- Video menunjukkan alur: buka dokumentasi → buat file PHP (contoh `api.php`) → copy kode → isi `target` dan `message` → jalankan dari browser → cek pesan masuk di WhatsApp.

**Transkrip hasil caption Indonesian (dirapikan, bukan transkripsi manual audio):**

> Kali ini akan menjelaskan cara mengirim pesan ke satu nomor menggunakan API Fonnte. Selalu buka dokumentasi untuk mempermudah penggunaan API. Perlu membuat file contoh `api.php` dan copy kode yang ada di dokumentasi. Semua `postfield` kecuali `target` dan `message` [dapat mengikuti contoh]. Message dengan tes dari PHP, kemudian buka di browser agar bisa tereksekusi. Kita cek WhatsApp—ada pesan masuk, berarti pengiriman berhasil. Pada video berikutnya kita akan mempelajari bagaimana mengirimkan pesan broadcast ke banyak nomor menggunakan API Fonnte.

### 2. How to send WhatsApp blasts/broadcasts using PHP

- URL: <https://www.youtube.com/watch?v=2YdwSRWgzOI>
- Channel: Fonnte
- Durasi yang tampil di YouTube: sekitar 1:08.
- Snapshot saat discovery: sekitar 12K views, 3 tahun lalu, label **Auto-dubbed**.
- Video menunjukkan pemisahan nomor dengan koma dan penggunaan delay 5–10 detik.

**Transkrip hasil caption Indonesian (dirapikan):**

> Kali ini akan menjelaskan cara mengirim pesan ke banyak nomor menggunakan API Fonnte. Selalu buka dokumentasi untuk mempermudah penggunaan API. Untuk mengirim pesan ke beberapa nomor, gunakan koma untuk memisahkan antar-nomor. Saya akan mengirim ke nomor saya sebanyak tiga kali untuk mempermudah praktik. Kita akan menambahkan delay untuk jeda pengiriman antar-pesan, sekitar 5 sampai 10 detik. Message dengan tes dari PHP, kemudian buka file `api.php` di browser agar bisa tereksekusi. Jika pesan masuk, berarti pengiriman berhasil. Pada video berikutnya kita akan mempelajari cara mengirimkan pesan dengan variabel menggunakan API Fonnte.

## Catatan kualitas tutorial

**Yang bagus:** dokumentasi dan video mengikuti mental model implementasi nyata; contoh PHP copy-paste; halaman API memiliki daftar parameter dan contoh response/error; tutorial menjelaskan prasyarat secara eksplisit.

**Yang bisa di-improve:**

1. Sediakan caption/transcript resmi yang mudah diakses di embedded video; saat discovery, YouTube menampilkan “Subtitles/closed captions unavailable” pada player dan panel transcript tidak memuat segment, walau caption otomatis dapat diambil dari track YouTube.
2. Tambahkan timestamp ke video yang tertaut dari setiap kartu fitur.
3. Sediakan quickstart non-PHP (Node.js, Python, cURL) di jalur onboarding utama.
4. Tandai artikel lama vs dokumentasi terkini untuk mencegah contoh API yang out-of-date.
5. Tambahkan error recovery, keamanan token, webhook authentication/signature, retry/idempotency, dan observability ke tutorial production.

## Metode discovery video

Video dibuka melalui `playwright-cli` pada halaman YouTube untuk melihat title, durasi, channel, views, label auto-dubbed, dan player state. Caption Indonesian kemudian diambil untuk kebutuhan transkripsi ringkas. Hasil di atas diperlakukan sebagai ringkasan caption otomatis sehingga mungkin mengandung salah dengar/typo.
