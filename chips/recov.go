package chips

import (
	"image"

	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/state"
)

var Recov300 = state.Chip{
	Index:      160,
	Name:       "Recov300",
	BaseDamage: 0,
	OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
		e.HP += 300
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
	},
}
