# 01 — Product Vision dan Strategy

## Working name

`[TBD]` — gunakan nama sementara sampai positioning tervalidasi.

## Vision

Membantu bisnis Indonesia mengelola komunikasi pelanggan di WhatsApp dengan cara yang **mudah, aman, dapat diukur, dan dapat diintegrasikan**.

## Mission

Menyederhanakan proses dari “menghubungkan WhatsApp” sampai “pesan berhasil diterima dan ditindaklanjuti”, tanpa menyembunyikan risiko, biaya, maupun status teknis.

## Positioning statement

Untuk bisnis kecil-menengah dan developer yang membutuhkan customer communication di WhatsApp, produk ini adalah platform automation dan inbox yang memberikan onboarding sederhana, delivery visibility, dan API production-grade. Berbeda dari gateway yang hanya fokus pada kirim pesan, produk ini membuat keseluruhan lifecycle—consent, send, delivery, reply, automation, billing, dan recovery—mudah dipahami.

## Strategi channel

### Official channel — default

Menggunakan WhatsApp Business Platform/Cloud API atau partner resmi. Diposisikan sebagai jalur production, dengan onboarding WABA, template/policy, dan billing yang transparan.

### Unofficial channel — optional

Connector berbasis WhatsApp Web/`whatsmeow`, jika secara legal dan operasional layak. Diposisikan sebagai experimental/advanced, menggunakan risk acknowledgement, capability matrix, kill switch, dan tanpa SLA official.

### Prinsip isolasi

Provider-specific behavior tidak boleh bocor ke core domain. Core platform berbicara dengan interface channel; setiap channel memiliki worker, status mapping, limit, error code, dan policy sendiri.

## Problem statement

1. Business owner sulit memahami setup dan status pengiriman.
2. Developer membutuhkan API/webhook yang dapat dipercaya.
3. Tim support membutuhkan inbox dan ownership percakapan.
4. Campaign tanpa consent dapat merusak reputasi nomor.
5. Pricing dan limit gateway sering membingungkan.
6. Disconnect, failed delivery, token, dan webhook error membutuhkan recovery flow yang jelas.

## Segmen awal

### Primary: SMB dengan 1–10 agent

- Ecommerce, jasa, edukasi, klinik, appointment, property, dan komunitas.
- Ingin inbox bersama, notifikasi, follow-up, dan automation.
- Tidak ingin membangun integrasi WhatsApp dari nol.

### Secondary: developer/agency

- Membangun aplikasi untuk banyak client.
- Membutuhkan API, webhook, logs, test environment, dan multi-workspace.

### Deferred: enterprise

- Membutuhkan SLA, DPA, SSO, audit, procurement, dan support khusus.
- Tidak menjadi target MVP sebelum reliability baseline terbukti.

## Technical/product signature

**Fibertia**—integrasi GoFiber + Inertia + Vue milik tim—dapat menjadi signature internal yang membuat dashboard server-driven, cepat, dan konsisten dengan developer experience Anda. Ini bukan alasan untuk memaksa semua public API menjadi Inertia; public developer API tetap REST/OpenAPI, sementara dashboard memakai Fibertia.

## Product pillars

1. **Fast activation** — first successful delivered message secepat mungkin.
2. **Operational clarity** — status dan error dapat ditindaklanjuti.
3. **Safe messaging** — consent, opt-out, quiet hours, approval, rate guard.
4. **Developer trust** — API versioning, idempotency, webhook signature, retry, request ID.
5. **Team productivity** — inbox, assignment, notes, template, automation.
6. **Transparent economics** — provider fee, platform fee, usage, dan overage dijelaskan.

## Differentiation

Bukan mengejar jumlah fitur terbanyak. Diferensiasi yang diprioritaskan:

- onboarding paling mudah;
- delivery timeline paling jelas;
- recovery dan debugging terbaik;
- safety control sebagai default;
- satu UX untuk business owner dan developer;
- official/unofficial provider dibedakan secara jujur.

## Non-goals

- Menjamin nomor tidak pernah banned.
- Menyembunyikan penggunaan unofficial channel.
- Menawarkan bulk spam/unsolicited messaging.
- Menggantikan CRM/ERP secara penuh pada MVP.
- Membangun AI agent sebelum data, consent, dan handoff dasar matang.

## Assumptions yang harus divalidasi

- SMB bersedia melalui onboarding official jika value dan instruksi jelas.
- Delivery status dan recovery lebih penting daripada jumlah fitur.
- User bersedia membayar untuk inbox/team/automation, bukan hanya endpoint send.
- Agency membutuhkan multi-workspace dan credential isolation.
- Unofficial connector memiliki demand, tetapi acceptable risk berbeda per segmen.

## Strategic bets

| Bet | Hypothesis | Validasi |
|---|---|---|
| Official-first | Trust meningkatkan conversion bisnis serius | Interview + pilot |
| Inbox + notification | Kombinasi ini lebih sticky daripada broadcast saja | Prototype + usage |
| Safety-by-default | Guardrail mengurangi failure/ticket tanpa menurunkan value | Pilot campaign |
| Developer center | API quality mempercepat integration dan retention | Developer pilot |
| Dual provider | Satu core platform dapat melayani dua segment | Architecture spike |
