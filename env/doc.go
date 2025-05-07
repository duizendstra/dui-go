// Package env provides utilities for loading environment variables
// directly into Go structs. Configuration is managed via struct tags.
//
// It aims to simplify application configuration by providing a type-safe
// way to consume environment variables with support for defaults and
// required fields, without requiring external dependencies.
//
// Currently supported field types for struct population are:
// string, int, int64, and bool.
//
// Struct tags used for configuration:
//   - `env:"ENV_VAR_NAME"`: Specifies the environment variable name.
//     Defaults to the uppercase field name if omitted or empty.
//   - `envDefault:"value"`: Provides a default string value if the
//     environment variable is not set or is empty. This string will be
//     parsed to the field's type.
//   - `envRequired:"true"`: Marks the environment variable as mandatory.
//     An error is returned if it's not set or is an empty string.
package env