// Package config provides a bridge between the application and environment variables
package config

import (
	"fmt"

	"sapasora/platform/config"
)

// NewAppConfig returns a new AppConfig instance.
// @wired:decorate
func NewAppConfig(cfg config.Config) (config.Config, error) {
	configs := []struct {
		name   string
		target any
	}{
		{name: "app", target: &App{}},
		{name: "cache", target: &Cache{}},
		{name: "database", target: &Database{}},
		{name: "inertia", target: &Inertia{}},
		{name: "mail", target: &Mail{}},
		{name: "pubsub", target: &Pubsub{}},
		{name: "scheduler", target: &Scheduler{}},
		{name: "security", target: &Security{}},
		{name: "server", target: &Server{}},
		{name: "session", target: &Session{}},
		{name: "storage", target: &Storage{}},
		{name: "telegram", target: &Telegram{}},
	}

	for _, item := range configs {
		if err := cfg.LoadStruct(item.name, item.target); err != nil {
			return nil, fmt.Errorf("load %s config: %w", item.name, err)
		}
	}

	return cfg, nil
}
