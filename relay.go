package hvac

import (
	"time"
)

type Relay struct {
	Pin       uint8     // gpio pin on relay control board
	PumpID    PumpID    // non-zero if devices is a pump
	BlowerID  BlowerID  // non-zero if device is a blower
	ChillerID ChillerID // non-zero if the device is a chiller

	Running   bool          // is the device currently running
	StartTime time.Time     // time of current start
	StopTime  time.Time     // time to stop running
	RunTime   time.Duration // total run-time of the device
}
