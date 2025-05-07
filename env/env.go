package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Process populates the fields of the struct pointed to by 'spec'
// using environment variables.
//
// Fields in the struct should be tagged with:
//   - `env:"ENV_VAR_NAME"`: Specifies the environment variable name.
//     If not provided, the uppercase field name is used. An empty tag (`env:""`)
//     also defaults to the uppercase field name.
//   - `envDefault:"value"`: Specifies the default value to use if the
//     environment variable is not set or is an empty string.
//     The default value is treated as a string and will be parsed
//     according to the field's type.
//   - `envRequired:"true"`: Specifies that the environment variable
//     must be set and non-empty. An error is returned if a required
//     variable is missing or empty. If a field is required, `envDefault`
//     is effectively ignored for that variable (as it *must* be provided from the environment).
//
// Supported field types are: string, int, int64, bool.
// Other types will result in an error during processing if a value needs to be set for them.
//
// Order of precedence for a field's value:
// 1. Environment variable (if set and non-empty).
// 2. Default value (if env var is not set or empty, and field is not required).
// 3. Go zero value (if env var not set or empty, no default, and not required).
//
// An empty string from an environment variable is considered "empty" for
// the purpose of 'required' checks. For non-string types, an empty environment
// variable (even if set) will cause parsing to use the default value if available,
// or error if required and no non-empty value can be determined. For string types,
// an empty string from the environment is a valid value.
func Process(spec any) error {
	val := reflect.ValueOf(spec)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("env: spec must be a pointer to a struct, got %T", spec)
	}
	if val.IsNil() {
		return fmt.Errorf("env: spec cannot be a nil pointer")
	}
	s := val.Elem()
	if s.Kind() != reflect.Struct {
		return fmt.Errorf("env: spec must be a pointer to a struct, got pointer to %s", s.Kind())
	}

	typ := s.Type()
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		fieldType := typ.Field(i)

		if !field.CanSet() {
			continue // Skip unexported fields
		}

		envVarName := strings.ToUpper(fieldType.Name)
		if tagName, ok := fieldType.Tag.Lookup("env"); ok && tagName != "" {
			envVarName = tagName
		}

		defaultValueStr := fieldType.Tag.Get("envDefault")
		required := fieldType.Tag.Get("envRequired") == "true"

		envValueStr, envExists := os.LookupEnv(envVarName)

		var valueToParse string
		shouldParseAndSet := false

		if required {
			if !envExists || envValueStr == "" {
				return fmt.Errorf("env: required environment variable %s is not set or is empty", envVarName)
			}
			valueToParse = envValueStr
			shouldParseAndSet = true
		} else { // Not required
			if envExists && envValueStr != "" { // Env var is set and non-empty
				valueToParse = envValueStr
				shouldParseAndSet = true
			} else if envExists && envValueStr == "" { // Env var is set but empty string
				if field.Kind() == reflect.String { // For strings, empty is a valid value if not required
					valueToParse = ""
					shouldParseAndSet = true
				} else if defaultValueStr != "" { // For non-strings, empty string from env, use default
					valueToParse = defaultValueStr
					shouldParseAndSet = true
				}
				// If non-string, empty from env, no default, not required: shouldParseAndSet remains false (use Go zero)
			} else if !envExists && defaultValueStr != "" { // Env var not set, but default exists
				valueToParse = defaultValueStr
				shouldParseAndSet = true
			}
			// If !envExists, no default, not required: shouldParseAndSet remains false (use Go zero)
		}

		if !shouldParseAndSet {
			continue // Leave field with its Go zero value
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(valueToParse)
		case reflect.Int, reflect.Int64:
			intValue, err := strconv.ParseInt(valueToParse, 0, field.Type().Bits())
			if err != nil {
				return fmt.Errorf("env: failed to parse int for %s (variable %s, value: '%s'): %w", fieldType.Name, envVarName, valueToParse, err)
			}
			field.SetInt(intValue)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(valueToParse)
			if err != nil {
				return fmt.Errorf("env: failed to parse bool for %s (variable %s, value: '%s'): %w", fieldType.Name, envVarName, valueToParse, err)
			}
			field.SetBool(boolValue)
		default:
			// Only error if we were actually trying to parse a value for an unsupported type.
			// If shouldParseAndSet was false, we wouldn't reach here for this field.
			return fmt.Errorf("env: unsupported type %s for field %s (variable %s)", field.Kind(), fieldType.Name, envVarName)
		}
	}
	return nil
}
