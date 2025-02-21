// Package env provides functionality for loading and validating environment variables.
// It allows you to define a list of expected environment variables, specifying whether
// they are mandatory, have a default value, or require validation.
//
// Typical usage:
//
//	loader := env.NewEnvLoader()
//	vars := []env.EnvVar{
//	  {Key: "API_KEY", Mandatory: true},
//	  {Key: "PORT", Mandatory: false, DefaultValue: "8080"},
//	}
//	envMap, err := loader.LoadEnv(ctx, vars)
//	if err != nil {
//	  // handle error (e.g., missing mandatory variable)
//	}
//
// The EnvLoader returns errors if mandatory variables are missing or validation fails.
// Using this approach simplifies startup configuration and ensures that your application
// clearly reports configuration issues through returned errors.
package env
