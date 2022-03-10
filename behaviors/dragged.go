package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type Dragged struct {
	PostDragParalyzeTime state.Ticks

	dragComplete         bool
	dragCompleteDuration state.Ticks
}

func (eb *Dragged) Flip() {
}

func (eb *Dragged) Clone() state.EntityBehavior {
	return &Dragged{
		eb.PostDragParalyzeTime,
		eb.dragComplete, eb.dragCompleteDuration,
	}
}

func (eb *Dragged) Step(e *state.Entity, s *state.State) {
	if eb.dragComplete {
		eb.dragCompleteDuration++
		if eb.dragCompleteDuration == 24 {
			if eb.PostDragParalyzeTime > 0 {
				e.SetBehavior(&Paralyzed{Duration: eb.PostDragParalyzeTime}, s)
			} else {
				e.SetBehavior(&Idle{}, s)
			}
			return
		}
		return
	}

	if e.SlideState.Slide.Direction == state.DirectionNone {
		eb.dragComplete = true
	}
}

func (eb *Dragged) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	var childNode draw.Node
	if eb.PostDragParalyzeTime > 0 {
		childNode = (&Paralyzed{Duration: 0}).Appearance(e, b)
	} else {
		childNode = draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.FlinchAnimation.Frames[eb.dragCompleteDuration])
	}

	rootNode.Children = append(rootNode.Children, childNode)
	return rootNode
}
