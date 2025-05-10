// Package env provides utilities for loading environment variables
// directly into Go structs in a type-safe manner. Configuration is managed via struct tags,
// eliminating the need for external dependencies for basic environment configuration.
//
// It aims to simplify application configuration by providing a type-safe
// way to consume environment variables with support for defaults and
// required fields.
//
// Currently supported field types for struct population are:
// string, int, int64, and bool.
//
// Struct tags used for configuration:
//   - : Specifies the environment variable name.
//     Defaults to the uppercase field name if omitted or empty.
//   - : Provides a default string value if the
//     environment variable is not set or is empty. This string will be
//     parsed to the field's type.
//   - : Marks the environment variable as mandatory.
//     An error is returned if it's not set or is an empty string.
//
// Example struct and usage:
//
//	type AppConfig struct {
//	    Host string
//	    Port int
//	}
//
//	func main() {
//	    os.Setenv("APP_PORT", "8080") // Example setup
//	    var cfg AppConfig
//	    if err := env.Process(&cfg); err != nil {
//	        log.Fatalf("Error loading config: %v", err)
//	    }
//	    // cfg.Host will be "localhost", cfg.Port will be 8080
//	}
package env
