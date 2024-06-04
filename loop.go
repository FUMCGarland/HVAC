package hvac

// LoopID is the unique identifier for a loop
type LoopID uint8

// A loop is a physical pipe that moves heated or cooled water to blowers or radiant devices around the campus. Loops connect blowers and pumps; some zones are radiant only and do not have blowers.
type Loop struct {
	Name        string
	ID          LoopID
	RadiantZone ZoneID
}
