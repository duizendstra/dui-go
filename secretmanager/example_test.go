package secretmanager_test

import (
	"context"
	"fmt"
	"io"

	"log/slog"

	"github.com/duizendstra/dui-go/secretmanager"
)

func ExampleNewClient() {
	// This example demonstrates creating a client with a configuration struct.
	// In a real application, you would be authenticated with GCP and might
	// configure your logger to write to os.Stdout or another destination.

	ctx := context.Background()

	// CORRECTED: For this testable example, we create a logger that writes to
	// io.Discard. This allows us to demonstrate the configuration syntax
	// without producing log output that would interfere with the test's
	// expected output.
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := secretmanager.Config{
		ProjectID: "a-dummy-project-id",
		Logger:    silentLogger,
	}

	client, err := secretmanager.NewClient(ctx, cfg)
	if err != nil {
		// This error is expected in a standard test environment that lacks credentials.
		// Since we're not logging to stdout, the user won't see this log message
		// during a normal `go test` run, but it's still useful for debugging.
		slog.Error("Note: client creation failed as expected without credentials", "error", err)
		fmt.Println("Client creation was attempted.")
		return
	}
	defer client.Close()

	// This part would only run in a real, authenticated environment.
	fmt.Println("Successfully created Secret Manager client.")

	// The `// Output:` block now correctly matches the expected output when
	// credentials are not present. If they *were* present, the output would be
	// "Successfully created Secret Manager client." and the test would fail,
	// correctly indicating an unexpected environment setup.
	// Output:
	// Client creation was attempted.
}
