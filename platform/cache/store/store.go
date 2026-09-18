// Package store provides a implementation of the CacheRepository interface for cache
package store

import (
	"crypto/rand"
	"fmt"
	"time"
)

type cacheItem struct {
	value     any
	expiredAt time.Time
	forever   bool
}

type lockInfo struct {
	token     string
	expiredAt time.Time
}

// generateLockToken creates a unique token for the lock
func generateLockToken(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to time-based token if random generation fails
		return fmt.Sprintf("lock_%d_%d", time.Now().UnixNano(), time.Now().UnixMilli())
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}
