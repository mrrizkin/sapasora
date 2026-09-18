// Package config provides the config tool for the application
package config

import (
	"time"
)

type Config interface {
	LoadStruct(name string, cfg any) error
	Set(key string, value any)
	Get(key string, defaultValue ...any) any
	GetString(key string, defaultValue ...string) string
	GetInt(key string, defaultValue ...int) int
	GetBool(key string, defaultValue ...bool) bool
	GetFloat64(key string, defaultValue ...float64) float64
	GetDuration(key string, defaultValue ...time.Duration) time.Duration
	Has(key string) bool
	All() map[string]any
}
