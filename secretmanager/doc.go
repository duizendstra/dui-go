// Package secretmanager provides a client for securely retrieving secrets
// from Google Cloud Secret Manager.
//
// It simplifies the interaction with the GCP API, providing a clean and focused
// interface for fetching the latest version of secrets.
//
// Usage:
//
//	import (
//		"context"
//		"fmt"
//		"log"
//		"os"
//
//		"github.com/your-org/goclient-secretmanager/secretmanager"
//	)
//
//	func main() {
//	    ctx := context.Background()
//	    projectID := os.Getenv("GCP_PROJECT_ID")
//
//	    // Create a new client
//	    client, err := secretmanager.NewClient(ctx, projectID)
//	    if err != nil {
//	        log.Fatalf("Failed to create secret manager client: %v", err)
//	    }
//	    defer client.Close()
//
//	    // Fetch a secret
//	    secretValue, err := client.GetSecret(ctx, "my-app-db-password")
//	    if err != nil {
//	        log.Fatalf("Failed to retrieve secret: %v", err)
//	    }
//
//	    fmt.Println("Successfully retrieved secret.")
//	}
//
// This package is thread-safe after initialization.
package secretmanager
