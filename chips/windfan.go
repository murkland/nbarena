package chips

import (
	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/state"
)

func makeWindFanOnUse(isFan bool) func(s *state.State, e *state.Entity, damage state.Damage) {
	return func(s *state.State, e *state.Entity, damage state.Damage) {
		x, y := e.TilePos.XY()
		dx, _ := e.Facing().XY()
		s.AttachEntity(&state.Entity{
			TilePos:     state.TilePosXY(x+dx, y),
			MaxLifeTime: 1440,

			HP:    40,
			MaxHP: 40,

			IsFlipped:            e.IsFlipped,
			IsAlliedWithAnswerer: e.IsAlliedWithAnswerer,

			Traits: state.EntityTraits{
				IgnoresTileEffects:   true,
				CannotFlinch:         true,
				CannotFlash:          true,
				StatusGuard:          true,
				IgnoresTileOwnership: true,
				CannotSlide:          true,
			},

			BehaviorState: state.EntityBehaviorState{
				Behavior: &behaviors.WindFan{Owner: e.ID(), IsFan: isFan},
			},
		})

	}
}

var Wind = &state.Chip{
	Index:      128,
	Name:       "Wind",
	BaseDamage: 0,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return nil
	},
}

var Fan = &state.Chip{
	Index:      129,
	Name:       "Fan",
	BaseDamage: 0,
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		return nil
	},
}
