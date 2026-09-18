# 01 — Product Brief

## Executive summary

Kita akan membangun SaaS untuk otomasi dan operasional komunikasi customer di WhatsApp. Produk harus lebih baik dari gateway yang hanya menyediakan endpoint send: user perlu tahu cara connect, apakah pesan benar-benar terkirim, kenapa gagal, siapa yang menangani conversation, berapa biaya, dan bagaimana memulihkan masalah.

## Vision

Membantu bisnis Indonesia mengelola komunikasi pelanggan di WhatsApp dengan cara yang mudah, aman, dapat diukur, dan dapat diintegrasikan.

## Mission

Menyederhanakan proses dari menghubungkan WhatsApp sampai pesan berhasil diterima dan ditindaklanjuti, tanpa menyembunyikan risiko, biaya, maupun status teknis.

## Problem

- Owner nonteknis sulit melakukan setup dan memahami error.
- Developer membutuhkan API/webhook production-grade.
- Tim support membutuhkan shared inbox, assignment, notes, dan audit.
- Campaign tanpa consent dapat merusak reputasi nomor.
- Pricing, quota, dan provider fee sering membingungkan.
- Disconnect, template rejected, token error, webhook failure, dan delivery failure membutuhkan recovery flow.

## Target awal

### Primary

SMB dengan 1–10 agent: ecommerce, jasa, edukasi, klinik, appointment, property, dan komunitas.

### Secondary

Developer dan agency yang membangun integrasi untuk banyak client.

### Deferred

Enterprise dengan kebutuhan SLA, DPA, SSO/SCIM, procurement, dan dedicated support.

## Primary wedge

**Shared inbox + transactional notification** untuk official channel.

Alasan:

- value mudah diukur;
- frequency penggunaan tinggi;
- lebih aman daripada memulai dari unsolicited broadcast;
- mendorong retention melalui inbox dan automation;
- cocok untuk production dan developer integration.

## Positioning

> Untuk bisnis kecil-menengah dan developer yang membutuhkan customer communication di WhatsApp, produk ini adalah platform automation dan inbox dengan onboarding sederhana, delivery visibility, dan API production-grade. Berbeda dari gateway yang hanya fokus pada kirim pesan, produk ini mengelola consent, send, delivery, reply, automation, billing, dan recovery dalam satu pengalaman.

## Product pillars

1. Fast activation.
2. Operational clarity.
3. Safe messaging.
4. Developer trust.
5. Team productivity.
6. Transparent economics.

## Official dan unofficial channel

### Official Cloud API

Production/default: Meta/WABA onboarding, template/policy, provider billing, dan capability sesuai official platform.

### Unofficial WhatsApp Web/whatsmeow

Experimental/advanced: QR/pairing, personal/group capability bila tersedia, tetapi risiko logout/ban, protocol change, maintenance, dan no official SLA. Tidak boleh disamarkan sebagai official.

## Non-goals

- Menjamin nomor tidak pernah banned.
- Menawarkan bulk spam.
- Menganggap AI sebagai pengganti agent tanpa guardrail.
- Menggantikan CRM/ERP sepenuhnya.
- Membangun seluruh feature list kompetitor sebelum core lifecycle stabil.

## Success definition

Produk berhasil jika user dapat:

1. memahami channel yang dipilih;
2. menyelesaikan setup;
3. mengirim test message;
4. melihat status sampai delivered/failed;
5. menerima dan menjawab conversation;
6. memperbaiki failure tanpa selalu menghubungi support;
7. mengintegrasikan use case yang sama melalui API/webhook.
