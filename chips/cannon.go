package chips

import (
	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/state"
)

var Cannon = &state.Chip{
	Index:      0,
	Name:       "Cannon",
	BaseDamage: 40,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Cannon{Style: behaviors.CannonStyleCannon, Damage: damage}
	},
}

var HiCannon = &state.Chip{
	Index:      1,
	Name:       "HiCannon",
	BaseDamage: 100,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Cannon{Style: behaviors.CannonStyleHiCannon, Damage: damage}
	},
}

var MCannon = &state.Chip{
	Index:      2,
	Name:       "M-Cannon",
	BaseDamage: 180,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Cannon{Style: behaviors.CannonStyleMCannon, Damage: damage}
	},
}
