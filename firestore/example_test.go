package firestore

import (
	"context"
	"fmt"
	"log"

	"github.com/duizendstra/dui-go/testutil"
)

// ExampleFirestoreKV demonstrates how to use a KV implementation. Here we use
// a MockFirestoreKV from the testutil package to avoid depending on a real Firestore.
// In production code, you would use NewKV(ctx, "your-project-id", "your-collection").
func ExampleFirestoreKV() {
	ctx := context.Background()

	// Use a mock KV for demonstration. In production, call NewKV(ctx, "your-project-id", "your-collection")
	kv := testutil.NewMockFirestoreKV()
	defer kv.Close()

	// Set a value
	if err := kv.Set(ctx, "demoKey", "demoValue"); err != nil {
		log.Fatalf("Set failed: %v", err)
	}

	// Get the value
	val, err := kv.Get(ctx, "demoKey")
	if err != nil {
		log.Fatalf("Get failed: %v", err)
	}
	fmt.Println("demoKey:", val)

	// Get a non-existent key returns empty string
	val, err = kv.Get(ctx, "missingKey")
	if err != nil {
		log.Fatalf("Get failed: %v", err)
	}
	fmt.Println("missingKey:", val)
}
