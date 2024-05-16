package hvac

type damperID uint8

// Damper is a type to be used in the future for controlling how "open" a blower is
type Damper struct {
	ID    damperID
	Name  string
	State int8
	Min   int8
	Max   int8
}
