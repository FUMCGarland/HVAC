package hvac

// Command is what the MQTT server sends to the relay controllers
type Command struct {
	TargetState bool
	RunTime     uint64 // seconds to run
	Source      string // manual, schedule, auto
}

// Response is what the Relay Controller sends to the MQTT server
type Response struct {
	CurrentState bool
	RanTime      uint64 // seconds actually ran
}

type PumpCommand Command
type BlowerCommand Command
type PumpResponse Response
type BlowerResponse Response

type DeviceID interface {
	CanEnable() error
	// Get() // pump & blower should have an interface...
	Start(uint64, string) error
	Stop(string) error
}

type MQTTRequest struct {
	Device  DeviceID
	Command Command
}
