package chips

import (
	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/state"
)

var WindRack = &state.Chip{
	Index:      79,
	Name:       "WindRack",
	BaseDamage: 140,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.WindRack{Damage: damage}
	},
}
