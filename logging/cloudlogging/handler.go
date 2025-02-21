// go-dui/logging/cloudlogging/handler.go
//
// Responsibilities of this file:
//   - Define a CloudLoggingHandler that translates slog Levels to Cloud Logging severities.
//   - Integrate trace and span info from the request context into logs.
//   - Optionally add source location details for error-level logs.
//
// Use this handler in conjunction with the WithCloudTraceContext middleware.

package cloudlogging

import (
	"context"
	"log/slog"
	"os"
	"runtime"
)

type CloudLoggingHandler struct {
	handler slog.Handler
}

func NewCloudLoggingHandler(component string) *CloudLoggingHandler {
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
		levelVar.Set(slog.LevelInfo)
	}

	baseHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: false,
		Level:     &levelVar,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			} else if a.Key == slog.SourceKey {
				a.Key = "logging.googleapis.com/sourceLocation"
			} else if a.Key == slog.LevelKey {
				a.Key = "severity"
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
			}
			return a
		},
	})

	handlerWithAttrs := baseHandler.WithAttrs([]slog.Attr{
		slog.String("component", component),
	})

	return &CloudLoggingHandler{handler: handlerWithAttrs}
}

func (h *CloudLoggingHandler) Handle(ctx context.Context, rec slog.Record) error {
	rec = rec.Clone()

	trace, _ := ctx.Value(traceKey{}).(string)
	if trace == "" {
		trace = "projects/unknown-project/traces/unknown-trace"
	}
	rec.Add("logging.googleapis.com/trace", slog.StringValue(trace))

	if spanID, ok := ctx.Value(spanIDKey{}).(string); ok {
		rec.Add("logging.googleapis.com/spanId", slog.StringValue(spanID))
	}

	if traceSampled, ok := ctx.Value(traceSampledKey{}).(bool); ok {
		rec.Add("logging.googleapis.com/trace_sampled", slog.BoolValue(traceSampled))
	} else {
		rec.Add("logging.googleapis.com/trace_sampled", slog.BoolValue(false))
	}

	if rec.Level >= slog.LevelError {
		pc, file, line, ok := runtime.Caller(3)
		if ok {
			funcName := runtime.FuncForPC(pc).Name()
			rec.Add("logging.googleapis.com/sourceLocation", slog.GroupValue(
				slog.String("file", file),
				slog.Int("line", line),
				slog.String("function", funcName),
			))
		}
	}

	return h.handler.Handle(ctx, rec)
}

func (h *CloudLoggingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CloudLoggingHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *CloudLoggingHandler) WithGroup(name string) slog.Handler {
	return &CloudLoggingHandler{handler: h.handler.WithGroup(name)}
}

func (h *CloudLoggingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}
