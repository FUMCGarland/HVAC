package datalogger

import (
	"context"
	"fmt"
	"log"
	"strings"
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

func DataLogger(ctx context.Context) {
	c := hvac.GetConfig()

	lj := &lumberjack.Logger{
		Filename:   c.DataLogFile,
		MaxSize:    100, // megabytes
		MaxBackups: 30,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
	log.SetOutput(lj)

	// clear and re-set the logger flags so the first line doesn't get the timestamp
	{
		flags := log.Flags()
		log.SetFlags(0)
		writeHeader(c)
		log.SetFlags(flags)
	}

	// every 5 minutes, write another line to the log file
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			writeLine(c)
		case <-ctx.Done():
			// start a new data file for each startup
			// _ = lj.Rotate()
			return
		}
	}
}

func writeHeader(c *hvac.Config) {
	var b strings.Builder

	b.WriteString("Date,Outside Temp,Outside Humidity")
	for _, k := range c.Blowers {
		b.WriteString(",")
		b.WriteString(k.Name)
		b.WriteString(" Running")
	}

	for _, k := range c.Pumps {
		b.WriteString(",")
		b.WriteString(k.Name)
		b.WriteString(" Running")
	}

	for _, k := range c.Chillers {
		b.WriteString(",")
		b.WriteString(k.Name)
		b.WriteString(" Running")
	}

	for _, k := range c.Rooms {
		b.WriteString(",")
		b.WriteString(k.Name)
		b.WriteString(" Temp")
		b.WriteString(",")
		b.WriteString(k.Name)
		b.WriteString(" Humidity")
		b.WriteString(",")
		b.WriteString(k.Name)
		b.WriteString(" Target")
	}

	for _, k := range c.Zones {
		b.WriteString(",")
		b.WriteString(k.Name)
		b.WriteString(" Zone Average Temp")
	}

	log.Println(b.String())
}

func writeLine(c *hvac.Config) {
	var b strings.Builder

	t, h := getOutsideTemp(c)
	b.WriteString(fmt.Sprintf(",%.2f,%d", t, h))
	for _, k := range c.Blowers {
		b.WriteString(",")
		b.WriteString(boolstr(k.Running))
	}

	for _, k := range c.Pumps {
		b.WriteString(",")
		b.WriteString(boolstr(k.Running))
	}

	for _, k := range c.Chillers {
		b.WriteString(",")
		b.WriteString(boolstr(k.Running))
	}

	for _, k := range c.Rooms {
		b.WriteString(",")
		b.WriteString(fmt.Sprintf("%.2f", k.Temperature))
		b.WriteString(",")
		b.WriteString(fmt.Sprintf("%d", k.Humidity))
		b.WriteString(",")
		b.WriteString(roomTarget(k, c))
	}

	for _, k := range c.Zones {
		b.WriteString(",")
		b.WriteString(fmt.Sprintf("%.2f", k.AverageTemp))
	}

	log.Println(b.String())
}

func boolstr(b bool) string {
	if b {
		return "True"
	}
	return "False"
}

func roomTarget(r *hvac.Room, c *hvac.Config) string {
	for _, k := range c.Zones {
		if k.ID == r.GetZoneIDInMode() {
			if c.SystemMode == hvac.SystemModeHeat {
				if r.Occupied {
					return fmt.Sprintf("%.2f", k.Targets.HeatingOccupiedTemp)
				} else {
					return fmt.Sprintf("%.2f", k.Targets.HeatingUnoccupiedTemp)
				}
			} else {
				if r.Occupied {
					return fmt.Sprintf("%.2f", k.Targets.CoolingOccupiedTemp)
				} else {
					return fmt.Sprintf("%.2f", k.Targets.CoolingUnoccupiedTemp)
				}
			}
		}
	}

	return "0"
}
