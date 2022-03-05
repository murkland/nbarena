package behaviors

import (
	"github.com/yumland/nbarena/bundle"
	"github.com/yumland/nbarena/draw"
	"github.com/yumland/nbarena/state"
)

type Cannon struct {
	Damage int
}

func (eb *Cannon) Clone() state.EntityBehavior {
	return &Cannon{
		eb.Damage,
	}
}

func (eb *Cannon) Step(e *state.Entity, s *state.State) {
	// TODO: Hitbox.
	if e.BehaviorElapsedTime() == 29 {
		e.SetBehavior(&Brace{})
	}
}

func (eb *Cannon) Interrupts(e *state.Entity) state.EntityBehaviorInterrupts {
	return state.EntityBehaviorInterrupts{}
}

func (eb *Cannon) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.CannonAnimation.Frames[e.BehaviorElapsedTime()]))

	cannonNode := &draw.OptionsNode{Layer: 9}
	cannonNode.Opts.GeoM.Translate(float64(16), float64(-24))
	rootNode.Children = append(rootNode.Children, cannonNode)
	cannonNode.Children = append(cannonNode.Children, draw.ImageWithFrame(b.CannonSprites.CannonImage, b.CannonSprites.Animation.Frames[e.BehaviorElapsedTime()]))

	return rootNode
}
