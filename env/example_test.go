package env_test // Use _test to ensure it's a black-box test

import (
	"fmt"
	"log"
	"os"

	"github.com/duizendstra/dui-go/env" // Adjust import path to your project
)

type AppConfig struct {
	Host        string `env:"APP_HOST" envDefault:"localhost"`
	Port        int    `env:"APP_PORT" envRequired:"true"`
	DebugMode   bool   `env:"APP_DEBUG_MODE" envDefault:"false"`
	ServiceName string `env:"SERVICE_NAME"` // Will look for SERVICE_NAME (uppercase field name)
	EmptyEnvVar string `env:"EMPTY_ENV_VAR"` // Env var set to empty, no default
	NotSetInt   int    `env:"NOT_SET_INT"`   // Go zero value expected
	NotSetBool  bool   `env:"NOT_SET_BOOL"`  // Go zero value expected
}

func ExampleProcess() {
	// Set up environment variables for the example
	os.Setenv("APP_PORT", "8080")
	os.Setenv("SERVICE_NAME", "MyAwesomeService")
	os.Setenv("APP_DEBUG_MODE", "true")
	os.Setenv("EMPTY_ENV_VAR", "") // Explicitly set to empty

	// Clean up environment variables after the test
	defer func() {
		os.Unsetenv("APP_HOST") // ensure it is not set if test runner has it
		os.Unsetenv("APP_PORT")
		os.Unsetenv("SERVICE_NAME")
		os.Unsetenv("APP_DEBUG_MODE")
		os.Unsetenv("EMPTY_ENV_VAR")
		os.Unsetenv("NOT_SET_INT")  // Ensure these are not set for the test
		os.Unsetenv("NOT_SET_BOOL")
	}()

	var cfg AppConfig
	err := env.Process(&cfg)
	if err != nil {
		log.Fatalf("Failed to process environment: %v", err)
	}

	fmt.Printf("Host: %s\n", cfg.Host)
	fmt.Printf("Port: %d\n", cfg.Port)
	fmt.Printf("Debug Mode: %t\n", cfg.DebugMode)
	fmt.Printf("Service Name: %s\n", cfg.ServiceName)
	fmt.Printf("Empty Env Var: '%s'\n", cfg.EmptyEnvVar)
	fmt.Printf("Not Set Int (Go zero value): %d\n", cfg.NotSetInt)
	fmt.Printf("Not Set Bool (Go zero value): %t\n", cfg.NotSetBool)

	// Output:
	// Host: localhost
	// Port: 8080
	// Debug Mode: true
	// Service Name: MyAwesomeService
	// Empty Env Var: ''
	// Not Set Int (Go zero value): 0
	// Not Set Bool (Go zero value): false
}

func ExampleProcess_requiredMissing() {
	// Ensure APP_REQUIRED_PORT is not set for this specific example run
	originalVal, originalExists := os.LookupEnv("APP_REQUIRED_PORT")
	os.Unsetenv("APP_REQUIRED_PORT")
	defer func() {
		if originalExists {
			os.Setenv("APP_REQUIRED_PORT", originalVal)
		} else {
			os.Unsetenv("APP_REQUIRED_PORT")
		}
	}()

	type Config struct {
		RequiredPort int `env:"APP_REQUIRED_PORT" envRequired:"true"`
	}

	var cfg Config
	err := env.Process(&cfg)
	if err != nil {
		fmt.Printf("Error: %v\n", err) // Output the error for demonstration
	} else {
		fmt.Println("No error, but expected one for missing required var.")
	}

	// Output:
	// Error: env: required environment variable APP_REQUIRED_PORT is not set or is empty
}
