package hvac

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/FUMCGarland/hvac/log"
	"github.com/go-co-op/gocron/v2"
	// "github.com/google/uuid"
)

func (s *OccupancySchedule) GetRecurringEntry(id uint8) *OccupancyRecurringEntry {
	for k := range s.Recurring {
		if s.Recurring[k].ID == id {
			return &s.Recurring[k]
		}
	}
	return nil
}

func (s *OccupancySchedule) AddRecurringEntry(e *OccupancyRecurringEntry) error {
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

	if existing := s.GetRecurringEntry(e.ID); existing != nil {
		err := fmt.Errorf("cannot reuse existing ID")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if e.Name == "" {
		e.Name = fmt.Sprintf("Unnamed %d", e.ID)
	}

	if err := buildRecurringJob(e); err != nil {
		log.Error(err.Error(), "entry", e)
		return err
	}

	occupancy.Recurring = append(occupancy.Recurring, *e)
	occupancy.writeToStore()
	return nil
}

func buildRecurringJob(e *OccupancyRecurringEntry) error {
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

	_, err := occScheduler.NewJob(
		gocron.WeeklyJob(
			1,
			func() []time.Weekday { return e.Weekdays },
			func() []gocron.AtTime { return attimes },
		),
		gocron.NewTask(
			func() {
				log.Debug("marking room as occupied", "e", e)
				for _, room := range e.Rooms {
					room.Get().Occupied = true
				}
			},
		),
		gocron.WithTags(roomStrings...),
		gocron.WithName(e.Name),
	)

	return err
}

func (s *OccupancySchedule) RemoveRecurringEntry(id uint8) {
	index := -1
	for k := range occupancy.Recurring {
		if s.Recurring[k].ID == id {
			index = k
			break
		}
	}
	if index == -1 {
		log.Info("unknown occupancy recurring entry ID", "id", id)
		return
	}
	log.Info("removing job from recurring schedule", "id", id)
	occScheduler.RemoveByTags(fmt.Sprintf("%d", id))
	s.Recurring = append(s.Recurring[:index], s.Recurring[index+1:]...)
	log.Debug("new schedule", "s", s.Recurring)
	_ = s.writeToStore()
}

// EditEntry updates an entry in the OccupancySchedule, keyed based on e.ID
func (s *OccupancySchedule) EditRecuringEntry(e *OccupancyRecurringEntry) error {
	index := -1
	for k := range occupancy.Recurring {
		if s.Recurring[k].ID == e.ID {
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

	s.Recurring[index] = *e
	if err := s.writeToStore(); err != nil {
		log.Error(err.Error(), "entry", e)
		return err
	}

	log.Info("removing job from schedule", "id", e.ID)
	occScheduler.RemoveByTags(fmt.Sprintf("%d", e.ID))

	if err := buildRecurringJob(e); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
