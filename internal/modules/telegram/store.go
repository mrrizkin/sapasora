package telegram

import "sync"

type Store[T any] struct {
	mu sync.Mutex
	m  map[string]T
}

func NewStore[T any]() *Store[T] {
	return &Store[T]{
		m: make(map[string]T),
	}
}

func (s *Store[T]) Get(key string) (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.m[key]
	return v, ok
}

func (s *Store[T]) Set(key string, value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

func (s *Store[T]) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)
}

func (s *Store[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m = make(map[string]T)
}

func (s *Store[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.m)
}

func (s *Store[T]) Keys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}

func (s *Store[T]) Values() []T {
	s.mu.Lock()
	defer s.mu.Unlock()
	values := make([]T, 0, len(s.m))
	for _, v := range s.m {
		values = append(values, v)
	}
	return values
}

func (s *Store[T]) Has(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.m[key]
	return ok
}
