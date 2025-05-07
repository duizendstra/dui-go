package cloudlogging

import (
	"bytes"
	"context"
	"encoding/json"
	// "io" // No longer needed directly
	"log/slog"
	"os"
	"strings"
	// "sync" // No longer needed
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewCloudLoggingHandler
func TestNewCloudLoggingHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("sets component attribute", func(t *testing.T) {
		testComponent := "my-test-component"
		var buf bytes.Buffer
		// Use the internal constructor with our buffer
		h := newCloudLoggingHandlerWithWriter(&buf, testComponent)
		logger := slog.New(h)

		logger.InfoContext(ctx, "test message for component")

		var logOutput map[string]interface{}
		outputString := buf.String()
		require.NotEmpty(t, outputString, "Log buffer is empty")

		err := json.Unmarshal([]byte(strings.TrimSpace(outputString)), &logOutput)
		require.NoError(t, err, "Failed to unmarshal log output: %s", outputString)

		assert.Equal(t, testComponent, logOutput["component"])
	})

	logLevelTests := []struct {
		name         string
		envLogLevel  string
		setEnv       bool
		expectedSlog slog.Level
		shouldLog    map[slog.Level]bool
	}{
		{
			name:         "default log level is INFO when LOG_LEVEL is not set",
			setEnv:       false,
			expectedSlog: slog.LevelInfo,
			shouldLog: map[slog.Level]bool{
				slog.LevelDebug: false,
				slog.LevelInfo:  true,
				slog.LevelWarn:  true,
				slog.LevelError: true,
			},
		},
		{
			name:         "LOG_LEVEL=DEBUG sets level to DEBUG",
			envLogLevel:  "DEBUG",
			setEnv:       true,
			expectedSlog: slog.LevelDebug,
			shouldLog: map[slog.Level]bool{
				slog.LevelDebug: true,
				slog.LevelInfo:  true,
			},
		},
		{
			name:         "LOG_LEVEL=INFO sets level to INFO",
			envLogLevel:  "INFO",
			setEnv:       true,
			expectedSlog: slog.LevelInfo,
			shouldLog: map[slog.Level]bool{
				slog.LevelDebug: false,
				slog.LevelInfo:  true,
			},
		},
		{
			name:         "LOG_LEVEL=WARN sets level to WARN",
			envLogLevel:  "WARN",
			setEnv:       true,
			expectedSlog: slog.LevelWarn,
			shouldLog: map[slog.Level]bool{
				slog.LevelInfo:  false,
				slog.LevelWarn:  true,
				slog.LevelError: true,
			},
		},
		{
			name:         "LOG_LEVEL=WARNING sets level to WARN",
			envLogLevel:  "WARNING",
			setEnv:       true,
			expectedSlog: slog.LevelWarn,
			shouldLog: map[slog.Level]bool{
				slog.LevelInfo: false,
				slog.LevelWarn: true,
			},
		},
		{
			name:         "LOG_LEVEL=ERROR sets level to ERROR",
			envLogLevel:  "ERROR",
			setEnv:       true,
			expectedSlog: slog.LevelError,
			shouldLog: map[slog.Level]bool{
				slog.LevelWarn:  false,
				slog.LevelError: true,
			},
		},
		{
			name:         "LOG_LEVEL=INVALID defaults to INFO",
			envLogLevel:  "INVALID_LEVEL",
			setEnv:       true,
			expectedSlog: slog.LevelInfo,
			shouldLog: map[slog.Level]bool{
				slog.LevelDebug: false,
				slog.LevelInfo:  true,
			},
		},
	}

	for _, tt := range logLevelTests {
		t.Run(tt.name, func(t *testing.T) {
			originalLogLevel := os.Getenv("LOG_LEVEL")
			if tt.setEnv {
				os.Setenv("LOG_LEVEL", tt.envLogLevel)
			} else {
				os.Unsetenv("LOG_LEVEL")
			}
			// Defer must be inside the loop for tt to be correct on defer execution
			defer func(original string, wasSet bool) {
				if wasSet {
					os.Setenv("LOG_LEVEL", original)
				} else {
					os.Unsetenv("LOG_LEVEL")
				}
			}(originalLogLevel, originalLogLevel != "" || tt.setEnv)

			var buf bytes.Buffer
			// Use the internal constructor with our buffer
			h := newCloudLoggingHandlerWithWriter(&buf, "test-log-level-component")
			logger := slog.New(h)

			assert.True(t, h.Enabled(ctx, tt.expectedSlog), "Handler should be enabled for its configured level %s", tt.expectedSlog)
			if tt.expectedSlog > slog.LevelDebug {
				assert.False(t, h.Enabled(ctx, tt.expectedSlog-1), "Handler should NOT be enabled for one level below its configured level %s", tt.expectedSlog-1)
			}

			for levelToTest, shouldBeLogged := range tt.shouldLog {
				buf.Reset() // Clear buffer for each log attempt
				logger.Log(ctx, levelToTest, "test log level message")
				outputString := buf.String()

				if shouldBeLogged {
					assert.NotEmpty(t, strings.TrimSpace(outputString), "Expected log for level %s, but got none. Log string: '%s'", levelToTest, strings.TrimSpace(outputString))
					assert.True(t, h.Enabled(ctx, levelToTest), "h.Enabled should be true for logged level %s", levelToTest)
				} else {
					assert.Empty(t, strings.TrimSpace(outputString), "Expected no log for level %s, but got: '%s'", levelToTest, strings.TrimSpace(outputString))
					assert.False(t, h.Enabled(ctx, levelToTest), "h.Enabled should be false for non-logged level %s", levelToTest)
				}
			}
		})
	}

	t.Run("replaces attributes for Cloud Logging format", func(t *testing.T) {
		attrTests := []struct {
			name        string
			level       slog.Level
			message     string
			attrs       []slog.Attr
			jsonKeyPath []string
		}{
			{
				name:        "message key becomes 'message'",
				level:       slog.LevelInfo,
				message:     "hello world",
				jsonKeyPath: []string{"message"},
			},
			{
				name:        "level key becomes 'severity' and INFO maps to 'INFO'",
				level:       slog.LevelInfo,
				message:     "info level test",
				jsonKeyPath: []string{"severity"},
			},
			{
				name:        "level key becomes 'severity' and DEBUG maps to 'DEBUG'",
				level:       slog.LevelDebug,
				message:     "debug level test",
				jsonKeyPath: []string{"severity"},
			},
			{
				name:        "level key becomes 'severity' and WARN maps to 'WARNING'",
				level:       slog.LevelWarn,
				message:     "warn level test",
				jsonKeyPath: []string{"severity"},
			},
			{
				name:        "level key becomes 'severity' and ERROR maps to 'ERROR'",
				level:       slog.LevelError,
				message:     "error level test",
				jsonKeyPath: []string{"severity"},
			},
			{
				name:        "custom LevelNotice maps to 'NOTICE'",
				level:       LevelNotice,
				message:     "notice level test",
				jsonKeyPath: []string{"severity"},
			},
			{
				name:        "custom LevelCritical maps to 'CRITICAL'",
				level:       LevelCritical,
				message:     "critical level test",
				jsonKeyPath: []string{"severity"},
			},
			{
				name:        "custom LevelAlert maps to 'ALERT'",
				level:       LevelAlert,
				message:     "alert level test",
				jsonKeyPath: []string{"severity"},
			},
			{
				name:        "custom LevelEmergency maps to 'EMERGENCY'",
				level:       LevelEmergency,
				message:     "emergency level test",
				jsonKeyPath: []string{"severity"},
			},
			// Add a test case specifically for source location transformation
			{
				name:        "source key becomes logging.googleapis.com/sourceLocation",
				level:       slog.LevelInfo, // Check transformation even at INFO level
				message:     "source transform test",
				jsonKeyPath: []string{"logging.googleapis.com/sourceLocation"}, // We expect this key
			},
		}

		originalLogLevel := os.Getenv("LOG_LEVEL")
		os.Setenv("LOG_LEVEL", "DEBUG") // Ensure all levels are processed
		defer os.Setenv("LOG_LEVEL", originalLogLevel)

		for _, tt := range attrTests {
			t.Run(tt.name, func(t *testing.T) {
				var buf bytes.Buffer
				// Use the internal constructor with our buffer
				hDebug := newCloudLoggingHandlerWithWriter(&buf, "test-attr-component-debug")
				loggerDebug := slog.New(hDebug)

				loggerDebug.LogAttrs(ctx, tt.level, tt.message, tt.attrs...)

				var logOutput map[string]interface{}
				trimmedOutput := strings.TrimSpace(buf.String())
				require.NotEmpty(t, trimmedOutput, "Expected log output for attribute test, but got none. Level: %s", tt.level.String())

				err := json.Unmarshal([]byte(trimmedOutput), &logOutput)
				require.NoError(t, err, "Failed to unmarshal log output: %s", trimmedOutput)

				current := logOutput
				var actualValue interface{} // Variable to store the final value found at the path

				// Traverse the jsonKeyPath
				found := true
				for i, key := range tt.jsonKeyPath {
					val, ok := current[key]
					if !ok {
						found = false
						break // Key not found at this level
					}
					if i == len(tt.jsonKeyPath)-1 {
						actualValue = val // Store the final value
					} else {
						var okCast bool
						current, okCast = val.(map[string]interface{})
						if !okCast {
							found = false
							break // Intermediate key is not a map
						}
					}
				}
				require.True(t, found, "Key path %v not found in log output. Full log: %s", tt.jsonKeyPath, trimmedOutput)

				// Perform assertions based on the first key in the path
				firstKey := tt.jsonKeyPath[0]
				if firstKey == "message" {
					assert.Equal(t, tt.message, actualValue, "Message content mismatch. Full log: %s", trimmedOutput)
				} else if firstKey == "severity" {
					var expectedSeverityString string
					switch tt.level {
					case slog.LevelDebug:
						expectedSeverityString = "DEBUG"
					case slog.LevelInfo:
						expectedSeverityString = "INFO"
					case LevelNotice:
						expectedSeverityString = "NOTICE"
					case slog.LevelWarn:
						expectedSeverityString = "WARNING"
					case slog.LevelError:
						expectedSeverityString = "ERROR"
					case LevelCritical:
						expectedSeverityString = "CRITICAL"
					case LevelAlert:
						expectedSeverityString = "ALERT"
					case LevelEmergency:
						expectedSeverityString = "EMERGENCY"
					default:
						expectedSeverityString = "DEFAULT"
					}
					assert.Equal(t, expectedSeverityString, actualValue, "Severity mismatch. Full log: %s", trimmedOutput)
				} else if firstKey == "logging.googleapis.com/sourceLocation" {
					// Check structure of source location
					sourceMap, ok := actualValue.(map[string]interface{})
					require.True(t, ok, "sourceLocation is not a map. Full log: %s", trimmedOutput)
					assert.Contains(t, sourceMap["file"], ".go", "Source file check. Full log: %s", trimmedOutput) // Basic check
					assert.Greater(t, sourceMap["line"], float64(0), "Source line check. Full log: %s", trimmedOutput)
					assert.NotEmpty(t, sourceMap["function"], "Source function check. Full log: %s", trimmedOutput)
				}
				// Add other specific value checks if needed for different jsonKeyPaths
			})
		}
	})

	t.Run("Handle method adds trace and source location", func(t *testing.T) {
		baseTestComponent := "handle-component"

		originalLogLevel := os.Getenv("LOG_LEVEL")
		os.Setenv("LOG_LEVEL", "DEBUG") // Ensure all levels are processed
		defer os.Setenv("LOG_LEVEL", originalLogLevel)

		t.Run("adds trace, span, and sampled info from context", func(t *testing.T) {
			var buf bytes.Buffer
			h := newCloudLoggingHandlerWithWriter(&buf, baseTestComponent)
			logger := slog.New(h)

			traceID := "test-trace-id-123"
			spanID := "test-span-id-456"
			projectID := "test-project"

			originalProjectIDEnv := os.Getenv("GOOGLE_CLOUD_PROJECT")
			os.Setenv("GOOGLE_CLOUD_PROJECT", projectID)
			defer os.Setenv("GOOGLE_CLOUD_PROJECT", originalProjectIDEnv)

			ctxWithTrace := context.WithValue(context.Background(), traceKey{}, "projects/"+projectID+"/traces/"+traceID)
			ctxWithTrace = context.WithValue(ctxWithTrace, spanIDKey{}, spanID)
			ctxWithTrace = context.WithValue(ctxWithTrace, traceSampledKey{}, true)

			logger.InfoContext(ctxWithTrace, "message with trace")

			var logOutput map[string]interface{}
			trimmedOutput := strings.TrimSpace(buf.String())
			require.NotEmpty(t, trimmedOutput, "Expected log output but got none.")
			err := json.Unmarshal([]byte(trimmedOutput), &logOutput)
			require.NoError(t, err, "Failed to unmarshal log output: %s", trimmedOutput)

			assert.Equal(t, "projects/"+projectID+"/traces/"+traceID, logOutput["logging.googleapis.com/trace"])
			assert.Equal(t, spanID, logOutput["logging.googleapis.com/spanId"])
			assert.Equal(t, true, logOutput["logging.googleapis.com/trace_sampled"])
		})

		t.Run("adds default trace info if not in context", func(t *testing.T) {
			var buf bytes.Buffer
			h := newCloudLoggingHandlerWithWriter(&buf, baseTestComponent)
			logger := slog.New(h)

			logger.InfoContext(context.Background(), "message without trace context")

			var logOutput map[string]interface{}
			trimmedOutput := strings.TrimSpace(buf.String())
			require.NotEmpty(t, trimmedOutput, "Expected log output but got none.")
			err := json.Unmarshal([]byte(trimmedOutput), &logOutput)
			require.NoError(t, err, "Failed to unmarshal log output: %s", trimmedOutput)

			assert.Equal(t, "projects/unknown-project/traces/unknown-trace", logOutput["logging.googleapis.com/trace"])
			_, spanExists := logOutput["logging.googleapis.com/spanId"] // SpanID should not be present if not in context
			assert.False(t, spanExists, "spanId should not exist when not in context")
			assert.Equal(t, false, logOutput["logging.googleapis.com/trace_sampled"])
		})

		// == Sub-test for Source Location (Checking ReplaceAttr works) ==
		t.Run("adds sourceLocation for error level logs via ReplaceAttr", func(t *testing.T) {
			var buf bytes.Buffer
			h := newCloudLoggingHandlerWithWriter(&buf, baseTestComponent)
			logger := slog.New(h)

			// The log call itself is the source
			logger.ErrorContext(ctx, "this is an error message") // Line ~405

			var logOutput map[string]interface{}
			trimmedOutput := strings.TrimSpace(buf.String())
			require.NotEmpty(t, trimmedOutput, "Expected log output but got none.")
			err := json.Unmarshal([]byte(trimmedOutput), &logOutput)
			require.NoError(t, err, "Failed to unmarshal log output: %s", trimmedOutput)

			// Check for sourceLocation (added by ReplaceAttr)
			sourceLocation, ok := logOutput["logging.googleapis.com/sourceLocation"].(map[string]interface{})
			require.True(t, ok, "sourceLocation not found or not a map. Full log: %s", trimmedOutput)

			// Check that the file is this test file.
			assert.Contains(t, sourceLocation["file"], "handler_test.go", "File should be handler_test.go. Full log: %s", trimmedOutput)
			assert.IsType(t, float64(0), sourceLocation["line"], "Line should be a number. Full log: %s", trimmedOutput)
			lineNum := sourceLocation["line"].(float64)
			assert.True(t, lineNum > 0, "Line number should be positive (%f). Full log: %s", lineNum, trimmedOutput)
			// Check that the function contains the name of this test function.
			assert.Contains(t, sourceLocation["function"], "TestNewCloudLoggingHandler", "Function name check. Full log: %s", trimmedOutput)
		})

		// == Sub-test for Source Location at INFO level (via ReplaceAttr) ==
		t.Run("adds sourceLocation for info level logs via ReplaceAttr", func(t *testing.T) {
			var buf bytes.Buffer
			h := newCloudLoggingHandlerWithWriter(&buf, baseTestComponent)
			logger := slog.New(h)

			logger.InfoContext(ctx, "this is an info message")

			var logOutput map[string]interface{}
			trimmedOutput := strings.TrimSpace(buf.String())
			require.NotEmpty(t, trimmedOutput, "Expected log output but got none.")
			err := json.Unmarshal([]byte(trimmedOutput), &logOutput)
			require.NoError(t, err, "Failed to unmarshal log output: %s", trimmedOutput)

			// In this revised approach, ReplaceAttr adds sourceLocation regardless of level.
			_, ok := logOutput["logging.googleapis.com/sourceLocation"]
			assert.True(t, ok, "sourceLocation SHOULD be present even for INFO logs with AddSource=true. Full log: %s", trimmedOutput)
		})
	})
}

// TestCloudLoggingHandler_WithAttrs_WithGroup
func TestCloudLoggingHandler_WithAttrs_WithGroup(t *testing.T) {
	ctx := context.Background()
	baseComponent := "with-attrs-groups"

	originalLogLevel := os.Getenv("LOG_LEVEL")
	os.Setenv("LOG_LEVEL", "DEBUG")
	defer os.Setenv("LOG_LEVEL", originalLogLevel)

	t.Run("WithAttrs adds attributes", func(t *testing.T) {
		var buf bytes.Buffer
		h := newCloudLoggingHandlerWithWriter(&buf, baseComponent)
		loggerWithTestWriter := slog.New(h).With("attr1", "value1", "attr2", 123)

		loggerWithTestWriter.InfoContext(ctx, "message from WithAttrs")

		var logOutput map[string]interface{}
		trimmedOutput := strings.TrimSpace(buf.String())
		require.NotEmpty(t, trimmedOutput, "Expected log output but got none.")
		err := json.Unmarshal([]byte(trimmedOutput), &logOutput)
		require.NoError(t, err, "Failed to unmarshal log output: %s", trimmedOutput)

		assert.Equal(t, baseComponent, logOutput["component"])
		assert.Equal(t, "value1", logOutput["attr1"])
		assert.Equal(t, float64(123), logOutput["attr2"])
	})

	t.Run("WithGroup adds a group", func(t *testing.T) {
		var buf bytes.Buffer
		h := newCloudLoggingHandlerWithWriter(&buf, baseComponent)
		loggerWithTestWriter := slog.New(h).WithGroup("myGroup").With("groupAttr", "groupValue")

		loggerWithTestWriter.InfoContext(ctx, "message from WithGroup")

		var logOutput map[string]interface{}
		trimmedOutput := strings.TrimSpace(buf.String())
		require.NotEmpty(t, trimmedOutput, "Expected log output but got none.")
		err := json.Unmarshal([]byte(trimmedOutput), &logOutput)
		require.NoError(t, err, "Failed to unmarshal log output: %s", trimmedOutput)

		assert.Equal(t, baseComponent, logOutput["component"])
		group, ok := logOutput["myGroup"].(map[string]interface{})
		require.True(t, ok, "myGroup not found or not a map. Full log: %s", trimmedOutput)
		assert.Equal(t, "groupValue", group["groupAttr"])
	})

	t.Run("Chained WithAttrs and WithGroup", func(t *testing.T) {
		var buf bytes.Buffer
		h := newCloudLoggingHandlerWithWriter(&buf, baseComponent)
		finalLogger := slog.New(h).With("outerAttr", "outer").
			WithGroup("group1").With("innerAttr1", "inner1").
			WithGroup("group2").With("innerAttr2", "inner2")

		finalLogger.InfoContext(ctx, "chained message")

		var logOutput map[string]interface{}
		trimmedOutput := strings.TrimSpace(buf.String())
		require.NotEmpty(t, trimmedOutput, "Expected log output but got none.")
		err := json.Unmarshal([]byte(trimmedOutput), &logOutput)
		require.NoError(t, err, "Failed to unmarshal log output: %s", trimmedOutput)

		assert.Equal(t, baseComponent, logOutput["component"])
		assert.Equal(t, "outer", logOutput["outerAttr"])

		group1, ok := logOutput["group1"].(map[string]interface{})
		require.True(t, ok, "group1 not found. Full log: %s", trimmedOutput)
		assert.Equal(t, "inner1", group1["innerAttr1"])

		group2, ok := group1["group2"].(map[string]interface{})
		require.True(t, ok, "group2 not found under group1. Full log: %s", trimmedOutput)
		assert.Equal(t, "inner2", group2["innerAttr2"])
	})
}
