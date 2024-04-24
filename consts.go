package hvac

import (
	"time"
)

const MaxPumpRunTime uint64 = 14400
const MaxBlowerRunTime uint64 = 14400
const PumpMinTimeBetweenRuns time.Duration = (5 * 60 * time.Second)
