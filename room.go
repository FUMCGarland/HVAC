package hvac

type RoomID uint16

type Room struct {
	ID          RoomID
	Name        string
	Zone        ZoneID
	Temperature uint8
	Humidity    uint8
	Occupied    bool
}

func (r RoomID) Get() *Room {
	for k := range c.Rooms {
		if c.Rooms[k].ID == r {
			return &c.Rooms[k]
		}
	}
	return nil
}

func (r *Room) SetTemp(temp uint8) {
	r.Temperature = temp

	if c.ControlMode != ControlTemp {
		return
	}

	// this is too naïve
	// check all rooms in zone, exclude zeros, build avg?
	// once a hour? set everything to zero to skip dead sensors

	// we could get into a state where there is a legit 7 degree difference between two rooms and trigger a battle
	// hold-down timers won't solve that
	// also, finish the hold-down timer logic

	zone := r.Zone.Get()
	switch c.SystemMode {
	case SystemModeHeat:
		// TODO deal with zone average first
		if (r.Occupied && temp < zone.Targets.HeatingOccupiedTemp-3) || (!r.Occupied && temp < zone.Targets.HeatingUnoccupiedTemp-3) {
			zone.ID.Start(defaultRunDuration, "temp")
			return
		}
		if (r.Occupied && temp > zone.Targets.HeatingOccupiedTemp+3) || (!r.Occupied && temp < zone.Targets.HeatingUnoccupiedTemp+3) {
			zone.ID.Stop("temp")
			return
		}
	case SystemModeCool:
		// TODO deal with zone average first
		if (r.Occupied && temp > zone.Targets.CoolingOccupiedTemp+3) || (!r.Occupied && temp > zone.Targets.CoolingUnoccupiedTemp+3) {
			zone.ID.Start(defaultRunDuration, "temp")
			return
		}
		if (r.Occupied && temp < zone.Targets.CoolingOccupiedTemp-3) || (!r.Occupied && temp > zone.Targets.CoolingUnoccupiedTemp-3) {
			zone.ID.Stop("temp")
			return
		}
	}
}

func (r *Room) SetHumidity(humidity uint8) {
	r.Humidity = humidity

	/* if c.ControlMode != ControlTemp {
		return
	} */
	// nothing to do other than accept the update
}
