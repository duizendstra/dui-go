package gcs

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

// Client provides a simplified interface for Google Cloud Storage.
type Client struct {
	gcsClient *storage.Client
}

// NewClient creates a new, authenticated client for GCS using Application
// Default Credentials.
func NewClient(ctx context.Context) (*Client, error) {
	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create underlying storage.Client: %w", err)
	}
	return &Client{gcsClient: gcsClient}, nil
}

// Upload streams data from an io.Reader to a new object in GCS.
// It is the caller's responsibility to handle the closing of the reader.
// This version uses a deferred Close for more robust resource management and
// named return values to correctly capture errors.
func (c *Client) Upload(ctx context.Context, bucket, object string, r io.Reader) (err error) {
	if bucket == "" {
		return fmt.Errorf("bucket name cannot be empty")
	}
	if object == "" {
		return fmt.Errorf("object name cannot be empty")
	}

	// Get a writer for the GCS object.
	writer := c.gcsClient.Bucket(bucket).Object(object).NewWriter(ctx)

	// Defer the Close call to ensure it runs even if io.Copy panics or returns an error.
	// This pattern ensures that the error from Close(), which is the definitive
	// success/fail signal for the upload, is properly handled and returned.
	defer func() {
		if closeErr := writer.Close(); closeErr != nil {
			// If we haven't already captured an error from io.Copy, use the Close error.
			// This prevents the more significant Close error from being masked.
			if err == nil {
				err = fmt.Errorf("failed to close GCS writer for object '%s': %w", object, closeErr)
			}
		}
	}()

	// Stream the data from the reader to the GCS writer.
	if _, copyErr := io.Copy(writer, r); copyErr != nil {
		return fmt.Errorf("failed to copy data to GCS object '%s': %w", object, copyErr)
	}

	// If io.Copy succeeds, the final result depends on writer.Close(), which is handled by the defer.
	return nil
}

// Close releases any resources held by the client. It should be called when
// the client is no longer needed.
func (c *Client) Close() error {
	return c.gcsClient.Close()
}
