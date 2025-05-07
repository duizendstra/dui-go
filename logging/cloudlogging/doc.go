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
// injects trace and span data if available, and can optionally include source
// location details for error-level and above logs.
//
// Typical usage:
//   http.Handle("/", WithCloudTraceContext(yourHandler))
//   logger := slog.New(NewCloudLoggingHandler("my-service"))
//   slog.SetDefault(logger)
//
// When logs are produced (e.g., slog.InfoContext), the handler references the
// trace information in the request context to annotate logs with trace fields,
// enabling better correlation and debugging in Cloud Logging and Cloud Trace.

package cloudlogging
