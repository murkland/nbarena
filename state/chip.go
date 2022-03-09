package state

type Chip struct {
	Index  int
	Name   string
	Damage int
	OnUse  func(s *State, e *Entity)
}

func (c Chip) Clone() Chip {
	// Chips are immutable.
	return c
}
