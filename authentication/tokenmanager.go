package authentication

import (
	"fmt"
	"sync"
	"time"

	"github.com/duizendstra-com/dui-go/cache"
)

// TokenFetcher is a function returning a new token and its expiry.
type TokenFetcher func() (string, time.Time, error)

// TokenManagerInterface defines the behavior for managing tokens.
type TokenManagerInterface interface {
	RegisterFetcher(key string, fetcher TokenFetcher)
	SetToken(key, token string, expiry time.Time)
	GetToken(key string) (string, error)
}

// cachedToken is stored in the cache.
type cachedToken struct {
	token  string
	expiry time.Time
}

// TokenManager manages tokens stored in a cache, refreshing them via fetchers when needed.
type TokenManager struct {
	mu       sync.Mutex
	c        cache.Cache
	fetchers map[string]TokenFetcher
}

// NewTokenManager returns a new TokenManager instance, storing tokens in the provided cache.
func NewTokenManager(c cache.Cache) *TokenManager {
	return &TokenManager{
		c:        c,
		fetchers: make(map[string]TokenFetcher),
	}
}

// RegisterFetcher associates a TokenFetcher with a given key. When GetToken sees a missing
// or expired token, it calls this fetcher to obtain a fresh one.
func (tm *TokenManager) RegisterFetcher(key string, fetcher TokenFetcher) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.fetchers[key] = fetcher
}

// SetToken manually stores a token and its expiry in the cache, bypassing the fetcher.
func (tm *TokenManager) SetToken(key, token string, expiry time.Time) {
	tm.c.Set(key, &cachedToken{
		token:  token,
		expiry: expiry,
	})
}

// GetToken retrieves a token for the key. If the token is present and not expired,
// it returns it. Otherwise, it fetches a new token from the registered TokenFetcher.
func (tm *TokenManager) GetToken(key string) (string, error) {
	val, ok := tm.c.Get(key)
	if ok {
		if ct, valid := val.(*cachedToken); valid && time.Now().Before(ct.expiry) {
			return ct.token, nil
		}
	}

	tm.mu.Lock()
	fetcher, hasFetcher := tm.fetchers[key]
	tm.mu.Unlock()

	if !hasFetcher {
		return "", fmt.Errorf("no fetcher registered for key: %s", key)
	}

	token, expiry, err := fetcher()
	if err != nil {
		return "", fmt.Errorf("failed to fetch token for key %s: %w", key, err)
	}

	tm.SetToken(key, token, expiry)
	return token, nil
}
