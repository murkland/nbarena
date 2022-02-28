package state

import (
	"github.com/yumland/clone"
	"github.com/yumland/pngsheet"
	"github.com/yumland/yumbattle/bundle"
	"github.com/yumland/yumbattle/draw"
)

type EntityBehaviorInterrupts struct {
	OnMove   bool
	OnCharge bool
}

type EntityBehavior interface {
	clone.Cloner[EntityBehavior]
	Appearance(e *Entity, b *bundle.Bundle) draw.Node
	Step(e *Entity)
	Interrupts(e *Entity) EntityBehaviorInterrupts
}

type IdleEntityBehavior struct {
}

func (eb *IdleEntityBehavior) Clone() EntityBehavior {
	return &IdleEntityBehavior{}
}

func (eb *IdleEntityBehavior) Step(e *Entity) {
}

func (eb *IdleEntityBehavior) Interrupts(e *Entity) EntityBehaviorInterrupts {
	return EntityBehaviorInterrupts{
		OnMove:   true,
		OnCharge: true,
	}
}

func (eb *IdleEntityBehavior) Appearance(e *Entity, b *bundle.Bundle) draw.Node {
	frame := b.MegamanSprites.IdleAnimation.Frames[0]
	return draw.ImageWithFrame(b.MegamanSprites.Image, frame)
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
		frame = b.MegamanSprites.MoveStartAnimation.Frames[e.behaviorElapsedTime]
	} else if e.behaviorElapsedTime < 6 {
		frame = b.MegamanSprites.MoveEndAnimation.Frames[e.behaviorElapsedTime-3]
	} else {
		frame = b.MegamanSprites.MoveEndAnimation.Frames[len(b.MegamanSprites.MoveEndAnimation.Frames)-1]
	}
	return draw.ImageWithFrame(b.MegamanSprites.Image, frame)
}

func (eb *MoveEntityBehavior) Interrupts(e *Entity) EntityBehaviorInterrupts {
	return EntityBehaviorInterrupts{}
}

type BusterEntityBehavior struct {
	IsPowerShot  bool
	cooldownTime Ticks
}

func (eb *BusterEntityBehavior) realElapsedTime(e *Entity) Ticks {
	t := e.behaviorElapsedTime
	if eb.IsPowerShot {
		t -= 5
	}
	return t
}

func (eb *BusterEntityBehavior) Clone() EntityBehavior {
	return &BusterEntityBehavior{
		eb.IsPowerShot,
		eb.cooldownTime,
	}
}

func (eb *BusterEntityBehavior) Step(e *Entity) {
	realElapsedTime := eb.realElapsedTime(e)
	eb.cooldownTime = 100

	if realElapsedTime == 5+eb.cooldownTime {
		e.SetBehavior(&IdleEntityBehavior{})
	}
}

func (eb *BusterEntityBehavior) Appearance(e *Entity, b *bundle.Bundle) draw.Node {
	realElapsedTime := eb.realElapsedTime(e)

	if realElapsedTime < 0 {
		frame := b.MegamanSprites.IdleAnimation.Frames[0]
		return draw.ImageWithFrame(b.MegamanSprites.Image, frame)
	}

	rootNode := &draw.OptionsNode{}

	megamanBusterAnimTime := int(realElapsedTime)
	if megamanBusterAnimTime >= len(b.MegamanSprites.BusterAnimation.Frames) {
		megamanBusterAnimTime = len(b.MegamanSprites.BusterAnimation.Frames) - 1
	}
	rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.BusterAnimation.Frames[megamanBusterAnimTime]))

	busterFrames := b.BusterSprites.Animations[0]
	busterAnimTime := int(realElapsedTime)
	if busterAnimTime >= len(busterFrames.Frames) {
		busterAnimTime = len(busterFrames.Frames) - 1
	}
	busterFrame := busterFrames.Frames[busterAnimTime]
	rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.BusterSprites.Image, busterFrame))

	muzzleFlashAnimTime := int(realElapsedTime) - 1
	if muzzleFlashAnimTime > 0 && muzzleFlashAnimTime < len(b.MuzzleFlashSprites.Animations[0].Frames) {
		muzzleFlashNode := &draw.OptionsNode{}
		muzzleFlashNode.Children = append(muzzleFlashNode.Children, draw.ImageWithFrame(b.MuzzleFlashSprites.Image, b.MuzzleFlashSprites.Animations[0].Frames[muzzleFlashAnimTime]))
		rootNode.Children = append(rootNode.Children, muzzleFlashNode)
	}

	return rootNode
}

func (eb *BusterEntityBehavior) Interrupts(e *Entity) EntityBehaviorInterrupts {
	realElapsedTime := eb.realElapsedTime(e)
	return EntityBehaviorInterrupts{
		OnMove: realElapsedTime >= 5,
	}
}
