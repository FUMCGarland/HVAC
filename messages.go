package hvac

import (
	"time"
)

// Command is what the MQTT server sends to the relay controllers
type Command struct {
	TargetState bool
	RunTime     time.Duration
	Source      string // manual, schedule, auto
}

// Response is what the Relay Controller sends to the MQTT server
type Response struct {
	CurrentState  bool
	RanTime       time.Duration // time actually ran
	TimeRemaining time.Duration
}

// types to make sure we aren't sending blower commands to a pump even though they look the same

// PumpCommand is what the controller sends to the relay-module to start or stop a pump
type PumpCommand Command
type BlowerCommand Command
type ChillerCommand Command
type ZoneCommand Command

// PumpResponse is what the relay-module sends back to the controller when a pump changes state
type PumpResponse Response
type BlowerResponse Response
type ChillerResponse Response
type ZoneResponse Response

// DeviceID is a generic interface for all devices (pumps, blowers, etc)
type DeviceID interface {
	canEnable() error
	// Get() device
	Start(time.Duration, string) error
	Stop(string)
}

// unused for now
type device interface {
	readFromStore() error
	writeToStore() error
}

// The MQTTRequest is a wrapper type which contains both the command and the ID of the device being controlled
type MQTTRequest struct {
	Device  DeviceID
	Command Command
}
