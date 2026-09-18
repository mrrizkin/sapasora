# Channel adapter contract

Track 4.2 adds the provider-neutral contract in
`internal/modules/channel`. It separates connection, health, capability,
send, contact/user, event, address, error, and inbound-event concerns so a
provider can implement only the capabilities it supports.

The package includes a deterministic in-memory fake and contract tests. The
contract uses normalized domain shapes and opaque event envelopes; it does not
import WhatsApp, Telegram, or any other provider payload type. Error mapping
uses stable categories and a provider-code catalog.

This change is intentionally additive. Existing WhatsApp and Telegram
implementations do **not** adopt the contract yet. The next step is a
provider-by-provider integration and migration, including capability mapping,
lifecycle behavior, and provider-specific normalization at each adapter
boundary.
