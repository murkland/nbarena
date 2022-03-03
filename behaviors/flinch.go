package behaviors

import (
	"github.com/yumland/yumbattle/bundle"
	"github.com/yumland/yumbattle/draw"
	"github.com/yumland/yumbattle/state"
)

type Flinch struct {
}

func (eb *Flinch) Clone() state.EntityBehavior {
	return &Flinch{}
}

func (eb *Flinch) Step(e *state.Entity, sh *state.StepHandle) {
	if e.BehaviorElapsedTime() == 24 {
		e.SetBehavior(&Idle{})
	}
}

func (eb *Flinch) Interrupts(e *state.Entity) state.EntityBehaviorInterrupts {
	return state.EntityBehaviorInterrupts{}
}

func (eb *Flinch) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	return draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.FlinchEndAnimation.Frames[int(e.BehaviorElapsedTime())])
}
