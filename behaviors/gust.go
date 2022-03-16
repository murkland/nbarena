package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/state/query"
)

type Gust struct {
	Owner state.EntityID
}

func (eb *Gust) Clone() state.EntityBehavior {
	return &Gust{
		eb.Owner,
	}
}

func (eb *Gust) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *Gust) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	return nil
}

func (eb *Gust) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime%4 == 1 {
		x, y := e.TilePos.XY()
		x += query.DXForward(e.IsFlipped)
		if !e.MoveDirectly(state.TilePosXY(x, y)) {
			e.IsPendingDestruction = true
			return
		}
	}

	var h state.Hit
	h.Element = state.ElementWind
	h.ForcedMovement = state.ForcedMovement{Type: state.ForcedMovementTypeSlide, Direction: e.Facing()}
	s.ApplyHit(s.Entities[eb.Owner], e.TilePos, h)
}

func (eb *Gust) Cleanup(e *state.Entity, s *state.State) {
}
