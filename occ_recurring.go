package hvac

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/FUMCGarland/hvac/log"
	"github.com/go-co-op/gocron/v2"
)

func (s *OccupancySchedule) GetRecurringEntry(id OccupancyRecurringID) *OccupancyRecurringEntry {
	for _, k := range s.Recurring {
		if k.ID == id {
			return k
		}
	}
	return nil
}

func (s *OccupancySchedule) getLowestUnusedRecurringID() OccupancyRecurringID {
	used := make([]OccupancyRecurringID, 0, len(s.Recurring))

	for _, k := range s.Recurring {
		used = append(used, k.ID)
	}
	slices.Sort(used)

	var i OccupancyRecurringID
	for i = 0; i < OccupancyRecurringID(len(used)); i++ {
		if !slices.Contains(used, i) {
			return i
		}
	}
	return i + OccupancyRecurringID(1)
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
		// return err
		e.ID = s.getLowestUnusedRecurringID()
	}

	if e.Name == "" {
		e.Name = fmt.Sprintf("Unnamed %d", e.ID)
	}

	if err := buildRecurringJob(e); err != nil {
		log.Error(err.Error(), "entry", e)
		return err
	}

	s.Recurring = append(s.Recurring, e)
	if err := s.writeToStore(); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func buildRecurringJob(e *OccupancyRecurringEntry) error {
	attimes := make([]gocron.AtTime, 0)

	times := strings.Split(e.StartTime, ";")
	for _, v := range times {
		log.Debug("time", "in time", v)
		units := strings.Split(v, ":")
		hour, err := strconv.ParseInt(units[0], 10, 8)
		if err != nil {
			log.Error(err.Error())
			return err
		}

		minute, err := strconv.ParseInt(units[1], 10, 8)
		if err != nil {
			log.Error(err.Error())
			return err
		}

		log.Debug("time", "configured time", v)
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
				log.Info("recurring occupancy start", "name", e.Name, "rooms", e.Rooms)
				zoneActivated := false
				zones := make([]ZoneID, 0)

				for _, room := range e.Rooms {
					r := room.Get()
					if r == nil {
						log.Warn("got nil room starting recurring occupancy, update the rules")
						continue
					}
					r.Occupied = true

					zoneActivated = false
					for _, k := range zones {
						if k == r.GetZoneIDInMode() {
							zoneActivated = true
						}
					}
					if !zoneActivated {
						log.Info("activating zone")
						r.GetZoneInMode().UpdateTemp() // recalculates the avg and runs if needed
						zones = append(zones, r.GetZoneIDInMode())
					}
				}
			},
		),
		gocron.WithTags(newTags...),
		gocron.WithName(e.Name),
	)
	if err != nil {
		return err
	}

	endtimes := make([]gocron.AtTime, 0)
	etimes := strings.Split(e.EndTime, ";")
	for _, v := range etimes {
		log.Debug("end time", "in time", v)
		units := strings.Split(v, ":")
		hour, err := strconv.ParseInt(units[0], 10, 8)
		if err != nil {
			log.Error(err.Error())
			return err
		}

		minute, err := strconv.ParseInt(units[1], 10, 8)
		if err != nil {
			log.Error(err.Error())
			return err
		}

		log.Debug("end time", "configured time", v)
		endtimes = append(endtimes, gocron.NewAtTime(uint(hour), uint(minute), 0))
	}

	_, err = occScheduler.NewJob(
		gocron.WeeklyJob(
			1,
			func() []time.Weekday { return e.Weekdays },
			func() []gocron.AtTime { return endtimes },
		),
		gocron.NewTask(
			func() {
				log.Info("recurring occupancy end", "name", e.Name, "rooms", e.Rooms)
				for _, room := range e.Rooms {
					r := room.Get()
					r.Occupied = false
					r.GetZoneInMode().UpdateTemp() // recalculates the avg and runs if needed
				}
				cleanOneTimeSchedule()

			},
		),
		gocron.WithTags(newTags...),
		gocron.WithName(fmt.Sprintf("%s end", e.Name)),
	)

	return err
}

func (s *OccupancySchedule) RemoveRecurringEntry(id OccupancyRecurringID) {
	index := -1
	for i, k := range s.Recurring {
		if k.ID == id {
			index = i
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
	for i, k := range s.Recurring {
		if k.ID == e.ID {
			index = i
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

	s.Recurring[index] = e
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
