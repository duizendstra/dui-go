package cloudlogging

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type CloudLoggingHandler struct {
	handler slog.Handler // The underlying handler (likely slog.JSONHandler with added attributes)
}

// newCloudLoggingHandlerWithWriter is an internal constructor that allows specifying the output writer.
// This is primarily for testing.
func newCloudLoggingHandlerWithWriter(out io.Writer, component string) *CloudLoggingHandler {
	var levelVar slog.LevelVar

	envLogLevel := os.Getenv("LOG_LEVEL")
	switch envLogLevel {
	case "DEBUG":
		levelVar.Set(slog.LevelDebug)
	case "INFO":
		levelVar.Set(slog.LevelInfo)
	case "WARN", "WARNING":
		levelVar.Set(slog.LevelWarn)
	case "ERROR":
		levelVar.Set(slog.LevelError)
	default:
		levelVar.Set(slog.LevelInfo) // Default to Info
	}

	// Configure the base JSON handler
	opts := &slog.HandlerOptions{
		AddSource: true, // <--- Enable source addition in the base handler
		Level:     &levelVar,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.MessageKey:
				a.Key = "message" // Rename message key
			case slog.LevelKey:
				a.Key = "severity" // Rename level key and map values
				level := a.Value.Any().(slog.Level)
				switch level {
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
				// The base handler added the source, now we rename and restructure it.
				if source, ok := a.Value.Any().(*slog.Source); ok && source != nil {
					// Create the structured sourceLocation attribute for Cloud Logging
					return slog.Attr{
						Key: "logging.googleapis.com/sourceLocation",
						// Use GroupValue to represent the nested structure
						Value: slog.GroupValue(
							slog.String("file", source.File),
							slog.Int("line", source.Line),
							slog.String("function", source.Function),
						),
					}
				}
				// If casting fails or source is nil, drop the attribute by returning an empty Attr
				return slog.Attr{}
			}
			return a
		},
	}

	// Create the base handler
	baseHandler := slog.NewJSONHandler(out, opts)

	// Prepend the component attribute using WithAttrs
	// This returns a handler, not necessarily a *slog.JSONHandler
	handlerWithComponent := baseHandler.WithAttrs([]slog.Attr{
		slog.String("component", component),
	})

	// Return our wrapper. The wrapped handler now includes the component attribute
	// and has the ReplaceAttr logic applied.
	return &CloudLoggingHandler{handler: handlerWithComponent}
}

// NewCloudLoggingHandler creates a new CloudLoggingHandler that writes to os.Stderr.
func NewCloudLoggingHandler(component string) *CloudLoggingHandler {
	return newCloudLoggingHandlerWithWriter(os.Stderr, component)
}

// Handle now primarily focuses on adding trace information, as source location
// is handled by the underlying JSONHandler and formatted by ReplaceAttr.
func (h *CloudLoggingHandler) Handle(ctx context.Context, rec slog.Record) error {
	// It's generally safer to clone if we are adding attributes, even if
	// the underlying handler might also clone.
	newRec := rec.Clone()

	// Add trace information
	trace, _ := ctx.Value(traceKey{}).(string)
	if trace == "" {
		trace = "projects/unknown-project/traces/unknown-trace"
	}
	newRec.AddAttrs(slog.String("logging.googleapis.com/trace", trace))

	if spanID, ok := ctx.Value(spanIDKey{}).(string); ok && spanID != "" {
		newRec.AddAttrs(slog.String("logging.googleapis.com/spanId", spanID))
	}

	traceSampled, ok := ctx.Value(traceSampledKey{}).(bool)
	newRec.AddAttrs(slog.Bool("logging.googleapis.com/trace_sampled", ok && traceSampled))

	// Source location logic is removed from here.
	// The underlying handler (created with AddSource: true) adds the source info,
	// and ReplaceAttr formats the key and structure.

	// Delegate to the wrapped handler (which includes component attr and ReplaceAttr logic)
	return h.handler.Handle(ctx, newRec)
}

func (h *CloudLoggingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Ensure the CloudLoggingHandler wrapper type is preserved
	return &CloudLoggingHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *CloudLoggingHandler) WithGroup(name string) slog.Handler {
	// Ensure the CloudLoggingHandler wrapper type is preserved
	return &CloudLoggingHandler{handler: h.handler.WithGroup(name)}
}

func (h *CloudLoggingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Delegate enabling check to the wrapped handler
	return h.handler.Enabled(ctx, level)
}
