package gcs_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/duizendstra/dui-go/gcs"
)

func ExampleNewClient() {
	// This example demonstrates creating a client. In a real application,
	// you would be authenticated with GCP.
	ctx := context.Background()

	// Configure the client. The logger is set to discard output for clean test runs.
	cfg := gcs.Config{
		BucketName: "my-test-bucket",
		Logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	client, err := gcs.NewClient(ctx, cfg)
	if err != nil {
		// This error is expected in a standard test environment that lacks credentials.
		log.Printf("Note: client creation failed as expected without credentials: %v", err)
		fmt.Println("Client creation was attempted.")
		return
	}
	defer client.Close()

	// The following shows how you would use the client if it were created successfully.
	// This part is not executed during a standard test run because of the return above.
	fmt.Println("Successfully created GCS client.")

	data := strings.NewReader("file content")
	if err := client.Upload(ctx, "my-object", data); err != nil {
		log.Printf("upload would fail in this test environment: %v", err)
	}

	// Output:
	// Client creation was attempted.
}