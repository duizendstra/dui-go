package secretmanager

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// Client provides a simplified interface for interacting with Google Secret Manager.
type Client struct {
	gcpClient *secretmanager.Client
	projectID string
}

// NewClient creates a new, authenticated client for Google Secret Manager.
// It requires a context for initialization and the GCP Project ID where secrets are stored.
func NewClient(ctx context.Context, projectID string) (*Client, error) {
	if projectID == "" {
		return nil, fmt.Errorf("GCP ProjectID is required")
	}

	gcpClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create underlying secretmanager.Client: %w", err)
	}

	return &Client{
		gcpClient: gcpClient,
		projectID: projectID,
	}, nil
}

// GetSecret fetches the latest version of a secret from Secret Manager.
// The secretID is the short name of the secret (e.g., "my-api-key").
func (c *Client) GetSecret(ctx context.Context, secretID string) (string, error) {
	if secretID == "" {
		return "", fmt.Errorf("secretID cannot be empty")
	}

	// Build the full resource name required by the API.
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", c.projectID, secretID)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := c.gcpClient.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version '%s': %w", name, err)
	}

	if result.Payload == nil || result.Payload.Data == nil {
		return "", fmt.Errorf("retrieved secret payload for '%s' is nil", name)
	}

	return string(result.Payload.Data), nil
}

// Close releases any resources held by the client. It should be called when
// the client is no longer needed.
func (c *Client) Close() error {
	return c.gcpClient.Close()
}
