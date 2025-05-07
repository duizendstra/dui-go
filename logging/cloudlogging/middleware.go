// dui-go/logging/cloudlogging/middleware.go
//
// Responsibilities of this file:
//   - Define an HTTP middleware (WithCloudTraceContext) that extracts X-Cloud-Trace-Context
//     header data and attaches trace info to the request context.
//
// If the GCP project ID cannot be determined, it defaults to "unknown-project".
// If the X-Cloud-Trace-Context header is missing or invalid, the trace is set to
// a fallback ("projects/unknown-project/traces/") without a valid trace ID.

package cloudlogging

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"cloud.google.com/go/compute/metadata"
)

type traceKey struct{}
type spanIDKey struct{}
type traceSampledKey struct{}

var reCloudTraceContext = regexp.MustCompile(
	`^([a-f\d]+)(?:/([a-f\d]+))?;o=(\d+)$`,
)

func deconstructXCloudTraceContext(s string) (traceID, spanID string, traceSampled bool) {
	matches := reCloudTraceContext.FindStringSubmatch(s)
	if len(matches) == 4 {
		traceID = matches[1]
		spanID = matches[2] // may be empty
		traceSampled = matches[3] == "1"
	}
	return
}

// WithCloudTraceContext extracts trace info from the X-Cloud-Trace-Context
// header and attaches it to the request context.
func WithCloudTraceContext(h http.Handler) http.Handler {
	projectID, err := metadata.ProjectIDWithContext(context.Background())
	if err != nil {
		// Error retrieving from metadata server, try environment variable fallback
		if envProjectID := os.Getenv("GOOGLE_CLOUD_PROJECT"); envProjectID != "" {
			projectID = envProjectID
		} else {
			projectID = "unknown-project" // Fallback if neither metadata nor env var is available
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceHeader := r.Header.Get("X-Cloud-Trace-Context")
		traceID, spanID, traceSampled := deconstructXCloudTraceContext(traceHeader)

		// Ensure traceID is not empty to prevent "projects/project-id/traces/"
		effectiveTraceID := traceID
		if effectiveTraceID == "" {
			effectiveTraceID = "unknown-trace-id" // Or some other placeholder if traceID is empty
		}
		trace := fmt.Sprintf("projects/%s/traces/%s", projectID, effectiveTraceID)

		ctx := r.Context()
		ctx = context.WithValue(ctx, traceKey{}, trace)
		if spanID != "" { // Only add spanIDKey if spanID is present
			ctx = context.WithValue(ctx, spanIDKey{}, spanID)
		}
		ctx = context.WithValue(ctx, traceSampledKey{}, traceSampled)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
