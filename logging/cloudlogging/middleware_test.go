package cloudlogging

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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

		// *** FIXED LOGIC ***
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

// TestDetermineProjectID remains the same as it was already correct.
func TestDetermineProjectID(t *testing.T) {
	originalFetcher := fetcher
	originalEnv, envSet := os.LookupEnv("GOOGLE_CLOUD_PROJECT")
	t.Cleanup(func() {
		SetProjectIDFetcher(originalFetcher)
		if envSet {
			os.Setenv("GOOGLE_CLOUD_PROJECT", originalEnv)
		} else {
			os.Unsetenv("GOOGLE_CLOUD_PROJECT")
		}
		resetDetermineProjectID()
	})

	t.Run("Metadata Success", func(t *testing.T) {
		resetDetermineProjectID()
		SetProjectIDFetcher(&mockProjectIDFetcher{id: "proj-meta", err: nil})
		os.Unsetenv("GOOGLE_CLOUD_PROJECT")
		if id := determineProjectID(); id != "proj-meta" {
			t.Errorf("got %q, want %q", id, "proj-meta")
		}
	})

	t.Run("Metadata Fails, Env Set", func(t *testing.T) {
		resetDetermineProjectID()
		SetProjectIDFetcher(&mockProjectIDFetcher{id: "", err: errors.New("metadata unavailable")})
		os.Setenv("GOOGLE_CLOUD_PROJECT", "proj-env")
		if id := determineProjectID(); id != "proj-env" {
			t.Errorf("got %q, want %q", id, "proj-env")
		}
	})

	t.Run("Metadata Fails, Env Unset", func(t *testing.T) {
		resetDetermineProjectID()
		SetProjectIDFetcher(&mockProjectIDFetcher{id: "", err: errors.New("metadata unavailable")})
		os.Unsetenv("GOOGLE_CLOUD_PROJECT")
		if id := determineProjectID(); id != "unknown-project" {
			t.Errorf("got %q, want %q", id, "unknown-project")
		}
	})
}

// TestWithCloudTraceContext now uses the corrected helper.
func TestWithCloudTraceContext(t *testing.T) {
	resetDetermineProjectID()
	SetProjectIDFetcher(&mockProjectIDFetcher{id: "p-test", err: nil})
	t.Cleanup(func() {
		resetDetermineProjectID()
	})

	testCases := []struct {
		name           string
		header         string
		expectedTrace  string
		expectedSpan   string
		expectSampled  bool
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
