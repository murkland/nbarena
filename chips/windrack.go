package chips

import (
	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/state"
)

var WindRack = &state.Chip{
	Index:      79,
	Name:       "WindRack",
	BaseDamage: 140,
	OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
		e.NextBehavior = &behaviors.WindRack{Damage: damage}
	},
}
