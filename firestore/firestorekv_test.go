package firestore

import (
	"context"
	"testing"
	"time"

	"github.com/duizendstra-com/dui-go/internal/testutil"
)

// TestMockFirestoreKV ensures that MockFirestoreKV behaves like a FirestoreKV in-memory.
// No external Firestore emulator or credentials are required for this test.
func TestMockFirestoreKV(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	kv := testutil.NewMockFirestoreKV() // Use the mock instead of real FirestoreKV
	defer kv.Close()

	// Test Set operation
	if err := kv.Set(ctx, "testKey", "testValue"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Test Get existing key
	val, err := kv.Get(ctx, "testKey")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "testValue" {
		t.Errorf("expected 'testValue', got %q", val)
	}

	// Test Get non-existent key
	val, err = kv.Get(ctx, "nonExistent")
	if err != nil {
		t.Fatalf("Get nonExistent failed: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty for nonExistent, got %q", val)
	}
}
