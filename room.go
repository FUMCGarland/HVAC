package hvac

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/FUMCGarland/hvac/log"
	"github.com/go-co-op/gocron/v2"
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
	CoolZone    ZoneID
	HeatZone    ZoneID
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

		// update the zone average and run logic on starting/stopping the zone
		zone := r.HeatZone.Get()
		zone.UpdateTemp()
	case SystemModeCool:
		if r.Temperature <= chillerLockoutTemp && !c.ChillerLockout {
			log.Warn("locking out chiller, room temp too low", "room", r.ID, "temp", r.Temperature)
			c.ChillerLockout = true
			for k := range c.Chillers {
				c.Chillers[k].ID.Stop("lockout")
			}
		}

		// update the zone average and run logic on starting/stopping the zone
		zone := r.CoolZone.Get()
		zone.UpdateTemp()
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
	if zone := room.GetZoneInMode(); zone != nil {
		log.Debug("running zone check")
		zone.UpdateTemp()
	}

	if !room.Occupied {
		return
	}

	log.Debug("scheduling room to be unoccupied")
	end := time.Now().Add(time.Hour * 2)
	_, err := occScheduler.NewJob(
		gocron.OneTimeJob(
			gocron.OneTimeJobStartDateTime(end),
		),
		gocron.NewTask(
			func() {
				log.Debug("clearing manual occupancy for room: %s", room.Name)
				room.Occupied = false
				room.GetZoneInMode().UpdateTemp() // recalculates the avg and runs if needed
			},
		),
		// gocron.WithTags(e.Name, scheduleTagOccupancy, scheduleTagOneTime),
		gocron.WithName(fmt.Sprintf("manual occupancy end")),
	)
	if err != nil {
		log.Error(err.Error())
	}
}

func (r Room) GetZoneInMode() *Zone {
	switch c.SystemMode {
	case SystemModeHeat:
		return r.HeatZone.Get()
	case SystemModeCool:
		return r.CoolZone.Get()
	}
	return nil
}

func (r Room) GetZoneIDInMode() ZoneID {
	switch c.SystemMode {
	case SystemModeHeat:
		return r.HeatZone
	case SystemModeCool:
		return r.CoolZone
	}
	return 0
}
