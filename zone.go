package hvac

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

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

	if zt.HeatingUnoccupiedTemp < 55 || zt.HeatingUnoccupiedTemp > 80 {
		return oor
	}
	if zt.HeatingOccupiedTemp < 55 || zt.HeatingOccupiedTemp > 80 {
		return oor
	}
	if zt.CoolingUnoccupiedTemp < 55 || zt.CoolingUnoccupiedTemp > 80 {
		return oor
	}
	if zt.CoolingOccupiedTemp < 55 || zt.CoolingOccupiedTemp > 80 {
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
