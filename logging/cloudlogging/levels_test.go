// dui-go/pkg/logging/cloudlogging/levels_test.go

package cloudlogging

import (
	"log/slog"
	"testing"
)

func TestStringToLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"DEBUG", slog.LevelDebug},
		{"INFO", slog.LevelInfo},
		{"NOTICE", LevelNotice},
		{"WARNING", slog.LevelWarn},
		{"ERROR", slog.LevelError},
		{"CRITICAL", LevelCritical},
		{"ALERT", LevelAlert},
		{"EMERGENCY", LevelEmergency},
		{"UNKNOWN", slog.LevelInfo}, // Default
	}

	for _, tt := range tests {
		got := StringToLevel(tt.input)
		if got != tt.expected {
			t.Errorf("StringToLevel(%q) = %d; want %d", tt.input, got, tt.expected)
		}
	}
}
