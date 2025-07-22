// dui-go/logging/cloudlogging/doc.go
//
// Package cloudlogging provides middleware and handlers to integrate Go's
// standard HTTP handling and logging with Google Cloud Logging and Cloud Trace.
//
// The WithCloudTraceContext middleware inspects the X-Cloud-Trace-Context
// header from incoming HTTP requests, extracting trace and span information.
// This trace context is stored in the request's context so that subsequent
// handlers and logging calls can produce logs correlated with a specific Cloud
// Trace span.
//
// The CloudLoggingHandler is a custom slog.Handler that formats logs in a manner
// compatible with Google Cloud Logging. It sets appropriate severity fields,
// injects trace and span data if available, and automatically includes structured
// source location details for all log entries.
//
// # Typical Usage (HTTP Services)
//
// For a standard web server, wrap your handler with the middleware. The trace
// context will be automatically propagated from incoming requests.
//
//	http.Handle("/", WithCloudTraceContext(yourHandler))
//	logger := slog.New(NewCloudLoggingHandler("my-service"))
//	slog.SetDefault(logger)
//
// # Usage in Google Cloud Run Jobs
//
// Cloud Run Jobs are not triggered by HTTP requests, so they do not have an
// incoming X-Cloud-Trace-Context header. To correlate all logs from a single
// job execution, you should manually create a trace ID and inject it into the
// context at the start of your job using the `WithTrace` helper function.
//
// This pattern ensures all logs produced during the job run can be easily
// filtered and viewed together in the Cloud Logging UI.
//
//	func main() {
//		// 1. Set up the production-ready logger.
//		logger := slog.New(cloudlogging.NewCloudLoggingHandler("my-cloud-run-job"))
//		slog.SetDefault(logger)
//
//		// 2. Create a unique trace ID for this execution.
//		// We can use the built-in execution ID from Cloud Run or generate a UUID.
//		executionID := os.Getenv("CLOUD_RUN_EXECUTION")
//		if executionID == "" {
//			executionID = uuid.NewString() // Requires github.com/google/uuid
//		}
//
//		// 3. Create a root context for this job and add the trace.
//		ctx := cloudlogging.WithTrace(context.Background(), executionID)
//
//		// 4. Run the job's business logic, passing the context down.
//		slog.InfoContext(ctx, "Job starting.")
//		// ... run job logic with ctx ...
//		slog.InfoContext(ctx, "Job finished successfully.")
//	}
package cloudlogging
