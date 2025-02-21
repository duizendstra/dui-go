package cloudlogging

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithCloudTraceContext(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trace, _ := r.Context().Value(traceKey{}).(string)
		spanID, _ := r.Context().Value(spanIDKey{}).(string)
		traceSampled, _ := r.Context().Value(traceSampledKey{}).(bool)

		if trace == "" {
			t.Error("expected trace to be set, got empty string")
		}
		if spanID != "123" {
			t.Errorf("expected spanID = 123, got %q", spanID)
		}
		if !traceSampled {
			t.Error("expected traceSampled = true, got false")
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("X-Cloud-Trace-Context", "abcdef1234567890abcdef1234567890/123;o=1")

	rr := httptest.NewRecorder()
	handler := WithCloudTraceContext(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned non-OK status: got %v want %v", status, http.StatusOK)
	}
}

func TestWithCloudTraceContext_NoHeader(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trace, _ := r.Context().Value(traceKey{}).(string)
		if trace == "" {
			t.Error("expected a default trace value, got empty string")
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	rr := httptest.NewRecorder()

	handler := WithCloudTraceContext(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned non-OK status: got %v want %v", status, http.StatusOK)
	}
}

// Test case for a malformed X-Cloud-Trace-Context header
func TestWithCloudTraceContext_MalformedHeader(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trace, _ := r.Context().Value(traceKey{}).(string)
		if trace == "" {
			t.Error("expected a default trace value, got empty string")
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("X-Cloud-Trace-Context", "not-a-valid-trace-header")

	rr := httptest.NewRecorder()
	handler := WithCloudTraceContext(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned non-OK status for malformed header: got %v want %v", status, http.StatusOK)
	}
}
