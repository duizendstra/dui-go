package authentication

import (
	"fmt"
	"time"

	"github.com/duizendstra-com/go-dui/cache"
)

// ExampleTokenManager_basic demonstrates creating and using a TokenManager.
func ExampleTokenManager_basic() {
	// Create a TokenManager with an in-memory cache.
	tm := NewTokenManager(cache.NewInMemoryCache())

	// Register a fetcher that returns a dynamic token valid for 5 seconds.
	tm.RegisterFetcher("my-service", func() (string, time.Time, error) {
		return "dynamic-token", time.Now().Add(5 * time.Second), nil
	})

	// First call => fetch from the fetcher
	token, err := tm.GetToken("my-service")
	if err != nil {
		fmt.Println("Error fetching token:", err)
		return
	}
	fmt.Println("Initial token:", token)

	// Second call => uses the cached token
	token2, err := tm.GetToken("my-service")
	if err != nil {
		fmt.Println("Error fetching second token:", err)
		return
	}
	fmt.Println("Cached token:", token2)

	// Output:
	// Initial token: dynamic-token
	// Cached token: dynamic-token
}

// ExampleTokenManager_error demonstrates error handling scenarios.
func ExampleTokenManager_error() {
	tm := NewTokenManager(cache.NewInMemoryCache())

	// Attempt to get a token for an unknown service => error
	_, err := tm.GetToken("unknown-service")
	if err != nil {
		fmt.Println("Expected error for unknown-service:", err)
	}

	// Register a fetcher that always fails
	tm.RegisterFetcher("failing-service", func() (string, time.Time, error) {
		return "", time.Time{}, fmt.Errorf("intentional failure")
	})

	_, err = tm.GetToken("failing-service")
	if err != nil {
		fmt.Println("Expected error for failing-service:", err)
	}

	// Output:
	// Expected error for unknown-service: no fetcher registered for key: unknown-service
	// Expected error for failing-service: failed to fetch token for key failing-service: intentional failure
}
