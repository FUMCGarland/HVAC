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
type PumpCommand Command
type BlowerCommand Command
type ChillerCommand Command
type ZoneCommand Command

type PumpResponse Response
type BlowerResponse Response
type ChillerResponse Response
type ZoneResponse Response

type DeviceID interface {
	canEnable() error
	// Get() device
	Start(time.Duration, string) error
	Stop(string)
}

type device interface {
	readFromStore() error
	writeToStore() error
}

type MQTTRequest struct {
	Device  DeviceID
	Command Command
}
