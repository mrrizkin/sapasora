// Package message contains the provider-neutral Track 7.1 message domain.
//
// It owns tenant/workspace-scoped message records, attachment metadata,
// delivery observations, and append-only message events. Provider send APIs,
// state transition policy, transports, migrations, and UI wiring remain outside
// this package.
package message
