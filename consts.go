package hvac

import (
	"time"
)

type DegF uint8

const defaultRunDuration time.Duration = (1 * time.Hour)

// TODO: implement this
const MaxPumpRunTime time.Duration = (6 * time.Hour)
const MinPumpRunTime time.Duration = (30 * time.Minute)
const MaxBlowerRunTime time.Duration = (6 * time.Hour)
const MinBlowerRunTime time.Duration = (30 * time.Minute)
const MaxChillerRunTime time.Duration = (6 * time.Hour)
const MinChillerRunTime time.Duration = (30 * time.Minute)

const pumpMinTimeBetweenRuns time.Duration = (5 * time.Minute)
const blowerMinTimeBetweenRuns time.Duration = (5 * time.Minute)
const chillerMinTimeBetweenRuns time.Duration = (5 * time.Minute)
const minZoneTemp DegF = 55
const maxZoneTemp DegF = 85

// this is the "I don't know what I'm doing" part of the project
const zoneHysterisisRange DegF = 3 // DegF above and below the zone's configured ranges
const chillerLockoutTemp DegF = 68 // DegF, if too cold, don't run chiller
const chillerRecoveryTemp DegF = 75
const boilerLockoutTemp DegF = 78
const boilerREcoveryTemp DegF = 68
