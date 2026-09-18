package config

import "time"

type Cache struct {
	Default string                `name:"default" env:"CACHE_DEFAULT,default=file"`
	Prefix  string                `name:"prefix"  env:"CACHE_PREFIX,default=app_cache"`
	Stores  map[string]CacheStore `name:"stores"`
	TTL     time.Duration         `name:"ttl"     env:"CACHE_TTL,default=300s"`
}

type CacheStore struct {
	Driver       string            `name:"driver"          env:"CACHE_DRIVER"`
	Path         string            `name:"path"            env:"CACHE_PATH"`
	Table        string            `name:"table"           env:"CACHE_TABLE"`
	Connection   string            `name:"connection"      env:"CACHE_CONNECTION"`
	LockConn     string            `name:"lock_connection" env:"CACHE_LOCK_CONNECTION"`
	PersistentID string            `name:"persistent_id"   env:"CACHE_PERSISTENT_ID"`
	SASLUser     string            `name:"sasl_user"       env:"CACHE_SASL_USER"`
	SASLPassword string            `name:"sasl_password"   env:"CACHE_SASL_PASS"`
	Servers      []MemcachedServer `name:"servers"`
	Key          string            `name:"key"             env:"CACHE_DYNAMODB_KEY"`
	Secret       string            `name:"secret"          env:"CACHE_DYNAMODB_SECRET"`
	Region       string            `name:"region"          env:"CACHE_DYNAMODB_REGION,default=us-east-1"`
	TableName    string            `name:"table_name"      env:"CACHE_DYNAMODB_TABLE,default=cache"`
	Endpoint     string            `name:"endpoint"        env:"CACHE_DYNAMODB_ENDPOINT"`
}

type MemcachedServer struct {
	Host   string `name:"host"   env:"MEMCACHED_HOST,default=127.0.0.1"`
	Port   int    `name:"port"   env:"MEMCACHED_PORT,default=11211"`
	Weight int    `name:"weight" env:"MEMCACHED_WEIGHT,default=100"`
}
