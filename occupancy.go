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
	// "github.com/google/uuid"
)

// OccupancyScheduleList is just a wrapper so we can hang methods on it
type OccupancyScheduleList struct {
	List []OccupancyScheduleEntry
}

var occupancy OccupancyScheduleList
var soccupancy gocron.Scheduler

// ScheduleEntry is the definition of a job to be run at specified times
type OccupancyScheduleEntry struct {
	ID        uint8
	Name      string
	StartTime string // "6:30"
	EndTime   string // "15:30"
	Weekdays  []time.Weekday
	Rooms     []RoomID
}

func init() {
	var err error
	soccupancy, err = gocron.NewScheduler()
	if err != nil {
		log.Error(err.Error())
	}
}

// GetOccupancySchedule returns the live schedule
func (c *Config) GetOccupancySchedule() (*OccupancyScheduleList, error) {
	if len(occupancy.List) == 0 {
		s, err := readOccupancyFromStore()
		if err != nil {
			log.Error(err.Error())
			return &occupancy, err
		}
		occupancy = *s
	}
	return &occupancy, nil
}

func (s *OccupancyScheduleList) writeToStore() error {
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

func readOccupancyFromStore() (*OccupancyScheduleList, error) {
	sl := OccupancyScheduleList{}
	sl.List = []OccupancyScheduleEntry{}

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

	if len(sl.List) == 0 {
		sl.List = make([]OccupancyScheduleEntry, 0)
	}

	for k := range sl.List {
		log.Info("loading occupancy entry", "entry", sl.List[k])
		if err := buildOccupancyJob(&sl.List[k]); err != nil {
			log.Error(err.Error())
			return &sl, err
		}
	}

	return &sl, nil
}

// GetEntry returns an entry in the OccupancyScheduleList by ID
func (s *OccupancyScheduleList) GetEntry(id uint8) *OccupancyScheduleEntry {
	for k := range s.List {
		if s.List[k].ID == id {
			return &s.List[k]
		}
	}
	return nil
}

// AddEntry adds a new entry to the list of jobs to run
func (s *OccupancyScheduleList) AddEntry(e *OccupancyScheduleEntry) error {
	log.Debug("adding entry", "e", e)

	if len(e.Weekdays) == 0 {
		err := fmt.Errorf("cannot schedule an entry not on any days")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if len(e.Rooms) == 0 {
		err := fmt.Errorf("cannot schedule an entry without any rooms")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if existing := s.GetEntry(e.ID); existing != nil {
		err := fmt.Errorf("cannot reuse existing ID")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if e.Name == "" {
		e.Name = fmt.Sprintf("Unnamed %d", e.ID)
	}

	if err := buildOccupancyJob(e); err != nil {
		log.Error(err.Error(), "entry", e)
		return err
	}

	occupancy.List = append(occupancy.List, *e)
	occupancy.writeToStore()
	return nil
}

// buildJob processes a ScheduleEntry and loads it into gocron
func buildOccupancyJob(e *OccupancyScheduleEntry) error {
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

	var roomStrings []string
	for _, r := range e.Rooms {
		roomStrings = append(roomStrings, fmt.Sprintf("%d", r))
	}

	_, err := soccupancy.NewJob(
		gocron.WeeklyJob(
			1,
			func() []time.Weekday { return e.Weekdays },
			func() []gocron.AtTime { return attimes },
		),
		gocron.NewTask(
			func() {
				log.Debug("starting scheduled entry", "e", e)
				for _, room := range e.Rooms {
					room.Get().Occupied = true
				}
			},
		),
		gocron.WithTags(roomStrings...),
		gocron.WithName(e.Name),
		/* 
		gocron.WithEventListeners(
			gocron.AfterJobRuns(
				func(jobID uuid.UUID, jobName string) {
					for _, room := range e.Rooms {
						room.Get().Occupied = false
					}
				},
			),
		),
		*/
	)

	return err
}

// RemoveEntry removes a ScheduleEntry from the list by ID
func (s *OccupancyScheduleList) RemoveEntry(id uint8) {
	index := -1
	for k := range occupancy.List {
		if s.List[k].ID == id {
			index = k
			break
		}
	}
	if index == -1 {
		log.Info("unknown occupancy schedule entry ID", "id", id)
		return
	}
	log.Info("removing job from occupancy schedule", "id", id)
	soccupancy.RemoveByTags(fmt.Sprintf("%d", id))
	s.List = append(s.List[:index], s.List[index+1:]...)
	log.Debug("new schedule", "s", s.List)
	_ = s.writeToStore()
}

// EditEntry updates an entry in the ScheduleList, keyed based on e.ID
func (s *OccupancyScheduleList) EditEntry(e *OccupancyScheduleEntry) error {
	index := -1
	for k := range occupancy.List {
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

	if len(e.Rooms) == 0 {
		err := fmt.Errorf("cannot schedule an entry without rooms")
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
	soccupancy.RemoveByTags(fmt.Sprintf("%d", e.ID))

	if err := buildOccupancyJob(e); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
