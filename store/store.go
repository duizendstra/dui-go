// store/store.go
package store

import "context"

// KV defines simple key-value operations.
// It's defined here in the consumer package (store) rather than in the producer package (firestore).
type KV interface {
	// Get retrieves the value for a given key. If not found, returns an empty string with no error.
	Get(ctx context.Context, key string) (string, error)
	// Set stores the value under the given key, overwriting existing values.
	Set(ctx context.Context, key, value string) error
	// Close releases resources held by the KV implementation.
	Close() error
}

// Store defines a generic interface for key-value data storage.
// Store implementations can use a KV to persist data.
type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Close() error
}
