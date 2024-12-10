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
	Name             string
	OneDegreeAdjTime time.Duration // how long does it take the zone to move by 1 degF
	Targets          ZoneTargets
	AverageTemp      DegF
	ID               ZoneID
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
			log.Debug("stopping blower", "zone", z, "blower", c.Blowers[k].ID)
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
// if the zone is running, it extends the time the zone is running if the new duration is longer than the current duration
// e.g. 55 minutes left, called for 60 minutes, extends to 60 minutes.
// e.g. 120 minutes left, called for 60 minutes, does nothing
func (z ZoneID) Start(d time.Duration, msg string) error {
	enabled := make([]DeviceID, 0)

	log.Debug("starting/extending zone", "zone", z)

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

// this will get smarter once we have more data to do some ML on
func (z ZoneID) estimateOneDegAdjTime() (time.Duration, error) {
	return (15 * time.Minute), nil
}

func (z *Zone) recalcAvgTemp() {
	var avgCnt uint8
	var avgTot DegF
	var avg DegF
	maxAge := time.Now().Add(0 - tempMaxAge)
	for k := range c.Rooms {
		// in the zone, not zero, and more recent than tempMaxAge
		if c.Rooms[k].GetZoneIDInMode() == z.ID && c.Rooms[k].Temperature != 0 && c.Rooms[k].LastUpdate.After(maxAge) {
			avgCnt++
			avgTot += c.Rooms[k].Temperature
		}

		if avgCnt > 1 {
			avg = avgTot / DegF(avgCnt)
		} else {
			avg = avgTot // save a div if 1, avoid NaN if 0
		}
	}
	if avg != 0 {
		z.AverageTemp = avg
	}
}

func (z *Zone) UpdateTemp() {
	zoneOccupied := false
	for k := range c.Rooms {
		if c.Rooms[k].GetZoneIDInMode() == z.ID && c.Rooms[k].Occupied {
			zoneOccupied = true
			break
		}
	}

	z.recalcAvgTemp()
	log.Debug("zone temp", "zone", z.ID, "avg", z.AverageTemp)

	switch c.SystemMode {
	case SystemModeHeat:
		if (zoneOccupied && z.AverageTemp < z.Targets.HeatingOccupiedTemp-zoneHysterisisRange) || (!zoneOccupied && z.AverageTemp < z.Targets.HeatingUnoccupiedTemp-zoneHysterisisRange) {
			if c.ControlMode == ControlTemp {
				log.Info("starting/extending zone", "zone", z.ID, "avg temp", z.AverageTemp, "targets", z.Targets)
				_ = z.ID.Start(defaultRunDuration, "temp")
			}
			return
		}
		if (zoneOccupied && z.AverageTemp > z.Targets.HeatingOccupiedTemp) || (!zoneOccupied && z.AverageTemp > z.Targets.HeatingUnoccupiedTemp) {
			if c.ControlMode == ControlTemp && z.ID.IsRunning() {
				log.Info("stopping zone", "zone", z.ID, "avg temp", z.AverageTemp, "targets", z.Targets)
				z.ID.Stop("temp")
			}
			return
		}
	case SystemModeCool:
		if (zoneOccupied && z.AverageTemp > z.Targets.CoolingOccupiedTemp+zoneHysterisisRange) || (!zoneOccupied && z.AverageTemp > z.Targets.CoolingUnoccupiedTemp+zoneHysterisisRange) {
			if c.ControlMode == ControlTemp {
				log.Info("starting/extending zone", "zone", z.ID, "avg temp", z.AverageTemp, "targets", z.Targets)
				_ = z.ID.Start(defaultRunDuration, "temp")
			}
			return
		}
		if (zoneOccupied && z.AverageTemp < z.Targets.CoolingOccupiedTemp) || (!zoneOccupied && z.AverageTemp < z.Targets.CoolingUnoccupiedTemp) {
			if c.ControlMode == ControlTemp && z.ID.IsRunning() {
				log.Info("stopping zone", "zone", z.ID, "avg temp", z.AverageTemp, "targets", z.Targets)
				z.ID.Stop("temp")
			}
			return
		}
	}
}

// A zone is runing of all devices in the zone are running, no matter how they were started
// XXX TODO this is not complete for radiant heating zones
func (z ZoneID) IsRunning() bool {
	var totalDevices uint8
	for k := range c.Blowers {
		totalDevices++
		if c.Blowers[k].Zone == z {
			if !c.Blowers[k].Running {
				return false
			}

			pumpid := c.Blowers[k].getPump(c.SystemMode)
			if pumpid != 0 {
				totalDevices++
				if !pumpid.Get().Running {
					return false
				}
			}

			chillerid := pumpid.getChiller()
			if chillerid != 0 {
				totalDevices++
				if !chillerid.Get().Running {
					return false
				}
			}
		}
	}

	if totalDevices == 0 { // radiant heating zones in cooling mode
		return false
	}

	return true
}
