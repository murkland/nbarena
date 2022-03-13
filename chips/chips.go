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
		OnUse: func(s *state.State, e *state.Entity) {
			e.NextBehavior = &behaviors.Cannon{Style: behaviors.CannonStyleCannon, Damage: 40}
		},
	},
	{
		Index:  1,
		Name:   "HiCannon",
		Damage: 100,
		OnUse: func(s *state.State, e *state.Entity) {
			e.NextBehavior = &behaviors.Cannon{Style: behaviors.CannonStyleHiCannon, Damage: 100}
		},
	},
	{
		Index:  2,
		Name:   "M-Cannon",
		Damage: 180,
		OnUse: func(s *state.State, e *state.Entity) {
			e.NextBehavior = &behaviors.Cannon{Style: behaviors.CannonStyleMCannon, Damage: 180}
		},
	},
	{
		Index:  3,
		Name:   "AirShot",
		Damage: 20,
		OnUse: func(s *state.State, e *state.Entity) {
			e.NextBehavior = &behaviors.AirShot{Damage: 20}
		},
	},
	{
		Index:  4,
		Name:   "Vulcan1",
		Damage: 10,
		OnUse: func(s *state.State, e *state.Entity) {
			e.NextBehavior = &behaviors.Vulcan{Shots: 3, Damage: 10}
		},
	},
	{
		Index:  5,
		Name:   "Vulcan2",
		Damage: 15,
		OnUse: func(s *state.State, e *state.Entity) {
			e.NextBehavior = &behaviors.Vulcan{Shots: 4, Damage: 15}
		},
	},
	{
		Index:  6,
		Name:   "Vulcan3",
		Damage: 20,
		OnUse: func(s *state.State, e *state.Entity) {
			e.NextBehavior = &behaviors.Vulcan{Shots: 5, Damage: 20}
		},
	},
	{
		Index:  7,
		Name:   "SuprVulc",
		Damage: 20,
		OnUse: func(s *state.State, e *state.Entity) {
			e.NextBehavior = &behaviors.Vulcan{Shots: 10, Damage: 20}
		},
	},
	{
		Index:  73,
		Name:   "WideBlde",
		Damage: 150,
		OnUse: func(s *state.State, e *state.Entity) {
			e.NextBehavior = &behaviors.Sword{Damage: 150, Style: behaviors.SwordStyleBlade, Range: behaviors.SwordRangeWide}
		},
	},
	{
		Index:  79,
		Name:   "WindRack",
		Damage: 140,
		OnUse: func(s *state.State, e *state.Entity) {
			e.NextBehavior = &behaviors.WindRack{Damage: 140}
		},
	},
	// {
	// 	Index: 162,
	// 	Name:  "AreaGrab",
	// 	OnUse: func(s *state.State, e *state.Entity) {
	// 	},
	// },
}
