package hvac

import (
	"fmt"
	"time"

	"github.com/FUMCGarland/hvac/log"
	"github.com/go-co-op/gocron/v2"
)

// GetOneTimeEntry returns an entry based on ID
func (s *OccupancySchedule) GetOneTimeEntry(id OccupancyOneTimeID) *OccupancyOneTimeEntry {
	for k := range s.OneTime {
		if s.OneTime[k].ID == id {
			return &s.OneTime[k]
		}
	}
	return nil
}

// AddOneTimeEntry adds a new entry to the one-time scheduler
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
	if err := occupancy.writeToStore(); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// buildOneTimeJob adds a job to the onetime scheduler
func buildOneTimeJob(e *OccupancyOneTimeEntry) error {
	var maxPreRunTime time.Duration
	for _, r := range e.Rooms {
		timeDiff, err := r.getPreRunTime()
		if err != nil {
			log.Error(err.Error())
			continue
		}
		if timeDiff > maxPreRunTime {
			maxPreRunTime = timeDiff
		}
	}
	log.Debug("Setting prerun offset", "prerun", maxPreRunTime, "job", e.Name)
	startWithPrerun := e.Start.Add(0 - maxPreRunTime)

	if e.Start.Before(time.Now()) {
		log.Info("not adding job in the past")
		occupancy.RemoveOneTimeEntry(e.ID)
		return nil // not an error, just info
	}

	_, err := occScheduler.NewJob(
		gocron.OneTimeJob(
			gocron.OneTimeJobStartDateTime(startWithPrerun),
		),
		gocron.NewTask(
			func() {
				log.Info("marking room as occupied", "e", e)
				zones := make([]ZoneID, 0)
				zoneActivated := false
				for _, room := range e.Rooms {
					r := room.Get()
					if r == nil {
						log.Warn("got nil room starting recurring occupancy, update the rules")
						continue
					}
					r.Occupied = true
					zoneActivated = false
					for k := range zones {
						if zones[k] == r.Zone {
							zoneActivated = true
						}
					}
					if !zoneActivated {
						log.Debug("activating zone")
						r.Zone.Get().UpdateTemp() // recalculates the avg and runs if needed
						zones = append(zones, r.Zone)
					}
				}
			},
		),
		gocron.WithTags(e.Name, scheduleTagOccupancy, scheduleTagOneTime),
		gocron.WithName(e.Name),
	)
	if err != nil {
		return err
	}

	_, err = occScheduler.NewJob(
		gocron.OneTimeJob(
			gocron.OneTimeJobStartDateTime(e.End),
		),
		gocron.NewTask(
			func() {
				log.Debug("marking room as unoccupied", "e", e)
				for _, room := range e.Rooms {
					r := room.Get()
					r.Occupied = false
					r.Zone.Get().UpdateTemp() // recalculates the avg and runs if needed
				}
				cleanOneTimeSchedule()

			},
		),
		gocron.WithTags(e.Name, scheduleTagOccupancy, scheduleTagOneTime),
		gocron.WithName(fmt.Sprintf("%s end", e.Name)),
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveOneTimeEntry remove an entry from the running scheduler and the configuration
func (s *OccupancySchedule) RemoveOneTimeEntry(id OccupancyOneTimeID) {
	index := -1
	for k := range occupancy.OneTime {
		if s.OneTime[k].ID == id {
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
	s.OneTime = append(s.OneTime[:index], s.OneTime[index+1:]...)
	log.Debug("new schedule", "s", s.OneTime)
	if err := s.writeToStore(); err != nil {
		log.Error(err.Error())
	}
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

	log.Debug("removing job from schedule", "id", e.ID)
	// need to add RemoveByName
	occScheduler.RemoveByTags(fmt.Sprintf("%d", e.ID))

	if err := buildOneTimeJob(e); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func cleanOneTimeSchedule() {
	log.Debug("cleaning one-time occupancy schedule")
	for k := range occupancy.OneTime {
		if occupancy.OneTime[k].Start.Before(time.Now()) {
			occupancy.RemoveOneTimeEntry(occupancy.OneTime[k].ID)
		}
	}
	if err := occupancy.writeToStore(); err != nil {
		log.Error(err.Error())
	}
}
