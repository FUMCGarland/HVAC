package hvac

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/FUMCGarland/hvac/log"
)

// ZoneID is a unique identifier for a zone
type ZoneID uint8

// A zone is a collection of rooms which are controlled together, either by radiant heat or blowers
type Zone struct {
	ID          ZoneID
	Name        string
	Targets     ZoneTargets
	AverageTemp DegF
}

// Each zone has four target temps, based on systemMode and room occupancy
type ZoneTargets struct {
	HeatingOccupiedTemp   DegF // 68 -- heat to 68 if someone is schedule to be in the zone
	HeatingUnoccupiedTemp DegF // 60 -- let it get down to 60 if no one is schedule to be in the zone
	CoolingOccupiedTemp   DegF // 74 -- cool to 74 if someone is sheduled to be in the zone
	CoolingUnoccupiedTemp DegF // 80 -- let it get up to 80 if no one is scheduled to be in the zone
}

// Get returns a populated zone for a given ZoneID
func (z ZoneID) Get() *Zone {
	for k := range c.Zones {
		if c.Zones[k].ID == z {
			return &c.Zones[k]
		}
	}
	return nil
}

// SetTargets sets a zone's target tempratures, called from the REST interface
func (z *Zone) SetTargets(c *Config, zt *ZoneTargets) error {
	oor := fmt.Errorf("zone temperature out of sane range")

	if zt.HeatingUnoccupiedTemp < minZoneTemp || zt.HeatingUnoccupiedTemp > maxZoneTemp {
		return oor
	}
	if zt.HeatingOccupiedTemp < minZoneTemp || zt.HeatingOccupiedTemp > maxZoneTemp {
		return oor
	}
	if zt.CoolingUnoccupiedTemp < minZoneTemp || zt.CoolingUnoccupiedTemp > maxZoneTemp {
		return oor
	}
	if zt.CoolingOccupiedTemp < minZoneTemp || zt.CoolingOccupiedTemp > maxZoneTemp {
		return oor
	}

	z.Targets.HeatingUnoccupiedTemp = zt.HeatingUnoccupiedTemp
	z.Targets.HeatingOccupiedTemp = zt.HeatingOccupiedTemp
	z.Targets.CoolingUnoccupiedTemp = zt.CoolingUnoccupiedTemp
	z.Targets.CoolingOccupiedTemp = zt.CoolingOccupiedTemp

	if err := z.writeToStore(); err != nil {
		return err
	}
	return nil
}

// TODO: write to a single file instead of a file-per-zone

func (z *Zone) writeToStore() error {
	path := path.Join(c.StateStore, fmt.Sprintf("zone-%d.json", z.ID))

	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	j, err := json.Marshal(*z)
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

func (z *Zone) readFromStore() error {
	path := path.Join(c.StateStore, fmt.Sprintf("zone-%d.json", z.ID))

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	var in Zone
	if err = json.Unmarshal(data, &in); err != nil {
		log.Error(err.Error())
		return err
	}

	z.Targets.HeatingUnoccupiedTemp = in.Targets.HeatingUnoccupiedTemp
	z.Targets.HeatingOccupiedTemp = in.Targets.HeatingOccupiedTemp
	z.Targets.CoolingUnoccupiedTemp = in.Targets.CoolingUnoccupiedTemp
	z.Targets.CoolingOccupiedTemp = in.Targets.CoolingOccupiedTemp

	return nil
}

// Stop shuts down all devices in a zone
func (z ZoneID) Stop(msg string) {
	log.Debug("stopping zone", "ID", z, "msg", msg)

	for k := range c.Blowers {
		if c.Blowers[k].Zone == z && c.Blowers[k].Running {
			log.Debug("stopping blower on zone", "zone", c.Blowers[k].ID)
			c.Blowers[k].ID.Stop(msg)
		}
	}
	// stopping the blowers will stop the pumps and chillers if necessary.

	if c.SystemMode == SystemModeHeat {
		// shut down the radiant loops for the zone
		for k := range c.Loops {
			if c.Loops[k].RadiantZone == z {
				pump := c.getPumpFromLoop(c.Loops[k].ID)
				if pump.Get().Running {
					pump.Stop(msg)
				}
			}
		}
	}
}

// Start starts up all devices in a zone
func (z ZoneID) Start(d time.Duration, msg string) error {
	enabled := make([]DeviceID, 0)

	for k := range c.Blowers {
		if c.Blowers[k].Zone == z {
			if err := c.Blowers[k].ID.Start(d, msg); err != nil {
				stopAll(enabled)
				return err
			}
			enabled = append(enabled, c.Blowers[k].ID)
			time.Sleep(time.Second) // let blowers start before attempting to start pump

			pumpid := c.Blowers[k].getPump(c.SystemMode)
			if pumpid != 0 {
				// don't skip if the pump is already running, we want to extend the running time if necessary
				if err := pumpid.Start(d, msg); err != nil {
					stopAll(enabled)
					return err
				}
				enabled = append(enabled, pumpid)
			}
			time.Sleep(time.Second) // let pumps start before attempting to start chiller

			chillerid := pumpid.getChiller()
			if chillerid != 0 {
				// don't skip if the chiller is already running, we want to extend the running time if necessary
				if err := chillerid.Start(d, msg); err != nil {
					stopAll(enabled)
					return err
				}
				enabled = append(enabled, chillerid)
			}
		}
	}
	return nil
}

func stopAll(enabled []DeviceID) {
	log.Debug("calling stopAll", "enabled", enabled)
	for k := range enabled {
		enabled[k].Stop("internal")
	}
}
