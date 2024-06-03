package hvac

import (
	"time"
)

type DegF float32 // uint8

// min/max durations for user-scheduled run-times, does not affect temp based
const (
	defaultRunDuration time.Duration = (time.Hour) // used in temp control mode, we just need a value to send to the relays
	MaxPumpRunTime     time.Duration = (6 * time.Hour)
	MinPumpRunTime     time.Duration = (30 * time.Minute)
	MaxBlowerRunTime   time.Duration = (6 * time.Hour)
	MinBlowerRunTime   time.Duration = (30 * time.Minute)
	MaxChillerRunTime  time.Duration = (6 * time.Hour)
	MinChillerRunTime  time.Duration = (30 * time.Minute)
)

const (
	pumpMinTimeBetweenRuns    time.Duration = (5 * time.Minute) // pumps must stop for minimum of 5 minutes between runs
	blowerMinTimeBetweenRuns  time.Duration = (5 * time.Minute) // blower must stop for minimum of 5 minutes between runs
	chillerMinTimeBetweenRuns time.Duration = (5 * time.Minute) // chiller must stop for minimum of 5 minutes between runs
)

const zoneHysterisisRange DegF = 3 // DegF above and below the zone's configured ranges e.g. 68 = 65-71

// lockout temps
const (
	chillerLockoutTemp  DegF = 60 // if a room is too cold, stop all chillers (let blowers/pumps run) to prevent freezeout
	chillerRecoveryTemp DegF = 68 // all rooms must be above this temp to unlock the chiller
	boilerLockoutTemp   DegF = 78 // if a room is warmer than this, stop all heater pumps (thereby stopping boilers)
	boilerRecoveryTemp  DegF = 68 // all rooms must fall below this to unlock the heater pumps
	minZoneTemp         DegF = 60 // the coolest we will accept as a user-defined value
	maxZoneTemp         DegF = 80 // the warmest we will accept as a user-defined value
)

const tempMaxAge time.Duration = (4 * time.Hour)

// MQTT topic strings
const (
	BlowersTopic         string = "blowers"
	ChillersTopic        string = "chillers"
	PumpsTopic           string = "pumps"
	RoomsTopic           string = "rooms"
	CurrentStateEndpoint string = "currentstate"
	TargetStateEndpoint  string = "targetstate"
	TempEndpoint         string = "temp"     // Shellies use their own topic, this is for other (potential) sensors
	HumidityEndpoint     string = "humidity" // Shellies use their own topic, this is for other (potential) sensors
)

// tags for scheduler, not fully utilized yet, will be helpful when we unify the two schedulers
const (
	scheduleTagOccupancy string = "occupancy"
	scheduleTagDevice    string = "device"
	scheduleTagRecurring string = "recurring"
	scheduleTagOneTime   string = "onetime"
)

// QoS is the MQTT Quality of Service value
// only 0 has been tested
const QoS byte = 0
