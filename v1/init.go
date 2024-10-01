package cplugin

import (
	"log/slog"
	"os"
	"sync"
)

var (
	pluginCache   = make(map[string]any)
	mutex         sync.RWMutex
	defaultLogger Logger
)

func init() {
	defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))
}

func SetLogger(l Logger) {
	defaultLogger = l
}

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
