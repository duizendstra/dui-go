// Package cache provides a key-value caching abstraction for storing and
// retrieving data at runtime. It defines a Cache interface representing
// common cache operations and includes a simple, in-memory implementation.
//
// Common Use Cases:
//   - Caching frequently accessed values to improve performance
//   - Storing lightweight, ephemeral data that does not need persistence
//
// Key Features:
//   - Thread-safe operations for concurrent access
//   - Bulk set and retrieval methods for convenience
//   - A flush method to clear the entire cache at once
//
// This package returns a concrete in-memory cache type, and an interface
// that describes its usage. Consumers can rely on the Cache interface to
// abstract away the implementation details if desired. Additional caching
// backends (e.g. persistent or distributed) could be introduced later by
// providing new implementations that satisfy the Cache interface.
package cache
