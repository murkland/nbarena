package state

type Chip struct {
	Index        int
	Name         string
	BaseDamage   int
	MakeBehavior func(damage Damage) EntityBehavior
}

func (c Chip) Clone() Chip {
	// Chips are immutable.
	return c
}
