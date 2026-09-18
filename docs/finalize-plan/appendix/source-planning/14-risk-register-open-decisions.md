# 14 — Risk Register dan Open Decisions

## Risk register

| ID | Risiko | Prob. | Impact | Mitigasi | Owner |
|---|---|---:|---:|---|---|
| R1 | Approval/eligibility official channel tertunda | M | H | Mulai dengan test assets/partner, siapkan checklist | Product/Business |
| R2 | Perubahan Meta policy/API | M | H | Adapter boundary, changelog monitoring, contract tests | Platform |
| R3 | Whatsmeow protocol break/ban | H | H | Experimental flag, isolated worker, kill switch, disclosure | Channel |
| R4 | Cross-tenant data leak | L | Critical | Tenant context, tests, review, optional RLS | Engineering |
| R5 | Credential/API key leak | M | Critical | Encryption, scoped key, redacted logs, rotation | Security |
| R6 | Duplicate outbound message | M | H | Idempotency, dedup, retry classification | Messaging |
| R7 | Webhook customer endpoint down | H | M | Durable retry, backoff, DLQ, replay | Platform |
| R8 | Campaign spam/abuse | M | H | Consent, approval, rate guard, opt-out | Product/Trust |
| R9 | Provider fee/unit economics buruk | M | H | Usage ledger, cost alert, pilot pricing | Finance |
| R10 | Scope terlalu besar | H | H | Wedge jelas, P0 gate, vertical slice | Product |
| R11 | Docs/API/UI tidak sinkron | M | M | OpenAPI/schema source of truth, CI check | Developer Experience |
| R12 | Sensitive data retention berlebihan | M | H | Retention default, deletion job, privacy review | Security/Legal |
| R13 | AI memberi jawaban salah/bocor | M | H | Assist-only, knowledge boundaries, human handoff | Product/AI |
| R14 | Support load tinggi | M | M | Diagnostics, health check, self-serve docs | Support |

## Open decisions

### Product

- [ ] Nama, brand, dan kategori produk.
- [ ] Wedge pertama: notification, inbox, atau keduanya.
- [ ] Negara pertama selain Indonesia.
- [ ] Apakah campaign masuk MVP atau P1.
- [ ] Apakah unofficial connector tersedia sejak private beta.

### Channel/business

- [ ] Direct Meta Tech Provider vs partner/BSP.
- [ ] Ownership billing Meta/provider.
- [ ] Official Cloud API vs provider abstraction pihak ketiga.
- [ ] Policy dan legal review untuk whatsmeow.
- [ ] Capability yang benar-benar diperlukan untuk groups/personal number.

### Technical

- [ ] Queue semantics di atas Valkey: Streams/worker library yang dipilih.
- [ ] Auth: session cookie vs managed auth provider.
- [ ] Cloud/provider deployment dengan Podman-compatible runtime.
- [ ] Payment provider dan tax/invoice.
- [ ] Message/media retention final.
- [ ] OpenAPI generation approach dengan Scalar.
- [ ] Multi-region requirement.
- [ ] Convention Fibertia: shared props, validation, partial reload, upload.
- [ ] Convention Zod: schema ownership, versioning, generation dari OpenAPI, dan error boundary.
- [ ] Prepared statement lifecycle pada repository sqlx.
- [ ] Nanoid generation: application layer vs database function.

### UX

- [ ] Bahasa default dan localization.
- [ ] Persona onboarding branching.
- [ ] Design system choice/build.
- [ ] Approval workflow campaign.
- [ ] Detail status yang visible untuk nontechnical user.

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
### Positive
- ...
### Negative
- ...

## Follow-up
- [ ] ...
```

## Current decisions

### D-001 — Go + Vue

**Status:** Proposed/accepted as implementation direction.

**Reason:** developer experience yang sudah ada, Go cocok untuk network/concurrency/provider workers, Vue cocok untuk dashboard SPA dan reuse pengalaman Whatsmeow gateway.

### D-002 — Official-first

**Status:** Accepted direction.

**Reason:** lebih tepat untuk production/trust. Unofficial tetap dapat didukung sebagai adapter terisolasi, bukan default.

### D-003 — VSA modular monolith

**Status:** Accepted direction.

**Reason:** feature/capability menjadi boundary yang mudah dipahami dan dikirim sebagai vertical slice. Ports/adapters tetap digunakan pada boundary eksternal, tetapi bukan struktur utama codebase.

### D-004 — GoFiber + Fibertia + Vue

**Status:** Accepted.

**Reason:** sesuai pengalaman developer dan menjadi signature teknis produk; Inertia mengurangi kebutuhan client-side router/global state pada dashboard.

### D-005 — Dual ID

**Status:** Accepted.

**Reason:** `BIGSERIAL` efisien untuk internal join/index, sedangkan `nanoid(21)` aman dan stabil untuk public reference. Pada external contract, field `id: string` diisi dari `public_id`; internal `id` dan field `public_id` tidak pernah keluar.

### D-006 — Tooling

**Status:** Accepted.

**Decision:** sqlx + prepared statements, goose, zerolog, Valkey, Podman, dan OpenAPI + Scalar.

### D-007 — Safety default

**Status:** Accepted direction.

**Reason:** consent, opt-out, rate guard, approval, dan observability adalah bagian dari value, bukan fitur tambahan.

### D-008 — Zod runtime validation

**Status:** Accepted.

**Decision:** semua frontend runtime boundary—request, response, Inertia props, shared props, storage, dan external payload—divalidasi dengan Zod. Invalid data harus fail early/loudly melalui error boundary dan telemetry ter-redact; tidak boleh silent error.
