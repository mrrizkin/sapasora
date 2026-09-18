package store

import (
	"errors"
	"sync"
	"time"
)

type Memory struct {
	data  sync.Map
	locks sync.Map // map[string]*lockInfo to track locks
}

func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) Has(key string) bool {
	v, ok := m.data.Load(key)
	if !ok {
		return false
	}
	item := v.(cacheItem)
	if !item.forever && time.Now().After(item.expiredAt) {
		m.data.Delete(key)
		return false
	}
	return true
}

func (m *Memory) Get(key string) (any, bool) {
	if !m.Has(key) {
		return nil, false
	}
	v, _ := m.data.Load(key)
	return v.(cacheItem).value, true
}

func (m *Memory) Pull(key string) (any, bool) {
	val, ok := m.Get(key)
	if ok {
		m.data.Delete(key)
	}
	return val, ok
}

func (m *Memory) Put(key string, value any, ttl time.Duration) {
	m.data.Store(key, cacheItem{
		value:     value,
		expiredAt: time.Now().Add(ttl),
	})
}

func (m *Memory) Add(key string, value any, ttl time.Duration) bool {
	if m.Has(key) {
		return false
	}
	m.Put(key, value, ttl)
	return true
}

func (m *Memory) Increment(key string, n int64) (int64, error) {
	val, ok := m.Get(key)
	if !ok {
		m.Put(key, n, time.Hour)
		return n, nil
	}
	switch v := val.(type) {
	case int:
		v += int(n)
		m.Put(key, v, time.Hour)
		return int64(v), nil
	case int64:
		v += n
		m.Put(key, v, time.Hour)
		return v, nil
	default:
		return 0, errors.New("value is not numeric")
	}
}

func (m *Memory) Decrement(key string, n int64) (int64, error) {
	return m.Increment(key, -n)
}

func (m *Memory) Forever(key string, value any) {
	m.data.Store(key, cacheItem{value: value, forever: true})
}

func (m *Memory) Forget(key string) bool {
	_, ok := m.data.Load(key)
	m.data.Delete(key)
	return ok
}

func (m *Memory) Flush() error {
	m.data.Range(func(k, v any) bool {
		m.data.Delete(k)
		return true
	})
	return nil
}

func (m *Memory) Remember(
	key string,
	ttl time.Duration,
	callback func() (any, error),
) (any, error) {
	if val, ok := m.Get(key); ok {
		return val, nil
	}
	val, err := callback()
	if err != nil {
		return nil, err
	}
	m.Put(key, val, ttl)
	return val, nil
}

func (m *Memory) RememberForever(key string, callback func() (any, error)) (any, error) {
	if val, ok := m.Get(key); ok {
		return val, nil
	}
	val, err := callback()
	if err != nil {
		return nil, err
	}
	m.Forever(key, val)
	return val, nil
}

func (m *Memory) Close() error {
	return nil
}

func (m *Memory) IsLocked(key string) bool {
	v, ok := m.locks.Load(key)
	if !ok {
		return false
	}
	lock := v.(*lockInfo)
	if time.Now().After(lock.expiredAt) {
		m.locks.Delete(key)
		return false
	}
	return true
}

func (m *Memory) AcquireLock(key string, ttl time.Duration, token string) (bool, string) {
	// Check if already locked
	if m.IsLocked(key) {
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
	actual, loaded := m.locks.LoadOrStore(key, lock)
	if loaded {
		// Another goroutine acquired the lock in the meantime
		existingLock := actual.(*lockInfo)
		if time.Now().After(existingLock.expiredAt) {
			// Existing lock is expired, try to replace it
			m.locks.Store(key, lock)
			return true, newToken
		}
		return false, ""
	}

	return true, newToken
}

func (m *Memory) ReleaseLock(key string, token string) bool {
	v, ok := m.locks.Load(key)
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
		m.locks.Delete(key)
		return false // lock was already expired
	}

	// Remove the lock
	m.locks.Delete(key)
	return true
}
