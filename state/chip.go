package state

type Chip struct {
	Index           int
	Name            string
	BehaviorFactory func() EntityBehavior
}

func (c Chip) Clone() Chip {
	// Chips are immutable.
	return c
}
