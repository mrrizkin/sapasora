package config

import "time"

type Server struct {
	ReadTimeout  time.Duration `name:"read_timeout"  env:"SERVER_READ_TIMEOUT,default=15s"`
	WriteTimeout time.Duration `name:"write_timeout" env:"SERVER_WRITE_TIMEOUT,default=30s"`
	IdleTimeout  time.Duration `name:"idle_timeout"  env:"SERVER_IDLE_TIMEOUT,default=60s"`
	BodyLimit    int           `name:"body_limit"    env:"SERVER_BODY_LIMIT,default=16777216"`
}
