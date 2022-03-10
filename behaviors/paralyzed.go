package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type Paralyzed struct {
	Duration state.Ticks
}

func (eb *Paralyzed) Flip() {
}

func (eb *Paralyzed) Clone() state.EntityBehavior {
	return &Paralyzed{}
}

func (eb *Paralyzed) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == eb.Duration {
		e.SetBehavior(&Idle{}, s)
	}
}

func (eb *Paralyzed) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	if (e.ElapsedTime()/2)%2 == 1 {
		rootNode.Opts.ColorM.Translate(1.0, 1.0, 0.0, 0.0)
	}
	rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.StuckAnimation.Frames[int(e.BehaviorState.ElapsedTime)%len(b.MegamanSprites.StuckAnimation.Frames)]))
	return rootNode
}
