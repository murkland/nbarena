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
	frame := b.MegamanSprites.IdleAnimation.Frames[0]
	return draw.ImageWithFrame(b.MegamanSprites.Image, frame)
}
