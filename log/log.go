package log

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/smallnest/ringbuffer"
)

// Currently log is just a wrapper around "log/slog", if we need more roubust logging facilities, we can add them here later
// e.g. logging to an MQTT stream
var l *slog.Logger
var lvl *slog.LevelVar

const bufsiz = 1048576

var buf *ringbuffer.RingBuffer

// Start initializes the logging interface and returns a "log/slog"
func Start() *slog.Logger {
	lvl = new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))

	buf = ringbuffer.New(bufsiz)

	l = log
	return log
}

// Get returns the logger itself so it can be used in other subsystems (e.g. MQTT)
func Get() *slog.Logger {
	return l
}

func Debug(title string, args ...interface{}) {
	l.Debug(title, args...)
	fmt.Fprintf(buf, title, args...)
}

func Info(title string, args ...interface{}) {
	l.Info(title, args...)
	fmt.Fprintf(buf, title, args...)
}

func Error(title string, args ...interface{}) {
	l.Error(title, args...)
	fmt.Fprintf(buf, title, args...)
}

func Warn(title string, args ...interface{}) {
	l.Warn(title, args...)
	fmt.Fprintf(buf, title, args...)
}

func Fatal(title string, args ...interface{}) {
	l.Error(title, args...)
	fmt.Fprintf(buf, title, args...)
	panic(title)
}

func EnableDebug() {
	lvl.Set(slog.LevelDebug)
	l.Debug("debugging enabled")
}

func ReadBuf() (string, error) {
	b := make([]byte, 1048576)

	count, err := buf.Read(b)
	if err != nil {
		Error("error", "err", err.Error())
		return "", err
	}
	o := string(b[:count])
	return o, nil
}
