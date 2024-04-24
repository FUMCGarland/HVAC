package hvac

type damperID uint8

type Damper struct {
	ID    damperID
	Name  string
	State int8
	Min   int8
	Max   int8
}
