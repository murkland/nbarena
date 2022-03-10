package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type Bubbled struct {
	Duration state.Ticks
}

func (eb *Bubbled) Flip() {
}

func (eb *Bubbled) Clone() state.EntityBehavior {
	return &Bubbled{}
}

func (eb *Bubbled) Step(e *state.Entity, s *state.State) {
	if e.BehaviorElapsedTime() == eb.Duration {
		e.SetBehavior(&Idle{}, s)
	}
}

func (eb *Bubbled) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	// TODO: Renber bubble.
	return draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.StuckAnimation.Frames[int(e.BehaviorElapsedTime())%len(b.MegamanSprites.StuckAnimation.Frames)])
}
