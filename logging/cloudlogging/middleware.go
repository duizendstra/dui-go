// go-dui/logging/cloudlogging/middleware.go
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
		fmt.Printf("Error retrieving project ID from metadata server: %v\n", err)
		// Try environment variable fallback
		if envProjectID := os.Getenv("GOOGLE_CLOUD_PROJECT"); envProjectID != "" {
			projectID = envProjectID
			fmt.Printf("Using project ID from GOOGLE_CLOUD_PROJECT: %s\n", projectID)
		} else {
			projectID = "unknown-project"
			fmt.Printf("No project ID found in metadata or environment. Using fallback: %s\n", projectID)
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceHeader := r.Header.Get("X-Cloud-Trace-Context")
		traceID, spanID, traceSampled := deconstructXCloudTraceContext(traceHeader)
		trace := fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)

		ctx := r.Context()
		ctx = context.WithValue(ctx, traceKey{}, trace)
		ctx = context.WithValue(ctx, spanIDKey{}, spanID)
		ctx = context.WithValue(ctx, traceSampledKey{}, traceSampled)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
