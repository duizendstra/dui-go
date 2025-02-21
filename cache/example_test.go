package cache

import (
	"fmt"
)

// ExampleInMemoryCache demonstrates basic usage of the InMemoryCache.
//
// To view this example, run:
//
//	go doc myapp/cache
//
// To run this example as a test, run:
//
//	go test myapp/cache -v
func ExampleInMemoryCache() {
	// Create a new in-memory cache
	c := NewInMemoryCache()

	// Set a single key-value pair
	c.Set("username", "alice")

	// Retrieve the value
	val, ok := c.Get("username")
	if ok {
		fmt.Println("username:", val)
	}

	// Set multiple key-value pairs at once
	c.SetAll(map[string]interface{}{
		"count":  42,
		"status": "active",
	})

	// Print all values currently in the cache
	all := c.GetAll()
	fmt.Println("All values:", all)

	// Flush the cache
	c.Flush()
	fmt.Println("All values after Flush:", c.GetAll())

	// Output:
	// username: alice
	// All values: map[count:42 status:active username:alice]
	// All values after Flush: map[]
}
