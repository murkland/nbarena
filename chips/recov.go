package chips

import (
	"image"

	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/state"
)

func makeRecovOnUse(hp int) func(s *state.State, e *state.Entity, damage state.Damage) {
	return func(s *state.State, e *state.Entity, damage state.Damage) {
		e.HP += hp
		if e.HP > e.MaxHP {
			e.HP = e.MaxHP
		}

		// TODO: Trigger antirecovery.

		s.AttachSound(&state.Sound{
			Type: bundle.SoundTypeRecov,
		})
		s.AttachDecoration(&state.Decoration{
			Type:      bundle.DecorationTypeRecov,
			TilePos:   e.TilePos,
			Offset:    image.Point{0, 0},
			IsFlipped: e.IsFlipped,
		})

		e.ChipUseLockoutTimeLeft = 30
	}
}

var Recov10 = state.Chip{
	Index:      153,
	Name:       "Recov10",
	BaseDamage: 0,
	OnUse:      makeRecovOnUse(10),
}

var Recov30 = state.Chip{
	Index:      154,
	Name:       "Recov30",
	BaseDamage: 0,
	OnUse:      makeRecovOnUse(30),
}

var Recov50 = state.Chip{
	Index:      155,
	Name:       "Recov50",
	BaseDamage: 0,
	OnUse:      makeRecovOnUse(50),
}

var Recov80 = state.Chip{
	Index:      156,
	Name:       "Recov80",
	BaseDamage: 0,
	OnUse:      makeRecovOnUse(80),
}

var Recov120 = state.Chip{
	Index:      157,
	Name:       "Recov120",
	BaseDamage: 0,
	OnUse:      makeRecovOnUse(120),
}

var Recov150 = state.Chip{
	Index:      158,
	Name:       "Recov150",
	BaseDamage: 0,
	OnUse:      makeRecovOnUse(150),
}

var Recov200 = state.Chip{
	Index:      159,
	Name:       "Recov200",
	BaseDamage: 0,
	OnUse:      makeRecovOnUse(200),
}

var Recov300 = state.Chip{
	Index:      160,
	Name:       "Recov300",
	BaseDamage: 0,
	OnUse:      makeRecovOnUse(300),
}
