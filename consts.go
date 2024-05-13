package hvac

import (
	"time"
)

type DegF uint8

const defaultRunDuration time.Duration = (time.Hour)

// min/max durations for user-scheduled run-times, does not affect temp based
const MaxPumpRunTime time.Duration = (6 * time.Hour)
const MinPumpRunTime time.Duration = (30 * time.Minute)
const MaxBlowerRunTime time.Duration = (6 * time.Hour)
const MinBlowerRunTime time.Duration = (30 * time.Minute)
const MaxChillerRunTime time.Duration = (6 * time.Hour)
const MinChillerRunTime time.Duration = (30 * time.Minute)

const pumpMinTimeBetweenRuns time.Duration = (5 * time.Minute)    // pumps must stop for minimum of 5 minutes between runs
const blowerMinTimeBetweenRuns time.Duration = (5 * time.Minute)  // blower must stop for minimum of 5 minutes between runs
const chillerMinTimeBetweenRuns time.Duration = (5 * time.Minute) // chiller must stop for minimum of 5 minutes between runs

const minZoneTemp DegF = 60 // the coolest we will accept as a user-defined value
const maxZoneTemp DegF = 80 // the warmest we will accept as a user-defined value

const zoneHysterisisRange DegF = 3 // DegF above and below the zone's configured ranges e.g. 68 = 65-71

const chillerLockoutTemp DegF = 68  // if a room is too cold, stop all chillers (let blowers/pumps run) to prevent freezeout
const chillerRecoveryTemp DegF = 75 // all rooms must be above this temp to unlock the chiller
const boilerLockoutTemp DegF = 78   // if a room is warmer than this, stop all heater pumps (thereby stopping boilers)
const boilerRecoveryTemp DegF = 68  // all rooms must fall below this to unlock the heater pumps
