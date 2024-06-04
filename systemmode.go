package hvac

import (
	"fmt"
	"os"
	"path"

	"github.com/FUMCGarland/hvac/log"
)

// SystemModeT is a convenience type to ensure clean code
type SystemModeT uint8

const (
	SystemModeHeat    SystemModeT = iota // the system is heating
	SystemModeCool                       // the system is cooling
	SystemModeUnknown                    // unused
)

// SystemMode is a wrapper type for clean JSON files passed from the REST interface
type SystemMode struct {
	Mode SystemModeT
}

// systemModeStrings is used to translate between the SystemModeT and a friendly name
var systemModeStrings = []string{"heat", "cool", "unknown"}

// ToString returns a friendly string for a SystemModeT
func (t SystemModeT) ToString() string {
	return systemModeStrings[t]
}

// systemModeFromString returns a SystemModeT that matches a string
func systemModeFromString(s string) SystemModeT {
	if s == "heat" {
		return SystemModeHeat
	}
	if s == "cool" {
		return SystemModeCool
	}
	return SystemModeUnknown
}

// SetSystemMode sets the system into the requested mode
func (c *Config) SetSystemMode(sm SystemModeT) error {
	// TODO enforce controlMode == off
	if sm != SystemModeHeat && sm != SystemModeCool {
		err := fmt.Errorf("unknown system mode")
		log.Error(err.Error())
		return err
	}
	log.Debug("setting system mode", "mode", sm, "string", sm.ToString())
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
	return systemModeFromString(string(data)), nil
}
