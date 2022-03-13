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

func (eb *Paralyzed) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *Paralyzed) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == eb.Duration-1 {
		e.NextBehavior = &Idle{}
	}
}

func (eb *Paralyzed) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	if (e.ElapsedTime()/2)%2 == 1 {
		rootNode.Opts.ColorM.Translate(1.0, 1.0, 0.0, 0.0)
	}
	rootNode.Children = append(rootNode.Children, draw.ImageWithAnimation(b.MegamanSprites.Image, b.MegamanSprites.StuckAnimation, int(e.BehaviorState.ElapsedTime)))
	return rootNode
}
