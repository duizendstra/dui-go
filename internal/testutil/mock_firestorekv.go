package testutil

import (
	"context"
	"sync"
)

// MockFirestoreKV is an in-memory mock that simulates the behavior of a FirestoreKV.
// It stores keys and values in a map and returns empty strings for missing keys,
// just like a FirestoreKV would if a document doesn't exist.
//
// This mock never simulates errors by default, but you can add logic to do so if needed.
type MockFirestoreKV struct {
	mu   sync.Mutex
	data map[string]string
}

// NewMockFirestoreKV creates a new MockFirestoreKV instance with an empty in-memory map.
func NewMockFirestoreKV() *MockFirestoreKV {
	return &MockFirestoreKV{data: make(map[string]string)}
}

// Get retrieves the value associated with the given key.
// If the key does not exist, returns an empty string and no error.
func (m *MockFirestoreKV) Get(ctx context.Context, key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	val := m.data[key]
	return val, nil
}

// Set stores the given value under the specified key.
func (m *MockFirestoreKV) Set(ctx context.Context, key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}

// Close is a no-op for MockFirestoreKV, present only to match the FirestoreKV interface.
func (m *MockFirestoreKV) Close() error {
	return nil
}
