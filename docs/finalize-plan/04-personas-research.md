# 04 — Personas dan Research Plan

## Persona utama

### Business owner / ops

Owner/ops SMB nonteknis dengan 1–5 channel/nomor.

**Jobs:** konfirmasi order/appointment, shared reply, campaign terkontrol, memahami failure.

**Pain:** setup Meta, template rejected, nomor disconnect, status tidak jelas, takut spam/ban.

**Success:** first value <15 menit dan issue umum dapat diselesaikan sendiri.

### Developer

Developer Go/Node/PHP atau agency implementer.

**Jobs:** kirim event bisnis, menerima webhook, replay/debug, integrasi idempotent.

**Pain:** docs drift, token manual, error generik, retry tidak jelas, webhook insecure.

**Success:** POC <1 jam dan production checklist jelas.

### Support lead

Mengelola 3–20 agent.

**Jobs:** assignment, notes, canned reply, SLA, performance.

**Pain:** double reply, percakapan hilang, tidak ada ownership/audit.

### Agency

Mengimplementasikan solusi untuk banyak client.

**Jobs:** multi-workspace, credential isolation, deployment monitoring, billing/support per client.

## JTBD prioritas

1. Connect official channel.
2. Send test message dan melihat delivered/failed.
3. Receive inbound message.
4. Reply dan assign conversation.
5. Create template/automation sederhana.
6. Integrate API/webhook.
7. Diagnose failed delivery.

## Research questions

### Market

- Use case harian apa yang paling bernilai?
- Apakah calon user sudah memiliki WABA/Meta Business?
- Official, BSP, atau personal QR mana yang mereka pilih dan mengapa?
- Biaya/risk apa yang dapat diterima?

### UX

- Bagian onboarding mana yang membingungkan?
- Status apa yang dibutuhkan owner vs developer?
- Apakah consent/campaign approval dipahami?
- Apakah official/unofficial mudah dibedakan?

### Technical

- Stack customer saat ini?
- SDK apa yang dibutuhkan?
- Payload/event wajib?
- Requirement latency, throughput, retention, retry?

## Research sample

- 5 business owner.
- 5 developer.
- 3 support lead.
- 2 agency.

## Interview guide

1. Ceritakan workflow WhatsApp sekarang.
2. Apa yang terjadi saat pesan/nomor gagal?
3. Siapa yang mengelola nomor dan credential?
4. Definisi pesan sukses?
5. Pernah mengalami report/ban/disconnect?
6. Fitur yang paling sering dipakai/tidak dipakai?
7. Biaya operasional saat ini?
8. Apa yang harus tersedia agar mau pindah?
9. Siapa yang approve campaign/template?
10. Apa yang tidak boleh terjadi?

## Prototype test

Prototype:

- connect channel wizard;
- send test message;
- inbox;
- failed message recovery;
- pricing calculator.

Measure:

- task completion;
- time on task;
- error count;
- confidence;
- official/unofficial understanding;
- support-assisted step.

## Pilot

5–10 workspace dengan use case invoice, appointment, lead follow-up, dan support. Gunakan synthetic/redacted data; jangan meminta token atau melakukan broadcast ke recipient tanpa consent.
