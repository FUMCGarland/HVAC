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

type ScheduleList struct {
	List []ScheduleEntry
}

var schedule ScheduleList
var sz gocron.Scheduler

type ScheduleEntry struct {
	ID        uint8
	Name      string
	Mode      SystemModeT
	StartTime string // "10:30" "18:00;22:30;24:00""
	Weekdays  []time.Weekday
	RunTime   uint64
	Pumps     []PumpID
	Blowers   []BlowerID
}

func init() {
	var err error
	sz, err = gocron.NewScheduler()
	if err != nil {
		log.Error(err.Error())
	}
}

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

	j, err := json.Marshal(&s)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("writing", "e", s, "json", j)
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
	sl.List = make([]ScheduleEntry, 0)

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

	for _, v := range sl.List { // use copy, not live data, see if this fixes the Weekdays corruption
		v.Weekdays = uniq(v.Weekdays)
		log.Info("loading schedule entry", "entry", v)
		if err := buildJob(&v); err != nil {
			log.Error(err.Error())
			return &sl, err
		}
	}

	return &sl, nil
}

func uniq(s []time.Weekday) []time.Weekday {
	seen := make(map[time.Weekday]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[time.Weekday(v)]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func (s *ScheduleList) GetEntry(id uint8) *ScheduleEntry {
	for k := range s.List {
		if s.List[k].ID == id {
			return &s.List[k]
		}
	}
	return nil
}

func (s *ScheduleList) AddEntry(e *ScheduleEntry) error {
	log.Info("adding entry", "e", e)

	if len(e.Weekdays) == 0 {
		err := fmt.Errorf("cannot schedule an entry not on any days")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if len(e.Blowers) == 0 && len(e.Pumps) == 0 {
		err := fmt.Errorf("cannot schedule an entry without any pumps or blowers")
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
	schedule.writeToStore()
	return nil
}

func buildJob(e *ScheduleEntry) error {
	attimes := make([]gocron.AtTime, 0)

	times := strings.Split(e.StartTime, ";")
	for _, v := range times {
		log.Info("time", "time", v)
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

	_, err := sz.NewJob(
		gocron.WeeklyJob(
			1,
			gocron.NewWeekdays(e.Weekdays[0], e.Weekdays[1:]...),
			gocron.NewAtTimes(attimes[0], attimes[1:]...),
		),
		gocron.NewTask(
			func() {
				log.Info("starting scheduled entry", "e", e)
				for _, blower := range e.Blowers {
					log.Info("starting blower", "blower", blower, "duration", e.RunTime)
					blower.Start(e.RunTime, "scheduled")
				}
				for _, pump := range e.Pumps {
					log.Info("starting pump", "pump", pump, "duration", e.RunTime)
					pump.Start(e.RunTime, "scheduled")
				}
			},
		),
		gocron.WithTags(fmt.Sprintf("%d", e.ID)),
		gocron.WithName(e.Name),
	)

	return err
}

func (s *ScheduleList) RemoveEntry(id uint8) {
	index := -1
	for k := range schedule.List {
		if s.List[k].ID == id {
			index = k
			break
		}
	}
	if index == -1 {
		log.Info("unknown schedule entry ID", "id", id)
		return
	}
	log.Info("removing job from schedule", "id", id)
	sz.RemoveByTags(fmt.Sprintf("%d", id))
	s.List = append(s.List[:index], s.List[index+1:]...)
	log.Info("new schedule", "s", s.List)
	_ = s.writeToStore()
}

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

	if len(e.Blowers) == 0 && len(e.Pumps) == 0 {
		err := fmt.Errorf("cannot schedule an entry without any pumps or blowers")
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

	log.Info("removing job from schedule", "id", e.ID)
	sz.RemoveByTags(fmt.Sprintf("%d", e.ID))

	if err := buildJob(e); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
