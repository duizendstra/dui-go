package authentication

import "time"

// nowFunc returns the current time. In tests, we can replace it for stable results.
var nowFunc = time.Now

// Token represents a generic token with associated metadata, including its expiry time.
// IsExpired checks whether the token has passed its expiration time.
type Token struct {
	Value   string
	Expires time.Time
}

// IsExpired returns true if the current time is after the token's Expires time.
func (t Token) IsExpired() bool {
	return nowFunc().After(t.Expires)
}
