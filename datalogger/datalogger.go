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
			return
		}
	}
}

func writeHeader(c *hvac.Config) {
	var b strings.Builder

	b.WriteString("Date,OutsideTemp,OusideHumidity")
	for k := range c.Blowers {
		b.WriteString(",")
		b.WriteString(c.Blowers[k].Name)
		b.WriteString(" Running")
	}

	for k := range c.Pumps {
		b.WriteString(",")
		b.WriteString(c.Pumps[k].Name)
		b.WriteString(" Running")
	}

	for k := range c.Chillers {
		b.WriteString(",")
		b.WriteString(c.Chillers[k].Name)
		b.WriteString(" Running")
	}

	for k := range c.Rooms {
		b.WriteString(",")
		b.WriteString(c.Rooms[k].Name)
		b.WriteString(" Temp")
		b.WriteString(",")
		b.WriteString(c.Rooms[k].Name)
		b.WriteString(" Humidity")
		b.WriteString(",")
		b.WriteString(c.Rooms[k].Name)
		b.WriteString(" Target")
	}

	log.Println(b.String())
}

func writeLine(c *hvac.Config) {
	var b strings.Builder

	b.WriteString(time.Now().String())
	t, h := getOutsideTemp(c)
	b.WriteString(fmt.Sprintf(",%.2f,%d", t, h))
	for k := range c.Blowers {
		b.WriteString(",")
		b.WriteString(boolstr(c.Blowers[k].Running))
	}

	for k := range c.Pumps {
		b.WriteString(",")
		b.WriteString(boolstr(c.Pumps[k].Running))
	}

	for k := range c.Chillers {
		b.WriteString(",")
		b.WriteString(boolstr(c.Chillers[k].Running))
	}

	for k := range c.Rooms {
		b.WriteString(",")
		b.WriteString(fmt.Sprintf("%d", c.Rooms[k].Temperature))
		b.WriteString(",")
		b.WriteString(fmt.Sprintf("%d", c.Rooms[k].Humidity))
		b.WriteString(",")
		b.WriteString(roomTarget(c.Rooms[k], c))
	}

	log.Println(b.String())
}

func boolstr(b bool) string {
	if b {
		return "True"
	}
	return "False"
}

func roomTarget(r hvac.Room, c *hvac.Config) string {
	for k := range c.Zones {
		if c.Zones[k].ID == r.Zone {
			if c.SystemMode == hvac.SystemModeHeat {
				if r.Occupied {
					return fmt.Sprintf("%d", c.Zones[k].Targets.HeatingOccupiedTemp)
				} else {
					return fmt.Sprintf("%d", c.Zones[k].Targets.HeatingUnoccupiedTemp)
				}
			} else {
				if r.Occupied {
					return fmt.Sprintf("%d", c.Zones[k].Targets.CoolingOccupiedTemp)
				} else {
					return fmt.Sprintf("%d", c.Zones[k].Targets.CoolingOccupiedTemp)
				}
			}
		}
	}

	return "0"
}
