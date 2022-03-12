package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

// This is a questionable implementation of drag: I'm pretty sure drag and slide move entities via completely different mechanisms, but nbarena merges them together.
// This probably results in some subtly incorrect behavior.

type Dragged struct {
	PostDragParalyzeTime state.Ticks
	Type                 state.DragType

	dragComplete         bool
	dragCompleteDuration state.Ticks
}

func (eb *Dragged) Flip() {
}

func (eb *Dragged) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *Dragged) Clone() state.EntityBehavior {
	return &Dragged{
		eb.PostDragParalyzeTime, eb.Type,
		eb.dragComplete, eb.dragCompleteDuration,
	}
}

func (eb *Dragged) Step(e *state.Entity, s *state.State) {
	if eb.dragComplete {
		eb.dragCompleteDuration++
		if eb.dragCompleteDuration == 24-1 {
			if eb.PostDragParalyzeTime > 0 {
				e.NextBehavior = &Paralyzed{Duration: eb.PostDragParalyzeTime}
			} else {
				e.NextBehavior = &Idle{}
			}
			return
		}
		return
	}

	// TODO: support big drags
	if e.SlideState.Direction == state.DirectionNone {
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
