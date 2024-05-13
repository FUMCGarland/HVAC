package hvac

import (
	"time"

	"github.com/FUMCGarland/hvac/log"
)

type RoomID uint16

type Room struct {
	ID          RoomID
	Name        string
	Zone        ZoneID
	Temperature DegF
	Humidity    uint8
	LastUpdate  time.Time
	Occupied    bool
}

func (r RoomID) Get() *Room {
	for k := range c.Rooms {
		if c.Rooms[k].ID == r {
			return &c.Rooms[k]
		}
	}
	return nil
}

func (r *Room) SetTemp(temp DegF) {
	r.Temperature = temp
	r.LastUpdate = time.Now()

	zone := r.Zone.Get()
	{
		var avgCnt DegF
		var avgTot DegF
		var avg DegF
		hourAgo := time.Now().Add(0 - time.Hour)
		for k := range c.Rooms {
			// in the zone, and not zero, and more recent than an hour ago
			if c.Rooms[k].Zone == zone.ID && c.Rooms[k].Temperature != 0 && c.Rooms[k].LastUpdate.After(hourAgo) {
				avgCnt++
				avgTot += c.Rooms[k].Temperature
			}
			avg = avgTot / avgCnt
		}
		if avg != 0 {
			temp = avg
		}
	}
	log.Info("room temp", "room", r.Name, "zone", zone.ID, "room temp", r.Temperature, "zone avg", temp)

	switch c.SystemMode {
	case SystemModeHeat:
		if r.Temperature >= boilerLockoutTemp {
			log.Info("locking out boiler, room temp too high")
			boilerLockout = true
		}
		if (r.Occupied && temp < zone.Targets.HeatingOccupiedTemp-zoneHysterisisRange) || (!r.Occupied && temp < zone.Targets.HeatingUnoccupiedTemp-zoneHysterisisRange) {
			if c.ControlMode == ControlTemp {
				zone.ID.Start(defaultRunDuration, "temp")
			}
			return
		}
		if (r.Occupied && temp > zone.Targets.HeatingOccupiedTemp+zoneHysterisisRange) || (!r.Occupied && temp > zone.Targets.HeatingUnoccupiedTemp+zoneHysterisisRange) {
			if c.ControlMode == ControlTemp {
				zone.ID.Stop("temp")
			}
			return
		}
	case SystemModeCool:
		if r.Temperature <= chillerLockoutTemp {
			log.Info("locking out chiller, room temp too low")
			chillerLockout = true
		}
		if (r.Occupied && temp > zone.Targets.CoolingOccupiedTemp+zoneHysterisisRange) || (!r.Occupied && temp > zone.Targets.CoolingUnoccupiedTemp+zoneHysterisisRange) {
			log.Info("starting zone (if in temp control mode)", "zone", zone.ID, "avg temp", temp)
			if c.ControlMode == ControlTemp {
				zone.ID.Start(defaultRunDuration, "temp")
			}
			return
		}
		if (r.Occupied && temp < zone.Targets.CoolingOccupiedTemp-zoneHysterisisRange) || (!r.Occupied && temp < zone.Targets.CoolingUnoccupiedTemp-zoneHysterisisRange) {
			log.Info("stopping zone (if in temp control mode)", "zone", zone.ID, "avg temp", temp)
			if c.ControlMode == ControlTemp {
				zone.ID.Stop("temp")
			}
			return
		}
	}
}

func (r *Room) SetHumidity(humidity uint8) {
	r.Humidity = humidity
	// nothing to do other than accept the update
}
