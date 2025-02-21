package env

import (
	"context"
	"fmt"
	"os"
)

// EnvVar represents a configuration for a single environment variable.
// - Key: the variable name.
// - Mandatory: if true, an error is returned if not set.
// - DefaultValue: used if optional and not set.
// - Validation: optional function to validate the value.
type EnvVar struct {
	Key          string
	Mandatory    bool
	DefaultValue string
	Validation   func(string) bool
}

// EnvLoader reads and validates environment variables. It no longer logs, only returns errors.
type EnvLoader struct{}

// NewEnvLoader constructs an EnvLoader. It does not require a logger.
func NewEnvLoader() *EnvLoader {
	return &EnvLoader{}
}

// handleEnvError returns an error about missing or invalid environment variables.
func (el *EnvLoader) handleEnvError(key string, message string) error {
	return fmt.Errorf(message, key)
}

// LoadEnv reads, validates, and returns environment variables based on the given configuration.
// Returns an error if a mandatory variable is missing or validation fails.
func (el *EnvLoader) LoadEnv(ctx context.Context, vars []EnvVar) (map[string]string, error) {
	envMap := make(map[string]string)

	for _, v := range vars {
		value, exists := os.LookupEnv(v.Key)

		if !exists && v.Mandatory {
			return nil, el.handleEnvError(v.Key,
				"Mandatory environment variable is missing: %s")
		}

		if !exists {
			value = v.DefaultValue
		}

		if v.Validation != nil && !v.Validation(value) {
			return nil, el.handleEnvError(v.Key,
				"Validation failed for environment variable: %s")
		}

		envMap[v.Key] = value
	}

	return envMap, nil
}
