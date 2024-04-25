package hvac

import (
	"time"
)

// TODO: implement this
// const MaxPumpRunTime uint64 = 14400
// const MaxBlowerRunTime uint64 = 14400

const pumpMinTimeBetweenRuns time.Duration = (5 * 60 * time.Second)
const blowerMinTimeBetweenRuns time.Duration = (5 * 60 * time.Second)
const minZoneTemp = 55
const maxZoneTemp = 85
