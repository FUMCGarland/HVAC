package datalogger

import (
	"context"
	"log"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/FUMCGarland/hvac"
)

// It's a bit of a kludge to use the log interface for this,
// but lumberjack looks like it will do the rotation and splitting for us

// TODO extend lumberjack and add the ability to add a custom log line on
// rotation

// log file format
// date, outsideTemp, blower[n].Running, Pump[n].Running, Chiller[n].Running, Room[n].Temp, Room[n].Humidity, Room[n].Target(fromZone)

func dataLogger(ctx context.Context, datadir string) {
	c := hvac.GetConfig()

	log.SetOutput(&lumberjack.Logger{
		Filename:   c.DataLogFile,
		MaxSize:    100, // megabytes
		MaxBackups: 30,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	writeHeader(c)
	for {
		select {
		case <-ticker.C:
			writeLine(c)
		case <-ctx.Done():
			break
		}
	}
}

func writeHeader(c *hvac.config) {
}

func writeLine(c *hvac.Config) {
}
