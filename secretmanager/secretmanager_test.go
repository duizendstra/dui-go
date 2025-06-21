package secretmanager

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewClient_Validation tests the constructor's input validation using a Config struct.
func TestNewClient_Validation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		config      Config
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Fails with empty project ID",
			config:      Config{ProjectID: ""},
			expectError: true,
			errorMsg:    "GCP ProjectID is required in the config",
		},
		{
			name:        "Succeeds with non-empty project ID",
			config:      Config{ProjectID: "a-valid-project-id"},
			expectError: false,
			errorMsg:    "",
		},
		{
			name: "Succeeds with a logger provided",
			config: Config{
				ProjectID: "a-valid-project-id",
				Logger:    slog.Default(),
			},
			expectError: false,
			errorMsg:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(ctx, tt.config)

			if tt.expectError {
				require.Error(t, err, "Expected an error but got nil")
				assert.Contains(t, err.Error(), tt.errorMsg, "Error message mismatch")
			} else {
				if err != nil {
					assert.NotContains(t, err.Error(), "is required in the config", "Should not have a validation error")
				}
			}
		})
	}
}

// TestGetSecret_Validation remains the same, testing the method's own validation.
func TestGetSecret_Validation(t *testing.T) {
	client := &Client{
		projectID: "dummy-project",
		gcpClient: nil,
		logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	t.Run("Fails with empty secret ID", func(t *testing.T) {
		_, err := client.GetSecret(context.Background(), "")
		require.Error(t, err, "Expected an error for empty secret ID, but got nil")
		assert.Equal(t, "secretID cannot be empty", err.Error())
	})
}