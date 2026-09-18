// Package config provides a bridge between the application and environment variables
package config

import "sapasora/platform/config"

// NewAppConfig returns a new AppConfig instance
// @wired:decorate
func NewAppConfig(cfg config.Config) config.Config {
	cfg.LoadStruct("app", &App{})
	cfg.LoadStruct("cache", &Cache{})
	cfg.LoadStruct("database", &Database{})
	cfg.LoadStruct("inertia", &Inertia{})
	cfg.LoadStruct("mail", &Mail{})
	cfg.LoadStruct("pubsub", &Pubsub{})
	cfg.LoadStruct("scheduler", &Scheduler{})
	cfg.LoadStruct("security", &Security{})
	cfg.LoadStruct("session", &Session{})
	cfg.LoadStruct("storage", &Storage{})
	cfg.LoadStruct("telegram", &Telegram{})
	return cfg
}
