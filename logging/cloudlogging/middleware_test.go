package cloudlogging

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// checkContextValues is a test helper that runs the middleware and verifies
// that the expected trace context values are correctly injected.
func checkContextValues(t *testing.T, req *http.Request, expectedTracePrefixProject, expectedTraceID, expectedSpanID string, expectSampled bool) {
	t.Helper()
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		trace := ctx.Value(traceKey{}).(string)
		spanID, spanOK := ctx.Value(spanIDKey{}).(string)
		sampled, sampledOK := ctx.Value(traceSampledKey{}).(bool)

		expectedTraceString := "projects/" + expectedTracePrefixProject + "/traces/" + expectedTraceID
		if trace != expectedTraceString {
			t.Errorf("Context trace: got %q, want %q", trace, expectedTraceString)
		}

		// Correctly check for span presence and value.
		if expectedSpanID != "" {
			if !spanOK || spanID != expectedSpanID {
				t.Errorf("Context spanID: got %q (present %t), want %q", spanID, spanOK, expectedSpanID)
			}
		} else {
			if spanOK {
				t.Errorf("Context spanID: should be absent, but was present with value %q", spanID)
			}
		}

		// Correctly check for sampled presence and value.
		if expectSampled {
			if !sampledOK || !sampled {
				t.Errorf("Context sampled: got %t (present %t), want true", sampled, sampledOK)
			}
		} else {
			if sampledOK {
				t.Errorf("Context sampled: should be absent, but was present with value %t", sampled)
			}
		}
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	middleware := WithCloudTraceContext(dummyHandler)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("dummyHandler returned non-OK status: got %v", status)
	}
}

// --- Middleware Tests ---

// TestDetermineProjectID specifically tests the project ID determination logic,
// covering precedence (metadata, environment variable, default) and the use of sync.Once.
func TestDetermineProjectID(t *testing.T) {
	originalFetcher := fetcher
	originalEnv, envSet := os.LookupEnv("GOOGLE_CLOUD_PROJECT")
	// Ensure original fetcher and env var are restored after the test.
	t.Cleanup(func() {
		SetProjectIDFetcher(originalFetcher)
		if envSet {
			os.Setenv("GOOGLE_CLOUD_PROJECT", originalEnv)
		} else {
			os.Unsetenv("GOOGLE_CLOUD_PROJECT")
		}
		resetDetermineProjectID() // Reset after all sub-tests in this function complete.
	})

	t.Run("Metadata Success", func(t *testing.T) {
		resetDetermineProjectID() // Reset for this specific sub-test.
		mockFetcher := &mockProjectIDFetcher{id: "proj-meta", err: nil}
		SetProjectIDFetcher(mockFetcher)
		os.Unsetenv("GOOGLE_CLOUD_PROJECT") // Ensure env var is not set.

		// First call should use the mock fetcher.
		id := determineProjectID()
		if id != "proj-meta" {
			t.Errorf("Expected project ID 'proj-meta', got %q", id)
		}
	})

	t.Run("Metadata Fails, Env Set", func(t *testing.T) {
		resetDetermineProjectID()
		mockFetcher := &mockProjectIDFetcher{id: "", err: errors.New("metadata unavailable")} // Simulate fetcher error.
		SetProjectIDFetcher(mockFetcher)
		os.Setenv("GOOGLE_CLOUD_PROJECT", "proj-env") // Set env var fallback.

		id := determineProjectID()
		if id != "proj-env" { // Should use the env var.
			t.Errorf("Expected project ID 'proj-env', got %q", id)
		}
	})

	t.Run("Metadata Fails, Env Unset", func(t *testing.T) {
		resetDetermineProjectID()
		mockFetcher := &mockProjectIDFetcher{id: "", err: errors.New("metadata unavailable")} // Simulate fetcher error.
		SetProjectIDFetcher(mockFetcher)
		os.Unsetenv("GOOGLE_CLOUD_PROJECT") // Ensure env var is also unset.

		id := determineProjectID()
		if id != "unknown-project" { // Should use the hardcoded default.
			t.Errorf("Expected project ID 'unknown-project', got %q", id)
		}
	})
}

// TestWithCloudTraceContext focuses on the middleware's ability to parse the
// X-Cloud-Trace-Context header and inject the correct values into the request context.
func TestWithCloudTraceContext(t *testing.T) {
	resetDetermineProjectID()
	SetProjectIDFetcher(&mockProjectIDFetcher{id: "p-test", err: nil})
	t.Cleanup(func() {
		resetDetermineProjectID()
	})

	testCases := []struct {
		name          string
		header        string
		expectedTrace string
		expectedSpan  string
		expectSampled bool
	}{
		{"Full Header, Sampled", "trace123/span456;o=1", "trace123", "span456", true},
		{"Full Header, Not Sampled", "traceABC/spanDEF;o=0", "traceABC", "spanDEF", false},
		{"Empty Span, Sampled", "traceXYZ/;o=1", "traceXYZ", "", true},
		{"No Span Part, Not Sampled", "traceOnly;o=0", "traceOnly", "", false},
		{"No Header", "", "unknown-trace", "", false},
		{"Malformed Header", "invalid-format", "unknown-trace", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tc.header != "" {
				req.Header.Set("X-Cloud-Trace-Context", tc.header)
			}
			checkContextValues(t, req, "p-test", tc.expectedTrace, tc.expectedSpan, tc.expectSampled)
		})
	}
}

// TestWithTrace specifically tests the WithTrace helper function, ensuring it
// correctly injects a formatted trace string into the context.
func TestWithTrace(t *testing.T) {
	// Setup: Ensure a deterministic project ID for the test.
	originalFetcher := fetcher
	SetProjectIDFetcher(&mockProjectIDFetcher{id: "test-project", err: nil})
	t.Cleanup(func() {
		SetProjectIDFetcher(originalFetcher)
		resetDetermineProjectID()
	})

	// The trace ID we want to inject.
	traceID := "my-job-trace-id-123"

	// Execution: Call the function under test.
	ctx := WithTrace(context.Background(), traceID)

	// Verification: Check the context for the correct values.
	// 1. Check that the trace value was added correctly.
	traceVal, ok := ctx.Value(traceKey{}).(string)
	require.True(t, ok, "traceKey should exist in the context")
	expectedTrace := "projects/test-project/traces/my-job-trace-id-123"
	assert.Equal(t, expectedTrace, traceVal)

	// 2. Check that span and sampled keys were NOT added.
	_, spanOK := ctx.Value(spanIDKey{}).(string)
	assert.False(t, spanOK, "spanIDKey should not be set by WithTrace")

	_, sampledOK := ctx.Value(traceSampledKey{}).(bool)
	assert.False(t, sampledOK, "traceSampledKey should not be set by WithTrace")
}
