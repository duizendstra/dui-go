// Package gcs provides a simplified client for interacting with
// Google Cloud Storage (GCS).
//
// It simplifies common operations like uploading objects and is designed to be
// used as a generic, reusable library in any Go application that needs to
// interact with GCS.
//
// Usage:
//
//	import (
//		"context"
//		"log"
//		"strings"
//
//		"github.com/your-org/goclient-gcs/gcs"
//	)
//
//	func main() {
//	    ctx := context.Background()
//
//	    // Create a new client. It uses Application Default Credentials.
//	    client, err := gcs.NewClient(ctx)
//	    if err != nil {
//	        log.Fatalf("Failed to create GCS client: %v", err)
//	    }
//	    defer client.Close()
//
//	    // Prepare data to upload. An io.Reader is used for flexibility.
//	    data := strings.NewReader("This is the content of my file.")
//	    bucketName := "my-gcs-bucket"
//	    objectName := "path/to/my-object.txt"
//
//	    // Upload the data
//	    if err := client.Upload(ctx, bucketName, objectName, data); err != nil {
//	        log.Fatalf("Failed to upload object: %v", err)
//	    }
//
//	    log.Println("Successfully uploaded object to GCS.")
//	}
//
// This package is thread-safe after initialization.
package gcs
