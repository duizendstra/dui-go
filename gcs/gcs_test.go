package gcs

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUpload_Validation tests the input validation for the Upload method.
// It is a pure unit test and does not require a live GCS client.
func TestUpload_Validation(t *testing.T) {
	// Create a dummy client for testing validation logic. The internal gcsClient
	// is nil because we only want to test the validation checks, which execute
	// before the client is used.
	client := &Client{gcsClient: nil}
	ctx := context.Background()
	dummyReader := strings.NewReader("some data")

	testCases := []struct {
		name        string
		bucket      string
		object      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Fails with empty bucket name",
			bucket:      "",
			object:      "my-object",
			expectError: true,
			errorMsg:    "bucket name cannot be empty",
		},
		{
			name:        "Fails with empty object name",
			bucket:      "my-bucket",
			object:      "",
			expectError: true,
			errorMsg:    "object name cannot be empty",
		},
		// CORRECTED: The success case is removed because it cannot be tested
		// in this unit test without causing a nil pointer dereference.
		// Testing a successful upload requires an integration test.
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Upload(ctx, tt.bucket, tt.object, dummyReader)

			// We only test the cases that are expected to fail validation.
			require.Error(t, err, "Expected a validation error but got nil")
			assert.Equal(t, tt.errorMsg, err.Error(), "Error message mismatch")
		})
	}
}

// Note: A complete test of the Upload method's success path would be an
// integration test, structured like this:
//
// func TestUpload_Integration(t *testing.T) {
//     // This test should be skipped unless an integration build tag is provided.
//     t.Skip("Skipping integration test for GCS Upload")
//
//     // 1. Initialize a real GCS client (or an emulator).
//     // 2. Define a test bucket and object name.
//     // 3. Call client.Upload().
//     // 4. Assert that the error returned is nil.
//     // 5. Optionally, download the object and verify its content.
//     // 6. Clean up the created object from the bucket.
// }
