package chips

import (
	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/state"
)

var Sword = &state.Chip{
	Index:      70,
	Name:       "Sword",
	BaseDamage: 80,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Sword{Damage: damage, Style: behaviors.SwordStyleSword, Range: behaviors.SwordRangeShort}
	},
}

var WideSwrd = &state.Chip{
	Index:      71,
	Name:       "WideSwrd",
	BaseDamage: 80,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Sword{Damage: damage, Style: behaviors.SwordStyleSword, Range: behaviors.SwordRangeWide}
	},
}

var LongSwrd = &state.Chip{
	Index:      72,
	Name:       "LongSwrd",
	BaseDamage: 100,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Sword{Damage: damage, Style: behaviors.SwordStyleSword, Range: behaviors.SwordRangeLong}
	},
}

var WideBlde = &state.Chip{
	Index:      73,
	Name:       "WideBlde",
	BaseDamage: 150,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Sword{Damage: damage, Style: behaviors.SwordStyleBlade, Range: behaviors.SwordRangeWide}
	},
}

var LongBlde = &state.Chip{
	Index:      74,
	Name:       "LongBlde",
	BaseDamage: 150,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Sword{Damage: damage, Style: behaviors.SwordStyleBlade, Range: behaviors.SwordRangeLong}
	},
}
