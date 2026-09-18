package cache

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

// TaggedCache provides operations on a set of tagged cache items
type TaggedCache struct {
	cache    *Cache
	tagNames []string
	tagSet   map[string]bool // for efficient lookups
}

// prefixedKey generates a cache key that includes the tag information
// Format: {prefix}:{tag1}_{tag2}_...:{original_key}
func (tc *TaggedCache) prefixedKey(key string) string {
	if len(tc.tagNames) == 0 {
		return key
	}

	tags := strings.Join(tc.tagNames, "_")
	if tc.cache.config.Prefix != "" {
		return fmt.Sprintf("%s:%s:%s", tc.cache.config.Prefix, tags, key)
	}
	return fmt.Sprintf("%s:%s", tags, key)
}

// Get retrieves an item from the cache by key
func (tc *TaggedCache) Get(key string) (any, bool) {
	prefixedKey := tc.prefixedKey(key)
	return tc.cache.Get(prefixedKey)
}

// Put stores an item in the cache for a given number of seconds
func (tc *TaggedCache) Put(key string, value any, ttl time.Duration) {
	prefixedKey := tc.prefixedKey(key)
	tc.cache.Put(prefixedKey, value, ttl)

	// Record the key in tag index so we can flush later
	tc.recordKeyInTagIndex(prefixedKey)
}

// Add stores an item in the cache if the key doesn't already exist
func (tc *TaggedCache) Add(key string, value any, ttl time.Duration) bool {
	prefixedKey := tc.prefixedKey(key)
	result := tc.cache.Add(prefixedKey, value, ttl)

	if result {
		// Record the key in tag index since it was added
		tc.recordKeyInTagIndex(prefixedKey)
	}

	return result
}

// Increment increments the value of an item in the cache
func (tc *TaggedCache) Increment(key string, n int64) (int64, error) {
	prefixedKey := tc.prefixedKey(key)
	return tc.cache.Increment(prefixedKey, n)
}

// Decrement decrements the value of an item in the cache
func (tc *TaggedCache) Decrement(key string, n int64) (int64, error) {
	prefixedKey := tc.prefixedKey(key)
	return tc.cache.Decrement(prefixedKey, n)
}

// Forever stores an item in the cache indefinitely
func (tc *TaggedCache) Forever(key string, value any) {
	prefixedKey := tc.prefixedKey(key)
	tc.cache.Forever(prefixedKey, value)

	// Record the key in tag index
	tc.recordKeyInTagIndex(prefixedKey)
}

// Forget removes an item from the cache
func (tc *TaggedCache) Forget(key string) bool {
	prefixedKey := tc.prefixedKey(key)
	return tc.cache.Forget(prefixedKey)
}

// Flush removes all items from the cache with these tags
func (tc *TaggedCache) Flush() error {
	// Get all keys associated with these tags and remove them
	for _, tagName := range tc.tagNames {
		tagIndexKey := tc.getTagIndexKey(tagName)
		if keys, exists := tc.cache.Get(tagIndexKey); exists {
			if keySlice, ok := keys.([]string); ok {
				for _, key := range keySlice {
					tc.cache.Forget(key)
				}
			}
		}

		// Clear the tag index itself
		tc.cache.Forget(tagIndexKey)
	}

	return nil
}

// Remember gets an item from the cache, or execute the callback and store the result
func (tc *TaggedCache) Remember(
	key string,
	ttl time.Duration,
	callback func() (any, error),
) (any, error) {
	prefixedKey := tc.prefixedKey(key)
	value, exists := tc.cache.Get(prefixedKey)
	if exists {
		return value, nil
	}

	callbackResult, err := callback()
	if err != nil {
		return nil, err
	}

	tc.cache.Put(prefixedKey, callbackResult, ttl)
	tc.recordKeyInTagIndex(prefixedKey)

	return callbackResult, nil
}

// RememberForever gets an item from the cache, or execute the callback and store the result forever
func (tc *TaggedCache) RememberForever(key string, callback func() (any, error)) (any, error) {
	prefixedKey := tc.prefixedKey(key)
	value, exists := tc.cache.Get(prefixedKey)
	if exists {
		return value, nil
	}

	callbackResult, err := callback()
	if err != nil {
		return nil, err
	}

	tc.cache.Forever(prefixedKey, callbackResult)
	tc.recordKeyInTagIndex(prefixedKey)

	return callbackResult, nil
}

// getTagIndexKey returns the cache key used to store the list of keys for a tag
func (tc *TaggedCache) getTagIndexKey(tagName string) string {
	if tc.cache.config.Prefix != "" {
		return fmt.Sprintf("%s:tag_index:%s", tc.cache.config.Prefix, tagName)
	}
	return fmt.Sprintf("tag_index:%s", tagName)
}

// recordKeyInTagIndex records the key in the tag index for later flushing
func (tc *TaggedCache) recordKeyInTagIndex(key string) {
	for _, tagName := range tc.tagNames {
		tagIndexKey := tc.getTagIndexKey(tagName)

		// Get existing keys for this tag
		existingKeys := []string{}
		if val, exists := tc.cache.Get(tagIndexKey); exists {
			if keys, ok := val.([]string); ok {
				existingKeys = keys
			}
		}

		// Add this key if it's not already in the list
		keyExists := slices.Contains(existingKeys, key)

		if !keyExists {
			newKeys := append(existingKeys, key)
			tc.cache.Forever(tagIndexKey, newKeys) // Store the updated list
		}
	}
}
