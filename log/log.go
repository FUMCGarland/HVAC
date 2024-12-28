package log

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

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

func Debug(title string, args ...any) {
	l.Debug(title, args...)
	// WriteToBuff(title, args...)
}

func Info(title string, args ...interface{}) {
	l.Info(title, args...)
	WriteToBuff(title, args...)
}

func Error(title string, args ...interface{}) {
	l.Error(title, args...)
	WriteToBuff(title, args...)
}

func Warn(title string, args ...interface{}) {
	l.Warn(title, args...)
	WriteToBuff(title, args...)
}

func Fatal(title string, args ...interface{}) {
	l.Error(title, args...)
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

func WriteToBuff(title string, args ...any) {
	var out strings.Builder

	if strings.Contains(title, "client disconnected") {
		return
	}

	now := time.Now()

	out.WriteString(now.Format(time.DateTime))
	out.WriteString(" : ")
	out.WriteString(title)
	out.WriteString(" ")
	pos := 0
	for _, v := range args {
		switch v := v.(type) {
		case string:
			out.WriteString(v)
			if pos%2 == 0 {
				out.WriteString("=")
			} else {
				out.WriteString("; ")
			}
		case []byte:
			out.WriteString(string(v))
		case int, uint8:
			out.WriteString(fmt.Sprintf("%d; ", v))
		default:
			out.WriteString(fmt.Sprintf("%v; ", v))
		}
		pos = pos + 1
	}
	out.WriteString("\n")
	built := []byte(out.String())
	if _, err := buf.TryWrite(built); err != nil {
		buf.Reset() // if the buffer is full, dump it
		Debug("reset log buffer")
	}
}
