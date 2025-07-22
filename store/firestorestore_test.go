package store

import (
	"context"
	"testing"
	"time"

	"github.com/duizendstra/dui-go/testutil"
)

// TestFirestoreStoreWithMock demonstrates how to test the FirestoreStore using a mock
// instead of a real Firestore connection. By substituting testutil.NewMockFirestoreKV()
// we avoid dependencies on external services.
func TestFirestoreStoreWithMock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockKV := testutil.NewMockFirestoreKV() // mock Firestore-like KV from testutil
	store := &FirestoreStore{kv: mockKV}

	// Test Set
	if err := store.Set(ctx, "testKey", "testValue"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Test Get existing key
	val, err := store.Get(ctx, "testKey")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "testValue" {
		t.Errorf("expected 'testValue', got %q", val)
	}

	// Test Get non-existing key
	val, err = store.Get(ctx, "nonexistentKey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty string for nonexistentKey, got %q", val)
	}
}
