package hvac

import (
	"time"
)

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
const minZoneTemp = 55
const maxZoneTemp = 85
