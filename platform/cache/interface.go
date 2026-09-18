package cache

import "time"

type CacheRepository interface {
	Has(key string) bool
	Get(key string) (any, bool)
	Pull(key string) (any, bool)
	Put(key string, value any, ttl time.Duration)
	Add(key string, value any, ttl time.Duration) bool
	Increment(key string, n int64) (int64, error)
	Decrement(key string, n int64) (int64, error)
	Forever(key string, value any)
	Forget(key string) bool
	Flush() error
	Remember(key string, ttl time.Duration, callback func() (any, error)) (any, error)
	RememberForever(key string, callback func() (any, error)) (any, error)
	Close() error

	// Lock interface
	AcquireLock(key string, ttl time.Duration, token string) (bool, string)
	ReleaseLock(key string, token string) bool
	IsLocked(key string) bool
}
