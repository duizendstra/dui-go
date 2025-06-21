package authentication

import (
	"errors"
	"testing"
	"time"

	"github.com/duizendstra/dui-go/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenManager(t *testing.T) {
	// Common setup for all sub-tests
	c := cache.NewInMemoryCache()
	tm := NewTokenManager(c)

	// This fetcher is used across multiple sub-tests.
	tm.RegisterFetcher("api-service", func() (string, time.Time, error) {
		return "fetched-token", time.Now().Add(time.Minute), nil
	})

	t.Run("it should call the fetcher when cache is empty", func(t *testing.T) {
		token, err := tm.GetToken("api-service")
		require.NoError(t, err)
		assert.Equal(t, "fetched-token", token)
	})

	t.Run("it should return the cached token on subsequent calls", func(t *testing.T) {
		// Ensure the token is already in the cache from the previous test or a new fetch.
		_, err := tm.GetToken("api-service")
		require.NoError(t, err)

		// This call should now hit the cache.
		token, err := tm.GetToken("api-service")
		require.NoError(t, err)
		assert.Equal(t, "fetched-token", token)
	})

	t.Run("it should re-fetch when the token is expired", func(t *testing.T) {
		// Manually set a token that is already expired.
		tm.SetToken("api-service", "manual-expired-token", time.Now().Add(-time.Second))

		// The test should now trigger a re-fetch, returning the fresh token.
		token, err := tm.GetToken("api-service")
		require.NoError(t, err)
		assert.Equal(t, "fetched-token", token)
	})

	t.Run("it should return an error for an unregistered service", func(t *testing.T) {
		_, err := tm.GetToken("unknown-service")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no fetcher registered for key: unknown-service")
	})

	t.Run("it should return an error when the fetcher fails", func(t *testing.T) {
		tm.RegisterFetcher("failing-service", func() (string, time.Time, error) {
			return "", time.Time{}, errors.New("fetch failed")
		})

		_, err := tm.GetToken("failing-service")
		require.Error(t, err)
		assert.Equal(t, "failed to fetch token for key failing-service: fetch failed", err.Error())
	})
}
