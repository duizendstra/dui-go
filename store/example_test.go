package store

import (
	"context"
	"fmt"
	"log"
)

// ExampleFirestoreStore demonstrates how to create and use a Firestore-backed Store.
// In practice, you would need valid GCP credentials and a Firestore collection.
// This example assumes that the Firestore environment is properly set up.
func ExampleFirestoreStore() {
	ctx := context.Background()
	projectID := "your-gcp-project-id" // Replace with your GCP project ID
	collection := "example-collection" // Replace with your Firestore collection name

	// Create the store
	s, err := NewFirestoreStore(ctx, projectID, collection)
	if err != nil {
		log.Fatalf("failed to create FirestoreStore: %v", err)
	}
	defer s.Close()

	// Set a value
	if err := s.Set(ctx, "myKey", "myValue"); err != nil {
		log.Fatalf("failed to set value: %v", err)
	}

	// Get the value
	value, err := s.Get(ctx, "myKey")
	if err != nil {
		log.Fatalf("failed to get value: %v", err)
	}

	fmt.Println("Retrieved value:", value)
}
