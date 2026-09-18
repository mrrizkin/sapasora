package cache

import (
	"context"
	"fmt"
	"time"

	"sapasora/platform/cache/config"
	"sapasora/platform/cache/store"

	"go.uber.org/fx"
)

type Cache struct {
	repo   CacheRepository
	config *config.Config
}

func NewCache(lc fx.Lifecycle, config *config.Config) (*Cache, error) {

	var repo CacheRepository
	switch config.Default {
	case "memory":
		repo = store.NewMemory()
	case "ristretto":
		repo = store.NewRistretto()
	case "redis":
		repo = store.NewRedis(
			fmt.Sprintf(
				"%s:%d",
				config.Stores["redis"].Servers[0].Host,
				config.Stores["redis"].Servers[0].Port,
			),
			config.Stores["redis"].Key, // password
			0,                          // DB number, can be made configurable
		)
	default:
		panic("Cache type not supported")
	}

	cache := &Cache{
		repo:   repo,
		config: config,
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return cache.Close()
		},
	})

	return cache, nil
}

func (c *Cache) Get(key string) (any, bool) {
	return c.repo.Get(key)
}

func (c *Cache) Has(key string) bool {
	return c.repo.Has(key)
}

func (c *Cache) Pull(key string) (any, bool) {
	return c.repo.Pull(key)
}

func (c *Cache) Put(key string, value any, ttl time.Duration) {
	c.repo.Put(key, value, ttl)
}

func (c *Cache) Add(key string, value any, ttl time.Duration) bool {
	return c.repo.Add(key, value, ttl)
}

func (c *Cache) Increment(key string, n int64) (int64, error) {
	return c.repo.Increment(key, n)
}

func (c *Cache) Decrement(key string, n int64) (int64, error) {
	return c.repo.Decrement(key, n)
}

func (c *Cache) Forever(key string, value any) {
	c.repo.Forever(key, value)
}

func (c *Cache) Forget(key string) bool {
	return c.repo.Forget(key)
}

func (c *Cache) Flush() error {
	return c.repo.Flush()
}

func (c *Cache) Remember(key string, ttl time.Duration, callback func() (any, error)) (any, error) {
	return c.repo.Remember(key, ttl, callback)
}

func (c *Cache) RememberForever(key string, callback func() (any, error)) (any, error) {
	return c.repo.RememberForever(key, callback)
}

// Tag returns a new TaggedCache instance for the given tag names
func (c *Cache) Tag(names ...string) *TaggedCache {
	tagSet := make(map[string]bool)
	for _, name := range names {
		tagSet[name] = true
	}

	return &TaggedCache{
		cache:    c,
		tagNames: names,
		tagSet:   tagSet,
	}
}

func (c *Cache) Close() error {
	return c.repo.Close()
}

func Get[T any](c *Cache, key string, defaultFn ...func() (T, error)) (T, bool) {
	if v, ok := c.Get(key); ok {
		return v.(T), true
	}
	if len(defaultFn) > 0 {
		if v, err := defaultFn[0](); err == nil {
			return v, true
		}
	}
	return zeroValue[T](), false
}

func Has(c *Cache, key string) bool {
	return c.Has(key)
}

func Forget(c *Cache, key string) bool {
	return c.Forget(key)
}

func Flush(c *Cache) error {
	return c.Flush()
}

func Remember[T any](
	c *Cache,
	key string,
	ttl time.Duration,
	callback func() (T, error),
) (T, error) {
	if val, ok := c.Get(key); ok {
		return val.(T), nil
	}
	val, err := callback()
	if err != nil {
		return zeroValue[T](), err
	}
	c.Put(key, val, ttl)
	return val, nil
}

func RememberForever[T any](c *Cache, key string, callback func() (T, error)) (T, error) {
	if val, ok := c.Get(key); ok {
		return val.(T), nil
	}
	val, err := callback()
	if err != nil {
		return zeroValue[T](), err
	}
	c.Forever(key, val)
	return val, nil
}

func Increment(c *Cache, key string, n ...int64) (int64, error) {
	if len(n) > 0 {
		return c.repo.Increment(key, n[0])
	}

	return c.repo.Increment(key, 1)
}

func Decrement(c *Cache, key string, n ...int64) (int64, error) {
	if len(n) > 0 {
		return c.repo.Decrement(key, n[0])
	}
	return c.repo.Decrement(key, 1)
}

func Add(c *Cache, key string, value any, ttl time.Duration) bool {
	return c.repo.Add(key, value, ttl)
}

func Lock(c *Cache, key string, ttl time.Duration) *DistributedLock {
	return NewDistributedLock(c, key, ttl)
}

// Block waits until the lock is acquired or timeout occurs
func Block(lock *DistributedLock, timeout time.Duration) error {
	return lock.Block(timeout)
}

// Release releases the acquired lock
func Release(lock *DistributedLock) error {
	return lock.Release()
}

// AcquireLock attempts to acquire the lock immediately
func AcquireLock(lock *DistributedLock) (bool, error) {
	return lock.Acquire()
}

func Tags(c *Cache, keys ...string) *TaggedCache {
	return c.Tag(keys...)
}

func zeroValue[T any]() T {
	var zero T
	return zero
}
