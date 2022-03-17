package chips

import (
	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/state"
)

var Recov10 = &state.Chip{
	Index:      153,
	Name:       "Recov10",
	BaseDamage: 0,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Recov{HP: 10}
	},
}

var Recov30 = &state.Chip{
	Index:      154,
	Name:       "Recov30",
	BaseDamage: 0,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Recov{HP: 30}
	},
}

var Recov50 = &state.Chip{
	Index:      155,
	Name:       "Recov50",
	BaseDamage: 0,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Recov{HP: 50}
	},
}

var Recov80 = &state.Chip{
	Index:      156,
	Name:       "Recov80",
	BaseDamage: 0,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Recov{HP: 80}
	},
}

var Recov120 = &state.Chip{
	Index:      157,
	Name:       "Recov120",
	BaseDamage: 0,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Recov{HP: 120}
	},
}

var Recov150 = &state.Chip{
	Index:      158,
	Name:       "Recov150",
	BaseDamage: 0,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Recov{HP: 150}
	},
}

var Recov200 = &state.Chip{
	Index:      159,
	Name:       "Recov200",
	BaseDamage: 0,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Recov{HP: 200}
	},
}

var Recov300 = &state.Chip{
	Index:      160,
	Name:       "Recov300",
	BaseDamage: 0,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return &behaviors.Recov{HP: 300}
	},
}
