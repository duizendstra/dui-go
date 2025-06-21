// Package gcs provides a simplified client for interacting with
// Google Cloud Storage (GCS).
//
// It simplifies common operations like uploading objects and is designed to be
// used as a generic, reusable library in any Go application that needs to
// interact with GCS. The client is configured for a specific bucket upon
// creation.
//
// Usage:
//
//	import (
//		"context"
//		"log"
//		"log/slog"
//		"os"
//		"strings"
//
//		"github.com/duizendstra/dui-go/gcs"
//	)
//
//	func main() {
//	    ctx := context.Background()
//
//	    // Configuration for the client.
//	    cfg := gcs.Config{
//	        BucketName: "my-gcs-bucket",
//	        Logger:     slog.New(slog.NewJSONHandler(os.Stdout, nil)),
//	    }
//
//	    // Create a new client. It uses Application Default Credentials.
//	    client, err := gcs.NewClient(ctx, cfg)
//	    if err != nil {
//	        log.Fatalf("Failed to create GCS client: %v", err)
//	    }
//	    defer client.Close()
//
//	    // Prepare data to upload.
//	    data := strings.NewReader("This is the content of my file.")
//	    objectName := "path/to/my-object.txt"
//
//	    // Upload the data to the configured bucket.
//	    if err := client.Upload(ctx, objectName, data); err != nil {
//	        log.Fatalf("Failed to upload object: %v", err)
//	    }
//
//	    log.Println("Successfully uploaded object to GCS.")
//	}
//
// This package is thread-safe after initialization.
package gcs