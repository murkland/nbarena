package behaviors

import (
	"github.com/murkland/nbarena/state"
)

type UseChip struct {
	state.EntityBehavior
	Chip state.Chip
}

func (eb *UseChip) Clone() state.EntityBehavior {
	return &UseChip{eb.EntityBehavior.Clone(), eb.Chip}
}

func UseNextChip(e *state.Entity) bool {
	if len(e.Chips) == 0 {
		return false
	}
	chip := e.Chips[len(e.Chips)-1]
	e.Chips = e.Chips[:len(e.Chips)-1]
	e.NextBehavior = &UseChip{EntityBehavior: chip.MakeBehavior(), Chip: chip}
	return true
}
