package chips

import (
	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/state"
)

var AirShot = &state.Chip{
	Index:      3,
	Name:       "AirShot",
	BaseDamage: 20,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.AirShot{Damage: damage}
	},
}
