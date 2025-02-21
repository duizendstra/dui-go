package env

import (
	"context"
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	loader := NewEnvLoader()
	ctx := context.Background()

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

	envMap, err := loader.LoadEnv(ctx, vars)
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
	os.Unsetenv("MANDATORY_VAR")
	varsMissingMandatory := []EnvVar{
		{Key: "MANDATORY_VAR", Mandatory: true},
	}

	_, err = loader.LoadEnv(ctx, varsMissingMandatory)
	if err == nil {
		t.Error("expected error for missing mandatory variable, got nil")
	}

	// Test validation failure
	t.Setenv("VALIDATED_VAR", "invalid_value") // does not start with "valid"
	varsValidationFail := []EnvVar{
		{Key: "VALIDATED_VAR", DefaultValue: "valid_value", Validation: validFunc},
	}

	_, err = loader.LoadEnv(ctx, varsValidationFail)
	if err == nil {
		t.Error("expected error for invalid VALIDATED_VAR value, got nil")
	} else {
		expected := "Validation failed for environment variable: VALIDATED_VAR"
		if msg := err.Error(); msg[:len(expected)] != expected {
			t.Errorf("expected error starting with %q, got %q", expected, msg)
		}
	}
}

func TestLoadEnv_NoVariables(t *testing.T) {
	loader := NewEnvLoader()
	ctx := context.Background()

	envMap, err := loader.LoadEnv(ctx, nil)
	if err != nil {
		t.Fatalf("expected no error with no vars, got %v", err)
	}
	if len(envMap) != 0 {
		t.Errorf("expected empty env map, got %d entries", len(envMap))
	}
}

func TestLoadEnv_OptionalDefault(t *testing.T) {
	loader := NewEnvLoader()
	ctx := context.Background()

	vars := []EnvVar{
		{Key: "NON_EXISTENT", DefaultValue: "default_value"},
	}
	envMap, err := loader.LoadEnv(ctx, vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := envMap["NON_EXISTENT"]; got != "default_value" {
		t.Errorf("expected NON_EXISTENT=default_value, got %q", got)
	}
}
