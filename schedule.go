package hvac

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/FUMCGarland/hvac/log"
	"github.com/go-co-op/gocron/v2"
)

// ScheduleList is just a wrapper so we can hang methods on it
// The scheduler is used in "schedule" control mode (not temp) and starts
// and stops devices based on the schedule, not room temp/occupancy
// this is akin to the old mode of operation using the scheduler from the 1980s...
type ScheduleList struct {
	List []ScheduleEntry
}

var schedule ScheduleList
var scheduler gocron.Scheduler

// ScheduleEntry is the definition of a job to be run at specified times
type ScheduleEntry struct {
	ID        uint8
	Name      string
	Mode      SystemModeT
	StartTime string // "10:30" "18:00;22:30;24:00""
	Weekdays  []time.Weekday
	RunTime   time.Duration
	Zones     []ZoneID
}

// init() considered harmful, just a singleton to set up the global schedular
func init() {
	var err error
	scheduler, err = gocron.NewScheduler()
	if err != nil {
		log.Error(err.Error())
	}
}

// GetSchedule returns the live schedule
func (c *Config) GetSchedule() (*ScheduleList, error) {
	if len(schedule.List) == 0 {
		s, err := readScheduleFromStore()
		if err != nil {
			log.Error(err.Error())
			return &schedule, err
		}
		schedule = *s
	}
	return &schedule, nil
}

func (s *ScheduleList) writeToStore() error {
	path := path.Join(c.StateStore, "schedule.json")

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

func readScheduleFromStore() (*ScheduleList, error) {
	sl := ScheduleList{}
	sl.List = []ScheduleEntry{}

	path := path.Join(c.StateStore, "schedule.json")

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error(err.Error())
		return &sl, nil // file doesn't exist, start with an empty one
	}

	if err = json.Unmarshal(data, &sl); err != nil {
		log.Error(err.Error(), "path", path)
		return &sl, err
	}

	if len(sl.List) == 0 {
		sl.List = make([]ScheduleEntry, 0)
	}

	for k := range sl.List {
		log.Debug("loading schedule entry", "entry", sl.List[k])
		if err := buildJob(&sl.List[k]); err != nil {
			log.Error(err.Error())
			return &sl, err
		}
	}

	return &sl, nil
}

// GetEntry returns an entry in the ScheduleList by ID
func (s *ScheduleList) GetEntry(id uint8) *ScheduleEntry {
	for k := range s.List {
		if s.List[k].ID == id {
			return &s.List[k]
		}
	}
	return nil
}

// AddEntry adds a new entry to the list of jobs to run
func (s *ScheduleList) AddEntry(e *ScheduleEntry) error {
	log.Debug("adding entry", "e", e)

	if len(e.Weekdays) == 0 {
		err := fmt.Errorf("cannot schedule an entry not on any days")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if len(e.Zones) == 0 {
		err := fmt.Errorf("cannot schedule an entry without any zones")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if existing := s.GetEntry(e.ID); existing != nil {
		err := fmt.Errorf("cannot reuse existing ID")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if e.Mode != SystemModeHeat && e.Mode != SystemModeCool {
		err := fmt.Errorf("invalid system mode")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if e.Name == "" {
		e.Name = fmt.Sprintf("Unnamed %d", e.ID)
	}

	if err := buildJob(e); err != nil {
		log.Error(err.Error(), "entry", e)
		return err
	}

	schedule.List = append(schedule.List, *e)
	if err := schedule.writeToStore(); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// buildJob processes a ScheduleEntry and loads it into gocron
func buildJob(e *ScheduleEntry) error {
	attimes := make([]gocron.AtTime, 0)

	times := strings.Split(e.StartTime, ";")
	for _, v := range times {
		log.Debug("time", "time", v)
		units := strings.Split(v, ":")
		hour, err := strconv.ParseInt(units[0], 10, 8)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		hour = hour % 24
		minute, err := strconv.ParseInt(units[1], 10, 8)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		minute = minute % 60
		attimes = append(attimes, gocron.NewAtTime(uint(hour), uint(minute), 0))
	}

	_, err := scheduler.NewJob(
		gocron.WeeklyJob(
			1,
			func() []time.Weekday { return e.Weekdays },
			func() []gocron.AtTime { return attimes },
		),
		gocron.NewTask(
			func() {
				log.Debug("starting scheduled entry", "e", e)
				for _, zone := range e.Zones {
					log.Debug("starting zone", "zone", zone, "duration", e.RunTime.Minutes())
					if err := zone.Start(e.RunTime, "scheduled"); err != nil {
						log.Error(err.Error())
					}
				}
			},
		),
		gocron.WithTags(fmt.Sprintf("%d", e.ID), scheduleTagDevice, e.Name),
		gocron.WithName(e.Name),
	)

	return err
}

// RemoveEntry removes a ScheduleEntry from the list by ID
func (s *ScheduleList) RemoveEntry(id uint8) {
	index := -1
	for k := range schedule.List {
		if s.List[k].ID == id {
			index = k
			break
		}
	}
	if index == -1 {
		log.Warn("unknown schedule entry ID", "id", id)
		return
	}
	log.Debug("removing job from schedule", "id", id)
	scheduler.RemoveByTags(fmt.Sprintf("%d", id))
	s.List = append(s.List[:index], s.List[index+1:]...)
	log.Debug("new schedule", "s", s.List)
	_ = s.writeToStore()
}

// EditEntry updates an entry in the ScheduleList, keyed based on e.ID
func (s *ScheduleList) EditEntry(e *ScheduleEntry) error {
	index := -1
	for k := range schedule.List {
		if s.List[k].ID == e.ID {
			index = k
			break
		}
	}
	if index == -1 {
		err := fmt.Errorf("cannot update non-existent entry")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if len(e.Weekdays) == 0 {
		err := fmt.Errorf("cannot schedule an entry not on any days")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if len(e.Zones) == 0 {
		err := fmt.Errorf("cannot schedule an entry without any zones")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if e.Mode != SystemModeHeat && e.Mode != SystemModeCool {
		err := fmt.Errorf("invalid system mode")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if e.Name == "" {
		e.Name = fmt.Sprintf("Unnamed %d", e.ID)
	}

	s.List[index] = *e
	if err := s.writeToStore(); err != nil {
		log.Error(err.Error(), "entry", e)
		return err
	}

	log.Debug("removing job from schedule", "id", e.ID)
	scheduler.RemoveByTags(fmt.Sprintf("%d", e.ID))

	if err := buildJob(e); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
