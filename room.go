package hvac

import (
	"strings"
	"time"

	"github.com/FUMCGarland/hvac/log"
)

// RoomID is the room number or other identifying number
// uint16 since we've got 3 floors and a uint8 would not be enough
type RoomID uint16

// Room is the basic data type for a physical space
// temperature, humidity, and occupancy are tracked per-room
// device control is per-zone. This will allow us to track how long different
// rooms take to bring to proper temperature so we can set start-times
// properly. When we start building controls for the damnpers and valves
// having per-room data will help with system tuning even more
type Room struct {
	ID          RoomID
	Name        string
	Zone        ZoneID
	Temperature DegF
	Humidity    uint8
	Battery     uint8
	LastUpdate  time.Time
	Occupied    bool
	ShellyID    string
}

// Get returns a full room struct for a RoomID
func (r RoomID) Get() *Room {
	for k := range c.Rooms {
		if c.Rooms[k].ID == r {
			return &c.Rooms[k]
		}
	}
	return nil
}

// SetTemp records the temperature as reported by the sensors, called from the MQTT subsystem
// if the zone average is out-of-range, the proper devices are enabled to bring the zone into temperature
func (r *Room) SetTemp(temp DegF) {
	r.Temperature = temp
	r.LastUpdate = time.Now()

	zone := r.Zone.Get()
	{
		var avgCnt uint8
		var avgTot DegF
		var avg DegF
		maxAge := time.Now().Add(0 - tempMaxAge)
		for k := range c.Rooms {
			// in the zone, and not zero, and more recent than an hour ago
			if c.Rooms[k].Zone == zone.ID && c.Rooms[k].Temperature != 0 && c.Rooms[k].LastUpdate.After(maxAge) {
				avgCnt++
				avgTot += c.Rooms[k].Temperature
			}
			avg = avgTot / DegF(avgCnt)
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
			for k := range c.Pumps {
				c.Pumps[k].ID.Stop("lockout")
			}
		}
		if (r.Occupied && temp < zone.Targets.HeatingOccupiedTemp-zoneHysterisisRange) || (!r.Occupied && temp < zone.Targets.HeatingUnoccupiedTemp-zoneHysterisisRange) {
			if c.ControlMode == ControlTemp {
				_ = zone.ID.Start(defaultRunDuration, "temp")
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
			for k := range c.Chillers {
				c.Chillers[k].ID.Stop("lockout")
			}
		}
		if (r.Occupied && temp > zone.Targets.CoolingOccupiedTemp+zoneHysterisisRange) || (!r.Occupied && temp > zone.Targets.CoolingUnoccupiedTemp+zoneHysterisisRange) {
			log.Info("starting zone (if in temp control mode)", "zone", zone.ID, "avg temp", temp)
			if c.ControlMode == ControlTemp {
				_ = zone.ID.Start(defaultRunDuration, "temp")
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

// SetHumidity records the humidity as reported by the sensors, called from MQTT subsystem
func (r *Room) SetHumidity(humidity uint8) {
	r.Humidity = humidity
}

// SetBattery records the battery status as reported by the sensors, called from MQTT subsystem
func (r *Room) SetBattery(battery uint8) {
	r.Battery = battery
}

// GetRoomIDFromShelly returns a RoomID based on an associated (case insensitive) shelly ID
func GetRoomIDFromShelly(shellyID string) RoomID {
	for k := range c.Rooms {
		if strings.EqualFold(c.Rooms[k].ShellyID, shellyID) {
			return c.Rooms[k].ID
		}
	}

	return 0
}
