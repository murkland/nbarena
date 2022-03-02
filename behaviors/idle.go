package behaviors

import (
	"github.com/yumland/yumbattle/bundle"
	"github.com/yumland/yumbattle/draw"
	"github.com/yumland/yumbattle/state"
)

type Idle struct {
}

func (eb *Idle) Clone() state.EntityBehavior {
	return &Idle{}
}

func (eb *Idle) Step(e *state.Entity, sh *state.StepHandle) {
}

func (eb *Idle) Interrupts(e *state.Entity) state.EntityBehaviorInterrupts {
	return state.EntityBehaviorInterrupts{
		OnMove:   true,
		OnCharge: true,
	}
}

func (eb *Idle) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	frames := b.MegamanSprites.IdleAnimation.Frames
	frame := frames[int(e.BehaviorElapsedTime())%len(frames)]
	return draw.ImageWithFrame(b.MegamanSprites.Image, frame)
}