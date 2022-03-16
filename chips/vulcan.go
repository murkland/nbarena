package chips

import (
	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/state"
)

var Vulcan1 = &state.Chip{
	Index:      4,
	Name:       "Vulcan1",
	BaseDamage: 10,
	OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
		e.NextBehavior = &behaviors.Vulcan{Shots: 3, Damage: damage, IsSuper: false}
	},
}

var Vulcan2 = &state.Chip{
	Index:      5,
	Name:       "Vulcan2",
	BaseDamage: 15,
	OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
		e.NextBehavior = &behaviors.Vulcan{Shots: 4, Damage: damage, IsSuper: false}
	},
}

var Vulcan3 = &state.Chip{
	Index:      6,
	Name:       "Vulcan3",
	BaseDamage: 20,
	OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
		e.NextBehavior = &behaviors.Vulcan{Shots: 5, Damage: damage, IsSuper: false}
	},
}

var SuprVulc = &state.Chip{
	Index:      7,
	Name:       "SuprVulc",
	BaseDamage: 20,
	OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
		e.NextBehavior = &behaviors.Vulcan{Shots: 10, Damage: damage, IsSuper: true}
	},
}
