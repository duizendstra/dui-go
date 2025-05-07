// Package authentication provides functionality for handling authentication tokens.
// It includes:
//   - Token: A simple type representing a token and its expiration time.
//   - TokenManager: A component for retrieving and caching tokens from external sources.
//
// Typical usage:
//
//	import (
//	    "fmt"
//	    "log"
//	    "time"
//	    "github.com/duizendstra/dui-go/cache" // Added for clarity
//	    "github.com/duizendstra/dui-go/authentication"
//	)
//
//	func main() {
//	    tm := authentication.NewTokenManager(cache.NewInMemoryCache()) // Corrected
//	    tm.RegisterFetcher("my-service", func() (string, time.Time, error) {
//	      // Fetch a fresh token and return it along with its expiry.
//	      return "fresh-token", time.Now().Add(1 * time.Hour), nil
//	    })
//
//	    token, err := tm.GetToken("my-service")
//	    if err != nil {
//	      log.Fatalf("failed to get token: %v", err)
//	    }
//	    fmt.Println("Token:", token)
//	}
package authentication
