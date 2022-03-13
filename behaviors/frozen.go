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

func (eb *Frozen) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *Frozen) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == eb.Duration-1 {
		e.NextBehavior = &Idle{}
	}
}

func (eb *Frozen) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	rootNode.Opts.ColorM.Translate(float64(0xa5)/float64(0xff), float64(0xa5)/float64(0xff), float64(0xff)/float64(0xff), 0.0)
	rootNode.Children = append(rootNode.Children, draw.ImageWithAnimation(b.MegamanSprites.Image, b.MegamanSprites.StuckAnimation, int(e.BehaviorState.ElapsedTime)))
	return rootNode
}
