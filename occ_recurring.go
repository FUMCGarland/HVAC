package hvac

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/FUMCGarland/hvac/log"
	"github.com/go-co-op/gocron/v2"
)

func (s *OccupancySchedule) GetRecurringEntry(id OccupancyRecurringID) *OccupancyRecurringEntry {
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
	if err := occupancy.writeToStore(); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func buildRecurringJob(e *OccupancyRecurringEntry) error {
	attimes := make([]gocron.AtTime, 0)

	var maxPreRunTime time.Duration
	for _, r := range e.Rooms {
		timeDiff, err := r.GetPreRunTime()
		if err != nil {
			log.Error(err.Error())
			continue
		}
		if timeDiff > maxPreRunTime {
			maxPreRunTime = timeDiff
		}
	}
	preRunHours := int64(maxPreRunTime.Hours())
	preRunMinutes := int64(maxPreRunTime.Minutes()) % 60

	times := strings.Split(e.StartTime, ";")
	for _, v := range times {
		log.Debug("time", "time", v)
		units := strings.Split(v, ":")
		hour, err := strconv.ParseInt(units[0], 10, 8)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		hour = (hour - preRunHours) % 24

		minute, err := strconv.ParseInt(units[1], 10, 8)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		minute = (minute - preRunMinutes) % 60
		if minute < 0 {
			hour = hour - 1
			minute = minute + 60
		}

		log.Debug("time", "configured time", v, "prerun hour", hour, "prerun minute", minute)
		attimes = append(attimes, gocron.NewAtTime(uint(hour), uint(minute), 0))
	}

	newTags := []string{e.Name, scheduleTagOccupancy, scheduleTagRecurring}
	for _, r := range e.Rooms {
		newTags = append(newTags, fmt.Sprintf("%d", r))
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
					r := room.Get()
					if r == nil {
						log.Warn("got nil room starting recurring occupancy, update the rules")
						continue
					}
					r.Occupied = true
				}
			},
		),
		gocron.WithTags(newTags...),
		gocron.WithName(e.Name),
	)

	return err
}

func (s *OccupancySchedule) RemoveRecurringEntry(id OccupancyRecurringID) {
	index := -1
	for k := range occupancy.Recurring {
		if s.Recurring[k].ID == id {
			index = k
			break
		}
	}
	if index == -1 {
		log.Warn("unknown occupancy recurring entry ID", "id", id)
		return
	}
	log.Debug("removing job from recurring schedule", "id", id)
	occScheduler.RemoveByTags(fmt.Sprintf("%d", id))
	s.Recurring = append(s.Recurring[:index], s.Recurring[index+1:]...)
	log.Debug("new schedule", "s", s.Recurring)
	if err := s.writeToStore(); err != nil {
		log.Error(err.Error())
	}
}

// EditEntry updates an entry in the OccupancySchedule, keyed based on e.ID
func (s *OccupancySchedule) EditRecurringEntry(e *OccupancyRecurringEntry) error {
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

	log.Debug("removing job from schedule", "id", e.ID)
	occScheduler.RemoveByTags(fmt.Sprintf("%d", e.ID))

	if err := buildRecurringJob(e); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
