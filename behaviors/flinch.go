package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type Flinch struct {
}

func (eb *Flinch) Clone() state.EntityBehavior {
	return &Flinch{}
}

func (eb *Flinch) Step(e *state.Entity, s *state.State) {
	if e.BehaviorElapsedTime() == 24 {
		e.SetBehavior(&Idle{})
	}
}

func (eb *Flinch) Interrupts(e *state.Entity) state.EntityBehaviorInterrupts {
	return state.EntityBehaviorInterrupts{}
}

func (eb *Flinch) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	return draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.FlinchAnimation.Frames[int(e.BehaviorElapsedTime())])
}