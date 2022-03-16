package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type Bubbled struct {
	Duration state.Ticks
}

func (eb *Bubbled) Clone() state.EntityBehavior {
	return &Bubbled{eb.Duration}
}

func (eb *Bubbled) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == eb.Duration-1 {
		e.NextBehavior = &Idle{}
	}
}

func (eb *Bubbled) Cleanup(e *state.Entity, s *state.State) {
}

func (eb *Bubbled) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	// TODO: Renber bubble.
	return draw.ImageWithAnimation(b.MegamanSprites.Image, b.MegamanSprites.StuckAnimation, int(e.BehaviorState.ElapsedTime))
}
