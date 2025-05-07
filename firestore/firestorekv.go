package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// KV represents simple key-value operations. Implementations should handle non-existent keys
// by returning an empty string and no error from Get.
type KV interface {
	// Get retrieves the value associated with a key. If the key does not exist,
	// it returns an empty string and no error.
	Get(ctx context.Context, key string) (string, error)

	// Set stores the value under the given key.
	Set(ctx context.Context, key, value string) error

	// Close releases any resources associated with this KV implementation.
	Close() error
}

// FirestoreKV provides a key-value abstraction using a Firestore collection.
// Documents are stored with a "value" field. Missing documents return empty strings from Get.
type FirestoreKV struct {
	client     *firestore.Client
	collection string
}

// NewKV creates a FirestoreKV instance using the specified projectID and collection.
// It returns an error if the Firestore client cannot be created.
func NewKV(ctx context.Context, projectID, collection string) (*FirestoreKV, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create Firestore client: %w", err)
	}
	return &FirestoreKV{client: client, collection: collection}, nil
}

// Get retrieves the value for the given key from Firestore. If the key does not exist,
// it returns an empty string and no error.
func (f *FirestoreKV) Get(ctx context.Context, key string) (string, error) {
	docRef := f.client.Collection(f.collection).Doc(key)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Key does not exist, return empty string and no error.
			return "", nil
		}
		return "", fmt.Errorf("firestore get error (key=%s): %w", key, err)
	}

	var data map[string]interface{}
	if err := docSnap.DataTo(&data); err != nil {
		return "", fmt.Errorf("firestore data decode error (key=%s): %w", key, err)
	}

	value, _ := data["value"].(string)
	return value, nil
}

// Set writes a value at the given key in Firestore. Overwrites existing values.
func (f *FirestoreKV) Set(ctx context.Context, key, value string) error {
	docRef := f.client.Collection(f.collection).Doc(key)
	_, err := docRef.Set(ctx, map[string]interface{}{"value": value}, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("firestore set error (key=%s): %w", key, err)
	}
	return nil
}

// Close releases Firestore resources. After calling Close, the FirestoreKV should no longer
// be used.
func (f *FirestoreKV) Close() error {
	return f.client.Close()
}
