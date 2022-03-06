package chips

import (
	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/state"
)

var Chips = []state.Chip{
	{
		Index: 0,
		Name:  "Cannon",
		BehaviorFactory: func() state.EntityBehavior {
			return &behaviors.Cannon{Style: behaviors.CannonStyleCannon, Damage: 40}
		},
	},
	{
		Index: 1,
		Name:  "HiCannon",
		BehaviorFactory: func() state.EntityBehavior {
			return &behaviors.Cannon{Style: behaviors.CannonStyleHiCannon, Damage: 100}
		},
	},
	{
		Index: 1,
		Name:  "M-Cannon",
		BehaviorFactory: func() state.EntityBehavior {
			return &behaviors.Cannon{Style: behaviors.CannonStyleMCannon, Damage: 180}
		},
	},
	{
		Index: 70,
		Name:  "Sword",
		BehaviorFactory: func() state.EntityBehavior {
			return &behaviors.Sword{Damage: 80}
		},
	},
}
