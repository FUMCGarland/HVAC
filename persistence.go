package hvac

import (
	"os"

	"github.com/FUMCGarland/hvac/log"
)

const (
	storenameSystemMode  = "SystemMode"  // filename for the system mode
	storenameControlMode = "ControlMode" // file name for the control mode
)

func (c *Config) loadFromStore() error {
	_, err := os.ReadDir(c.StateStore)
	if err != nil {
		log.Fatal("state storage directory does not exist", "dir", c.StateStore)
	}

	sm, err := c.readSystemMode()
	if err != nil {
		log.Error("unable to load last system mode")
		sm = SystemModeHeat
	}
	if err := c.SetSystemMode(sm); err != nil {
		log.Fatal(err.Error())
	}
	cm, err := c.readControlMode()
	if err != nil {
		log.Error("unable to load last system control mode")
		cm = ControlOff
	}
	// starts the scheduler or temp handler if set
	if err := c.SetControlMode(cm); err != nil {
		log.Fatal(err.Error())
	}

	for k := range c.Pumps {
		if err := c.Pumps[k].readFromStore(); err != nil {
			log.Error(err.Error())
			return err
		}
	}

	for k := range c.Blowers {
		if err := c.Blowers[k].readFromStore(); err != nil {
			log.Error(err.Error())
			return err
		}
	}

	for k := range c.Zones {
		if err := c.Zones[k].readFromStore(); err != nil {
			log.Error(err.Error())
			return err
		}
	}

	for k := range c.Chillers {
		if err := c.Chillers[k].readFromStore(); err != nil {
			log.Error(err.Error())
			return err
		}
	}

	for k := range c.Rooms {
		if err := c.Rooms[k].readFromStore(); err != nil {
			log.Error(err.Error())
			return err
		}
	}

	s, err := readScheduleFromStore()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	schedule = *s // pointless?

	o, err := readOccupancyFromStore()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	occupancy = *o // pointless?
	return nil
}

func (c *Config) WriteToStore() error {
	if _, err := os.ReadDir(c.StateStore); err != nil {
		log.Fatal("state storage directory does not exist", "dir", c.StateStore)
	}

	if err := c.writeSystemMode(); err != nil {
		log.Error(err.Error())
		return err
	}

	if err := c.writeControlMode(); err != nil {
		log.Error(err.Error())
		return err
	}

	for k := range c.Pumps {
		if err := c.Pumps[k].writeToStore(); err != nil {
			log.Error(err.Error())
			return err
		}
	}

	for k := range c.Blowers {
		if err := c.Blowers[k].writeToStore(); err != nil {
			log.Error(err.Error())
			return err
		}
	}

	for k := range c.Zones {
		if err := c.Zones[k].writeToStore(); err != nil {
			log.Error(err.Error())
			return err
		}
	}

	for k := range c.Chillers {
		if err := c.Chillers[k].writeToStore(); err != nil {
			log.Error(err.Error())
			return err
		}
	}

	for k := range c.Rooms {
		if err := c.Rooms[k].writeToStore(); err != nil {
			log.Error(err.Error())
			return err
		}
	}

	if err := (&schedule).writeToStore(); err != nil {
		log.Error(err.Error())
		return err
	}

	if err := (&occupancy).writeToStore(); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
