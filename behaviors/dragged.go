package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type Dragged struct {
	Direction            state.Direction
	IsBig                bool
	PostDragParalyzeTime state.Ticks

	dragComplete         bool
	dragCompleteDuration state.Ticks
}

func (eb *Dragged) Flip() {
	eb.Direction = eb.Direction.FlipH()
}

func (eb *Dragged) Clone() state.EntityBehavior {
	return &Dragged{
		eb.Direction, eb.IsBig, eb.PostDragParalyzeTime,
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

	if e.BehaviorElapsedTime()%4 == 0 {
		if !e.StartMove(getNextDragEndTilePos(e.TilePos, eb.Direction), s.Field) {
			eb.dragComplete = true
		}
	} else if e.BehaviorElapsedTime()%4 == 2 {
		e.FinishMove()
		if !eb.IsBig {
			eb.dragComplete = true
		}
	}
}

func (eb *Dragged) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	var childNode draw.Node
	if eb.PostDragParalyzeTime > 0 {
		childNode = (&Paralyzed{Duration: 0}).Appearance(e, b)
	} else {
		childNode = draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.FlinchAnimation.Frames[eb.dragCompleteDuration])
	}

	rootNode := &draw.OptionsNode{}
	if !eb.dragComplete {
		// TODO: Render this correctly on the other side of the screen.
		dx, dy := eb.Direction.XY()
		offset := (int(e.BehaviorElapsedTime())+2+4)%4 - 2
		dx *= offset
		dy *= offset

		rootNode.Opts.GeoM.Translate(float64(dx*state.TileRenderedWidth/4), float64(dy*(state.TileRenderedHeight/4)))
	}
	rootNode.Children = append(rootNode.Children, childNode)
	return rootNode
}

func getNextDragEndTilePos(tilePos state.TilePos, direction state.Direction) state.TilePos {
	x, y := tilePos.XY()
	dx, dy := direction.XY()
	return state.TilePosXY(x+dx, y+dy)
}
