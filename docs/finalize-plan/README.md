# Finalize Plan — WhatsApp Engagement SaaS

Dokumen handoff untuk tim produk, design, engineering, security, support, dan business.

## Apa yang dibangun

SaaS **WhatsApp Engagement Platform** yang membantu bisnis mengelola customer communication melalui:

- shared inbox;
- transactional notification;
- templates;
- campaigns yang aman;
- automation;
- analytics;
- public API dan webhook;
- official WhatsApp Cloud API sebagai jalur production/default;
- optional unofficial WhatsApp Web connector berbasis `whatsmeow` untuk experimental/advanced use case.

Produk ini bukan clone Fonnte. Fokus diferensiasi adalah onboarding, delivery visibility, recovery, safety, developer experience, dan billing yang lebih baik.

## Navigation

### Product dan discovery

- [01 — Product brief](01-product-brief.md)
- [02 — Discovery findings](02-discovery-findings.md)
- [03 — Branding dan positioning](03-branding-positioning.md)
- [04 — Personas dan research plan](04-personas-research.md)

### Scope dan experience

- [05 — Product scope dan roadmap](05-product-scope-roadmap.md)
- [06 — UX dan information architecture](06-ux-information-architecture.md)
- [07 — Backlog dan acceptance criteria](07-backlog-acceptance-criteria.md)

### Technical handoff

- [08 — VSA system architecture](08-vsa-system-architecture.md)
- [09 — GoFiber, Fibertia, Vue, dan engineering conventions](09-technical-stack.md)
- [10 — Domain model, data, dan ID contract](10-domain-data-id-contract.md)
- [11 — API, webhook, dan provider contract](11-api-webhook-provider-contract.md)
- [12 — Security, compliance, dan reliability](12-security-compliance-reliability.md)
- [13 — Infrastructure dan observability](13-infrastructure-observability.md)

### Business dan delivery

- [14 — Pricing, billing, dan unit economics](14-pricing-billing-unit-economics.md)
- [15 — Delivery, testing, dan QA](15-delivery-testing-qa.md)
- [16 — Metrics dan launch plan](16-metrics-launch-plan.md)
- [17 — Risk register, decisions, dan open questions](17-risk-decisions-open-questions.md)
- [18 — Team handoff checklist](18-team-handoff-checklist.md)

## Non-negotiable decisions

1. **VSA adalah architecture utama.** Ports/adapters hanya dipakai di dependency boundary.
2. Backend menggunakan **Go + GoFiber**.
3. Dashboard menggunakan **Inertia + Vue 3 + TypeScript** melalui library internal **Fibertia**.
4. UI component menggunakan **Nuxt UI**.
5. Tidak menggunakan `vue-router` dan Pinia sebagai default.
6. Frontend menggunakan **Zod** untuk runtime validation request, response, Inertia props, shared props, storage, dan external payload.
7. Database menggunakan PostgreSQL + **sqlx**, dengan prepared statements pada query hot path/repository construction bila terukur bermanfaat.
8. Migration menggunakan **goose**.
9. Logging menggunakan **zerolog**.
10. Cache/queue menggunakan **Valkey**.
11. API schema/documentation menggunakan **OpenAPI + Scalar**.
12. Local container workflow menggunakan **Podman**.
13. ID internal adalah `BIGSERIAL`; ID external adalah `id: string` dari `public_id` NanoID 21 karakter.
14. Internal `id` dan field `public_id` tidak pernah keluar dari backend.
15. Official channel adalah default production; `whatsmeow` adalah optional experimental dan tidak boleh diposisikan sebagai official.

## Roadmap ringkas

```text
Validate → Foundation → Internal Alpha → Private Pilot → Public Beta → GA → Scale
```

Detail gate, scope, dan exit criteria ada di [05 — Product scope dan roadmap](05-product-scope-roadmap.md).

## Cara membaca dokumen ini

- Product/design/support mulai dari 01–07.
- Engineering mulai dari 08–13.
- Business/operations mulai dari 14–18.
- Semua anggota team wajib membaca 10 — ID contract dan 12 — security sebelum implementasi.
- Checklist handoff di 18 harus dipakai saat kickoff dan release review.

## Source completeness

Folder `appendix/source-discovery/` berisi salinan lengkap seluruh dokumen discovery Fonnte yang dibuat sebelumnya.

Folder `appendix/source-planning/` berisi salinan lengkap seluruh dokumen planning sebelumnya.

Kedua folder tersebut dipertahankan sebagai audit trail agar tidak ada informasi discovery/planning yang hilang saat dokumen dirapikan.
