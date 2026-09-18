package scheduler

import (
	"time"

	"sapasora/platform/config"
)

// LoadConfig loads scheduler configuration from the application config
func LoadConfig(cfg config.Config) SchedulerConfig {
	return SchedulerConfig{
		Workers:    cfg.GetInt("scheduler.workers", 3),
		Queue:      cfg.GetString("scheduler.queue", "default"),
		MaxRetries: cfg.GetInt("scheduler.max_retries", 25),
		Timeout:    cfg.GetDuration("scheduler.timeout", 30*time.Minute),
		DBStorage:  cfg.GetBool("scheduler.db_storage", false),
		DBDriver:   cfg.GetString("scheduler.db_driver", ""),
		DBDSN:      cfg.GetString("scheduler.db_dsn", ""),
	}
}
