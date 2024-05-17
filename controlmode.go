package hvac

import (
	"fmt"
	"github.com/FUMCGarland/hvac/log"
	"os"
	"path"
)

// ControlModeT is the real ControlMode type, used everywhere except the REST request; just simplify the REST...
type ControlModeT uint8

const (
	ControlManual   ControlModeT = iota // Manual Mode
	ControlSchedule                     // Schedule Individual devices
	ControlTemp                         // Thermostatic mode
	ControlOff                          // Everything off
)

var systemControlModeStrings = []string{"manual", "schedule", "temp", "off"}

// ControlMode is the JSON wireformat for the REST request, dumbly named, probably unnecessary
type ControlMode struct {
	ControlMode ControlModeT
}

// ToString returns a friendly name for the ControlModeT
func (t ControlModeT) ToString() string {
	return systemControlModeStrings[t]
}

// Take a friendly name and return the ControlModeT
func ControlModeFromString(s string) ControlModeT {
	switch s {
	case "manual":
		return ControlManual
	case "schedule":
		return ControlSchedule
	case "temp":
		return ControlTemp
	case "off":
		return ControlOff
	}
	log.Error("unknown system control mode string", "mode", s)
	return ControlManual
}

// SetControlMode is called from LoadConfig at startup and when the mode is changed manually
func (c *Config) SetControlMode(cm ControlModeT) error {
	if cm > ControlOff {
		err := fmt.Errorf("unknown system control mode")
		log.Error(err.Error())
		return err
	}

	switch cm {
	case ControlManual:
		log.Info("stopping schedulers")
		StopAll()
		_ = scheduler.StopJobs()
		// occScheduler.StopJobs()
		occScheduler.Start() // during development, won't hurt since the room's occupied flag is ignored if not in temp mode
		log.Info("control mode manual")
	case ControlSchedule:
		// occScheduler.StopJobs()
		occScheduler.Start() // during development, won't hurt since the room's occupied flag is ignored if not in temp mode
		log.Info("starting scheduler")
		scheduler.Start()
	case ControlTemp:
		log.Info("stopping scheduler")
		_ = scheduler.StopJobs()
		log.Info("starting temp mode")
		occScheduler.Start()
	case ControlOff:
		log.Info("stopping schedulers")
		_ = scheduler.StopJobs()
		// occScheduler.StopJobs()
		occScheduler.Start() // during development, won't hurt since the room's occupied flag is ignored if not in temp mode
		StopAll()
		log.Info("control mode off")
	}

	c.ControlMode = cm
	if err := c.writeControlMode(); err != nil {
		log.Error("unable to write cm to disk", "error", err.Error())
		return err
	}
	return nil
}

// write the current mode to disk so that if we restart we come back correctly
func (c *Config) writeControlMode() error {
	path := path.Join(c.StateStore, storenameControlMode)

	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if _, err := fd.WriteString(c.ControlMode.ToString()); err != nil {
		log.Error(err.Error())
		return err
	}

	if err := fd.Close(); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// readControlMode called at startup to get the last mode from disk
func (c *Config) readControlMode() (ControlModeT, error) {
	path := path.Join(c.StateStore, storenameControlMode)

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error(err.Error())
		return ControlManual, err
	}
	return ControlModeFromString(string(data)), nil
}
