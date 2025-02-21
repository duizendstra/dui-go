package authentication

import (
	"errors"
	"testing"
	"time"

	"github.com/duizendstra-com/go-dui/cache"
)

func TestTokenManager(t *testing.T) {
	// Use an in-memory cache for simplicity
	c := cache.NewInMemoryCache()
	tm := NewTokenManager(c)

	// 1) Register a fetcher returning a static token
	tm.RegisterFetcher("api-service", func() (string, time.Time, error) {
		return "fetched-token", time.Now().Add(time.Minute), nil
	})

	// 2) Initially, no token => fetcher is called
	token, err := tm.GetToken("api-service")
	if err != nil {
		t.Fatalf("unexpected error fetching token: %v", err)
	}
	if token != "fetched-token" {
		t.Errorf("expected 'fetched-token', got %q", token)
	}

	// 3) Second call => cached token
	token2, err := tm.GetToken("api-service")
	if err != nil {
		t.Fatalf("unexpected error second time: %v", err)
	}
	if token2 != "fetched-token" {
		t.Errorf("expected 'fetched-token', got %q", token2)
	}

	// 4) Manually set a token that expires now => force re-fetch
	tm.SetToken("api-service", "manual-token", time.Now())
	time.Sleep(10 * time.Millisecond) // ensure it is expired

	token3, err := tm.GetToken("api-service")
	if err != nil {
		t.Fatalf("unexpected error after expired token: %v", err)
	}
	if token3 != "fetched-token" {
		t.Errorf("expected 'fetched-token' after re-fetch, got %q", token3)
	}

	// 5) Unknown service => error
	_, err = tm.GetToken("unknown-service")
	if err == nil {
		t.Fatal("expected error for unknown-service")
	}

	// 6) Failing fetcher => error
	tm.RegisterFetcher("failing-service", func() (string, time.Time, error) {
		return "", time.Time{}, errors.New("fetch failed")
	})

	_, err = tm.GetToken("failing-service")
	if err == nil || err.Error() != "failed to fetch token for key failing-service: fetch failed" {
		t.Errorf("expected 'fetch failed' error, got %v", err)
	}
}
