package hvac

import (
	"time"
)

// Command is what the MQTT server sends to the relay controllers
type Command struct {
	TargetState bool
	RunTime     time.Duration // the time to run the device/zone
	Source      string        // a string describing the reasonx:, manual, schedule, auto
}

// Response is what the Relay Controller sends to the MQTT server
type Response struct {
	CurrentState  bool
	RanTime       time.Duration // time actually ran (when task completed)
	TimeRemaining time.Duration // time left on the clock (for periodic check-ins)
}

// DeviceID is a generic interface for all devices (pumps, blowers, etc)
type DeviceID interface {
	canEnable() error
	Start(time.Duration, string) error
	Stop(string)
}

// The MQTTRequest is a wrapper type which contains both the command and the ID of the device being controlled
type MQTTRequest struct {
	DeviceID DeviceID
	Command  Command
}
