package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type Idle struct {
}

func (eb *Idle) Clone() state.EntityBehavior {
	return &Idle{}
}

func (eb *Idle) Step(e *state.Entity, s *state.State) {
}

func (eb *Idle) Interrupts(e *state.Entity) state.EntityBehaviorInterrupts {
	return state.EntityBehaviorInterrupts{
		WithChipUse: true,
		WithMove:    true,
		WithCharge:  true,
	}
}

func (eb *Idle) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	frames := b.MegamanSprites.IdleAnimation.Frames
	frame := frames[int(e.BehaviorElapsedTime())%len(frames)]
	return draw.ImageWithFrame(b.MegamanSprites.Image, frame)
}
