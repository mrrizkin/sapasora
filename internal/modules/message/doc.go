// Package message contains the provider-neutral Track 7.1/7.2 message domain.
//
// It owns tenant/workspace-scoped message records, attachment metadata,
// delivery observations, strict message state transitions, and append-only
// message events. Provider send APIs, transports, migrations, and UI wiring
// remain outside this package.
package message
