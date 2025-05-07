package store

import (
	"context"
	"fmt"

	"github.com/duizendstra/dui-go/firestore"
)

// kvInterface matches the methods we need from a Firestore-like KV.
// It allows us to inject either a real FirestoreKV or a mock implementation for testing.
type kvInterface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Close() error
}

// FirestoreStore implements the Store interface using a Firestore-based KV (from the firestore package).
// This allows storing key-value data in a Firestore collection without changing the store's interface.
type FirestoreStore struct {
	kv kvInterface
}

// NewFirestoreStore creates a Store implementation backed by Firestore.
// It uses firestore.NewKV to connect to a Firestore project and collection.
//
// Example usage:
//
//	ctx := context.Background()
//	s, err := NewFirestoreStore(ctx, "my-gcp-project", "my-collection")
//	if err != nil {
//	  log.Fatalf("failed to create FirestoreStore: %v", err)
//	}
//	defer s.Close()
//
//	if err := s.Set(ctx, "foo", "bar"); err != nil {
//	  log.Fatalf("Set failed: %v", err)
//	}
//
//	val, err := s.Get(ctx, "foo")
//	// ...
func NewFirestoreStore(ctx context.Context, projectID, collection string) (Store, error) {
	realKV, err := firestore.NewKV(ctx, projectID, collection)
	if err != nil {
		return nil, fmt.Errorf("failed to create FirestoreKV for store: %w", err)
	}
	return &FirestoreStore{kv: realKV}, nil
}

// Get retrieves the value for a given key from Firestore. Returns an empty string if not found.
func (s *FirestoreStore) Get(ctx context.Context, key string) (string, error) {
	return s.kv.Get(ctx, key)
}

// Set stores the value for a given key, overwriting any existing value.
func (s *FirestoreStore) Set(ctx context.Context, key, value string) error {
	return s.kv.Set(ctx, key, value)
}

// Close releases any resources associated with the Firestore store.
func (s *FirestoreStore) Close() error {
	return s.kv.Close()
}
