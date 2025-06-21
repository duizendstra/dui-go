package secretmanager

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// Config holds the configuration for the Secret Manager client.
type Config struct {
	ProjectID string
	// Logger is an optional structured logger. If nil, logging will be disabled.
	Logger *slog.Logger
}

// Client provides a simplified interface for interacting with Google Secret Manager.
type Client struct {
	gcpClient *secretmanager.Client
	projectID string
	logger    *slog.Logger
}

// NewClient creates a new, authenticated client for Google Secret Manager.
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("GCP ProjectID is required in the config")
	}

	// Use the provided logger, or default to a silent one if not provided.
	logger := cfg.Logger
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	logger.Info("Initializing Google Secret Manager client...")
	gcpClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create underlying secretmanager.Client: %w", err)
	}

	return &Client{
		gcpClient: gcpClient,
		projectID: cfg.ProjectID,
		logger:    logger,
	}, nil
}

// GetSecret fetches the latest version of a secret from Secret Manager.
func (c *Client) GetSecret(ctx context.Context, secretID string) (string, error) {
	if secretID == "" {
		return "", fmt.Errorf("secretID cannot be empty")
	}

	c.logger.Debug("Fetching secret from Google Secret Manager", "secret_id", secretID)
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", c.projectID, secretID)
	req := &secretmanagerpb.AccessSecretVersionRequest{Name: name}

	result, err := c.gcpClient.AccessSecretVersion(ctx, req)
	if err != nil {
		c.logger.Error("Failed to access secret version", "name", name, "error", err)
		return "", fmt.Errorf("failed to access secret version '%s': %w", name, err)
	}

	if result.Payload == nil || result.Payload.Data == nil {
		return "", fmt.Errorf("retrieved secret payload for '%s' is nil", name)
	}

	c.logger.Info("Successfully fetched secret", "secret_id", secretID)
	return string(result.Payload.Data), nil
}

// Close releases any resources held by the client.
func (c *Client) Close() error {
	return c.gcpClient.Close()
}