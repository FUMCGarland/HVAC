package hvac

import (
	"fmt"
	// "time"

	"github.com/FUMCGarland/hvac/log"
	"github.com/go-co-op/gocron/v2"
)

func (s *OccupancySchedule) GetOneTimeEntry(id uint8) *OccupancyOneTimeEntry {
	for k := range s.OneTime {
		if s.OneTime[k].ID == id {
			return &s.OneTime[k]
		}
	}
	return nil
}

func (s *OccupancySchedule) AddOneTimeEntry(e *OccupancyOneTimeEntry) error {
	log.Debug("adding entry", "e", e)

	if len(e.Rooms) == 0 {
		err := fmt.Errorf("cannot schedule an event without any rooms")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if existing := s.GetOneTimeEntry(e.ID); existing != nil {
		err := fmt.Errorf("cannot reuse existing ID")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if e.Name == "" {
		e.Name = fmt.Sprintf("Unnamed %d", e.ID)
	}

	if err := buildOneTimeJob(e); err != nil {
		log.Error(err.Error(), "entry", e)
		return err
	}

	occupancy.OneTime = append(occupancy.OneTime, *e)
	occupancy.writeToStore()
	return nil
}

func buildOneTimeJob(e *OccupancyOneTimeEntry) error {
	var roomStrings []string
	for _, r := range e.Rooms {
		roomStrings = append(roomStrings, fmt.Sprintf("%d", r))
	}

	_, err := occScheduler.NewJob(
		gocron.OneTimeJob(
			gocron.OneTimeJobStartDateTime(e.Start),
		),
		gocron.NewTask(
			func() {
				log.Debug("marking room as occupied", "e", e)
				for _, room := range e.Rooms {
					room.Get().Occupied = true
				}
			},
		),
		gocron.WithTags(e.Name),
		gocron.WithName(e.Name),
	)

	_, err = occScheduler.NewJob(
		gocron.OneTimeJob(
			gocron.OneTimeJobStartDateTime(e.End),
		),
		gocron.NewTask(
			func() {
				log.Debug("marking room as unoccupied", "e", e)
				for _, room := range e.Rooms {
					room.Get().Occupied = false
				}
			},
		),
		gocron.WithTags(e.Name),
		gocron.WithName(e.Name),
	)

	return err
}

func (s *OccupancySchedule) RemoveOneTimeEntry(id uint8) {
	index := -1
	for k := range occupancy.OneTime {
		if s.OneTime[k].ID == id {
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
	s.OneTime = append(s.OneTime[:index], s.OneTime[index+1:]...)
	log.Debug("new schedule", "s", s.OneTime)
	_ = s.writeToStore()
}

// EditEntry updates an entry in the OccupancySchedule, keyed based on e.ID
func (s *OccupancySchedule) EditOneTimeEntry(e *OccupancyOneTimeEntry) error {
	index := -1
	for k := range occupancy.OneTime {
		if s.OneTime[k].ID == e.ID {
			index = k
			break
		}
	}
	if index == -1 {
		err := fmt.Errorf("cannot update non-existent entry")
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

	s.OneTime[index] = *e
	if err := s.writeToStore(); err != nil {
		log.Error(err.Error(), "entry", e)
		return err
	}

	log.Info("removing job from schedule", "id", e.ID)
	occScheduler.RemoveByTags(fmt.Sprintf("%d", e.ID))

	if err := buildOneTimeJob(e); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
