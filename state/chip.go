package state

type Chip struct {
	Index      int
	Name       string
	BaseDamage int
	OnUse      func(s *State, e *Entity, damage Damage)
}

func (c Chip) Clone() Chip {
	// Chips are immutable.
	return c
}
