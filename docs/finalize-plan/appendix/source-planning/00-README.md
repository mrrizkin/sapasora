# Planning — SaaS WhatsApp Automation

## Tujuan folder

Folder ini berisi rencana produk, UX, bisnis, arsitektur teknis, delivery, dan risiko untuk membangun SaaS WhatsApp automation versi improved.

## Arah produk

Produk yang direncanakan adalah **WhatsApp Engagement Platform** untuk bisnis kecil-menengah dan developer, dengan:

- onboarding yang mudah;
- official WhatsApp Cloud API sebagai jalur production/default;
- optional unofficial connector berbasis `whatsmeow` untuk use case experimental/personal/group;
- shared inbox, campaign, automation, template, API, webhook, dan analytics;
- safety, consent, observability, dan recovery sebagai fitur inti.

## Technology direction

- Backend: **Go + GoFiber**
- Full-stack request/rendering: **Inertia + Vue 3 + TypeScript**, menggunakan library internal **Fibertia**
- Runtime validation frontend: **Zod** untuk request, response, props, dan runtime boundaries
- UI components: **Nuxt UI**
- State/navigation: Inertia navigation dan page props; tidak menggunakan `vue-router` atau Pinia sebagai default
- Database: PostgreSQL + **sqlx** dengan prepared statements pada repository
- IDs: internal `BIGSERIAL` + public `nanoid(21)`; external contract hanya `id: text` berisi `public_id`
- Cache/queue: **Valkey**
- Logging: **zerolog**
- Migration: **goose**
- API schema/documentation UI: **Scalar** berbasis OpenAPI
- Object storage: S3-compatible
- Local/runtime tooling: **Podman** + Compose-compatible workflow
- Deployment: containerized, CI/CD, managed infrastructure terlebih dahulu

## Dokumen

1. [Product vision dan strategy](01-product-vision-strategy.md)
2. [Persona, JTBD, dan research plan](02-persona-jtbd-research.md)
3. [Scope MVP dan roadmap](03-scope-mvp-roadmap.md)
4. [UX, IA, dan user flow](04-ux-information-architecture.md)
5. [System architecture](05-system-architecture.md)
6. [Technical stack Go + Vue](06-technical-stack-go-vue.md)
7. [Domain model dan data](07-domain-model-data.md)
8. [API, webhook, dan provider contract](08-api-webhook-provider-contract.md)
9. [Security, compliance, dan reliability](09-security-compliance-reliability.md)
10. [Infrastructure, deployment, dan observability](10-infrastructure-deployment-observability.md)
11. [Pricing, billing, dan unit economics](11-pricing-billing-unit-economics.md)
12. [Delivery process, testing, dan QA](12-delivery-testing-qa.md)
13. [Backlog epics dan user stories](13-backlog-epics-user-stories.md)
14. [Risk register dan open decisions](14-risk-register-open-decisions.md)
15. [Metrics, analytics, dan launch plan](15-metrics-analytics-launch.md)
16. [ID convention dan public contract](16-id-convention-public-contract.md)
17. [Frontend schema validation dengan Zod](17-frontend-schema-validation-zod.md)

## Status planning

| Area | Status |
|---|---|
| Problem/market hypothesis | Draft |
| Target segment | Draft — perlu interview |
| Official API strategy | Direction decided: official-first |
| Unofficial connector | Optional/experimental, bukan default |
| Backend/frontend stack | Accepted direction: GoFiber + Fibertia + Vue 3 + TypeScript |
| MVP scope | Draft |
| UX flow | Draft |
| Technical architecture | Draft — VSA modular monolith |
| Pricing | Hypothesis |
| Compliance/legal | Needs specialist review |

## Decision rule

Jika ada konflik antara kecepatan launch dan keamanan platform/channel, pilih keamanan. Jangan mengiklankan unofficial connector sebagai official API dan jangan menyamakan SLA, capability, atau reliability keduanya.

## Out of scope untuk fase awal

- Meniru seluruh fitur Fonnte.
- Omnichannel penuh.
- AI autonomous agent tanpa guardrail.
- Marketplace integrasi besar.
- Enterprise SSO/SCIM sebelum ada demand.
- Arbitrary bulk messaging tanpa consent/safety control.
