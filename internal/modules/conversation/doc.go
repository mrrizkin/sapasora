// Package conversation defines the provider-neutral conversation domain.
//
// Conversation resources are tenant/workspace scoped and intentionally stop at
// the domain and repository boundary. HTTP, UI, persistence migrations,
// authorization, message delivery, and provider adapters remain integration
// concerns. Diagnostics omit participant contact details and other PII.
package conversation
