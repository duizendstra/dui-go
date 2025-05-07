package testutil

import (
	"context"
	"testing"
)

func TestMockFirestoreKV(t *testing.T) {
	mkv := NewMockFirestoreKV()

	ctx := context.Background()
	if err := mkv.Set(ctx, "testKey", "testValue"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := mkv.Get(ctx, "testKey")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "testValue" {
		t.Errorf("expected 'testValue', got %q", val)
	}

	val, err = mkv.Get(ctx, "nonExistent")
	if err != nil {
		t.Fatalf("Get nonExistent failed: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty for nonExistent, got %q", val)
	}
}
