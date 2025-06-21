// Package secretmanager provides a client for securely retrieving secrets
// from Google Cloud Secret Manager.
//
// It simplifies the interaction with the GCP API, providing a clean and focused
// interface for fetching the latest version of secrets. The client can be
// configured with an optional slog.Logger for structured logging.
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
//
//	    // Configuration for the client
//	    cfg := secretmanager.Config{
//	        ProjectID: os.Getenv("GCP_PROJECT_ID"),
//	        // Logger is optional. If nil, no logs will be produced.
//	        Logger:    slog.New(slog.NewJSONHandler(os.Stdout, nil)),
//	    }
//
//	    // Create a new client
//	    client, err := secretmanager.NewClient(ctx, cfg)
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