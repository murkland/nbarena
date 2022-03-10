package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type Frozen struct {
	Duration state.Ticks
}

func (eb *Frozen) Flip() {
}

func (eb *Frozen) Clone() state.EntityBehavior {
	return &Frozen{}
}

func (eb *Frozen) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == eb.Duration {
		e.SetBehavior(&Idle{}, s)
	}
}

func (eb *Frozen) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	rootNode.Opts.ColorM.Translate(float64(0xa5)/float64(0xff), float64(0xa5)/float64(0xff), float64(0xff)/float64(0xff), 0.0)
	rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.StuckAnimation.Frames[int(e.BehaviorState.ElapsedTime)%len(b.MegamanSprites.StuckAnimation.Frames)]))
	return rootNode
}
