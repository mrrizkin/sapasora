package store

import (
	"errors"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

type Ristretto struct {
	cache *ristretto.Cache[string, any]
	locks sync.Map // map[string]*lockInfo to track locks (in-memory for this process)
}

func NewRistretto() *Ristretto {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config[string, any]{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}

	return &Ristretto{
		cache: ristrettoCache,
	}
}

func (r *Ristretto) Has(key string) bool {
	_, ok := r.cache.Get(key)
	return ok
}

func (r *Ristretto) Get(key string) (any, bool) {
	return r.cache.Get(key)
}

func (r *Ristretto) Pull(key string) (any, bool) {
	v, ok := r.cache.Get(key)
	if !ok {
		return nil, false
	}

	r.cache.Del(key)
	return v, true
}

func (r *Ristretto) Put(key string, value any, ttl time.Duration) {
	r.cache.SetWithTTL(key, value, 0, ttl)
}

func (r *Ristretto) Add(key string, value any, ttl time.Duration) bool {
	if r.Has(key) {
		return false
	}
	r.Put(key, value, ttl)
	return true
}

func (r *Ristretto) Increment(key string, n int64) (int64, error) {
	val, ok := r.cache.Get(key)
	if !ok {
		r.Put(key, n, time.Hour)
		return n, nil
	}
	switch v := val.(type) {
	case int:
		v += int(n)
		r.Put(key, v, time.Hour)
		return int64(v), nil
	case int64:
		v += n
		r.Put(key, v, time.Hour)
		return v, nil
	default:
		return 0, errors.New("value is not numeric")
	}
}

func (r *Ristretto) Decrement(key string, n int64) (int64, error) {
	return r.Increment(key, -n)
}

func (r *Ristretto) Forever(key string, value any) {
	r.cache.SetWithTTL(key, value, 0, 0)
}

func (r *Ristretto) Forget(key string) bool {
	_, ok := r.cache.Get(key)
	r.cache.Del(key)
	return ok
}

func (r *Ristretto) Flush() error {
	r.cache.Clear()
	return nil
}

func (r *Ristretto) Remember(
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

func (r *Ristretto) RememberForever(key string, callback func() (any, error)) (any, error) {
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

func (r *Ristretto) Close() error {
	r.cache.Close()
	return nil
}

func (r *Ristretto) IsLocked(key string) bool {
	// For Ristretto, we'll also maintain an in-memory lock tracking similar to Memory store
	// Since Ristretto is an in-memory cache, this works for process-level locks
	v, ok := r.locks.Load(key)
	if !ok {
		return false
	}
	lock := v.(*lockInfo)
	if time.Now().After(lock.expiredAt) {
		r.locks.Delete(key)
		return false
	}
	return true
}

func (r *Ristretto) AcquireLock(key string, ttl time.Duration, token string) (bool, string) {
	// Check if already locked
	if r.IsLocked(key) {
		return false, ""
	}

	newToken := token
	if newToken == "" {
		// Generate a new token if not provided
		newToken = generateLockToken(16)
	}

	lock := &lockInfo{
		token:     newToken,
		expiredAt: time.Now().Add(ttl),
	}

	// Use CompareAndSwap to ensure atomic operation
	actual, loaded := r.locks.LoadOrStore(key, lock)
	if loaded {
		// Another goroutine acquired the lock in the meantime
		existingLock := actual.(*lockInfo)
		if time.Now().After(existingLock.expiredAt) {
			// Existing lock is expired, try to replace it
			r.locks.Store(key, lock)
			return true, newToken
		}
		return false, ""
	}

	return true, newToken
}

func (r *Ristretto) ReleaseLock(key string, token string) bool {
	v, ok := r.locks.Load(key)
	if !ok {
		return false // lock doesn't exist
	}

	lock := v.(*lockInfo)

	// Check if the token matches
	if lock.token != token {
		return false // wrong token, can't release
	}

	// Check if lock is still valid (not expired)
	if time.Now().After(lock.expiredAt) {
		r.locks.Delete(key)
		return false // lock was already expired
	}

	// Remove the lock
	r.locks.Delete(key)
	return true
}
