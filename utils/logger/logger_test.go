package logger

import (
	"context"
	"testing"
)

func setup(t *testing.T) Logger {
	logger, err := New(CONSOLE_HANDLER, Options{
		Development:  true,
		Debug:        true,
		EnableCaller: true,
		SkipCaller:   3,
	})
	if err != nil {
		t.Fatal(err)
	}
	return logger
}

func setupWithSentry(t *testing.T) Logger {
	logger, err := New(CONSOLE_HANDLER, Options{
		Development:  false,
		Debug:        false,
		EnableCaller: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	return logger
}

func TestLog_Error(t *testing.T) {
	logger := setup(t)
	logger.Error("error example", "test", 2, "test2", 2.5)
}

func TestLog_Debug(t *testing.T) {
	logger := setup(t)
	logger.Debug("test", "msg", "hello")
}

func TestLog_Info(t *testing.T) {
	logger := setup(t)
	logger.Info("error example", "test", 2, "test2", 2.5)
}

func TestLog_Warn(t *testing.T) {
	logger := setup(t)
	logger.Warn("error example", "test", 2, "test2", 2.5)
}

func TestLog_ErrorContext(t *testing.T) {
	logger := setup(t)
	logger.ErrorContext(context.TODO(), "error example", "test", 2, "test2", 2.5)
}

func TestLog_DebugContext(t *testing.T) {
	logger := setup(t)
	logger.DebugContext(context.TODO(), "error example", "test", 2, "test2", 2.5)
}

func TestLog_InfoContext(t *testing.T) {
	logger := setup(t)
	logger.InfoContext(context.TODO(), "error example", "test", 2, "test2", 2.5)
}

func TestLog_WarnContext(t *testing.T) {
	logger := setup(t)
	logger.WarnContext(context.TODO(), "error example", "test", 2, "test2", 2.5)
}

func TestLog_Fatal(t *testing.T) {
	logger := setup(t)
	logger.Fatal("error example", "test", 2, "test2", 2.5)
}

func TestLog_FatalContext(t *testing.T) {
	logger := setup(t)
	logger.FatalContext(context.TODO(), "error example", "test", 2, "test2", 2.5)
}
