package hvac

import (
	"fmt"
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
	LastUpdate  time.Time
	Name        string
	ShellyID    string
	Temperature DegF
	ID          RoomID
	Zone        ZoneID
	Humidity    uint8
	Battery     uint8
	Occupied    bool
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

	// determine if the ZONE is occupied
	zoneOccupied := false
	for k := range c.Rooms {
		if c.Rooms[k].Zone == zone.ID && c.Rooms[k].Occupied {
			zoneOccupied = true
			break
		}
	}

	// calculate the average temperature of all rooms in the zone
	{
		var avgCnt uint8
		var avgTot DegF
		var avg DegF
		maxAge := time.Now().Add(0 - tempMaxAge)
		for k := range c.Rooms {
			// in the zone, not zero, and more recent than tempMaxAge
			if c.Rooms[k].Zone == zone.ID && c.Rooms[k].Temperature != 0 && c.Rooms[k].LastUpdate.After(maxAge) {
				avgCnt++
				avgTot += c.Rooms[k].Temperature
			}
			avg = avgTot / DegF(avgCnt)
		}
		if avg != 0 {
			zone.AverageTemp = avg
		}
	}
	log.Debug("room temp", "room", r.Name, "zone", zone.ID, "room temp", r.Temperature, "zone avg", zone.AverageTemp)

	switch c.SystemMode {
	case SystemModeHeat:
		if r.Temperature >= boilerLockoutTemp && !c.BoilerLockout {
			log.Warn("locking out boiler, room temp too high")
			c.BoilerLockout = true
			for k := range c.Pumps {
				c.Pumps[k].ID.Stop("lockout")
			}
		}
		if (zoneOccupied && zone.AverageTemp < zone.Targets.HeatingOccupiedTemp-zoneHysterisisRange) || (!zoneOccupied && zone.AverageTemp < zone.Targets.HeatingUnoccupiedTemp-zoneHysterisisRange) {
			if c.ControlMode == ControlTemp {
				log.Info("starting zone", "zone", zone.ID, "avg temp", zone.AverageTemp)
				_ = zone.ID.Start(defaultRunDuration, "temp")
			}
			return
		}
		if (zoneOccupied && zone.AverageTemp > zone.Targets.HeatingOccupiedTemp+zoneHysterisisRange) || (!zoneOccupied && zone.AverageTemp > zone.Targets.HeatingUnoccupiedTemp+zoneHysterisisRange) {
			if c.ControlMode == ControlTemp {
				log.Info("stopping zone", "zone", zone.ID, "avg temp", zone.AverageTemp)
				zone.ID.Stop("temp")
			}
			return
		}
	case SystemModeCool:
		if r.Temperature <= chillerLockoutTemp && !c.ChillerLockout {
			log.Warn("locking out chiller, room temp too low", "room", r.ID, "temp", r.Temperature)
			c.ChillerLockout = true
			for k := range c.Chillers {
				c.Chillers[k].ID.Stop("lockout")
			}
		}
		if (zoneOccupied && zone.AverageTemp > zone.Targets.CoolingOccupiedTemp+zoneHysterisisRange) || (!zoneOccupied && zone.AverageTemp > zone.Targets.CoolingUnoccupiedTemp+zoneHysterisisRange) {
			if c.ControlMode == ControlTemp {
				log.Info("starting zone", "zone", zone.ID, "avg temp", zone.AverageTemp)
				_ = zone.ID.Start(defaultRunDuration, "temp")
			}
			return
		}
		if (zoneOccupied && zone.AverageTemp < zone.Targets.CoolingOccupiedTemp-zoneHysterisisRange) || (!zoneOccupied && zone.AverageTemp < zone.Targets.CoolingUnoccupiedTemp-zoneHysterisisRange) {
			if c.ControlMode == ControlTemp {
				log.Info("stopping zone", "zone", zone.ID, "avg temp", zone.AverageTemp)
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

func (r RoomID) getPreRunTime() (time.Duration, error) {
	room := r.Get()
	if room == nil {
		err := fmt.Errorf("invalid room")
		return 0, err
	}
	zid := room.Zone
	z := zid.Get()
	if z == nil {
		err := fmt.Errorf("invalid zone")
		return 0, err
	}

	var tempDiff DegF
	if c.SystemMode == SystemModeHeat {
		tempDiff = z.Targets.HeatingOccupiedTemp - z.Targets.HeatingUnoccupiedTemp
	} else {
		tempDiff = z.Targets.CoolingUnoccupiedTemp - z.Targets.CoolingOccupiedTemp
	}

	if z.OneDegreeAdjTime == 0 {
		var err error
		z.OneDegreeAdjTime, err = zid.estimateOneDegAdjTime()
		if err != nil {
			log.Error(err.Error())
			return 0, err
		}
	}
	t := time.Duration(tempDiff) * z.OneDegreeAdjTime
	log.Debug("Zone occupancy range", "tempDiff", tempDiff, "zone 1 degree adjustment time", z.OneDegreeAdjTime, "time required", t)
	return t, nil
}
