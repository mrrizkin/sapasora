# 06 — Catatan Discovery

## Scope dan metode

- Target: <https://fonnte.com/>
- Metode browser: `playwright-cli` global.
- Halaman yang dibuka: homepage, blog tutorial, knowledge base, tutorial PHP API, tutorial PHP webhook, YouTube demo, register, support, terms, disclaimer, Ask Fai.
- Caption video: track caption Indonesian diambil setelah video ditinjau melalui browser; hasil dirapikan untuk ringkasan.
- Tidak melakukan register, login, pengiriman pesan, scan QR, pembelian, atau submit ticket.

## Urutan observasi

1. **Homepage** — menemukan positioning, CTA, fungsi, feature map, integrasi, testimonials, angka social proof, pricing, FAQ, dan disclaimer footer.
2. **Knowledge base** — menemukan struktur kategori Getting Started, Dashboard, How To, API, Device API, Webhook, AI, dan Misc.
3. **First-time guide** — mengonfirmasi alur add device → connect → send → optional autoreply/chatbot.
4. **Dashboard docs** — memeriksa device management, phone book, send message, autoreply, recurring.
5. **API docs** — memeriksa endpoint send, token, parameter, response, error, schedule/delay/sequence/rotator.
6. **Webhook docs/tutorial** — memeriksa payload, autoread, public URL, reply flow, dan attachment.
7. **AI docs** — memeriksa flow, custom data, history, dan quota/cost model.
8. **Pricing toggle** — menguji IDR/USD dan monthly/annual melalui click/evaluate pada DOM.
9. **Tutorial archive** — memeriksa 4 halaman arsip dan kategori integrasi.
10. **YouTube** — membuka video API dan broadcast; memeriksa title, duration, channel, views, auto-dubbed, player captions/transcript state.
11. **Register/support** — memeriksa field signup dan support ticket tanpa mengirim data.
12. **Terms/disclaimer** — membaca dependency WhatsApp Web, retention, support boundaries, no-SLA/no-refund, dan risk disclosure.
13. **Ask Fai** — membuka contextual help dari docs; diarahkan ke Aksita.ai dan menjawab prompt “Apa itu fonnte?”.

## Temuan penting yang perlu diverifikasi lebih lanjut

- Angka `2000++ devices`, `100K++ messages`, dan `<100ms` adalah klaim homepage.
- Homepage FAQ menyebut batas kecepatan 10 pesan/detik (hingga 600/menit) selama tidak terkena banned.
- File limitation mendokumentasikan maksimum 10 MB, sementara contoh error API menyebut di bawah 4 MB.
- Tutorial dan dokumentasi menyatakan POC; belum production-ready.
- Token adalah credential yang dapat mengirim pesan atas nama device; kebocoran token berdampak langsung.
- Webhook harus public dan autoread on. Status signature/authentication webhook perlu dikonfirmasi dari dokumentasi terbaru.
- Default AI tidak mengingat history; history dibuat dengan flow loop dan menambah cost.
- Paket dimiliki per-device, tidak share antar-device; upgrade/downgrade dapat terminate paket lama.
- YouTube player menampilkan captions unavailable/transcript panel tanpa segments saat ditinjau, tetapi automatic Indonesian caption track tersedia dan dipakai untuk ringkasan.

## Artefak browser

`playwright-cli` menyimpan snapshot dan console log sementara di `.playwright-cli/`. Artefak ini adalah jejak sesi discovery, bukan bagian dari kontrak API.

## Batasan discovery

- Tidak menguji dashboard authenticated.
- Tidak menguji deliverability, latency, reconnect, rate limit, banned behavior, webhook retry, atau deletion lifecycle.
- Tidak memverifikasi angka pricing terhadap invoice/payment.
- Tidak memverifikasi legal status, official WhatsApp policy, atau klaim customer testimonial secara independen.
- Tidak menganggap caption otomatis sebagai transkripsi audio yang sempurna.
