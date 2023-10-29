package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

type (
	HandleType  uint8
	Environment uint8
)

const (
	CONSOLE_HANDLER HandleType = iota
	TEXT_HANDLER
	JSON_HANDLER
)

const (
	DEVELOPMENT Environment = iota
	PRODUCTION
	RELEASE
)

const (
	_defaultSentryFlushTimeout = 1 * time.Second
)

type Log struct {
	skipCaller int
	slog       *slog.Logger
}

type Options struct {
	Development  bool // Development add development details of machine
	Debug        bool // Debug show debug devel message
	EnableCaller bool // EnableCaller show caller in line code
	SkipCaller   int  // SkipCaller skip caller level of CallerFrames https://github.com/golang/go/issues/59145#issuecomment-1481920720
}

type Logger interface {
	Debug(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	Warn(msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	Fatal(msg string, args ...any)
	FatalContext(ctx context.Context, msg string, args ...any)
	Log(ctx context.Context, level slog.Level, msg string, args ...any)
}

func New(
	handler HandleType,
	loggerOption Options,
) (Logger, error) {
	log := new(Log)
	logger := slog.Default()
	slogHandlerOpt := new(slog.HandlerOptions)
	slogHandlerOpt.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		return a
	}

	if loggerOption.Debug {
		slogHandlerOpt.Level = slog.LevelDebug
	}

	if loggerOption.EnableCaller {
		slogHandlerOpt.AddSource = true
	}

	switch handler {
	case JSON_HANDLER:
		logger = slog.New(slog.NewJSONHandler(os.Stderr, slogHandlerOpt))
	case TEXT_HANDLER:
		logger = slog.New(slog.NewTextHandler(os.Stderr, slogHandlerOpt))
	case CONSOLE_HANDLER:
		logger = slog.New(NewConsoleHandler(slogHandlerOpt))
	}

	if loggerOption.Development {
		buildInfo, _ := debug.ReadBuildInfo()
		logger = logger.With(slog.Group("debug_info",
			slog.String("go_version", buildInfo.GoVersion),
			slog.Int("pid", os.Getpid()),
			slog.String("os", runtime.GOOS),
			slog.String("os_arch", runtime.GOARCH),
		))
	}

	log.slog = logger
	log.skipCaller = loggerOption.SkipCaller

	return log, nil
}

func (l *Log) Debug(msg string, keyValues ...any) {
	l.Log(context.Background(), slog.LevelDebug, msg, keyValues...)
}

func (l *Log) DebugContext(ctx context.Context, msg string, keyValues ...any) {
	l.Log(ctx, slog.LevelDebug, msg, keyValues...)
}

func (l *Log) Info(msg string, keyValues ...any) {
	l.Log(context.Background(), slog.LevelInfo, msg, keyValues...)
}

func (l *Log) InfoContext(ctx context.Context, msg string, keyValues ...any) {
	l.Log(ctx, slog.LevelInfo, msg, keyValues...)
}

func (l *Log) Warn(msg string, keyValues ...any) {
	l.Log(context.Background(), slog.LevelWarn, msg, keyValues...)
}

func (l *Log) WarnContext(ctx context.Context, msg string, keyValues ...any) {
	l.Log(ctx, slog.LevelWarn, msg, keyValues...)
}

func (l *Log) Error(msg string, keyValues ...any) {
	l.Log(context.Background(), slog.LevelError, msg, keyValues...)
}

func (l *Log) ErrorContext(ctx context.Context, msg string, keyValues ...any) {
	l.Log(ctx, slog.LevelError, msg, keyValues...)
}

func (l *Log) Fatal(msg string, keyValues ...any) {
	defer os.Exit(1)
	l.Log(context.Background(), slog.LevelError, msg, keyValues...)
}

func (l *Log) FatalContext(ctx context.Context, msg string, keyValues ...any) {
	defer os.Exit(1)
	l.Log(ctx, slog.LevelError, msg, keyValues...)
}

func (l *Log) Log(ctx context.Context, level slog.Level, msg string, keyValues ...any) {
	var pcs [1]uintptr
	runtime.Callers(l.skipCaller, pcs[:])
	rec := slog.NewRecord(time.Now(), level, msg, pcs[0])
	rec.Add(keyValues...)

	_ = l.slog.Handler().Handle(ctx, rec)
}

func (e Environment) String() string {
	switch e {
	case DEVELOPMENT:
		return "development"
	case PRODUCTION:
		return "production"
	case RELEASE:
		return "release"
	}
	return ""
}
