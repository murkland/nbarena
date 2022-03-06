package chips

import (
	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/state"
)

var Chips = []state.Chip{
	{
		Index:  0,
		Name:   "Cannon",
		Damage: 40,
		BehaviorFactory: func() state.EntityBehavior {
			return &behaviors.Cannon{Style: behaviors.CannonStyleCannon, Damage: 40}
		},
	},
	{
		Index:  1,
		Name:   "HiCannon",
		Damage: 100,
		BehaviorFactory: func() state.EntityBehavior {
			return &behaviors.Cannon{Style: behaviors.CannonStyleHiCannon, Damage: 100}
		},
	},
	{
		Index:  1,
		Name:   "M-Cannon",
		Damage: 180,
		BehaviorFactory: func() state.EntityBehavior {
			return &behaviors.Cannon{Style: behaviors.CannonStyleMCannon, Damage: 180}
		},
	},
	{
		Index:  70,
		Name:   "Sword",
		Damage: 80,
		BehaviorFactory: func() state.EntityBehavior {
			return &behaviors.Sword{Damage: 80}
		},
	},
}
