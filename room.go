package hvac

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
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

	log.Debug("room temp", "room", r.ID, "temp", r.Temperature)

	// lock-outs for extreme room temps
	switch c.SystemMode {
	case SystemModeHeat:
		if r.Temperature >= boilerLockoutTemp && !c.BoilerLockout {
			log.Warn("locking out boiler, room temp too high")
			c.BoilerLockout = true
			for k := range c.Pumps {
				c.Pumps[k].ID.Stop("lockout")
			}
		}
	case SystemModeCool:
		if r.Temperature <= chillerLockoutTemp && !c.ChillerLockout {
			log.Warn("locking out chiller, room temp too low", "room", r.ID, "temp", r.Temperature)
			c.ChillerLockout = true
			for k := range c.Chillers {
				c.Chillers[k].ID.Stop("lockout")
			}
		}
	}

	// update the zone average and run logic on starting/stopping the zone
	zone := r.Zone.Get()
	zone.UpdateTemp()
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
	return 0, nil

	// disable for now
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
	log.Debug("Zone occupancy range", "tempDiff", tempDiff, "zone 1 degree adjustment time", z.OneDegreeAdjTime, "pre run minutes required", t)
	return t, nil
}

func (r *Room) writeToStore() error {
	path := path.Join(c.StateStore, fmt.Sprintf("room-%d.json", r.ID))
	log.Debug("writing room data", "file", path)

	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	j, err := json.Marshal(*r)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if _, err := fd.Write(j); err != nil {
		log.Error(err.Error())
		return err
	}

	if err := fd.Close(); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (r *Room) readFromStore() error {
	path := path.Join(c.StateStore, fmt.Sprintf("room-%d.json", r.ID))

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error(err.Error()) // log and ignore
		return nil
	}

	var in Room
	if err = json.Unmarshal(data, &in); err != nil {
		log.Error(err.Error())
		return err
	}

	r.Temperature = in.Temperature
	r.Humidity = in.Humidity
	r.Battery = in.Battery

	return nil
}

func (r RoomID) ToogleOccupancy() {
	room := r.Get()
	if room == nil {
		return
	}

	room.Occupied = !room.Occupied
	log.Debug("manually set room occupancy", "room", r, "state", room.Occupied)
	if zone := room.Zone.Get(); zone != nil {
		log.Debug("running zone check")
		zone.recalcAvgTemp()
	}
}
