package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testConfig struct {
	Name           string `env:"TEST_NAME" envDefault:"DefaultName"`
	Port           int    `env:"TEST_PORT" envRequired:"true"`
	Timeout        int64  `env:"TEST_TIMEOUT" envDefault:"5000"`
	Enabled        bool   `env:"TEST_ENABLED" envDefault:"true"` // Default "true" string for bool
	AutoStart      bool   `env:"AUTO_START" envDefault:"false"`// Default "false" string for bool
	NonTagged      string // Should look for NONTAGGED
	MissingVar     string `env:"MISSING_VAR" envDefault:"DefaultMissing"`
	EmptyVarSet    string `env:"EMPTY_VAR_SET"` // Set to "" explicitly, no default
	EmptyVarDefault string `env:"EMPTY_VAR_DEFAULT" envDefault:"DefaultForEmpty"` // Not set, or set to "", should use default
	NotSetInt      int    `env:"NOT_SET_INT"`   // Expect Go zero value
	NotSetBool     bool   `env:"NOT_SET_BOOL"`  // Expect Go zero value
	EmptyTagField  string `env:"" envDefault:"EmptyTagDefault"` // Uppercase field name: EMPTYTAGFIELD
}

type requiredOnlyConfig struct {
	MustExist    string `env:"MUST_EXIST" envRequired:"true"`
	MustExistToo int    `env:"MUST_EXIST_TOO" envRequired:"true"`
}

type badTypeConfig struct {
	BadField float32 `env:"BAD_FIELD"` // float32 is unsupported
}

// setupEnvVars sets environment variables and returns a cleanup function.
func setupEnvVars(t *testing.T, vars map[string]string) func() {
	t.Helper()
	originalVars := make(map[string]struct {
		value  string
		exists bool
	})

	for k, v := range vars {
		originalValue, exists := os.LookupEnv(k)
		originalVars[k] = struct {
			value  string
			exists bool
		}{originalValue, exists}
		err := os.Setenv(k, v)
		require.NoError(t, err)
	}

	return func() {
		for k := range vars { // Iterate over the keys that were originally in 'vars' to restore them
			original := originalVars[k]
			if original.exists {
				err := os.Setenv(k, original.value)
				require.NoError(t, err)
			} else {
				err := os.Unsetenv(k)
				require.NoError(t, err)
			}
		}
	}
}

func TestProcess_Success(t *testing.T) {
	cleanup := setupEnvVars(t, map[string]string{
		"TEST_NAME":    "ActualName",
		"TEST_PORT":    "8080",
		"TEST_TIMEOUT": "30",
		"TEST_ENABLED": "false", // Override default
		"AUTO_START":   "true",  // Override default
		"NONTAGGED":    "NonTaggedValue",
		"EMPTY_VAR_SET": "",      // Explicitly empty
		"EMPTYTAGFIELD": "FromEnv",
	})
	defer cleanup()

	var cfg testConfig
	err := Process(&cfg)
	require.NoError(t, err)

	assert.Equal(t, "ActualName", cfg.Name)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, int64(30), cfg.Timeout)
	assert.Equal(t, false, cfg.Enabled)
	assert.Equal(t, true, cfg.AutoStart)
	assert.Equal(t, "NonTaggedValue", cfg.NonTagged)
	assert.Equal(t, "DefaultMissing", cfg.MissingVar) // Uses default as MISSING_VAR is not set
	assert.Equal(t, "", cfg.EmptyVarSet)              // Explicit empty string from env
	assert.Equal(t, "DefaultForEmpty", cfg.EmptyVarDefault) // Uses default as EMPTY_VAR_DEFAULT is not set
	assert.Equal(t, 0, cfg.NotSetInt)                  // Go zero value
	assert.Equal(t, false, cfg.NotSetBool)             // Go zero value
	assert.Equal(t, "FromEnv", cfg.EmptyTagField)
}

func TestProcess_DefaultsAndZeroValues(t *testing.T) {
	varsToManage := []string{
		"TEST_NAME", "TEST_PORT", "TEST_TIMEOUT", "TEST_ENABLED", "AUTO_START",
		"NONTAGGED", "MISSING_VAR", "EMPTY_VAR_SET", "EMPTY_VAR_DEFAULT",
		"NOT_SET_INT", "NOT_SET_BOOL", "EMPTYTAGFIELD",
	}
	originalValues := make(map[string]struct{ v string; e bool })
	for _, k := range varsToManage {
		v, e := os.LookupEnv(k)
		originalValues[k] = struct{ v string; e bool }{v, e}
		if k != "TEST_PORT" { // TEST_PORT needs to be set for this specific test to pass Process
			os.Unsetenv(k)
		}
	}
	// Set the required TEST_PORT for this test
	err := os.Setenv("TEST_PORT", "9090")
	require.NoError(t, err)


	defer func() {
		for _, k := range varsToManage {
			ov := originalValues[k]
			if ov.e {
				os.Setenv(k, ov.v)
			} else {
				os.Unsetenv(k)
			}
		}
	}()


	var cfg testConfig
	err = Process(&cfg)
	require.NoError(t, err)

	assert.Equal(t, "DefaultName", cfg.Name)
	assert.Equal(t, 9090, cfg.Port) // Set value
	assert.Equal(t, int64(5000), cfg.Timeout)
	assert.Equal(t, true, cfg.Enabled)
	assert.Equal(t, false, cfg.AutoStart)
	assert.Equal(t, "", cfg.NonTagged) // Go zero as it was unset
	assert.Equal(t, "DefaultMissing", cfg.MissingVar)
	assert.Equal(t, "", cfg.EmptyVarSet)
	assert.Equal(t, "DefaultForEmpty", cfg.EmptyVarDefault)
	assert.Equal(t, 0, cfg.NotSetInt)
	assert.Equal(t, false, cfg.NotSetBool)
	assert.Equal(t, "EmptyTagDefault", cfg.EmptyTagField)
}

func TestProcess_RequiredMissing(t *testing.T) {
	t.Run("TEST_PORT required and not set", func(t *testing.T) {
		// Explicitly manage TEST_PORT and any other var for testConfig for this subtest
		originalTestPort, testPortExists := os.LookupEnv("TEST_PORT")
		originalTestName, testNameExists := os.LookupEnv("TEST_NAME")

		os.Unsetenv("TEST_PORT") // Ensure TEST_PORT is not set
		// Set other vars if they influence testConfig beyond the required field
		// For testConfig, TEST_NAME doesn't need to be set for this specific check, but good to control
		os.Unsetenv("TEST_NAME")


		defer func() {
			if testPortExists {
				os.Setenv("TEST_PORT", originalTestPort)
			} else {
				os.Unsetenv("TEST_PORT")
			}
			if testNameExists {
				os.Setenv("TEST_NAME", originalTestName)
			} else {
				os.Unsetenv("TEST_NAME")
			}
		}()

		var cfg testConfig
		err := Process(&cfg)
		require.Error(t, err, "Process should error because TEST_PORT is required but not set")
		if err != nil { // Check for nil before accessing Error()
			assert.Contains(t, err.Error(), "required environment variable TEST_PORT is not set or is empty")
		}
	})

	t.Run("MUST_EXIST required and set to empty", func(t *testing.T) {
		originalMustExist, mustExistExists := os.LookupEnv("MUST_EXIST")
		originalMustExistToo, mustExistTooExists := os.LookupEnv("MUST_EXIST_TOO")

		os.Setenv("MUST_EXIST", "")      // Set to empty, which is invalid for required
		os.Setenv("MUST_EXIST_TOO", "123") // Set this required field to a valid value

		defer func() {
			if mustExistExists {
				os.Setenv("MUST_EXIST", originalMustExist)
			} else {
				os.Unsetenv("MUST_EXIST")
			}
			if mustExistTooExists {
				os.Setenv("MUST_EXIST_TOO", originalMustExistToo)
			} else {
				os.Unsetenv("MUST_EXIST_TOO")
			}
		}()

		var cfg requiredOnlyConfig
		err := Process(&cfg)
		require.Error(t, err, "Process should error because MUST_EXIST is required but set to empty")
		if err != nil {
			assert.Contains(t, err.Error(), "required environment variable MUST_EXIST is not set or is empty")
		}
	})

	t.Run("MUST_EXIST_TOO required and not set", func(t *testing.T) {
		originalMustExist, mustExistExists := os.LookupEnv("MUST_EXIST")
		originalMustExistToo, mustExistTooExists := os.LookupEnv("MUST_EXIST_TOO")

		os.Setenv("MUST_EXIST", "ValidValue") // Set this required field to a valid value
		os.Unsetenv("MUST_EXIST_TOO")       // Ensure this required field is not set

		defer func() {
			if mustExistExists {
				os.Setenv("MUST_EXIST", originalMustExist)
			} else {
				os.Unsetenv("MUST_EXIST")
			}
			if mustExistTooExists {
				os.Setenv("MUST_EXIST_TOO", originalMustExistToo)
			} else {
				os.Unsetenv("MUST_EXIST_TOO")
			}
		}()

		var cfg requiredOnlyConfig
		err := Process(&cfg)
		require.Error(t, err, "Process should error because MUST_EXIST_TOO is required but not set")
		if err != nil {
			assert.Contains(t, err.Error(), "required environment variable MUST_EXIST_TOO is not set or is empty")
		}
	})
}

func TestProcess_ParseErrors(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectedMsg string
	}{
		{
			name: "int parse error",
			envVars: map[string]string{"TEST_PORT": "not-an-int"},
			expectedMsg: "failed to parse int for Port (variable TEST_PORT, value: 'not-an-int')",
		},
		{
			name: "bool parse error",
			envVars: map[string]string{"TEST_ENABLED": "not-a-bool"},
			expectedMsg: "failed to parse bool for Enabled (variable TEST_ENABLED, value: 'not-a-bool')",
		},
		{
			name: "int64 parse error",
			envVars: map[string]string{"TEST_TIMEOUT": "not-int64"},
			expectedMsg: "failed to parse int for Timeout (variable TEST_TIMEOUT, value: 'not-int64')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Base valid values for required fields in testConfig
			baseConfigVars := map[string]string{
				"TEST_PORT": "0", // Required
				// TEST_ENABLED and TEST_TIMEOUT have defaults, so not strictly needed here for Process to run
			}
			// Overlay test-specific vars
			for k, v := range tt.envVars {
				baseConfigVars[k] = v
			}

			cleanup := setupEnvVars(t, baseConfigVars)
			defer cleanup()

			var cfg testConfig
			err := Process(&cfg)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedMsg)
		})
	}
}

func TestProcess_UnsupportedType(t *testing.T) {
	var cfg badTypeConfig // BadField is float32
	// If BAD_FIELD is not set, Process should not error for this field
	err := Process(&cfg)
	require.NoError(t, err)

	// If BAD_FIELD is set, then it should error during attempt to parse/set
	cleanup := setupEnvVars(t, map[string]string{"BAD_FIELD": "1.23"})
	defer cleanup()

	err = Process(&cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type float32 for field BadField")
}

func TestProcess_NotAPointer(t *testing.T) {
	var cfg testConfig // Not a pointer
	err := Process(cfg)
	require.Error(t, err)
	assert.Equal(t, "env: spec must be a pointer to a struct, got env.testConfig", err.Error())
}

func TestProcess_NilPointer(t *testing.T) {
	var cfg *testConfig = nil // Nil pointer
	err := Process(cfg)
	require.Error(t, err)
	assert.Equal(t, "env: spec cannot be a nil pointer", err.Error())
}

func TestProcess_PointerToNonStruct(t *testing.T) {
	var i int
	err := Process(&i) // Pointer to int
	require.Error(t, err)
	assert.Equal(t, "env: spec must be a pointer to a struct, got pointer to int", err.Error())
}

func TestProcess_EmptyEnvTagDefaultsToFieldName(t *testing.T) {
	type Config struct {
		MyField string `env:"" envDefault:"helloFromEmptyTag"`
	}
	var cfg Config

	originalMyField, myFieldExists := os.LookupEnv("MYFIELD")
	os.Unsetenv("MYFIELD") // Ensure MYFIELD is not set to test default
	defer func() {
		if myFieldExists {
			os.Setenv("MYFIELD", originalMyField)
		} else {
			os.Unsetenv("MYFIELD")
		}
	}()

	err := Process(&cfg)
	require.NoError(t, err)
	assert.Equal(t, "helloFromEmptyTag", cfg.MyField)

	// Test when MYFIELD is set from env
	cleanup := setupEnvVars(t, map[string]string{"MYFIELD": "worldFromEmptyTag"})
	defer cleanup() // This will restore MYFIELD to its state before this setupEnvVars call
	
	err = Process(&cfg) // Re-process with the new env var
	require.NoError(t, err)
	assert.Equal(t, "worldFromEmptyTag", cfg.MyField)
}

func TestProcess_RequiredFieldIgnoresDefault(t *testing.T) {
	type Config struct {
		ImportantValue string `env:"IMPORTANT_VAL" envRequired:"true" envDefault:"shouldBeIgnored"`
	}
	var cfg Config

	t.Run("required var missing", func(t *testing.T) {
		originalImpVal, impValExists := os.LookupEnv("IMPORTANT_VAL")
		os.Unsetenv("IMPORTANT_VAL")
		defer func() {
			if impValExists { os.Setenv("IMPORTANT_VAL", originalImpVal) } else { os.Unsetenv("IMPORTANT_VAL") }
		}()

		err := Process(&cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "required environment variable IMPORTANT_VAL is not set or is empty")
	})

	t.Run("required var set", func(t *testing.T) {
		originalImpVal, impValExists := os.LookupEnv("IMPORTANT_VAL")
		os.Setenv("IMPORTANT_VAL", "actualValueFromEnv")
		defer func() {
			if impValExists { os.Setenv("IMPORTANT_VAL", originalImpVal) } else { os.Unsetenv("IMPORTANT_VAL") }
		}()
		
		err := Process(&cfg)
		require.NoError(t, err)
		assert.Equal(t, "actualValueFromEnv", cfg.ImportantValue)
	})
}

func TestProcess_EmptyStringFromEnvForNonString(t *testing.T) {
	type Config struct {
		MyInt  int  `env:"MY_INT" envDefault:"123"`
		MyBool bool `env:"MY_BOOL" envDefault:"true"`
	}
	var cfg Config

	originalMyInt, myIntExists := os.LookupEnv("MY_INT")
	originalMyBool, myBoolExists := os.LookupEnv("MY_BOOL")

	os.Setenv("MY_INT", "")  // Set to empty
	os.Setenv("MY_BOOL", "") // Set to empty
	defer func() {
		if myIntExists { os.Setenv("MY_INT", originalMyInt) } else { os.Unsetenv("MY_INT") }
		if myBoolExists { os.Setenv("MY_BOOL", originalMyBool) } else { os.Unsetenv("MY_BOOL") }
	}()

	err := Process(&cfg)
	require.NoError(t, err)
	assert.Equal(t, 123, cfg.MyInt, "Expected default int value when env var is empty")
	assert.Equal(t, true, cfg.MyBool, "Expected default bool value when env var is empty")
}

func TestProcess_UnexportedField(t *testing.T) {
    type Config struct {
        Exported   string `env:"EXPORTED_FIELD" envDefault:"exported"`
        unexported string // No tags, unexported
    }
    var cfg Config
    cfg.unexported = "initial" // Set initial value for unexported field

    // Explicitly manage env vars for this test
    originalExported, exportedExists := os.LookupEnv("EXPORTED_FIELD")
    originalUnexportedEnv, unexportedEnvExists := os.LookupEnv("UNEXPORTEDFIELD") // Default name for unexported

    os.Setenv("EXPORTED_FIELD", "from_env")
    os.Setenv("UNEXPORTEDFIELD", "env_unexported") // This should not affect the unexported field because it's unexported

    defer func() {
        if exportedExists { os.Setenv("EXPORTED_FIELD", originalExported) } else { os.Unsetenv("EXPORTED_FIELD") }
        if unexportedEnvExists { os.Setenv("UNEXPORTEDFIELD", originalUnexportedEnv) } else { os.Unsetenv("UNEXPORTEDFIELD") }
    }()

    err := Process(&cfg)
    require.NoError(t, err)
    assert.Equal(t, "from_env", cfg.Exported)
    assert.Equal(t, "initial", cfg.unexported, "Unexported field should not be modified by Process")
}

