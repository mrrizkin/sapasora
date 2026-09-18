package config

// Provider contains provider lifecycle settings shared by WhatsApp and
// Telegram startup workers.
type Provider struct {
	StartupConcurrency int `name:"startup_concurrency" env:"PROVIDER_STARTUP_CONCURRENCY,default=4"`
}
