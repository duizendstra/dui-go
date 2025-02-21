package cache

import "sync"

// InMemoryCache provides an in-memory implementation of the Cache interface.
// It stores data in a simple map protected by a mutex, ensuring safe concurrent
// access. It does not support expiration, persistence, or advanced features.
//
// This type is suitable for scenarios where cached data is small, does not
// need to persist between application restarts, and is frequently updated.
//
// Example:
//
//	c := NewInMemoryCache()
//	c.Set("foo", "bar")
//	if val, ok := c.Get("foo"); ok {
//	    fmt.Println(val) // prints "bar"
//	}
type InMemoryCache struct {
	mu   sync.Mutex
	data map[string]interface{}
}

// NewInMemoryCache returns a new, empty InMemoryCache instance.
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]interface{}),
	}
}

// Set associates a value with the given key, overwriting any existing value.
func (c *InMemoryCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

// Get retrieves the value associated with the given key.
// If the key exists, it returns (value, true). Otherwise, (nil, false).
func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.data[key]
	return val, ok
}

// SetAll stores multiple key-value pairs at once. Existing keys are overwritten.
func (c *InMemoryCache) SetAll(values map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range values {
		c.data[k] = v
	}
}

// GetAll returns a snapshot of all current key-value pairs in the cache.
// The returned map is a copy. Modifying it does not change the underlying cache.
func (c *InMemoryCache) GetAll() map[string]interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	copyMap := make(map[string]interface{}, len(c.data))
	for k, v := range c.data {
		copyMap[k] = v
	}
	return copyMap
}

// Flush removes all entries from the cache, leaving it empty.
func (c *InMemoryCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]interface{})
}
