package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yumland/clone"
	"github.com/yumland/pngsheet"
	"github.com/yumland/yumbattle/bundle"
	"github.com/yumland/yumbattle/draw"
)

type EntityBehavior interface {
	clone.Cloner[EntityBehavior]
	Appearance(e *Entity, b *bundle.Bundle) draw.Node
	Step(e *Entity)
}

type IdleEntityBehavior struct {
}

func (eb *IdleEntityBehavior) Clone() EntityBehavior {
	return &IdleEntityBehavior{}
}

func (eb *IdleEntityBehavior) Step(e *Entity) {
}

func (eb *IdleEntityBehavior) Appearance(e *Entity, b *bundle.Bundle) draw.Node {
	frame := b.Megaman.IdleAnimation.Frames[0]
	return draw.ImageWithOrigin(b.Megaman.Sprites.SubImage(frame.Rect).(*ebiten.Image), frame.Origin)
}

const moveEndlagTicks = 7

type MoveEntityBehavior struct {
}

func (eb *MoveEntityBehavior) Clone() EntityBehavior {
	return &MoveEntityBehavior{}
}

func (eb *MoveEntityBehavior) Step(e *Entity) {
	if e.behaviorElapsedTime == 3 {
		e.tilePos = e.futureTilePos
	}
	if e.behaviorElapsedTime == 6+moveEndlagTicks {
		e.SetBehavior(&IdleEntityBehavior{})
	}
}

func (eb *MoveEntityBehavior) Appearance(e *Entity, b *bundle.Bundle) draw.Node {
	var frame *pngsheet.Frame
	if e.behaviorElapsedTime < 3 {
		frame = b.Megaman.MoveStartAnimation.Frames[e.behaviorElapsedTime]
	} else if e.behaviorElapsedTime < 6 {
		frame = b.Megaman.MoveEndAnimation.Frames[e.behaviorElapsedTime-3]
	} else {
		frame = b.Megaman.MoveEndAnimation.Frames[len(b.Megaman.MoveEndAnimation.Frames)-1]
	}
	return draw.ImageWithOrigin(b.Megaman.Sprites.SubImage(frame.Rect).(*ebiten.Image), frame.Origin)
}
