package hvac

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/FUMCGarland/hvac/log"
	"github.com/go-co-op/gocron/v2"
	// "github.com/google/uuid"
)

type OccupancySchedule struct {
	Recurring []OccupancyRecurringEntry
	OneTime   []OccupancyOneTimeEntry
}

var occupancy OccupancySchedule
var occScheduler gocron.Scheduler

type OccupancyOneTimeID uint8
type OccupancyRecurringID uint8

// ScheduleEntry is the definition of a job to be run at specified times
type OccupancyRecurringEntry struct {
	ID        OccupancyRecurringID
	Name      string
	StartTime string // "6:30"
	EndTime   string // "15:30"
	Weekdays  []time.Weekday
	Rooms     []RoomID
}

// ScheduleEntry is the definition of a job to be run at specified times
type OccupancyOneTimeEntry struct {
	ID    OccupancyOneTimeID
	Name  string
	Start time.Time
	End   time.Time
	Rooms []RoomID
}

func init() {
	var err error
	occScheduler, err = gocron.NewScheduler()
	if err != nil {
		log.Error(err.Error())
	}
}

// GetOccupancySchedule returns the live schedule
func (c *Config) GetOccupancySchedule() (*OccupancySchedule, error) {
	if len(occupancy.Recurring) == 0 && len(occupancy.OneTime) == 0 {
		s, err := readOccupancyFromStore()
		if err != nil {
			log.Error(err.Error())
			return &occupancy, err
		}
		occupancy = *s
	}

	return &occupancy, nil
}

func (s *OccupancySchedule) writeToStore() error {
	path := path.Join(c.StateStore, "occupancy.json")

	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	j, err := json.Marshal(s)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if _, err := fd.Write(j); err != nil {
		log.Error(err.Error())
		return err
	}

	if err := fd.Close(); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func readOccupancyFromStore() (*OccupancySchedule, error) {
	sl := OccupancySchedule{}
	// sl.Recurring = []OccupancyRecurringEntry{}
	// sl.OneTime = []OccupancyOneTimeEntry{}

	path := path.Join(c.StateStore, "occupancy.json")

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error(err.Error())
		return &sl, nil // file doesn't exist, start with an empty one
	}

	if err = json.Unmarshal(data, &sl); err != nil {
		log.Error(err.Error(), "path", path)
		return &sl, err
	}

	if len(sl.Recurring) == 0 {
		sl.Recurring = make([]OccupancyRecurringEntry, 0)
	}

	if len(sl.OneTime) == 0 {
		sl.OneTime = make([]OccupancyOneTimeEntry, 0)
	}

	for k := range sl.Recurring {
		log.Info("loading recurring occupancy entry", "entry", sl.Recurring[k])
		if err := buildRecurringJob(&sl.Recurring[k]); err != nil {
			log.Error(err.Error())
			return &sl, err
		}
	}

	for k := range sl.OneTime {
		log.Info("loading onetime occupancy entry", "entry", sl.OneTime[k])
		if err := buildOneTimeJob(&sl.OneTime[k]); err != nil {
			log.Error(err.Error())
			return &sl, err
		}
	}

	return &sl, nil
}
