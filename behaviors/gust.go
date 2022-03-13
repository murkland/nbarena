package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/state/query"
)

type Gust struct {
	Direction state.Direction
}

func (eb *Gust) Flip() {
}

func (eb *Gust) Clone() state.EntityBehavior {
	return &Gust{
		eb.Direction,
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

	for _, target := range query.TangibleEntitiesAt(s, e.TilePos) {
		if target.IsAlliedWithAnswerer == e.IsAlliedWithAnswerer || target.FlashingTimeLeft > 0 {
			continue
		}

		var h state.Hit
		h.Traits.Element = state.ElementWind
		h.Traits.SlideDirection = eb.Direction
		target.Hit.Merge(h)
	}
}

func (eb *Gust) Cleanup(e *state.Entity, s *state.State) {
}
