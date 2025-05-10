// Package cache provides a key-value caching abstraction for storing and
// retrieving data at runtime. It defines a Cache interface representing
// common cache operations and includes a simple, in-memory implementation.
//
// Common Use Cases:
//   - Caching frequently accessed values to improve application performance.
//   - Storing lightweight, ephemeral data that does not require persistent storage.
//
// Key Features:
//   - Thread-safe operations for concurrent access in the provided InMemoryCache.
//   - Basic Get, Set, SetAll, GetAll, and Flush operations.
//   - A flexible Cache interface allowing for different backend implementations.
//
// Typical Usage:
//
//	import "github.com/duizendstra/dui-go/cache"
//	import "fmt"
//
//	func main() {
//	    c := cache.NewInMemoryCache()
//	    c.Set("myKey", "myValue")
//	    if val, ok := c.Get("myKey"); ok {
//	        fmt.Println(val) // Output: myValue
//	    }
//	}
//
// This package returns a concrete in-memory cache type (InMemoryCache), and the Cache
// interface that describes its usage. Consumers can rely on the Cache interface to
// abstract away implementation details if desired. Additional caching backends
// (e.g., persistent or distributed) could be introduced later by providing new
// implementations that satisfy the Cache interface.
//
// For testing code that depends on this cache.Cache interface, consider using the
// MockCache from the internal/testutil package.
package cache
