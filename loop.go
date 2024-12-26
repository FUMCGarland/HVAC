package hvac

import (
	"time"
)

// LoopID is the unique identifier for a loop
type LoopID uint8

// A loop is a physical pipe that moves heated or cooled water to blowers or radiant devices around the campus. Loops connect blowers and pumps; some zones are radiant only and do not have blowers.
type Loop struct {
	Name        string
	ID          LoopID
	RadiantZone ZoneID
}

func (l LoopID) Start(d time.Duration, msg string) error {
	pump := c.getPumpFromLoop(l)
	if pump != 0 {
		return pump.Start(d, msg)
	}
	return nil
}
