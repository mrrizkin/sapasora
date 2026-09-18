# 17 — Risk Register, Decisions, dan Open Questions

## Risk register

| ID | Risk | Prob. | Impact | Mitigation | Owner |
|---|---|---:|---:|---|---|
| R1 | Official approval/eligibility tertunda | M | H | Test assets/partner, checklist | Product/Business |
| R2 | Meta policy/API berubah | M | H | Integration boundary, monitoring, contract tests | Platform |
| R3 | Whatsmeow protocol break/ban | H | H | Isolated worker, flag, kill switch, disclosure | Channel |
| R4 | Cross-tenant leak | L | Critical | Workspace context, RLS option, tests | Engineering |
| R5 | Credential leak | M | Critical | Encryption, scoped key, rotation, redaction | Security |
| R6 | Duplicate outbound | M | H | Idempotency/dedup/retry policy | Messaging |
| R7 | Customer webhook down | H | M | Durable retry/backoff/DLQ/replay | Platform |
| R8 | Spam/abuse | M | H | Consent, approval, rate guard, opt-out | Trust |
| R9 | Provider/unit economics buruk | M | H | Usage ledger, cost alerts, pilot | Finance |
| R10 | Scope terlalu besar | H | H | Wedge/P0/vertical slices | Product |
| R11 | Docs/API/UI drift | M | M | OpenAPI + Zod contract tests | DX |
| R12 | Retention berlebihan | M | H | Retention/delete/privacy review | Security/Legal |
| R13 | AI hallucination/leak | M | H | Assist-only, boundaries, handoff | AI |
| R14 | Support load tinggi | M | M | Health check, diagnostics, docs | Support |

## Accepted decisions

- VSA modular monolith sebagai architecture utama.
- Go + GoFiber.
- Inertia + Vue 3 + TypeScript dengan Fibertia.
- Nuxt UI.
- Inertia navigation; tanpa vue-router/Pinia default.
- Zod runtime validation dan fail early/loudly.
- PostgreSQL + sqlx/prepared statements terukur.
- Internal `BIGSERIAL`; external `id: string` dari `public_id nanoid(21)`; keduanya tidak bercampur.
- goose, zerolog, Valkey, Podman, OpenAPI + Scalar.
- Official-first; whatsmeow optional/experimental.
- Safety default.

## Open decisions

### Product/business

- Nama/brand.
- Wedge final: notification, inbox, atau keduanya.
- Negara berikutnya.
- Campaign di MVP atau P1.
- Direct Meta Tech Provider vs BSP.
- Billing ownership Meta/provider.
- Legal/Meta review untuk whatsmeow.
- Group/personal capability yang benar-benar diperlukan.

### Technical

- Valkey queue implementation/semantics.
- Auth session vs managed provider.
- Cloud/deployment provider.
- Payment/tax/invoice.
- Retention final.
- OpenAPI generation vs manually maintained Zod strategy.
- Fibertia conventions: shared props, validation, partial reload, upload.
- Prepared statement lifecycle.
- NanoID generation package/convention.
- Multi-region requirement.

### UX

- Bahasa/localization.
- Persona branching.
- Design system tokens.
- Campaign approval.
- Detail status owner vs developer.

## Decision record template

```markdown
# ADR-XXX — Judul

Tanggal: YYYY-MM-DD
Status: Proposed | Accepted | Rejected | Superseded
Owner: ...

## Context
...

## Decision
...

## Alternatives
- ...

## Consequences
...

## Follow-up
- [ ] ...
```

## Questions before implementation

1. Apakah kita direct Meta Tech Provider atau mulai melalui BSP?
2. Segment pertama mana yang akan diwawancarai dan dipilotkan?
3. Apakah official onboarding bisa diuji sekarang dengan test asset?
4. Apakah whatsmeow masuk private beta atau ditunda sampai legal review?
5. Apa minimum feature yang membuat owner membayar?
6. Apa minimum API contract untuk developer pilot?
7. Berapa target workspace pilot dan siapa owner support-nya?
8. Nama sementara apa yang dipakai di repo/domain?
