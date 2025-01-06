package hvac

import (
	"fmt"
	"time"

	"github.com/FUMCGarland/hvac/log"
	"github.com/go-co-op/gocron/v2"
)

// GetOneTimeEntry returns an entry based on ID
func (s *OccupancySchedule) GetOneTimeEntry(id OccupancyOneTimeID) *OccupancyOneTimeEntry {
	for _, k := range s.OneTime {
		if k.ID == id {
			return k
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

	s.OneTime = append(s.OneTime, e)
	if err := s.writeToStore(); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// buildOneTimeJob adds a job to the onetime scheduler
func buildOneTimeJob(e *OccupancyOneTimeEntry) error {
	_, err := occScheduler.NewJob(
		gocron.OneTimeJob(
			gocron.OneTimeJobStartDateTime(e.Start),
		),
		gocron.NewTask(
			func() {
				log.Info("marking room as occupied", "e", e)
				zones := make([]ZoneID, 0)
				for _, room := range e.Rooms {
					r := room.Get()
					if r == nil {
						log.Warn("got invalid room starting recurring occupancy, update the rules")
						continue
					}
					r.Occupied = true

					// only run the zone temp recalc once per zone
					zoneActivated := false
					roomZone := r.GetZoneIDInMode()
					for _, checkzone := range zones {
						if checkzone == roomZone {
							zoneActivated = true
						}
					}
					if !zoneActivated {
						log.Debug("activating zone")
						roomZone.Get().UpdateTemp() // recalculates the avg and runs if needed
						zones = append(zones, roomZone)
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
				log.Info("marking room as unoccupied", "e", e)
				for _, room := range e.Rooms {
					r := room.Get()
					r.Occupied = false
					// recalc for every active room in every zone; we want to be sure to shut down on last-out
					r.GetZoneInMode().UpdateTemp()
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
	for i, k := range s.OneTime {
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
	s.OneTime = append(s.OneTime[:index], s.OneTime[index+1:]...)
	log.Debug("new schedule", "s", s.OneTime)
	if err := s.writeToStore(); err != nil {
		log.Error(err.Error())
	}
}

// EditEntry updates an entry in the OccupancySchedule, keyed based on e.ID
func (s *OccupancySchedule) EditOneTimeEntry(e *OccupancyOneTimeEntry) error {
	index := -1
	for i, k := range s.OneTime {
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

	if len(e.Rooms) == 0 {
		err := fmt.Errorf("cannot schedule an entry without rooms")
		log.Error(err.Error(), "entry", e)
		return err
	}

	if e.Name == "" {
		e.Name = fmt.Sprintf("Unnamed %d", e.ID)
	}

	s.OneTime[index] = e
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

	// remove from our schedule
	for _, k := range occupancy.OneTime {
		if k.ID != 0 && k.Start.Before(time.Now()) {
			occupancy.RemoveOneTimeEntry(k.ID)
		}
	}
	if err := occupancy.writeToStore(); err != nil {
		log.Error(err.Error())
	}

	// remove from the gocron schedule
	cleanOldJobs()
}
