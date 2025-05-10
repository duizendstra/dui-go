// Package authentication provides functionality for handling authentication tokens,
// simplifying their lifecycle management including retrieval, caching, and automatic refresh.
// It includes:
//   - Token: A simple type representing a token and its expiration time.
//   - TokenManager: A thread-safe component for retrieving and caching tokens from external sources.
//
// Typical usage:
//
//	import (
//	    "fmt"
//	    "log"
//	    "time"
//	    "github.com/duizendstra/dui-go/cache" // For providing a cache implementation
//	    "github.com/duizendstra/dui-go/authentication"
//	)
//
//	func main() {
//	    // Use a cache implementation, e.g., from dui-go/cache
//	    inMemCache := cache.NewInMemoryCache()
//	    tm := authentication.NewTokenManager(inMemCache)
//
//	    // Register a fetcher for a specific token key
//	    tm.RegisterFetcher("my-service-token", func() (string, time.Time, error) {
//	      // In a real scenario, fetch a fresh token from an identity provider or service
//	      // and return the token string, its server-side expiry time, and any error.
//	      return "fresh-secure-token-value", time.Now().Add(1 * time.Hour), nil
//	    })
//
//	    // Get the token; TokenManager handles caching and fetching
//	    token, err := tm.GetToken("my-service-token")
//	    if err != nil {
//	      log.Fatalf("failed to get token: %v", err)
//	    }
//	    fmt.Println("Acquired Token:", token)
//	}
package authentication
