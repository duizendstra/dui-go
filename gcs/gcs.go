package gcs

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"cloud.google.com/go/storage"
)

// Config holds the configuration for the GCS client.
type Config struct {
	BucketName string
	// Logger is an optional structured logger. If nil, logging is disabled.
	Logger *slog.Logger
}

// Client provides a simplified interface for Google Cloud Storage.
// It is configured to work with a specific bucket upon creation.
type Client struct {
	gcsClient  *storage.Client
	bucketName string
	logger     *slog.Logger
}

// NewClient creates a new, authenticated client for GCS using the provided
// configuration, preparing it to work with the specified bucket.
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.BucketName == "" {
		return nil, fmt.Errorf("GCS BucketName is required in the config")
	}

	logger := cfg.Logger
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create underlying storage.Client: %w", err)
	}

	logger.InfoContext(ctx, "Initialized Google Cloud Storage client", "bucket", cfg.BucketName)
	return &Client{
		gcsClient:  gcsClient,
		bucketName: cfg.BucketName,
		logger:     logger,
	}, nil
}

// Upload streams data from an io.Reader to a new object in the configured GCS bucket.
// It is the caller's responsibility to handle the closing of the reader.
func (c *Client) Upload(ctx context.Context, object string, r io.Reader) (err error) {
	if object == "" {
		return fmt.Errorf("object name cannot be empty")
	}

	c.logger.DebugContext(ctx, "Uploading object to GCS", "object", object, "bucket", c.bucketName)

	writer := c.gcsClient.Bucket(c.bucketName).Object(object).NewWriter(ctx)

	defer func() {
		if closeErr := writer.Close(); closeErr != nil {
			if err == nil {
				err = fmt.Errorf("failed to close GCS writer for object '%s': %w", object, closeErr)
			}
		}
	}()

	if _, copyErr := io.Copy(writer, r); copyErr != nil {
		return fmt.Errorf("failed to copy data to GCS object '%s': %w", object, copyErr)
	}

	c.logger.DebugContext(ctx, "Successfully uploaded object", "object", object, "bucket", c.bucketName)
	return nil
}

// Close releases any resources held by the client. It should be called when
// the client is no longer needed.
func (c *Client) Close() error {
	return c.gcsClient.Close()
}