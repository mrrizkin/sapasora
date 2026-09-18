// Package channel defines the provider-neutral channel domain, adapter
// boundary, and channel lifecycle state machine.
//
// The contracts and lifecycle state machine are additive. Existing WhatsApp
// and Telegram implementations are intentionally not required to implement
// them yet; provider-by-provider integration and migration remain separate
// work.
package channel
