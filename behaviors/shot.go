package behaviors

import (
	"image"
	"math/rand"

	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/state/query"
)

type Shot struct {
	Owner                   state.EntityID
	Damage                  state.Damage
	Hit                     state.Hit
	ExplosionDecorationType bundle.DecorationType
}

func (eb *Shot) Clone() state.EntityBehavior {
	return &Shot{
		eb.Owner,
		eb.Damage,
		eb.Hit,
		eb.ExplosionDecorationType,
	}
}

func (eb *Shot) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *Shot) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	return nil
}

func (eb *Shot) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime%2 == 1 {
		x, y := e.TilePos.XY()
		x += query.DXForward(e.IsFlipped)
		if !e.MoveDirectly(state.TilePosXY(x, y)) {
			e.IsPendingDestruction = true
			return
		}
	}

	h := eb.Hit
	h.AddDamage(eb.Damage)
	if s.ApplyHit(s.Entities[eb.Owner], e.TilePos, h) {
		if eb.ExplosionDecorationType != bundle.DecorationTypeNone {
			rand := rand.New(s.RandSource)

			xOff := rand.Intn(state.TileRenderedWidth / 4)
			yOff := -rand.Intn(state.TileRenderedHeight)

			s.AttachDecoration(&state.Decoration{
				Type:      eb.ExplosionDecorationType,
				TilePos:   e.TilePos,
				Offset:    image.Point{xOff, yOff},
				IsFlipped: e.IsFlipped,
			})
		}

		e.IsPendingDestruction = true
		return
	}
}

func (eb *Shot) Cleanup(e *state.Entity, s *state.State) {
}

func MakeShotEntity(owner *state.Entity, pos state.TilePos, shot *Shot) *state.Entity {
	shot.Owner = owner.ID()

	return &state.Entity{
		TilePos: pos,

		IsFlipped:            owner.IsFlipped,
		IsAlliedWithAnswerer: owner.IsAlliedWithAnswerer,

		Traits: state.EntityTraits{
			CanStepOnHoleLikeTiles: true,
			IgnoresTileEffects:     true,
			CannotFlinch:           true,
			IgnoresTileOwnership:   true,
			CannotSlide:            true,
			Intangible:             true,
		},

		BehaviorState: state.EntityBehaviorState{
			Behavior: shot,
		},
	}
}
