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
			e.SetBehavior(&behaviors.Cannon{Style: behaviors.CannonStyleCannon, Damage: 40}, s)
		},
	},
	{
		Index:  1,
		Name:   "HiCannon",
		Damage: 100,
		OnUse: func(s *state.State, e *state.Entity) {
			e.SetBehavior(&behaviors.Cannon{Style: behaviors.CannonStyleHiCannon, Damage: 100}, s)
		},
	},
	{
		Index:  1,
		Name:   "M-Cannon",
		Damage: 180,
		OnUse: func(s *state.State, e *state.Entity) {
			e.SetBehavior(&behaviors.Cannon{Style: behaviors.CannonStyleMCannon, Damage: 180}, s)
		},
	},
	{
		Index:  73,
		Name:   "WideBlde",
		Damage: 150,
		OnUse: func(s *state.State, e *state.Entity) {
			e.SetBehavior(&behaviors.Sword{Damage: 150, Style: behaviors.SwordStyleBlade, Range: behaviors.SwordRangeWide}, s)
		},
	},
	{
		Index:  79,
		Name:   "WindRack",
		Damage: 140,
		OnUse: func(s *state.State, e *state.Entity) {
			e.SetBehavior(&behaviors.WindRack{Damage: 140}, s)
		},
	},
	{
		Index: 162,
		Name:  "AreaGrab",
		OnUse: func(s *state.State, e *state.Entity) {
		},
	},
}
