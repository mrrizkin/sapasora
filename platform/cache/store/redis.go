package store

import (
	"context"
	"sapasora/platform/support/console"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
	mu     sync.RWMutex // protects access to client
}

func NewRedis(addr, password string, db int) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Redis{
		client: rdb,
	}
}

func (r *Redis) Has(key string) bool {
	ctx := context.Background()
	_, err := r.client.Get(ctx, key).Result()
	return err == nil
}

func (r *Redis) Get(key string) (any, bool) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, false
	}
	return val, true
}

func (r *Redis) Pull(key string) (any, bool) {
	ctx := context.Background()
	pipe := r.client.TxPipeline()

	getCmd := pipe.Get(ctx, key)
	_ = pipe.Del(ctx, key)
	_, err := pipe.Exec(ctx)

	if err != nil || getCmd.Err() != nil {
		return nil, false
	}

	val, getErr := getCmd.Result()
	if getErr != nil {
		return nil, false
	}

	return val, true
}

func (r *Redis) Put(key string, value any, ttl time.Duration) {
	ctx := context.Background()
	err := r.client.SetEx(ctx, key, value, ttl).Err()
	if err != nil {
		// FIXME: handle this error appropriately
		console.Error("Error setting key %s in Redis: %v", key, err)
	}
}

func (r *Redis) Add(key string, value any, ttl time.Duration) bool {
	ctx := context.Background()
	ok, err := r.client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return false
	}
	return ok
}

func (r *Redis) Increment(key string, n int64) (int64, error) {
	ctx := context.Background()
	return r.client.IncrBy(ctx, key, n).Result()
}

func (r *Redis) Decrement(key string, n int64) (int64, error) {
	ctx := context.Background()
	return r.client.DecrBy(ctx, key, n).Result()
}

func (r *Redis) Forever(key string, value any) {
	ctx := context.Background()
	err := r.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		// FIXME: handle this error appropriately
		console.Error("Error setting key %s permanently in Redis: %v\n", key, err)
	}
}

func (r *Redis) Forget(key string) bool {
	ctx := context.Background()
	deleted, err := r.client.Del(ctx, key).Result()
	if err != nil {
		console.Error("Error deleting key %s from Redis: %v\n", key, err)
		return false
	}
	return deleted > 0
}

func (r *Redis) Flush() error {
	ctx := context.Background()
	return r.client.FlushDB(ctx).Err()
}

func (r *Redis) Remember(
	key string,
	ttl time.Duration,
	callback func() (any, error),
) (any, error) {
	if val, ok := r.Get(key); ok {
		return val, nil
	}
	val, err := callback()
	if err != nil {
		return nil, err
	}
	r.Put(key, val, ttl)
	return val, nil
}

func (r *Redis) RememberForever(key string, callback func() (any, error)) (any, error) {
	if val, ok := r.Get(key); ok {
		return val, nil
	}
	val, err := callback()
	if err != nil {
		return nil, err
	}
	r.Forever(key, val)
	return val, nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}

// IsLocked checks if a key is currently locked in Redis
func (r *Redis) IsLocked(key string) bool {
	ctx := context.Background()
	lockKey := fmt.Sprintf("lock:%s", key)

	val, err := r.client.Get(ctx, lockKey).Result()
	if err != nil {
		return false // Key doesn't exist or other error
	}

	// The value should contain the lock token
	return val != ""
}

// AcquireLock attempts to acquire a distributed lock in Redis
func (r *Redis) AcquireLock(key string, ttl time.Duration, token string) (bool, string) {
	ctx := context.Background()
	lockKey := fmt.Sprintf("lock:%s", key)

	newToken := token
	if newToken == "" {
		newToken = generateLockToken(32)
	}

	// Use SET with NX (Not eXists) and EX (EXpire) options
	// This is atomic and will only set the key if it doesn't already exist
	ok, err := r.client.SetNX(ctx, lockKey, newToken, ttl).Result()
	if err != nil {
		return false, ""
	}

	if ok {
		return true, newToken
	}

	// If the lock was not acquired, return false with empty token
	return false, ""
}

// ReleaseLock releases a lock in Redis if the token matches
func (r *Redis) ReleaseLock(key string, token string) bool {
	ctx := context.Background()
	lockKey := fmt.Sprintf("lock:%s", key)

	script := `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
	`

	result, err := r.client.Eval(ctx, script, []string{lockKey}, token).Result()
	if err != nil {
		return false
	}

	// If the result is 1, it means the key was deleted (lock was released)
	// If the result is 0, it means the token didn't match or key was already expired
	return result == int64(1)
}
