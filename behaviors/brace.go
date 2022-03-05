package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type Brace struct {
}

func (eb *Brace) Clone() state.EntityBehavior {
	return &Brace{}
}

func (eb *Brace) Step(e *state.Entity, s *state.State) {
	if e.BehaviorElapsedTime() == 4 {
		e.SetBehavior(&Idle{})
	}
}

func (eb *Brace) Interrupts(e *state.Entity) state.EntityBehaviorInterrupts {
	return state.EntityBehaviorInterrupts{}
}

func (eb *Brace) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	return draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.BraceAnimation.Frames[int(e.BehaviorElapsedTime())])
}