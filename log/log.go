package log

import (
	"log/slog"
	"os"
)

// Currently log is just a wrapper around "log/slog", if we need more roubust logging facilities, we can add them here later
// e.g. logging to an MQTT stream
var l *slog.Logger
var lvl *slog.LevelVar

// Start initializes the logging interface and returns a "log/slog"
func Start() *slog.Logger {
	lvl = new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))
	// create a log destination on MQTT?

	l = log
	return log
}

// Get returns the logger itself so it can be used in other subsystems (e.g. MQTT)
func Get() *slog.Logger {
	return l
}

func Debug(title string, args ...interface{}) {
	l.Debug(title, args...)
}

func Info(title string, args ...interface{}) {
	l.Info(title, args...)
}

func Error(title string, args ...interface{}) {
	l.Error(title, args...)
}

func Warn(title string, args ...interface{}) {
	l.Warn(title, args...)
}

func Fatal(title string, args ...interface{}) {
	l.Error(title, args...)
	panic(title)
}

func EnableDebug() {
	lvl.Set(slog.LevelDebug)
	l.Debug("debugging enabled")
}
