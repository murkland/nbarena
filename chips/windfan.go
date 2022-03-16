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
			TilePos: state.TilePosXY(x+dx, y),

			IsFlipped:            e.IsFlipped,
			IsAlliedWithAnswerer: e.IsAlliedWithAnswerer,

			Traits: state.EntityTraits{
				CanStepOnHoleLikeTiles: true,
				IgnoresTileEffects:     true,
				CannotFlinch:           true,
				IgnoresTileOwnership:   true,
				CannotSlide:            true,
				Intangible:             true,
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
	OnUse:      makeWindFanOnUse(false),
}

var Fan = &state.Chip{
	Index:      129,
	Name:       "Fan",
	BaseDamage: 0,
	OnUse:      makeWindFanOnUse(true),
}
