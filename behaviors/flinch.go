package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type Flinch struct {
}

func (eb *Flinch) Flip() {
}

func (eb *Flinch) Clone() state.EntityBehavior {
	return &Flinch{}
}

func (eb *Flinch) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *Flinch) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == 24-1 {
		e.NextBehavior = &Idle{}
	}
}

func (eb *Flinch) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	return draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.FlinchAnimation.Frames[int(e.BehaviorState.ElapsedTime)])
}
