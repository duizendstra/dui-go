package gcs

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewClient_Validation tests the constructor's input validation.
func TestNewClient_Validation(t *testing.T) {
	ctx := context.Background()

	t.Run("Fails with empty bucket name", func(t *testing.T) {
		cfg := Config{BucketName: ""}
		_, err := NewClient(ctx, cfg)
		require.Error(t, err, "Expected an error for empty bucket name but got nil")
		assert.Equal(t, "GCS BucketName is required in the config", err.Error())
	})

	t.Run("Succeeds with non-empty bucket name", func(t *testing.T) {
		// This will still fail if credentials are not available, but it will
		// pass the initial validation check, which is what we're testing.
		cfg := Config{BucketName: "a-valid-bucket"}
		_, err := NewClient(ctx, cfg)
		if err != nil {
			assert.NotContains(t, err.Error(), "is required in the config")
		}
	})
}

// TestUpload_Validation tests the input validation for the Upload method.
// It is a pure unit test and does not require a live GCS client.
func TestUpload_Validation(t *testing.T) {
	// Create a dummy client for testing validation logic.
	client := &Client{
		gcsClient:  nil,
		bucketName: "dummy-bucket", // The bucket is now part of the client.
		logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	ctx := context.Background()
	dummyReader := strings.NewReader("some data")

	t.Run("Fails with empty object name", func(t *testing.T) {
		err := client.Upload(ctx, "", dummyReader)
		require.Error(t, err, "Expected a validation error but got nil")
		assert.Equal(t, "object name cannot be empty", err.Error())
	})
}