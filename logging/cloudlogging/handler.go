package cloudlogging

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// cloudLoggingReplaceAttr is the function that modifies slog attributes
// for Cloud Logging compatibility.
func cloudLoggingReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.MessageKey:
		a.Key = "message"
	case slog.LevelKey:
		a.Key = "severity"
		levelVal := a.Value.Any().(slog.Level)
		switch levelVal {
		case slog.LevelDebug:
			a.Value = slog.StringValue("DEBUG")
		case slog.LevelInfo:
			a.Value = slog.StringValue("INFO")
		case LevelNotice:
			a.Value = slog.StringValue("NOTICE")
		case slog.LevelWarn:
			a.Value = slog.StringValue("WARNING")
		case slog.LevelError:
			a.Value = slog.StringValue("ERROR")
		case LevelCritical:
			a.Value = slog.StringValue("CRITICAL")
		case LevelAlert:
			a.Value = slog.StringValue("ALERT")
		case LevelEmergency:
			a.Value = slog.StringValue("EMERGENCY")
		default:
			a.Value = slog.StringValue("DEFAULT")
		}
	case slog.SourceKey:
		// The slog handler provides a slog.Source object. We transform it.
		if source, ok := a.Value.Any().(*slog.Source); ok {
			a.Key = "logging.googleapis.com/sourceLocation"
			a.Value = slog.GroupValue(
				slog.String("file", source.File),
				slog.Int("line", source.Line),
				slog.String("function", source.Function),
			)
		}
	}
	return a
}

// CloudLoggingHandler is a custom slog.Handler for Google Cloud Logging.
// It embeds an underlying slog.Handler to automatically delegate methods like
// WithAttrs, WithGroup, and Enabled.
type CloudLoggingHandler struct {
	slog.Handler
}

// NewCloudLoggingHandlerForTest creates a CloudLoggingHandler that writes to the
// provided io.Writer instead of os.Stderr. This is useful for capturing and
// asserting on log output in tests.
func NewCloudLoggingHandlerForTest(out io.Writer, component string) *CloudLoggingHandler {
	var levelVar slog.LevelVar
	envLogLevel := os.Getenv("LOG_LEVEL")
	levelVar.Set(StringToLevel(envLogLevel))

	baseHandler := slog.NewJSONHandler(out, &slog.HandlerOptions{
		AddSource:   true, // Use the robust, built-in source capture.
		Level:       &levelVar,
		ReplaceAttr: cloudLoggingReplaceAttr,
	})

	// Add the component as a permanent attribute to the handler.
	handlerWithAttrs := baseHandler.WithAttrs([]slog.Attr{
		slog.String("component", component),
	})

	return &CloudLoggingHandler{Handler: handlerWithAttrs}
}

// NewCloudLoggingHandler creates a new CloudLoggingHandler that writes to os.Stderr
// for use in production.
func NewCloudLoggingHandler(component string) *CloudLoggingHandler {
	return NewCloudLoggingHandlerForTest(os.Stderr, component)
}

// Handle processes a log record, adding Cloud Trace context before passing it
// to the embedded JSON handler.
func (h *CloudLoggingHandler) Handle(ctx context.Context, rec slog.Record) error {
	rec = rec.Clone()

	// Determine the project ID using the cached, testable function.
	projectID := determineProjectID()

	traceVal, _ := ctx.Value(traceKey{}).(string)
	spanIDVal, _ := ctx.Value(spanIDKey{}).(string)
	sampledVal, _ := ctx.Value(traceSampledKey{}).(bool)

	// If no trace is in the context, create a default one.
	if traceVal == "" {
		traceVal = "projects/" + projectID + "/traces/unknown-trace"
	}

	rec.AddAttrs(slog.String("logging.googleapis.com/trace", traceVal))
	if spanIDVal != "" {
		rec.AddAttrs(slog.String("logging.googleapis.com/spanId", spanIDVal))
	}
	rec.AddAttrs(slog.Bool("logging.googleapis.com/trace_sampled", sampledVal))

	// Delegate the final handling to the embedded handler.
	return h.Handler.Handle(ctx, rec)
}
