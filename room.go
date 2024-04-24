package hvac

type RoomID uint16

type Room struct {
	ID          RoomID
	Name        string
	Zone        ZoneID
	Temperature uint8
}

func (r RoomID) Get() *Room {
	for k := range c.Rooms {
		if c.Rooms[k].ID == r {
			return &c.Rooms[k]
		}
	}
	return nil
}
