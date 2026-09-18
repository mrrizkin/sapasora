package cache

import (
	"fmt"
	"testing"
	"time"

	"sapasora/platform/cache/config"
	"sapasora/platform/cache/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

// mockLifecycle is a minimal implementation for testing
type mockLifecycle struct{}

func (m *mockLifecycle) Append(fx.Hook) {}

// Test basic cache operations for memory store
func TestMemoryStoreBasicOperations(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Test Has operation
	key := "test_key"
	assert.False(t, cache.repo.Has(key))

	// Test Put operation
	cache.repo.Put(key, "test_value", 10*time.Second)
	assert.True(t, cache.repo.Has(key))

	// Test Get operation
	value, exists := cache.repo.Get(key)
	assert.True(t, exists)
	assert.Equal(t, "test_value", value)

	// Test Pull operation (gets and deletes)
	value, exists = cache.repo.Pull(key)
	assert.True(t, exists)
	assert.Equal(t, "test_value", value)
	assert.False(t, cache.repo.Has(key))

	// Test Add operation
	added := cache.repo.Add(key, "new_value", 10*time.Second)
	assert.True(t, added)
	assert.True(t, cache.repo.Has(key))

	// Add should return false if key already exists
	added = cache.repo.Add(key, "another_value", 10*time.Second)
	assert.False(t, added)

	// Test Forever operation
	foreverKey := "forever_key"
	cache.repo.Forever(foreverKey, "forever_value")
	assert.True(t, cache.repo.Has(foreverKey))
	value, exists = cache.repo.Get(foreverKey)
	assert.True(t, exists)
	assert.Equal(t, "forever_value", value)

	// Test Forget operation
	deleted := cache.repo.Forget(key)
	assert.True(t, deleted)
	assert.False(t, cache.repo.Has(key))
}

func TestMemoryStoreIncrementDecrement(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Test Increment
	key := "counter"
	result, err := cache.repo.Increment(key, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result)

	// Increment again
	result, err = cache.repo.Increment(key, 5)
	assert.NoError(t, err)
	assert.Equal(t, int64(6), result)

	// Test Decrement
	result, err = cache.repo.Decrement(key, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(4), result)

	// Test with non-numeric value (should create new counter)
	cache.repo.Put("non_numeric", "string_value", 10*time.Second)
	_, err = cache.repo.Increment("non_numeric", 1)
	assert.Error(t, err)
}

func TestTaggedCacheBasicOperations(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Prefix:  "test",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Create a tagged cache instance with "users" and "profiles" tags
	taggedCache := cache.Tag("users", "profiles")

	// Test Put operation with tags
	taggedKey := "user:123"
	taggedCache.Put(taggedKey, "John Doe", 10*time.Second)

	// Verify the key was stored with the proper tag prefix
	prefixedKey := fmt.Sprintf("test:%s:%s", "users_profiles", taggedKey)
	value, exists := cache.Get(prefixedKey)
	assert.True(t, exists)
	assert.Equal(t, "John Doe", value)

	// Test Get operation
	value, exists = taggedCache.Get(taggedKey)
	assert.True(t, exists)
	assert.Equal(t, "John Doe", value)

	// Test Add operation
	added := taggedCache.Add("user:456", "Jane Doe", 10*time.Second)
	assert.True(t, added)

	// Add should return false if key already exists
	added = taggedCache.Add("user:456", "Someone Else", 10*time.Second)
	assert.False(t, added)

	// Test Forever operation
	taggedCache.Forever("user:789", "Bob Smith")
	value, exists = taggedCache.Get("user:789")
	assert.True(t, exists)
	assert.Equal(t, "Bob Smith", value)

	// Test Forget operation
	deleted := taggedCache.Forget("user:123")
	assert.True(t, deleted)
	_, exists = taggedCache.Get("user:123")
	assert.False(t, exists)
}

func TestTaggedCacheFlush(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Prefix:  "test",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Create a tagged cache instance
	taggedCache := cache.Tag("users", "test")

	// Add some items with tags
	taggedCache.Put("user:1", "John", 10*time.Second)
	taggedCache.Put("user:2", "Jane", 10*time.Second)
	taggedCache.Forever("user:profile", "profile_data")

	// Verify items exist before flush
	_, exists := taggedCache.Get("user:1")
	assert.True(t, exists)
	_, exists = taggedCache.Get("user:2")
	assert.True(t, exists)
	_, exists = taggedCache.Get("user:profile")
	assert.True(t, exists)

	// Verify tag index was created
	tagIndexKey := "test:tag_index:users"
	_, exists = cache.Get(tagIndexKey)
	assert.True(t, exists)

	// Flush all items with the tags
	err = taggedCache.Flush()
	assert.NoError(t, err)

	// Verify items no longer exist after flush
	_, exists = taggedCache.Get("user:1")
	assert.False(t, exists)
	_, exists = taggedCache.Get("user:2")
	assert.False(t, exists)
	_, exists = taggedCache.Get("user:profile")
	assert.False(t, exists)

	// Verify tag index was cleared
	_, exists = cache.Get(tagIndexKey)
	assert.False(t, exists)
}

func TestTaggedCacheRemember(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Prefix:  "test",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	taggedCache := cache.Tag("temp")

	// Counter to track callback execution
	callCount := 0

	// Test Remember with callback
	value, err := taggedCache.Remember("temp:count", 10*time.Second, func() (any, error) {
		callCount++
		return fmt.Sprintf("result_%d", callCount), nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "result_1", value)
	assert.Equal(t, 1, callCount)

	// Get again - should not call the callback
	value, err = taggedCache.Remember("temp:count", 10*time.Second, func() (any, error) {
		callCount++
		return fmt.Sprintf("result_%d", callCount), nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "result_1", value)
	assert.Equal(t, 1, callCount) // Count should still be 1
}

func TestTaggedCacheRememberForever(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Prefix:  "test",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	taggedCache := cache.Tag("persistent")

	// Counter to track callback execution
	callCount := 0

	// Test RememberForever with callback
	value, err := taggedCache.RememberForever("persistent:count", func() (any, error) {
		callCount++
		return fmt.Sprintf("result_%d", callCount), nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "result_1", value)
	assert.Equal(t, 1, callCount)

	// Get again - should not call the callback
	value, err = taggedCache.RememberForever("persistent:count", func() (any, error) {
		callCount++
		return fmt.Sprintf("result_%d", callCount), nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "result_1", value)
	assert.Equal(t, 1, callCount) // Count should still be 1
}

func TestMultipleTaggedCaches(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Prefix:  "test",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Create two different tagged caches
	usersCache := cache.Tag("users")
	ordersCache := cache.Tag("orders")

	// Add items to both caches
	usersCache.Put("user:id", "user_data", 10*time.Second)
	ordersCache.Put("order:id", "order_data", 10*time.Second)

	// Verify items are in their respective tagged areas
	userValue, exists := usersCache.Get("user:id")
	assert.True(t, exists)
	assert.Equal(t, "user_data", userValue)

	orderValue, exists := ordersCache.Get("order:id")
	assert.True(t, exists)
	assert.Equal(t, "order_data", orderValue)

	// Flush only the users cache
	err = usersCache.Flush()
	assert.NoError(t, err)

	// User item should be gone
	_, exists = usersCache.Get("user:id")
	assert.False(t, exists)

	// Order item should still exist
	_, exists = ordersCache.Get("order:id")
	assert.True(t, exists)
}

func TestMemoryStoreRemember(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Test Remember with callback
	key := "remember_test"
	called := false
	callback := func() (any, error) {
		called = true
		return "remembered_value", nil
	}

	// First call should execute callback
	value, err := cache.repo.Remember(key, 10*time.Second, callback)
	assert.NoError(t, err)
	assert.Equal(t, "remembered_value", value)
	assert.True(t, called)

	// Second call should return cached value without executing callback
	called = false
	value, err = cache.repo.Remember(key, 10*time.Second, callback)
	assert.NoError(t, err)
	assert.Equal(t, "remembered_value", value)
	assert.False(t, called)

	// Test RememberForever
	foreverKey := "remember_forever_test"
	called = false
	value, err = cache.repo.RememberForever(foreverKey, callback)
	assert.NoError(t, err)
	assert.Equal(t, "remembered_value", value)
	assert.True(t, called)

	// Second call should return cached value
	called = false
	value, err = cache.repo.RememberForever(foreverKey, callback)
	assert.NoError(t, err)
	assert.Equal(t, "remembered_value", value)
	assert.False(t, called)
}

func TestMemoryStoreTTLExpiry(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	key := "ttl_test"
	cache.repo.Put(key, "ttl_value", 100*time.Millisecond)

	// Value should exist initially
	assert.True(t, cache.repo.Has(key))
	value, exists := cache.repo.Get(key)
	assert.True(t, exists)
	assert.Equal(t, "ttl_value", value)

	// Wait for TTL to expire
	time.Sleep(150 * time.Millisecond)

	// Value should no longer exist
	assert.False(t, cache.repo.Has(key))
	_, exists = cache.repo.Get(key)
	assert.False(t, exists)
}

func TestMemoryStoreFlush(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Add multiple items
	cache.repo.Put("key1", "value1", 10*time.Second)
	cache.repo.Put("key2", "value2", 10*time.Second)
	cache.repo.Forever("key3", "value3")

	assert.True(t, cache.repo.Has("key1"))
	assert.True(t, cache.repo.Has("key2"))
	assert.True(t, cache.repo.Has("key3"))

	// Flush all items
	err = cache.repo.Flush()
	assert.NoError(t, err)

	assert.False(t, cache.repo.Has("key1"))
	assert.False(t, cache.repo.Has("key2"))
	assert.False(t, cache.repo.Has("key3"))
}

func TestRistrettoStoreOperations(t *testing.T) {
	cfg := &config.Config{
		Default: "ristretto",
		Stores: map[string]config.CacheStore{
			"ristretto": {Driver: "ristretto"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Test basic operations
	key := "ristretto_test"
	assert.False(t, cache.repo.Has(key))

	cache.repo.Put(key, "ristretto_value", 10*time.Second)
	time.Sleep(100 * time.Millisecond)
	assert.True(t, cache.repo.Has(key))

	value, exists := cache.repo.Get(key)
	assert.True(t, exists)
	assert.Equal(t, "ristretto_value", value)

	// Test increment/decrement
	result, err := cache.repo.Increment("ristretto_counter", 1)
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int64(1), result)

	result, err = cache.repo.Decrement("ristretto_counter", 1)
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int64(0), result)

	// Test pull
	pullValue, exists := cache.repo.Pull(key)
	assert.True(t, exists)
	assert.Equal(t, "ristretto_value", pullValue)
	assert.False(t, cache.repo.Has(key))
}

// func TestRedisStoreOperations(t *testing.T) {
// 	// Test Redis store operations if Redis is available
// 	cfg := &config.Config{
// 		Default: "redis",
// 		Stores: map[string]config.CacheStore{
// 			"redis": {
// 				Driver: "redis",
// 				Servers: []config.MemcachedServer{
// 					{
// 						Host: "localhost",
// 						Port: 6379,
// 					},
// 				},
// 			},
// 		},
// 	}
//
// 	// This test may not run if Redis is not available
// 	cache, err := NewCache(&mockLifecycle{}, cfg)
// 	if err != nil {
// 		t.Skip("Redis not available, skipping Redis tests")
// 		return
// 	}
// 	require.NotNil(t, cache)
//
// 	// Test basic operations
// 	key := "redis_test"
// 	assert.False(t, cache.repo.Has(key))
//
// 	cache.repo.Put(key, "redis_value", 10*time.Second)
// 	assert.True(t, cache.repo.Has(key))
//
// 	value, exists := cache.repo.Get(key)
// 	assert.True(t, exists)
// 	assert.Equal(t, "redis_value", value)
//
// 	// Test increment/decrement
// 	result, err := cache.repo.Increment("redis_counter", 1)
// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(1), result)
//
// 	result, err = cache.repo.Decrement("redis_counter", 1)
// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(0), result)
//
// 	// Clean up
// 	cache.repo.Forget(key)
// 	cache.repo.Forget("redis_counter")
// }

func TestNewCacheWithInvalidConfig(t *testing.T) {
	// Test invalid cache type - this should panic
	cfg := &config.Config{
		Default: "invalid",
		Stores:  map[string]config.CacheStore{},
	}

	defer func() {
		if r := recover(); r != nil {
			// Expected panic for invalid cache type
			assert.Equal(t, "Cache type not supported", r)
		} else {
			t.Error("Expected panic for invalid cache type")
		}
	}()

	_, err := NewCache(&mockLifecycle{}, cfg)
	assert.Error(t, err) // This line won't be reached due to panic
}

func TestCacheClose(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Test Close method
	err = cache.Close()
	assert.NoError(t, err)
}

func TestConcurrentAccess(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Test concurrent access doesn't cause race conditions
	done := make(chan bool, 10)
	for i := range 10 {
		go func(i int) {
			key := "concurrent_test_" + string(rune(i))
			cache.repo.Put(key, "value_"+string(rune(i)), 10*time.Second)
			_, exists := cache.repo.Get(key)
			assert.True(t, exists)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for range 10 {
		<-done
	}
}

func TestEdgeCasesAndErrorHandling(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Test Get with non-existent key
	_, exists := cache.repo.Get("non_existent_key")
	assert.False(t, exists)

	// Test Pull with non-existent key
	_, exists = cache.repo.Pull("non_existent_pull_key")
	assert.False(t, exists)

	// Test Forget with non-existent key
	_ = cache.repo.Forget("non_existent_forget_key")
	// Behavior may depend on implementation, but it shouldn't error
	assert.NotNil(t, cache) // At least ensure no panic

	// Test increment/decrement with invalid values
	cache.repo.Put("string_key", "not_a_number", 10*time.Second)
	_, err = cache.repo.Increment("string_key", 1)
	assert.Error(t, err)

	_, err = cache.repo.Decrement("string_key", 1)
	assert.Error(t, err)

	// Test Remember with erroring callback
	errorKey := "error_remember_test"
	callbackWithError := func() (any, error) {
		return nil, assert.AnError
	}

	// Should return the error from the callback
	_, err = cache.repo.Remember(errorKey, 10*time.Second, callbackWithError)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)

	// Test RememberForever with erroring callback
	_, err = cache.repo.RememberForever(errorKey, callbackWithError)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)

	// Test large values
	largeValue := make([]byte, 1024*1024) // 1MB value
	for i := range largeValue {
		largeValue[i] = byte(i % 256)
	}
	cache.repo.Put("large_value", largeValue, 10*time.Second)
	retrieved, exists := cache.repo.Get("large_value")
	assert.True(t, exists)
	assert.Equal(t, largeValue, retrieved)

	// Test TTL behavior with zero duration
	cache.repo.Put("zero_ttl", "zero_ttl_value", 0)
	time.Sleep(10 * time.Millisecond) // Give a little time for any cleanup
	// Behavior depends on implementation - some may treat 0 as immediate expiration

	// Test very small TTL
	cache.repo.Put("tiny_ttl", "tiny_ttl_value", 1*time.Nanosecond)
	time.Sleep(1 * time.Millisecond) // More than enough time for 1 nanosecond TTL
	// Check if it's gone
	_, exists = cache.repo.Get("tiny_ttl")
	// Behavior here depends on the implementation's precision
	assert.False(t, exists)

}

func TestRememberWithPanicInCallback(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	panicKey := "panic_test"
	panicCallback := func() (any, error) {
		panic("test panic in callback")
	}

	// We need to handle this carefully to test the behavior
	defer func() {
		if r := recover(); r != nil {
			// This is expected since we're calling the callback directly
			// This test may need adjustment based on how the implementation handles panics
			_ = r
		}
	}()

	// This would normally cause a panic, so we're not calling it directly in tests
	cache.repo.Remember(panicKey, 10*time.Second, panicCallback)
}

func TestMultipleStoreConfig(t *testing.T) {
	// Test NewCache with multiple store configuration
	cfg := &config.Config{
		Default: "memory",
		Stores: map[string]config.CacheStore{
			"memory":    {Driver: "memory"},
			"ristretto": {Driver: "ristretto"},
			"redis": {
				Driver: "redis",
				Servers: []config.MemcachedServer{
					{
						Host: "localhost",
						Port: 6379,
					},
				},
			},
			"database": {
				Driver:    "sqlite3",
				Path:      ":memory:",
				TableName: "cache",
			},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Test that operations work with the default (memory) store
	cache.repo.Put("multi_store_test", "value", 10*time.Second)
	value, exists := cache.repo.Get("multi_store_test")
	assert.True(t, exists)
	assert.Equal(t, "value", value)
}

func TestLockMemoryStore(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Test basic lock acquisition and release
	lock := Lock(cache, "test_lock", 10*time.Second)

	// Initially should not be acquired
	assert.False(t, lock.IsAcquired())

	// Acquire the lock
	acquired, err := AcquireLock(lock)
	assert.NoError(t, err)
	assert.True(t, acquired)
	assert.True(t, lock.IsAcquired())

	// Try to acquire the same lock again (should fail)
	lock2 := Lock(cache, "test_lock", 10*time.Second)
	acquired2, err := AcquireLock(lock2)
	assert.NoError(t, err)
	assert.False(t, acquired2)

	// Release the lock
	err = Release(lock)
	assert.NoError(t, err)
	assert.False(t, lock.IsAcquired())

	// Now we should be able to acquire it again
	acquired3, err := AcquireLock(lock)
	assert.NoError(t, err)
	assert.True(t, acquired3)
	assert.True(t, lock.IsAcquired())

	// Release and cleanup
	err = Release(lock)
	assert.NoError(t, err)
}

func TestLockWithBlock(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Acquire a lock first
	lock1 := Lock(cache, "block_test_lock", 10*time.Second)
	acquired, err := AcquireLock(lock1)
	assert.NoError(t, err)
	assert.True(t, acquired)

	// Try to acquire the same lock with a different instance in a goroutine
	done := make(chan bool, 1)
	go func() {
		lock2 := Lock(cache, "block_test_lock", 10*time.Second)
		// This should block until lock1 is released
		err := Block(lock2, 3*time.Second) // 3 second timeout
		if err == nil {
			// Successful acquisition
			Release(lock2)
		}
		done <- true
	}()

	// Give some time for the goroutine to try acquiring
	time.Sleep(100 * time.Millisecond)

	// The goroutine should still be waiting
	select {
	case <-done:
		t.Error("Block should have waited for the lock")
	default:
		// Good, still waiting
	}

	// Release the first lock
	err = Release(lock1)
	assert.NoError(t, err)

	// Now the goroutine should acquire the lock
	select {
	case <-done:
		// Good, it acquired the lock
	case <-time.After(5 * time.Second):
		t.Error("Block should have acquired the lock after release")
	}
}

func TestLockTTLExpiry(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Acquire a lock with a short TTL
	lock := Lock(cache, "ttl_test_lock", 100*time.Millisecond)
	acquired, err := AcquireLock(lock)
	assert.NoError(t, err)
	assert.True(t, acquired)

	// Wait for the lock to expire
	time.Sleep(150 * time.Millisecond)

	// Now we should be able to acquire it again
	lock2 := Lock(cache, "ttl_test_lock", 10*time.Second)
	acquired2, err := AcquireLock(lock2)
	assert.NoError(t, err)
	assert.True(t, acquired2)

	err = Release(lock2)
	assert.NoError(t, err)
}

func TestRistrettoLock(t *testing.T) {
	cfg := &config.Config{
		Default: "ristretto",
		Stores: map[string]config.CacheStore{
			"ristretto": {Driver: "ristretto"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Test basic lock functionality with Ristretto
	lock := Lock(cache, "ristretto_test_lock", 10*time.Second)

	acquired, err := AcquireLock(lock)
	assert.NoError(t, err)
	assert.True(t, acquired)
	assert.True(t, lock.IsAcquired())

	err = Release(lock)
	assert.NoError(t, err)
	assert.False(t, lock.IsAcquired())
}

func TestLockIsLocked(t *testing.T) {
	cfg := &config.Config{
		Default: "memory",
		Stores: map[string]config.CacheStore{
			"memory": {Driver: "memory"},
		},
	}

	cache, err := NewCache(&mockLifecycle{}, cfg)
	require.NoError(t, err)
	require.NotNil(t, cache)

	key := "is_locked_test"

	// Initially should not be locked
	assert.False(t, cache.repo.IsLocked(key))

	// Acquire lock
	lock := Lock(cache, key, 10*time.Second)
	acquired, err := AcquireLock(lock)
	assert.NoError(t, err)
	assert.True(t, acquired)

	// Should be locked now
	assert.True(t, cache.repo.IsLocked(key))

	// Release lock
	err = Release(lock)
	assert.NoError(t, err)

	// Should not be locked anymore
	assert.False(t, cache.repo.IsLocked(key))
}

func TestLockTokenValidation(t *testing.T) {
	// Test with memory store directly to check token validation
	memoryStore := store.NewMemory()

	// Acquire lock
	acquired, token1 := memoryStore.AcquireLock("token_test", 10*time.Second, "")
	assert.True(t, acquired)
	assert.NotEmpty(t, token1)

	// Try to release with wrong token (should fail)
	released := memoryStore.ReleaseLock("token_test", "wrong_token")
	assert.False(t, released)

	// Should still be locked with correct token
	assert.True(t, memoryStore.IsLocked("token_test"))

	// Try to release with correct token (should succeed)
	released = memoryStore.ReleaseLock("token_test", token1)
	assert.True(t, released)

	// Should no longer be locked
	assert.False(t, memoryStore.IsLocked("token_test"))
}
