package cloudlogging

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"cloud.google.com/go/compute/metadata"
)

// --- Interface for Project ID Fetching (for testability) ---
type projectIDFetcher interface {
	ProjectID(ctx context.Context) (string, error)
}

// metadataProjectIDFetcher implements the fetcher using the real metadata service.
type metadataProjectIDFetcher struct{}

func (f *metadataProjectIDFetcher) ProjectID(ctx context.Context) (string, error) {
	return metadata.ProjectIDWithContext(ctx)
}

// --- Middleware Globals and Setup ---
type traceKey struct{}
type spanIDKey struct{}
type traceSampledKey struct{}

var (
	determinedProjectID string
	projectIDOnce       sync.Once
	// fetcher is the current implementation for getting the project ID.
	// It defaults to the real metadata service but can be replaced for tests.
	fetcher projectIDFetcher = &metadataProjectIDFetcher{}
)

// SetProjectIDFetcher allows tests to replace the default project ID fetcher with a mock.
func SetProjectIDFetcher(f projectIDFetcher) {
	fetcher = f
	// Reset the sync.Once so that the new fetcher will be used on the next call.
	projectIDOnce = sync.Once{}
	determinedProjectID = ""
}

// determineProjectID gets the GCP Project ID, caching the result for performance.
// It checks the metadata service first, then the GOOGLE_CLOUD_PROJECT env var,
// and finally falls back to "unknown-project".
func determineProjectID() string {
	projectIDOnce.Do(func() {
		projID, err := fetcher.ProjectID(context.Background())
		if err == nil && projID != "" {
			determinedProjectID = projID
			return
		}
		if envProjectID := os.Getenv("GOOGLE_CLOUD_PROJECT"); envProjectID != "" {
			determinedProjectID = envProjectID
			return
		}
		determinedProjectID = "unknown-project"
	})
	return determinedProjectID
}

// --- Trace Context Parsing ---

// deconstructXCloudTraceContext uses robust string splitting to parse the trace header.
func deconstructXCloudTraceContext(headerValue string) (traceID, spanID string, traceSampled bool) {
	if headerValue == "" {
		return
	}

	parts := strings.Split(headerValue, ";o=")
	if len(parts) != 2 {
		return
	}
	tracePart := parts[0]
	optionsPart := parts[1]

	traceSampled = optionsPart == "1"

	spanParts := strings.Split(tracePart, "/")
	if len(spanParts) > 0 {
		traceID = spanParts[0]
	}
	if len(spanParts) > 1 {
		spanID = spanParts[1]
	}

	return
}

// --- Middleware and Helper Functions ---

// WithCloudTraceContext is an HTTP middleware that extracts trace information
// from the "X-Cloud-Trace-Context" header and injects it into the request's context.
func WithCloudTraceContext(h http.Handler) http.Handler {
	projectID := determineProjectID()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceHeader := r.Header.Get("X-Cloud-Trace-Context")
		traceID, spanID, traceSampled := deconstructXCloudTraceContext(traceHeader)

		resolvedTraceID := traceID
		if resolvedTraceID == "" {
			resolvedTraceID = "unknown-trace"
		}

		trace := fmt.Sprintf("projects/%s/traces/%s", projectID, resolvedTraceID)

		ctx := r.Context()
		ctx = context.WithValue(ctx, traceKey{}, trace)
		// Only set the spanID and sampled keys if they have meaningful values.
		if spanID != "" {
			ctx = context.WithValue(ctx, spanIDKey{}, spanID)
		}
		if traceSampled {
			ctx = context.WithValue(ctx, traceSampledKey{}, true)
		}

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithTrace creates a new context annotated with a Google Cloud Trace ID.
// It is intended for non-HTTP applications (like Cloud Run Jobs) where a trace
// is not propagated via incoming request headers.
//
// The traceID provided will be used to formulate the full trace string, e.g.,
// "projects/your-project-id/traces/your-trace-id". The project ID is determined
// automatically.
func WithTrace(ctx context.Context, traceID string) context.Context {
	projectID := determineProjectID()
	trace := fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)
	return context.WithValue(ctx, traceKey{}, trace)
}
