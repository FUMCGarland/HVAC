package hvac

import (
	// "fmt"
	"os"
	// "path"

	"github.com/FUMCGarland/hvac/log"
)

const (
	storenameSystemMode  = "SystemMode"
	storenameControlMode = "ControlMode"
)

// it would make sense to use badgerdb since it is already built in to mochi-mqtt, but it seems like overkill for this task

func (c *Config) loadFromStore() error {
	_, err := os.ReadDir(c.StateStore)
	if err != nil {
		log.Fatal("state storage directory does not exist", "dir", c.StateStore)
	}

	sm, err := c.readSystemMode()
	if err != nil {
		log.Error("unable to load last system mode")
	}
	c.SetSystemMode(sm)
	scm, err := c.readControlMode()
	if err != nil {
		log.Error("unable to load last system control mode")
	}
	// starts the scheduler or temp handler if set
	c.SetControlMode(scm)

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

	s, err := readScheduleFromStore()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	schedule = *s
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

	if err := (&schedule).writeToStore(); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
