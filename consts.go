package hvac

import (
	"time"
)

type DegF uint8

const defaultRunDuration time.Duration = (time.Hour)

const MaxPumpRunTime time.Duration = (6 * time.Hour)
const MinPumpRunTime time.Duration = (30 * time.Minute)
const MaxBlowerRunTime time.Duration = (6 * time.Hour)
const MinBlowerRunTime time.Duration = (30 * time.Minute)
const MaxChillerRunTime time.Duration = (6 * time.Hour)
const MinChillerRunTime time.Duration = (30 * time.Minute)

const pumpMinTimeBetweenRuns time.Duration = (5 * time.Minute)
const blowerMinTimeBetweenRuns time.Duration = (5 * time.Minute)
const chillerMinTimeBetweenRuns time.Duration = (5 * time.Minute)

const minZoneTemp DegF = 60 // the coolest we will accept as a user-defined value
const maxZoneTemp DegF = 80 // the warmest we will accept as a user-defined value

const zoneHysterisisRange DegF = 3 // DegF above and below the zone's configured ranges

const chillerLockoutTemp DegF = 68  // DegF, if too cold, don't run chiller
const chillerRecoveryTemp DegF = 75 // all rooms must be above this temp to unlock the chiller
const boilerLockoutTemp DegF = 78   // if a room is warmer than this, stop all heater pumps
const boilerRecoveryTemp DegF = 68  // all rooms must fall below this to unlock the heater pumps
