package store

import (
	"context"
	"fmt"
	"log"

	"github.com/duizendstra/dui-go/internal/testutil" // Required for the mock KV
)

// ExampleFirestoreStore demonstrates how to create and use a Firestore-backed Store.
//
// For actual use, NewFirestoreStore requires valid GCP credentials and a Firestore
// collection. The projectID and collection variables below would need to be set to your
// actual GCP project and Firestore collection.
//
// This runnable example uses a mock Key-Value store internally for demonstration purposes
// to allow `go test` to pass without real GCP credentials.
func ExampleFirestoreStore() {
	ctx := context.Background()
	// projectID := "your-gcp-project-id" // Replace with your GCP project ID for actual use
	// collection := "example-collection" // Replace with your Firestore collection name for actual use

	// How you would typically create the store for production:
	// s, err := NewFirestoreStore(ctx, projectID, collection)
	// if err != nil {
	//	 log.Fatalf("failed to create FirestoreStore for production: %v", err)
	// }
	// defer s.Close() // Ensure to close the production store

	// For this example to be runnable without GCP, we'll demonstrate
	// the Store's methods using a FirestoreStore with a mock KV backend.
	// In a real application, you would use the 's' from NewFirestoreStore above.
	mockKv := testutil.NewMockFirestoreKV()
	exampleStore := &FirestoreStore{kv: mockKv} // Manually create FirestoreStore with the mock KV

	// Set a value using the example store
	if err := exampleStore.Set(ctx, "myKey", "myValue"); err != nil {
		log.Fatalf("failed to set value using example store: %v", err)
	}

	// Get the value using the example store
	value, err := exampleStore.Get(ctx, "myKey")
	if err != nil {
		log.Fatalf("failed to get value using example store: %v", err)
	}

	fmt.Println("Retrieved value:", value)

	// Remember to close the store when done.
	// If 's' was a real store, you'd call s.Close().
	// For our exampleStore with a mock, Close is a no-op but good practice to show.
	if err := exampleStore.Close(); err != nil {
		log.Printf("Error closing example store: %v", err)
	}

	// Output:
	// Retrieved value: myValue
}
