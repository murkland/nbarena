package chips

import (
	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/state"
)

var Chips = []*state.Chip{
	{
		Index:      0,
		Name:       "Cannon",
		BaseDamage: 40,
		MakeBehavior: func(damage state.Damage) state.EntityBehavior {
			return &behaviors.Cannon{Style: behaviors.CannonStyleCannon, Damage: damage}
		},
	},
	{
		Index:      1,
		Name:       "HiCannon",
		BaseDamage: 100,
		MakeBehavior: func(damage state.Damage) state.EntityBehavior {
			return &behaviors.Cannon{Style: behaviors.CannonStyleHiCannon, Damage: damage}
		},
	},
	{
		Index:      2,
		Name:       "M-Cannon",
		BaseDamage: 180,
		MakeBehavior: func(damage state.Damage) state.EntityBehavior {
			return &behaviors.Cannon{Style: behaviors.CannonStyleMCannon, Damage: damage}
		},
	},
	{
		Index:      3,
		Name:       "AirShot",
		BaseDamage: 20,
		MakeBehavior: func(damage state.Damage) state.EntityBehavior {
			return &behaviors.AirShot{Damage: damage}
		},
	},
	{
		Index:      4,
		Name:       "Vulcan1",
		BaseDamage: 10,
		MakeBehavior: func(damage state.Damage) state.EntityBehavior {
			return &behaviors.Vulcan{Shots: 3, Damage: damage}
		},
	},
	{
		Index:      5,
		Name:       "Vulcan2",
		BaseDamage: 15,
		MakeBehavior: func(damage state.Damage) state.EntityBehavior {
			return &behaviors.Vulcan{Shots: 4, Damage: damage}
		},
	},
	{
		Index:      6,
		Name:       "Vulcan3",
		BaseDamage: 20,
		MakeBehavior: func(damage state.Damage) state.EntityBehavior {
			return &behaviors.Vulcan{Shots: 5, Damage: damage}
		},
	},
	{
		Index:      7,
		Name:       "SuprVulc",
		BaseDamage: 20,
		MakeBehavior: func(damage state.Damage) state.EntityBehavior {
			return &behaviors.Vulcan{Shots: 10, Damage: damage}
		},
	},
	{
		Index:      73,
		Name:       "WideBlde",
		BaseDamage: 150,
		MakeBehavior: func(damage state.Damage) state.EntityBehavior {
			return &behaviors.Sword{Damage: damage, Style: behaviors.SwordStyleBlade, Range: behaviors.SwordRangeWide}
		},
	},
	{
		Index:      79,
		Name:       "WindRack",
		BaseDamage: 140,
		MakeBehavior: func(damage state.Damage) state.EntityBehavior {
			return &behaviors.WindRack{Damage: damage}
		},
	},
	{
		Index:      160,
		Name:       "Recov300",
		BaseDamage: 0,
		MakeBehavior: func(damage state.Damage) state.EntityBehavior {
			return &behaviors.Recov{HP: 300}
		},
	},
	{
		Index: 162,
		Name:  "AreaGrab",
		MakeBehavior: func(damage state.Damage) state.EntityBehavior {
			return &behaviors.AreaGrab{}
		},
	},
}
