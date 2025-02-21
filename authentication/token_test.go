package authentication

import (
	"testing"
	"time"
)

func TestTokenIsExpired(t *testing.T) {
	fixedTime := time.Unix(10000, 0) // arbitrary stable reference time

	cases := []struct {
		name     string
		token    Token
		expected bool
	}{
		{
			name: "Not expired if expiry is in the future",
			token: Token{
				Value:   "valid-token",
				Expires: fixedTime.Add(time.Minute),
			},
			expected: false,
		},
		{
			name: "Expired if expiry is in the past",
			token: Token{
				Value:   "old-token",
				Expires: fixedTime.Add(-time.Minute),
			},
			expected: true,
		},
		{
			name: "Not expired if expiry equals current time",
			token: Token{
				Value:   "just-expired",
				Expires: fixedTime,
			},
			expected: false,
		},
	}

	originalNow := nowFunc
	nowFunc = func() time.Time { return fixedTime }
	defer func() { nowFunc = originalNow }()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.token.IsExpired()
			if got != c.expected {
				t.Errorf("expected %v, got %v for token %q", c.expected, got, c.token.Value)
			}
		})
	}
}
