package secretmanager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewClient_Validation tests the input validation of the NewClient constructor.
func TestNewClient_Validation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		projectID   string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Fails with empty project ID",
			projectID:   "",
			expectError: true,
			errorMsg:    "GCP ProjectID is required",
		},
		{
			name: "Succeeds with non-empty project ID",
			// Note: This will attempt a real connection and will fail if not authenticated.
			// The purpose of this test case is to ensure no error is returned for valid input,
			// even if the underlying connection fails later. The validation itself passes.
			projectID:   "a-valid-project-id",
			expectError: false, // We don't expect a validation error.
			errorMsg:    "",
		},
	}

	// CORRECTED: The loop now correctly uses `range tests`.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test intentionally ignores the actual client returned, as it might
			// be nil if authentication fails. We are only testing the validation error.
			_, err := NewClient(ctx, tt.projectID)

			if tt.expectError {
				require.Error(t, err, "Expected an error but got nil")
				assert.Contains(t, err.Error(), tt.errorMsg, "Error message mismatch")
			} else {
				// If we expect success, the validation itself should not fail.
				// The underlying GCP client creation might fail if not authenticated,
				// so we only check that our specific validation error is not present.
				if err != nil {
					assert.NotContains(t, err.Error(), "GCP ProjectID is required", "Should not have a validation error")
				}
			}
		})
	}
}

// TestGetSecret_Validation tests the input validation for the GetSecret method.
func TestGetSecret_Validation(t *testing.T) {
	// We can create a client that won't be used to make calls, just to test methods.
	client := &Client{
		projectID: "dummy-project",
		gcpClient: nil, // No real client needed for this validation test.
	}

	t.Run("Fails with empty secret ID", func(t *testing.T) {
		_, err := client.GetSecret(context.Background(), "")
		require.Error(t, err, "Expected an error for empty secret ID, but got nil")
		assert.Equal(t, "secretID cannot be empty", err.Error())
	})
}
