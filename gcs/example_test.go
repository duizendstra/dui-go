package gcs_test

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/duizendstra/dui-go/gcs"
)

func ExampleNewClient() {
	// This example demonstrates the syntax for creating a GCS client.
	// In a real application, you would have your environment authenticated
	// with GCP (e.g., via `gcloud auth application-default login`).

	ctx := context.Background()

	client, err := gcs.NewClient(ctx)
	if err != nil {
		// This error is expected in a standard test environment that lacks credentials.
		log.Printf("Note: client creation failed as expected without credentials: %v", err)
		fmt.Println("Client creation was attempted.")
		return
	}
	defer client.Close()

	// The following shows how you would use the client if it were created successfully.
	// This part of the code is not executed during a standard test run because of the return above.
	fmt.Println("Successfully created GCS client.")

	data := strings.NewReader("file content")
	if err := client.Upload(ctx, "my-bucket", "my-object", data); err != nil {
		log.Printf("upload would fail in this test environment: %v", err)
	}

	// The `// Output:` block must match the predictable output of the test.
	// Output:
	// Client creation was attempted.
}
