package behaviors

import (
	"github.com/yumland/pngsheet"
	"github.com/yumland/yumbattle/bundle"
	"github.com/yumland/yumbattle/draw"
	"github.com/yumland/yumbattle/state"
)

const moveEndlagTicks = 7

type Teleport struct {
}

func (eb *Teleport) Clone() state.EntityBehavior {
	return &Teleport{}
}

func (eb *Teleport) Step(e *state.Entity, sh *state.StepHandle) {
	if e.BehaviorElapsedTime() == 3 {
		e.FinishMove()
	}

	if e.BehaviorElapsedTime() == 6+moveEndlagTicks {
		e.SetBehavior(&Idle{})
	}
}

func (eb *Teleport) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	var frame *pngsheet.Frame
	if e.BehaviorElapsedTime() < 3 {
		frame = b.MegamanSprites.MoveStartAnimation.Frames[e.BehaviorElapsedTime()]
	} else if e.BehaviorElapsedTime() < 6 {
		frame = b.MegamanSprites.MoveEndAnimation.Frames[e.BehaviorElapsedTime()-3]
	} else {
		frame = b.MegamanSprites.MoveEndAnimation.Frames[len(b.MegamanSprites.MoveEndAnimation.Frames)-1]
	}
	return draw.ImageWithFrame(b.MegamanSprites.Image, frame)
}

func (eb *Teleport) Interrupts(e *state.Entity) state.EntityBehaviorInterrupts {
	return state.EntityBehaviorInterrupts{}
}
