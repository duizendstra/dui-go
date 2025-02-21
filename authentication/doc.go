// Package authentication provides functionality for handling authentication tokens.
// It includes:
//   - Token: A simple type representing a token and its expiration time.
//   - TokenManager: A component for retrieving and caching tokens from external sources.
//
// Typical usage:
//
//	tm := NewTokenManager(cache.NewInMemoryCache())
//	tm.RegisterFetcher("my-service", func() (string, time.Time, error) {
//	  // Fetch a fresh token and return it along with its expiry.
//	  return "fresh-token", time.Now().Add(1*time.Hour), nil
//	})
//
//	token, err := tm.GetToken("my-service")
//	if err != nil {
//	  log.Fatalf("failed to get token: %v", err)
//	}
//	fmt.Println("Token:", token)
//
// # Error Handling
//
// If no fetcher is registered for a given key or if the fetcher fails, GetToken returns an error.
// This allows you to gracefully handle missing or invalid tokens at runtime.
//
// # Testing
//
// This package does not specify how fetchers should obtain tokens, allowing you to integrate
// your own logic or use provided integrations like EasyflorTokenFetcher. For testing, you can use
// a mocking approach by depending on the TokenManagerInterface or by overriding fetchers with
// your own test logic.
package authentication
