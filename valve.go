package hvac

type ValveID uint8

// Valve is a type for future use when we build controllers to adjust the flow on the loops
type Valve struct {
	Name  string
	ID    ValveID
	State int8
	Min   int8
	Max   int8
}
