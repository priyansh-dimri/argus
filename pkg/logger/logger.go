package logger

import (
	"io"
	"log/slog"
	"os"
)

var Log *slog.Logger

func InitLogger() {
	Initialize(os.Stdout)
}

func Initialize(w io.Writer) {
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	Log = slog.New(handler)
	slog.SetDefault(Log)
}

func Info(msg string, args ...any) {
	if Log == nil {
		InitLogger()
	}
	Log.Info(msg, args...)
}

func Error(msg string, err error, args ...any) {
	if Log == nil {
		InitLogger()
	}
	allArgs := append([]any{"error", err}, args...)
	Log.Error(msg, allArgs...)
}

func With(args ...any) *slog.Logger {
	if Log == nil {
		InitLogger()
	}
	return Log.With(args...)
}
