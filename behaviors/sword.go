package behaviors

import (
	"github.com/yumland/yumbattle/bundle"
	"github.com/yumland/yumbattle/draw"
	"github.com/yumland/yumbattle/state"
)

type Sword struct {
	Rows      int
	Cols      int
	Damage    int
	AnimIndex int
}

func (eb *Sword) Clone() state.EntityBehavior {
	return &Sword{
		eb.Rows,
		eb.Cols,
		eb.Damage,
		eb.AnimIndex,
	}
}

func (eb *Sword) Step(e *state.Entity, sh *state.StepHandle) {
	// TODO: Everything.
	if e.BehaviorElapsedTime() == 22 {
		e.SetBehavior(&Idle{})
	}
}

func (eb *Sword) Interrupts(e *state.Entity) state.EntityBehaviorInterrupts {
	return state.EntityBehaviorInterrupts{}
}

func (eb *Sword) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	if e.BehaviorElapsedTime() < 21 {
		rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.SlashAnimation.Frames[e.BehaviorElapsedTime()]))
		rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.SwordSprites.Image, b.SwordSprites.Animations[eb.AnimIndex].Frames[e.BehaviorElapsedTime()]))
	} else {
		rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.IdleAnimation.Frames[0]))
	}
	return rootNode
}
