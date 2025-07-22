package testutil

import (
	"sync"

	"github.com/duizendstra/dui-go/cache"
)

// Compile-time check that MockCache implements cache.Cache
var _ cache.Cache = (*MockCache)(nil)

// MockCache is a mock implementation of the cache.Cache interface, designed
// for testing. It records calls to its methods, storing keys and values in an
// internal map.
//
// It is safe to call these methods concurrently, but reading the recorded calls
// while concurrent operations are ongoing is not recommended.
type MockCache struct {
	mu       sync.Mutex
	data     map[string]interface{}
	GetCalls []string
	SetCalls []struct {
		Key   string
		Value interface{}
	}
	SetAllCalls []map[string]interface{}
	FlushCalls  int
}

// NewMockCache returns a new, empty MockCache instance. This mock can be used
// in tests to verify key-value operations that rely on the cache.Cache interface.
func NewMockCache() *MockCache {
	return &MockCache{
		data: make(map[string]interface{}),
	}
}

// Set records the call and stores a value under the given key.
func (m *MockCache) Set(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SetCalls = append(m.SetCalls, struct {
		Key   string
		Value interface{}
	}{Key: key, Value: value})
	m.data[key] = value
}

// Get records the call and retrieves the value associated with the given key.
func (m *MockCache) Get(key string) (interface{}, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.GetCalls = append(m.GetCalls, key)
	val, ok := m.data[key]
	return val, ok
}

// SetAll records the call and stores multiple key-value pairs at once.
func (m *MockCache) SetAll(values map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SetAllCalls = append(m.SetAllCalls, values)
	for k, v := range values {
		m.data[k] = v
	}
}

// GetAll returns a copy of all key-value pairs currently stored.
func (m *MockCache) GetAll() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	copyMap := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		copyMap[k] = v
	}
	return copyMap
}

// Flush records the call and removes all entries from the cache.
func (m *MockCache) Flush() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.FlushCalls++
	m.data = make(map[string]interface{})
}
