// Package arr provides utilities for collections
package arr

func List[T any](items ...T) []T {
	return items
}

func Keys[T comparable, U any](collection map[T]U) []T {
	keys := make([]T, 0, len(collection))
	for key := range collection {
		keys = append(keys, key)
	}
	return keys
}

func Merge[T any](collections ...[]T) []T {
	result := make([]T, 0)
	for _, collection := range collections {
		result = append(result, collection...)
	}
	return result
}

// Pluck extracts values of type V from a slice of type T using a key function
func Pluck[T any, V any](collection []T, getKey func(T) V) []V {
	result := make([]V, 0, len(collection))
	for _, item := range collection {
		result = append(result, getKey(item))
	}
	return result
}

// KeyBy converts a slice of type T into a map where the key is of type K
// If multiple items have the same key, the last one will be kept
func KeyBy[T any, K comparable](collection []T, getKey func(T) K) map[K]T {
	result := make(map[K]T)
	for _, item := range collection {
		key := getKey(item)
		result[key] = item
	}
	return result
}

// GroupBy converts a slice of type T into a map where the key is of type K
// and the value is a slice of all items with that key
func GroupBy[T any, K comparable](collection []T, getKey func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range collection {
		key := getKey(item)
		result[key] = append(result[key], item)
	}
	return result
}

// Filter returns a new slice containing only the elements that satisfy the predicate
func Filter[T any](collection []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range collection {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map applies a transformation function to each element in the slice
func Map[T any, R any](collection []T, transform func(T) R) []R {
	result := make([]R, len(collection))
	for i, item := range collection {
		result[i] = transform(item)
	}
	return result
}

// Reduce reduces a slice to a single value using an accumulator function
func Reduce[T any, R any](collection []T, initial R, accumulator func(R, T) R) R {
	result := initial
	for _, item := range collection {
		result = accumulator(result, item)
	}
	return result
}

// First returns the first element that satisfies the predicate, along with a boolean indicating if an element was found
func First[T any](collection []T, predicate func(T) bool) (T, bool) {
	for _, item := range collection {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}
