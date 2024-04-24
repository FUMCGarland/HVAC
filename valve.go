package hvac

type valveID uint8

type Valve struct {
	ID    valveID
	Name  string
	State int8
	Min   int8
	Max   int8
}
