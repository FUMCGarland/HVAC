package hvac

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/FUMCGarland/hvac/log"
)

type ZoneID uint8

type Zone struct {
	ID      ZoneID
	Name    string
	Targets ZoneTargets
}

type ZoneTargets struct {
	HeatingOccupiedTemp   uint8 // 68 -- heat to 68 if someone is schedule to be in the zone
	HeatingUnoccupiedTemp uint8 // 60 -- let it get down to 60 if no one is schedule to be in the zone
	CoolingOccupiedTemp   uint8 // 74 -- cool to 74 if someone is sheduled to be in the zone
	CoolingUnoccupiedTemp uint8 // 80 -- let it get up to 80 if no one is scheduled to be in the zone
}

func (z ZoneID) Get() *Zone {
	for k := range c.Zones {
		if c.Zones[k].ID == z {
			return &c.Zones[k]
		}
	}
	return nil
}

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

func (z ZoneID) Stop(msg string) {
	// stop the blowers the pumps/chiller will cascade if necessary
	for k := range c.Blowers {
		if c.Blowers[k].Zone == z && c.Blowers[k].Running {
			c.Blowers[k].ID.Stop(msg)
		}
	}

	if c.SystemMode == SystemModeHeat {
		// shut down the radiant loops for the zone
		for k := range c.Loops {
			if c.Loops[k].RadiantZone == z {
				pump := c.GetPumpFromLoop(c.Loops[k].ID)
				if pump.Get().Running {
					pump.Stop(msg)
				}
			}
		}
	}
}

func (z ZoneID) Start(d time.Duration, msg string) error {
	enabled := make([]DeviceID, 0)

	for k := range c.Blowers {
		if c.Blowers[k].Zone == z {
			if err := c.Blowers[k].ID.Start(d, msg); err != nil {
				stopAll(enabled)
				return err
			}
			enabled = append(enabled, c.Blowers[k].ID)

			pumpid := c.Blowers[k].getPump(c.SystemMode)
			if pumpid != 0 {
				time.Sleep(1 * time.Second) // let blower start before attempting to start pump
				if err := pumpid.Start(d, msg); err != nil {
					stopAll(enabled)
					return err
				}
				enabled = append(enabled, pumpid)
			}

			chillerid := pumpid.getChiller()
			if chillerid != 0 {
				time.Sleep(1 * time.Second) // let pump start before attempting to start chiller
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
	for _, dev := range enabled {
		switch dev := dev.(type) {
		case BlowerID:
			dev.Stop("internal")
		case ChillerID:
			dev.Stop("internal")
		case PumpID:
			dev.Stop("internal")
		}
	}
}
