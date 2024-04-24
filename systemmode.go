package hvac

import (
	"fmt"
	"os"
	"path"

	"github.com/FUMCGarland/hvac/log"
)

type SystemModeT uint8

const (
	SystemModeHeat SystemModeT = iota
	SystemModeCool
	SystemModeUnknown
)

type SystemMode struct {
	Mode SystemModeT
}

var systemModeStrings = []string{"heat", "cool", "unknown"}

func (t SystemModeT) ToString() string {
	return systemModeStrings[t]
}

func SystemModeFromString(s string) SystemModeT {
	if s == "heat" {
		return SystemModeHeat
	}
	if s == "cool" {
		return SystemModeCool
	}
	return SystemModeUnknown
}

func (c *Config) SetSystemMode(sm SystemModeT) error {
	if sm != SystemModeHeat && sm != SystemModeCool {
		err := fmt.Errorf("unknown system mode")
		log.Error(err.Error())
		return err
	}
	log.Info("setting system mode", "mode", sm, "string", sm.ToString())
	c.SystemMode = sm
	if err := c.writeSystemMode(); err != nil {
		log.Error("unable to write system mode to disk", "error", err.Error())
		return err
	}
	return nil
}

func (c *Config) writeSystemMode() error {
	path := path.Join(c.StateStore, storenameSystemMode)

	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if _, err := fd.WriteString(c.SystemMode.ToString()); err != nil {
		log.Error(err.Error())
		return err
	}

	if err := fd.Close(); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (c *Config) readSystemMode() (SystemModeT, error) {
	path := path.Join(c.StateStore, storenameSystemMode)

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error(err.Error())
		return SystemModeUnknown, err
	}
	return SystemModeFromString(string(data)), nil
}
