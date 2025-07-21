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
type CloudLoggingHandler struct {
	handler slog.Handler
}

// newCloudLoggingHandlerWithWriter is an internal constructor that allows specifying the output writer.
// This is primarily for testing.
func newCloudLoggingHandlerWithWriter(out io.Writer, component string) *CloudLoggingHandler {
	var levelVar slog.LevelVar
	envLogLevel := os.Getenv("LOG_LEVEL")
	levelVar.Set(StringToLevel(envLogLevel))

	baseHandler := slog.NewJSONHandler(out, &slog.HandlerOptions{
		AddSource:   true, // Use the robust, built-in source capture.
		Level:       &levelVar,
		ReplaceAttr: cloudLoggingReplaceAttr,
	})

	handlerWithAttrs := baseHandler.WithAttrs([]slog.Attr{
		slog.String("component", component),
	})

	return &CloudLoggingHandler{handler: handlerWithAttrs}
}

// NewCloudLoggingHandler creates a new CloudLoggingHandler that writes to os.Stderr.
func NewCloudLoggingHandler(component string) *CloudLoggingHandler {
	return newCloudLoggingHandlerWithWriter(os.Stderr, component)
}

// Handle processes a log record, adding Cloud Trace context before passing it
// to the underlying JSON handler.
func (h *CloudLoggingHandler) Handle(ctx context.Context, rec slog.Record) error {
	rec = rec.Clone()

	// Determine the project ID using the cached, testable function from middleware.go
	projectID := determineProjectID()

	traceVal, _ := ctx.Value(traceKey{}).(string)
	spanIDVal, _ := ctx.Value(spanIDKey{}).(string)
	sampledVal, _ := ctx.Value(traceSampledKey{}).(bool)

	if traceVal == "" {
		traceVal = "projects/" + projectID + "/traces/unknown-trace"
	}

	rec.AddAttrs(slog.String("logging.googleapis.com/trace", traceVal))
	if spanIDVal != "" {
		rec.AddAttrs(slog.String("logging.googleapis.com/spanId", spanIDVal))
	}
	rec.AddAttrs(slog.Bool("logging.googleapis.com/trace_sampled", sampledVal))

	return h.handler.Handle(ctx, rec)
}

// WithAttrs delegates to the wrapped handler, returning a new CloudLoggingHandler.
func (h *CloudLoggingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CloudLoggingHandler{handler: h.handler.WithAttrs(attrs)}
}

// WithGroup delegates to the wrapped handler, returning a new CloudLoggingHandler.
func (h *CloudLoggingHandler) WithGroup(name string) slog.Handler {
	return &CloudLoggingHandler{handler: h.handler.WithGroup(name)}
}

// Enabled delegates to the wrapped handler.
func (h *CloudLoggingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}