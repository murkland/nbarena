package chips

import (
	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/state"
)

var AirShot = &state.Chip{
	Index:      3,
	Name:       "AirShot",
	BaseDamage: 20,
	OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
		e.NextBehavior = &behaviors.AirShot{Damage: damage}
	},
}
