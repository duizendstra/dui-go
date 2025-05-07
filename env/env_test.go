package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnv(t *testing.T) {
	loader := NewEnvLoader()

	t.Setenv("MANDATORY_VAR", "mandatory_value")
	t.Setenv("OPTIONAL_VAR", "optional_value")

	// Define a validation function that only allows values that start with "valid"
	validFunc := func(val string) bool {
		return len(val) >= 5 && val[:5] == "valid"
	}

	vars := []EnvVar{
		{Key: "MANDATORY_VAR", Mandatory: true},
		{Key: "OPTIONAL_VAR", DefaultValue: "default_opt"},
		{Key: "MISSING_OPTIONAL", DefaultValue: "default_missing"},
		{Key: "VALIDATED_VAR", DefaultValue: "valid_value", Validation: validFunc},
	}

	t.Setenv("VALIDATED_VAR", "valid_data")

	envMap, err := loader.LoadEnv(vars) // ctx removed
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check mandatory var
	if got := envMap["MANDATORY_VAR"]; got != "mandatory_value" {
		t.Errorf("expected MANDATORY_VAR=mandatory_value, got %q", got)
	}

	// Check optional var
	if got := envMap["OPTIONAL_VAR"]; got != "optional_value" {
		t.Errorf("expected OPTIONAL_VAR=optional_value, got %q", got)
	}

	// Check missing optional var uses default
	if got := envMap["MISSING_OPTIONAL"]; got != "default_missing" {
		t.Errorf("expected MISSING_OPTIONAL=default_missing, got %q", got)
	}

	// Check validated var
	if got := envMap["VALIDATED_VAR"]; got != "valid_data" {
		t.Errorf("expected VALIDATED_VAR=valid_data, got %q", got)
	}

	// Test missing mandatory variable
	os.Unsetenv("MANDATORY_VAR") // Unset for this specific sub-test
	varsMissingMandatory := []EnvVar{
		{Key: "MANDATORY_VAR", Mandatory: true},
	}

	_, err = loader.LoadEnv(varsMissingMandatory) // ctx removed
	if err == nil {
		t.Error("expected error for missing mandatory variable, got nil")
	}
	t.Setenv("MANDATORY_VAR", "mandatory_value") // Reset for subsequent tests if any in other files

	// Test validation failure
	t.Setenv("VALIDATED_VAR", "invalid_value") // does not start with "valid"
	varsValidationFail := []EnvVar{
		{Key: "VALIDATED_VAR", DefaultValue: "valid_value", Validation: validFunc},
	}

	_, err = loader.LoadEnv(varsValidationFail) // ctx removed
	if err == nil {
		t.Error("expected error for invalid VALIDATED_VAR value, got nil")
	} else {
		assert.ErrorContains(t, err, "Validation failed for environment variable: VALIDATED_VAR")
	}
	t.Setenv("VALIDATED_VAR", "valid_data") // Reset for subsequent tests
}

func TestLoadEnv_NoVariables(t *testing.T) {
	loader := NewEnvLoader()

	envMap, err := loader.LoadEnv(nil) // ctx removed
	if err != nil {
		t.Fatalf("expected no error with no vars, got %v", err)
	}
	if len(envMap) != 0 {
		t.Errorf("expected empty env map, got %d entries", len(envMap))
	}
}

func TestLoadEnv_OptionalDefault(t *testing.T) {
	loader := NewEnvLoader()

	vars := []EnvVar{
		{Key: "NON_EXISTENT", DefaultValue: "default_value"},
	}
	envMap, err := loader.LoadEnv(vars) // ctx removed
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := envMap["NON_EXISTENT"]; got != "default_value" {
		t.Errorf("expected NON_EXISTENT=default_value, got %q", got)
	}
}
