package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// DistributedLock represents a distributed lock
type DistributedLock struct {
	cache      *Cache
	key        string
	ttl        time.Duration
	token      string
	acquired   bool
	tokenMutex sync.Mutex // protects the token and acquired fields
}

// NewDistributedLock creates a new lock instance
func NewDistributedLock(cache *Cache, key string, ttl time.Duration) *DistributedLock {
	return &DistributedLock{
		cache: cache,
		key:   key,
		ttl:   ttl,
	}
}

// Acquire attempts to acquire the lock
func (l *DistributedLock) Acquire() (bool, error) {
	l.tokenMutex.Lock()
	defer l.tokenMutex.Unlock()

	if l.acquired {
		return true, nil // already acquired
	}

	token, err := generateLockToken()
	if err != nil {
		return false, err
	}

	success, newToken := l.cache.repo.AcquireLock(l.key, l.ttl, token)
	if success {
		l.token = newToken
		l.acquired = true
		return true, nil
	}

	return false, nil
}

// Release releases the lock if it's currently held
func (l *DistributedLock) Release() error {
	l.tokenMutex.Lock()
	defer l.tokenMutex.Unlock()

	if !l.acquired || l.token == "" {
		return errors.New("lock not acquired or already released")
	}

	if l.cache.repo.ReleaseLock(l.key, l.token) {
		l.acquired = false
		l.token = ""
		return nil
	}

	return errors.New(
		"failed to release lock - may have been expired or acquired by another process",
	)
}

// IsAcquired returns whether this lock instance currently holds the lock
func (l *DistributedLock) IsAcquired() bool {
	l.tokenMutex.Lock()
	defer l.tokenMutex.Unlock()
	return l.acquired
}

// generateLockToken creates a unique token for the lock
func generateLockToken() (string, error) {
	// Create a unique token using current time and some randomness
	return fmt.Sprintf("lock_%d_%d", time.Now().UnixNano(), time.Now().UnixMilli()), nil
}

// Block waits until the lock is acquired or timeout occurs
func (l *DistributedLock) Block(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-time.After(10 * time.Millisecond):
			acquired, err := l.Acquire()
			if err != nil {
				return err
			}
			if acquired {
				return nil
			}

			if time.Now().After(deadline) {
				return errors.New("timeout waiting for lock")
			}
		case <-time.After(time.Until(deadline)):
			return errors.New("timeout waiting for lock")
		}
	}
}

// ContextBlock waits until the lock is acquired or context is cancelled
func (l *DistributedLock) ContextBlock(ctx context.Context) error {
	for {
		acquired, err := l.Acquire()
		if err != nil {
			return err
		}
		if acquired {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Millisecond):
			// Continue loop
		}
	}
}
