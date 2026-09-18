package config

import "time"

type Scheduler struct {
	Workers    int           `name:"workers"      env:"SCHEDULER_WORKERS,default=3"`
	Queue      string        `name:"queue"        env:"SCHEDULER_QUEUE,default=default"`
	MaxRetries int           `name:"max_retries"  env:"SCHEDULER_MAX_RETRIES,default=25"`
	Timeout    time.Duration `name:"timeout"      env:"SCHEDULER_TIMEOUT,default=30m"`
	DBStorage  bool          `name:"db_storage"   env:"SCHEDULER_DB_STORAGE,default=false"`
	DBDriver   string        `name:"db_driver"    env:"SCHEDULER_DB_DRIVER"`
	DBDSN      string        `name:"db_dsn"       env:"SCHEDULER_DB_DSN"`
}
