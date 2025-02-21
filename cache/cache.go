package cache

// Cache defines a key-value cache interface.
//
// Implementations of Cache should handle concurrency safely if they expect
// concurrent access, and should document any specific behaviors or limitations
// regarding data retention, expiration, or storage size.
type Cache interface {
	// Set associates a value with the given key, overwriting any existing value.
	Set(key string, value interface{})

	// Get retrieves the value associated with the given key.
	// If the key exists, it returns (value, true).
	// Otherwise, it returns (nil, false).
	Get(key string) (interface{}, bool)

	// SetAll stores multiple key-value pairs at once. Existing keys are overwritten.
	SetAll(values map[string]interface{})

	// GetAll returns a snapshot of all current key-value pairs in the cache.
	// The returned map is a copy, so modifications to it do not affect the
	// underlying cache. Implementations should ensure this operation is safe
	// to call concurrently with others.
	GetAll() map[string]interface{}

	// Flush removes all entries from the cache, leaving it empty.
	Flush()
}
