package logger

import (
	"log/slog"
	"testing"
)

func TestHandler_Handle(t *testing.T) {
	logger := slog.New(NewConsoleHandler(&slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	}))

	if logger == nil {
		t.Fatal("logger is nil")
	}

	logger.Debug("logger initiated")
}
