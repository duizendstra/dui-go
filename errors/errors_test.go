// go-dui/errors/errors_test.go
//
// This test verifies the core functionality of the APIError type and its methods.
// It checks creating new errors, adding details, comparing errors with Is(),
// and ensuring that messages are formatted correctly.

package errors

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(400, "bad request")
	if err.Code != 400 {
		t.Errorf("expected code=400, got %d", err.Code)
	}
	if err.Message != "bad request" {
		t.Errorf("expected message='bad request', got %q", err.Message)
	}
	if len(err.Details) != 0 {
		t.Errorf("expected no details, got %d", len(err.Details))
	}
}

func TestNewFallbackMessage(t *testing.T) {
	err := New(500, "")
	if err.Message != "unknown error" {
		t.Errorf("expected fallback message='unknown error', got %q", err.Message)
	}
}

func TestWithDetails(t *testing.T) {
	err := New(404, "not found")
	err.WithDetails(ErrorDetail{Reason: "DB_MISS", Message: "No record in database"})
	if len(err.Details) != 1 {
		t.Fatalf("expected 1 detail, got %d", len(err.Details))
	}
	if err.Details[0].Reason != "DB_MISS" {
		t.Errorf("expected reason='DB_MISS', got %q", err.Details[0].Reason)
	}
	if err.Details[0].Message != "No record in database" {
		t.Errorf("expected message='No record in database', got %q", err.Details[0].Message)
	}
}

func TestErrorFormatting(t *testing.T) {
	err := New(500, "internal error")
	msg := err.Error()
	if msg != `APIError: code=500, message="internal error"` {
		t.Errorf("unexpected error format: %q", msg)
	}

	err.WithDetails(ErrorDetail{Reason: "CONFIG", Message: "missing config"})
	msg = err.Error()
	if msg != `APIError: code=500, message="internal error", details=[CONFIG: missing config]` {
		t.Errorf("unexpected error format with details: %q", msg)
	}
}

func TestIs(t *testing.T) {
	errA := New(400, "bad request")
	errB := New(400, "bad request")
	errC := New(400, "different message")
	errD := New(500, "bad request")

	if !errors.Is(errA, errB) {
		t.Error("errA should match errB (same code and message)")
	}
	if errors.Is(errA, errC) {
		t.Error("errA should not match errC (different messages)")
	}
	if errors.Is(errA, errD) {
		t.Error("errA should not match errD (different codes)")
	}

	// Non-APIError should not match
	if errors.Is(errA, errors.New("some other error")) {
		t.Error("APIError should not match a non-APIError")
	}
}

func TestPredefinedErrors(t *testing.T) {
	if ErrBadRequest.Code != 400 || ErrBadRequest.Message != "bad request" {
		t.Errorf("ErrBadRequest not as expected: code=%d msg=%q", ErrBadRequest.Code, ErrBadRequest.Message)
	}
	if ErrServerError.Code != 500 || ErrServerError.Message != "internal server error" {
		t.Errorf("ErrServerError not as expected: code=%d msg=%q", ErrServerError.Code, ErrServerError.Message)
	}
}
