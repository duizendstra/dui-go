package cloudlogging

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewCloudLoggingHandler_LevelFromEnv verifies that the log level is correctly set
// from the LOG_LEVEL environment variable and that the Enabled() method reflects this.
func TestNewCloudLoggingHandler_LevelFromEnv(t *testing.T) {
	testCases := []struct {
		name          string
		envValue      string
		testLevel     slog.Level
		expectEnabled bool
	}{
		{"Env DEBUG enables DEBUG", "DEBUG", slog.LevelDebug, true},
		{"Env INFO disables DEBUG", "INFO", slog.LevelDebug, false},
		{"Env WARN enables WARN", "WARN", slog.LevelWarn, true},
		{"Env WARN disables INFO", "WARN", slog.LevelInfo, false},
		{"Env ERROR enables ERROR", "ERROR", slog.LevelError, true},
		{"Env ERROR disables WARN", "ERROR", slog.LevelWarn, false},
		{"Env unset defaults to INFO", "", slog.LevelInfo, true},
		{"Env invalid defaults to INFO", "INVALID", slog.LevelInfo, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Backup and restore the environment variable INSIDE the sub-test
			// to prevent state from leaking between test cases.
			originalLevel, levelSet := os.LookupEnv("LOG_LEVEL")
			t.Cleanup(func() {
				if levelSet {
					os.Setenv("LOG_LEVEL", originalLevel)
				} else {
					os.Unsetenv("LOG_LEVEL")
				}
			})

			// Set the environment for this specific sub-test.
			if tc.envValue == "" {
				os.Unsetenv("LOG_LEVEL")
			} else {
				os.Setenv("LOG_LEVEL", tc.envValue)
			}

			// Create the handler which will read the environment variable.
			var buf bytes.Buffer
			handler := NewCloudLoggingHandlerForTest(&buf, "test-component")

			// Assert that the Enabled method works as expected.
			assert.Equal(t, tc.expectEnabled, handler.Enabled(context.Background(), tc.testLevel))
		})
	}
}

// TestCloudLoggingHandler_Handle verifies the core logic of the Handle method.
func TestCloudLoggingHandler_Handle(t *testing.T) {
	originalFetcher := fetcher
	t.Cleanup(func() {
		SetProjectIDFetcher(originalFetcher)
		resetDetermineProjectID()
	})

	t.Run("adds trace, span, and sampled info from context", func(t *testing.T) {
		var buf bytes.Buffer
		h := NewCloudLoggingHandlerForTest(&buf, "test-component")
		logger := slog.New(h)

		ctx := context.WithValue(context.Background(), traceKey{}, "projects/test-proj/traces/trace-abc")
		ctx = context.WithValue(ctx, spanIDKey{}, "span-123")
		ctx = context.WithValue(ctx, traceSampledKey{}, true)

		logger.InfoContext(ctx, "message with trace")

		var logOutput map[string]interface{}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &logOutput))

		assert.Equal(t, "projects/test-proj/traces/trace-abc", logOutput["logging.googleapis.com/trace"])
		assert.Equal(t, "span-123", logOutput["logging.googleapis.com/spanId"])
		assert.Equal(t, true, logOutput["logging.googleapis.com/trace_sampled"])
	})

	t.Run("adds default trace info if not in context", func(t *testing.T) {
		resetDetermineProjectID()
		SetProjectIDFetcher(&mockProjectIDFetcher{id: "", err: errors.New("metadata unavailable")})
		originalEnv, envSet := os.LookupEnv("GOOGLE_CLOUD_PROJECT")
		os.Unsetenv("GOOGLE_CLOUD_PROJECT")
		t.Cleanup(func() {
			if envSet {
				os.Setenv("GOOGLE_CLOUD_PROJECT", originalEnv)
			}
		})

		var buf bytes.Buffer
		h := NewCloudLoggingHandlerForTest(&buf, "test-component")
		logger := slog.New(h)

		logger.InfoContext(context.Background(), "message without trace")

		var logOutput map[string]interface{}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &logOutput))

		assert.Equal(t, "projects/unknown-project/traces/unknown-trace", logOutput["logging.googleapis.com/trace"])
		_, hasSpan := logOutput["logging.googleapis.com/spanId"]
		assert.False(t, hasSpan)
		assert.Equal(t, false, logOutput["logging.googleapis.com/trace_sampled"])
	})

	t.Run("adds source location for all levels", func(t *testing.T) {
		var buf bytes.Buffer
		h := NewCloudLoggingHandlerForTest(&buf, "test-component")
		logger := slog.New(h)

		logger.Info("a message that should have source")

		var logOutput map[string]interface{}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &logOutput))

		sourceLocRaw, exists := logOutput["logging.googleapis.com/sourceLocation"]
		require.True(t, exists)

		sourceLocMap, ok := sourceLocRaw.(map[string]interface{})
		require.True(t, ok)

		assert.Contains(t, sourceLocMap["file"], "handler_test.go")
		assert.Contains(t, sourceLocMap["function"], "TestCloudLoggingHandler_Handle")
		assert.Greater(t, sourceLocMap["line"], float64(0))
	})
}

// TestCloudLoggingReplaceAttr tests the attribute replacement logic directly.
func TestCloudLoggingReplaceAttr(t *testing.T) {
	tests := []struct {
		name         string
		inputAttr    slog.Attr
		expectedAttr slog.Attr
	}{
		{"Message Key", slog.String(slog.MessageKey, "hello"), slog.String("message", "hello")},
		{"Level Debug", slog.Any(slog.LevelKey, slog.LevelDebug), slog.String("severity", "DEBUG")},
		{"Level Info", slog.Any(slog.LevelKey, slog.LevelInfo), slog.String("severity", "INFO")},
		{"Level Notice", slog.Any(slog.LevelKey, LevelNotice), slog.String("severity", "NOTICE")},
		{"Level Warn", slog.Any(slog.LevelKey, slog.LevelWarn), slog.String("severity", "WARNING")},
		{"Level Error", slog.Any(slog.LevelKey, slog.LevelError), slog.String("severity", "ERROR")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cloudLoggingReplaceAttr(nil, tt.inputAttr)
			assert.Equal(t, tt.expectedAttr.Key, got.Key)
			assert.True(t, tt.expectedAttr.Value.Equal(got.Value))
		})
	}
}
