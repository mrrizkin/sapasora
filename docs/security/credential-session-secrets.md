# Credential/session secret foundation

`platform/secret` is the provider-neutral P4.4 foundation for encrypted credentials and provider-session material.

## Contract

- `Engine` uses AES-256-GCM. The envelope records the algorithm, key version, creation time, expiry metadata, nonce, and ciphertext.
- Envelope metadata is authenticated as additional data. Ciphertext, key, expiry, or version tampering fails authentication (or key lookup) without returning plaintext.
- `Keyring` resolves the current key and historical keys by opaque version. `ReEncrypt` decrypts with the recorded version and seals with the current version, preserving expiry metadata.
- `SecretReference` is the DTO-safe shape: ID, kind, key version, and timestamps only. `SecretRecord` is the persistence shape: a reference plus ciphertext envelope and reveal/revoke metadata. Neither contains plaintext.
- `SecretStore.ClaimReveal` is an atomic persistence contract. A durable implementation must claim a reveal transactionally so concurrent requests cannot both obtain a one-time reveal. `Replace` must compare-and-swap reveal/revoke metadata so a stale rotation cannot clear a lifecycle decision.
- `Service.Reveal` returns an ephemeral `RevealedSecret` with redacting string/JSON representations. A creation or rotation flow may explicitly hand off the value with `Store`/`Rotate` followed by `Reveal`; the persisted claim succeeds once only. Provider code should call `Take`/`Bytes` only at the provider boundary and wipe the value immediately after use.
- `AccessAuditEvent` is the complete audit payload: secret ID, kind, action, outcome, and timestamp only. The optional `AccessAuditor` contract records store, reveal, rotate, and revoke outcomes without receiving plaintext, ciphertext, envelopes, or provider metadata. `NewService` discovers an auditor implemented by the store; `NewServiceWithAudit` permits a separate sink.
- `SecretPurger` is an optional store capability exposed by `SecretService.Purge`. It removes only records revoked or expired at the supplied clock time. Existing `SecretStore` implementations remain valid without implementing it.

## Integration boundary

This package is additive and is not wired into the existing WhatsApp, Telegram, or `platform/session` providers. Provider integration requires an explicit durable `SecretStore` persistence contract, tenant/owner authorization, an `AccessAuditor` (either as the store's optional capability or a separate `NewServiceWithAudit` sink), a transactional `SecretPurger` retention implementation, and provider-specific lifecycle handling. Existing provider sessions must not be migrated implicitly or persisted through the in-memory store. Audit implementations must persist only the five `AccessAuditEvent` fields; purge jobs must run only after the owning provider disconnect/delete policy authorizes removal.

Production keyrings should resolve keys from a KMS or secret manager. `StaticKeyring` and `MemoryStore` are intended for tests/local development only. Application logs and diagnostics should log only `SecretReference` or the package's redacting representations; never log values returned by `RevealedSecret.Bytes`/`Take`.
