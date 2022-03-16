package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type WindFan struct {
	Owner state.EntityID
	IsFan bool
}

func (eb *WindFan) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *WindFan) Clone() state.EntityBehavior {
	return &WindFan{eb.Owner, eb.IsFan}
}

func (eb *WindFan) Cleanup(e *state.Entity, s *state.State) {
}

func (eb *WindFan) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime%13 == 0 {
		// TODO: Should spawn gusts on edge of opponent area.
		y := 2 - ((int(e.BehaviorState.ElapsedTime) / 13) % 3) + 1
		x := 1
		if eb.IsFan {
			x = state.TileCols - 2
		}

		if e.IsAlliedWithAnswerer {
			x = state.TileCols - x - 1
		}

		gustStyle := GustStyleWind
		if eb.IsFan {
			gustStyle = GustStyleFan
		}

		isFlipped := eb.IsFan
		if e.IsAlliedWithAnswerer {
			isFlipped = !isFlipped
		}

		s.AttachEntity(&state.Entity{
			TilePos: state.TilePosXY(x, y),

			IsFlipped:            isFlipped,
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
				Behavior: &Gust{eb.Owner, gustStyle, true},
			},
		})
	}
}

func (eb *WindFan) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	img := b.WindFanSprites.WindImage
	if eb.IsFan {
		img = b.WindFanSprites.FanImage
	}

	rootNode := &draw.OptionsNode{}
	rootNode.Children = append(rootNode.Children, draw.ImageWithAnimation(img, b.WindFanSprites.Animation, int(e.BehaviorState.ElapsedTime)))
	return rootNode
}
