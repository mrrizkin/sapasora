// Package config provides the cache configuration
package config

import (
	"sapasora/platform/config"
	"time"
)

// Config represents the entire cache configuration
type Config struct {
	Default string
	Prefix  string
	Stores  map[string]CacheStore
	TTL     time.Duration
}

// CacheStore describes each cache store backend configuration
type CacheStore struct {
	Driver       string
	Path         string
	Table        string
	Connection   string
	LockConn     string
	PersistentID string
	SASLUser     string
	SASLPassword string
	Servers      []MemcachedServer
	Key          string
	Secret       string
	Region       string
	TableName    string
	Endpoint     string
}

// MemcachedServer represents a memcached server entry
type MemcachedServer struct {
	Host   string
	Port   int
	Weight int
}

func NewConfig(cfg config.Config) *Config {
	stores := make(map[string]CacheStore)

	stores["file"] = CacheStore{
		Driver: "file",
		Path:   cfg.GetString("cache.file.path"),
	}

	stores["memory"] = CacheStore{
		Driver: "memory",
	}

	stores["redis"] = CacheStore{
		Driver:     "redis",
		Connection: cfg.GetString("cache.redis.connection"),
		LockConn:   cfg.GetString("cache.redis.lock_connection"),
		Key:        cfg.GetString("cache.redis.key"),
		Secret:     cfg.GetString("cache.redis.secret"),
		Servers: []MemcachedServer{
			{
				Host:   cfg.GetString("cache.redis.host"),
				Port:   cfg.GetInt("cache.redis.port"),
				Weight: cfg.GetInt("cache.redis.weight"),
			},
		},
	}

	stores["memcached"] = CacheStore{
		Driver: "memcached",
		Servers: []MemcachedServer{
			{
				Host:   cfg.GetString("cache.memcached.host"),
				Port:   cfg.GetInt("cache.memcached.port"),
				Weight: cfg.GetInt("cache.memcached.weight"),
			},
		},
	}

	return &Config{
		Default: cfg.GetString("cache.default"),
		Prefix:  cfg.GetString("cache.prefix"),
		Stores:  stores,
		TTL:     cfg.GetDuration("cache.ttl"),
	}
}
