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
	IsBig                bool

	dragCompleteDuration state.Ticks
}

func (eb *Dragged) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *Dragged) Clone() state.EntityBehavior {
	return &Dragged{
		eb.PostDragParalyzeTime, eb.IsBig,
		eb.dragCompleteDuration,
	}
}

func (eb *Dragged) Step(e *state.Entity, s *state.State) {
	// The implementation of drag hijacks normal slide processing: it relies on a precise interaction between how sliding is discarded in step.resolveOne if drag is active.
	//
	// We use the slide mechanism here because it's the only way nbarena is able to render players being partially on a tile.
	if e.SlideState.Direction == state.DirectionNone {
		eb.dragCompleteDuration++
		if eb.dragCompleteDuration == 24-1 {
			if eb.PostDragParalyzeTime > 0 {
				e.NextBehavior = &Paralyzed{Duration: eb.PostDragParalyzeTime}
			} else {
				e.NextBehavior = &Idle{}
			}
		}
		return
	}

	if !eb.IsBig && e.BehaviorState.ElapsedTime == 4 {
		e.SlideState = state.SlideState{}
	} else if e.BehaviorState.ElapsedTime%4 == 0 {
		x, y := e.TilePos.XY()
		dx, dy := e.SlideState.Direction.XY()

		if !e.StartMove(state.TilePosXY(x+dx, y+dy), s) {
			e.SlideState = state.SlideState{}
		}
	} else if e.BehaviorState.ElapsedTime%4 == 2 {
		e.FinishMove(s)
	}
}

func (eb *Dragged) Cleanup(e *state.Entity, s *state.State) {
}

func (eb *Dragged) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	var childNode draw.Node
	if eb.PostDragParalyzeTime > 0 {
		childNode = (&Paralyzed{Duration: 0}).Appearance(e, b)
	} else {
		// TODO: The timing is correct but the animation is wrong.
		childNode = draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.FlinchAnimation.Frames[eb.dragCompleteDuration])
	}

	rootNode.Children = append(rootNode.Children, childNode)
	return rootNode
}
