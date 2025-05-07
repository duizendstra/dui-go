package firestore

// import (
// 	"context"
// 	"testing"
// 	"time"
//
// 	"github.com/duizendstra/dui-go/internal/testutil" // Was used for the removed test
// )

// Note: The test for MockFirestoreKV that was previously in this file has been removed
// as it is redundant. The MockFirestoreKV is already tested in
// its own package: internal/testutil/mock_firestorekv_test.go.

// This file (firestore/firestorekv_test.go) is intended for integration tests
// of the actual firestore.FirestoreKV implementation against a Firestore emulator
// or a real Firestore instance (with appropriate build tags).
// Such tests are not yet implemented for v0.0.1.

// Example for future integration test structure:
/*
func TestFirestoreKV_Integration(t *testing.T) {
	// Requires FIREBASE_EMULATOR_HOST to be set, or other GCP credentials.
	// Add appropriate build tags (e.g., //go:build integration)

	// projectID := os.Getenv("TEST_FIRESTORE_PROJECT_ID")
	// if projectID == "" {
	// 	t.Skip("TEST_FIRESTORE_PROJECT_ID not set, skipping integration test")
	// }
	// collection := "test-kv-collection"

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// kv, err := NewKV(ctx, projectID, collection)
	// if err != nil {
	// 	t.Fatalf("Failed to create FirestoreKV for integration test: %v", err)
	// }
	// defer kv.Close()

	// t.Run("SetAndGet", func(t *testing.T) {
	// 	key := "integrationTestKey"
	// 	value := "integrationTestValue-" + time.Now().Format(time.RFC3339Nano)

	// 	err := kv.Set(ctx, key, value)
	// 	if err != nil {
	// 		t.Fatalf("Set failed: %v", err)
	// 	}

	// 	retrievedVal, err := kv.Get(ctx, key)
	// 	if err != nil {
	// 		t.Fatalf("Get failed: %v", err)
	// 	}
	// 	if retrievedVal != value {
	// 		t.Errorf("expected %q, got %q", value, retrievedVal)
	// 	}
	// })

	// t.Run("GetNonExistent", func(t *testing.T) {
	// 	key := "nonExistentKey-" + time.Now().Format(time.RFC3339Nano)
	// 	retrievedVal, err := kv.Get(ctx, key)
	// 	if err != nil {
	// 		t.Fatalf("Get for non-existent key failed: %v", err)
	// 	}
	// 	if retrievedVal != "" {
	// 		t.Errorf("expected empty string for non-existent key, got %q", retrievedVal)
	// 	}
	// })
}
*/
