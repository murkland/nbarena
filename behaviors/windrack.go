package behaviors

import (
	"image"

	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type WindRack struct {
	Damage state.Damage
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
		s.AttachDecoration(&state.Decoration{
			Type:      bundle.DecorationTypeWindSlash,
			TilePos:   e.TilePos,
			Offset:    image.Point{0, -16},
			IsFlipped: e.IsFlipped,
		})
	} else if e.BehaviorState.ElapsedTime == 9 {
		x, y := e.TilePos.XY()
		dx, _ := e.Facing().XY()

		for _, pos := range []state.TilePos{
			state.TilePosXY(x+dx, y),
			state.TilePosXY(x+dx, y+1),
			state.TilePosXY(x+dx, y-1),
		} {
			var h state.Hit
			h.ForcedMovement = state.ForcedMovement{Type: state.ForcedMovementTypeBigDrag, Direction: e.Facing()}
			h.Element = state.ElementWind
			h.CanCounter = true
			h.Flinch = true
			h.RemovesFullSynchro = true
			h.AddDamage(eb.Damage)
			s.ApplyHit(e, pos, h)
		}

		for i := 1; i <= 3; i++ {
			s.AttachEntity(&state.Entity{
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
					Behavior: &Gust{e.ID(), GustStyleNone, false},
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
