package behaviors

import (
	"image"

	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/state/query"
)

type WindRack struct {
	Damage state.Damage
}

func (eb *WindRack) Flip() {
}

func (eb *WindRack) Clone() state.EntityBehavior {
	return &WindRack{
		eb.Damage,
	}
}

func (eb *WindRack) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{
		CanBeCountered: true,
	}
}

func (eb *WindRack) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == 0 {
		s.AddDecoration(&state.Decoration{
			Type:      bundle.DecorationTypeWindSlash,
			TilePos:   e.TilePos,
			Offset:    image.Point{0, -16},
			IsFlipped: e.IsFlipped,
		})
	} else if e.BehaviorState.ElapsedTime == 9 {
		x, y := e.TilePos.XY()
		dx := query.DXForward(e.IsFlipped)

		var entities []*state.Entity
		entities = append(entities, query.TangibleEntitiesAt(s, state.TilePosXY(x+dx, y))...)
		entities = append(entities, query.TangibleEntitiesAt(s, state.TilePosXY(x+dx, y+1))...)
		entities = append(entities, query.TangibleEntitiesAt(s, state.TilePosXY(x+dx, y-1))...)

		for _, target := range entities {
			if target.FlashingTimeLeft == 0 {
				var h state.Hit
				h.Traits.Drag = state.DragTypeBig
				h.Traits.SlideDirection = e.Facing()
				h.Traits.Element = state.ElementWind
				h.AddDamage(eb.Damage)
				target.Hit.Merge(h)
			}
		}

		for i := 1; i <= 3; i++ {
			s.AddEntity(&state.Entity{
				TilePos: state.TilePosXY(x+dx, i),

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
					Behavior: &Gust{e.Facing()},
				},
			})
		}
	} else if e.BehaviorState.ElapsedTime == 27-1 {
		e.NextBehavior = &Idle{}
	}
}

func (eb *WindRack) Cleanup(e *state.Entity, s *state.State) {
}

func (eb *WindRack) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	rootNode.Children = append(rootNode.Children, draw.ImageWithAnimation(b.MegamanSprites.Image, b.MegamanSprites.SlashAnimation, int(e.BehaviorState.ElapsedTime)))

	swordNode := &draw.OptionsNode{Layer: 6}
	rootNode.Children = append(rootNode.Children, swordNode)
	swordNode.Children = append(swordNode.Children, draw.ImageWithAnimation(b.WindRackSprites.Image, b.WindRackSprites.Animations[0], int(e.BehaviorState.ElapsedTime)))

	return rootNode
}
