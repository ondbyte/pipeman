package nodd

import "sync"

// SyncMap is a concurrency-safe typed map
type SyncMap[K comparable, V any] struct {
	m sync.Map
}

func (sm *SyncMap[K, V]) Store(key K, value V) {
	sm.m.Store(key, value)
}

func (sm *SyncMap[K, V]) Load(key K) (V, bool) {
	val, ok := sm.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return val.(V), true
}

func (sm *SyncMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	actual, loaded := sm.m.LoadOrStore(key, value)
	return actual.(V), loaded
}

func (sm *SyncMap[K, V]) Delete(key K) {
	sm.m.Delete(key)
}

func (sm *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	sm.m.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}

// SyncSlice is a concurrency-safe slice
type SyncSlice[T any] struct {
	mu    *sync.RWMutex
	slice []T
}

func NewSyncSlice[T any](entries ...T) *SyncSlice[T] {
	return &SyncSlice[T]{
		slice: entries,
		mu:    &sync.RWMutex{},
	}
}

func (s *SyncSlice[T]) Append(val T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slice = append(s.slice, val)
}

func (s *SyncSlice[T]) Get(index int) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if index < 0 || index >= len(s.slice) {
		var zero T
		return zero, false
	}
	return s.slice[index], true
}

func (s *SyncSlice[T]) Set(index int, val T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.slice) {
		return false
	}
	s.slice[index] = val
	return true
}

func (s *SyncSlice[T]) Len() int {
	if s == nil {
		return 0
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.slice)
}

func (s *SyncSlice[T]) Slice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copySlice := make([]T, len(s.slice))
	copy(copySlice, s.slice)
	return copySlice
}

func (s *SyncSlice[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slice = nil
}
