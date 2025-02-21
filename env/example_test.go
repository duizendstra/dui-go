package env

import (
	"context"
	"fmt"
	"os"
)

func ExampleEnvLoader_LoadEnv() {
	// Set environment variables for mandatory fields so the example passes.
	os.Setenv("API_KEY", "secret-api-key")
	os.Setenv("DATABASE_DSN", "user:pass@tcp(localhost:3306)/dbname")

	// Define the list of environment variables to load.
	vars := []EnvVar{
		{Key: "API_KEY", Mandatory: true},                           // Mandatory
		{Key: "DATABASE_DSN", Mandatory: true},                      // Mandatory
		{Key: "PORT", Mandatory: false, DefaultValue: "8080"},       // Optional with default
		{Key: "ENV", Mandatory: false, DefaultValue: "development"}, // Optional with default
	}

	// Create an EnvLoader without a logger.
	loader := NewEnvLoader()

	// Load the environment variables based on the defined configuration.
	ctx := context.Background()
	envMap, err := loader.LoadEnv(ctx, vars)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the successfully loaded environment variables.
	fmt.Println("Loaded env:", envMap)

	// Output:
	// Loaded env: map[API_KEY:secret-api-key DATABASE_DSN:user:pass@tcp(localhost:3306)/dbname ENV:development PORT:8080]
}
