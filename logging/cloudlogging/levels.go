// dui-go/pkg/logging/cloudlogging/levels.go
//
// Package cloudlogging provides integration between Go's slog package and
// Google Cloud Logging's severity levels. This file defines custom log levels
// that align with Google Cloud Logging severities, as well as a utility
// function to map string representations of log levels to these custom
// slog.Level values.
//
// Dependencies:
//   - log/slog: Provides the base logging level types and constants.
//   - This package is intended to be used by a custom slog handler that
//     formats logs for consumption in Google Cloud Logging, ensuring the
//     severity fields map correctly.
//
// Responsibilities of this file:
//   - Define custom slog levels for notice, critical, alert, and emergency
//     to align with Google Cloud Logging severity conventions.
//   - Provide a conversion function from a string (e.g., from an environment
//     variable or configuration) to the corresponding slog.Level.
//   - Serve as a unified reference for severity level mappings, ensuring that
//     the rest of the logging codebase relies on these centralized definitions
//     for consistency and maintainability.

package cloudlogging

import (
	"log/slog"
	"strings" // Required for case-insensitive comparison
)

// Custom log levels for Cloud Logging.
// These constants are offsets from slog's built-in levels (e.g., slog.LevelInfo, slog.LevelError).
// They exist because Google Cloud Logging defines severities such as "NOTICE", "CRITICAL",
// "ALERT", and "EMERGENCY", which are not part of slog's built-in levels.
// By defining these custom levels, we can seamlessly integrate Cloud Logging severities
// into slog-based logging.
const (
	// LevelNotice defines a log level that sits just above INFO and corresponds
	// to the "NOTICE" severity in Cloud Logging. This may be used for normal
	// but significant events that do not yet warrant a WARNING.
	LevelNotice = slog.LevelInfo + 1

	// LevelCritical defines a log level above ERROR that corresponds to the
	// "CRITICAL" severity in Cloud Logging. This may be used for critical issues
	// that require immediate attention.
	LevelCritical = slog.LevelError + 1

	// LevelAlert defines a log level above CRITICAL corresponding to the "ALERT"
	// severity in Cloud Logging. This indicates a condition that must be dealt
	// with immediately. Alerts should trigger immediate action.
	LevelAlert = slog.LevelError + 2

	// LevelEmergency defines a log level above ALERT corresponding to the
	// "EMERGENCY" severity in Cloud Logging. This indicates a system-wide
	// emergency, and human intervention is almost certainly required.
	LevelEmergency = slog.LevelError + 3
)

// StringToLevel takes a string representation of a log severity (e.g., from
// user input, environment variables, or configuration) and returns the
// corresponding slog.Level. This includes both standard slog levels and
// the custom Cloud Logging levels defined above.
//
// The function ensures that external configurations can easily match a string
// severity like "NOTICE" or "CRITICAL" to the appropriate slog level, allowing
// flexible configuration of log thresholds or severity mappings. The comparison
// is case-insensitive.
//
// If the provided string does not match any known severity, it defaults to
// slog.LevelInfo for safety. This default ensures that logs are produced with
// at least a known level, rather than failing silently.
func StringToLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "NOTICE":
		return LevelNotice
	case "WARN", "WARNING": // Handle both common abbreviations
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	case "CRITICAL":
		return LevelCritical
	case "ALERT":
		return LevelAlert
	case "EMERGENCY":
		return LevelEmergency
	default:
		return slog.LevelInfo // Default fallback
	}
}
