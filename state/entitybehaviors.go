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
	Step(e *Entity, sh *StepHandle)
	Interrupts(e *Entity) EntityBehaviorInterrupts
}

type IdleEntityBehavior struct {
}

func (eb *IdleEntityBehavior) Clone() EntityBehavior {
	return &IdleEntityBehavior{}
}

func (eb *IdleEntityBehavior) Step(e *Entity, sh *StepHandle) {
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

func (eb *MoveEntityBehavior) Step(e *Entity, sh *StepHandle) {
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
	isJammed     bool
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
		eb.isJammed,
		eb.cooldownTime,
	}
}

// Buster cooldown time:
var busterCooldownDurations = [][]Ticks{
	// d = 1, 2, 3, 4, 5, 6
	{5, 9, 13, 17, 21, 25}, // Lv1
	{4, 8, 11, 15, 18, 21}, // Lv2
	{4, 7, 10, 13, 16, 18}, // Lv3
	{3, 5, 7, 9, 11, 13},   // Lv4
	{3, 4, 5, 6, 7, 8},     // Lv5
}

func (eb *BusterEntityBehavior) Step(e *Entity, sh *StepHandle) {
	realElapsedTime := eb.realElapsedTime(e)

	if realElapsedTime == 5+eb.cooldownTime {
		e.SetBehavior(&IdleEntityBehavior{})
	}

	if realElapsedTime == 1 {
		_, d := findNearestEntity(sh.state, e.id, e.tilePos, e.isAlliedWithAnswerer, e.isFlipped, horizontalDistance)
		eb.cooldownTime = busterCooldownDurations[0][d]

		x, y := e.tilePos.XY()
		if e.isFlipped {
			x--
		} else {
			x++
		}
		sh.SpawnEntity(&Entity{
			behavior: &busterShotEntityBehavior{
				isPowerShot: eb.IsPowerShot,
			},

			tilePos:                TilePosXY(x, y),
			hp:                     0,
			canStepOnHoleLikeTiles: true,
			ignoresTileEffects:     true,
			cannotFlinch:           true,
		})
	}
}

func (eb *BusterEntityBehavior) Appearance(e *Entity, b *bundle.Bundle) draw.Node {
	realElapsedTime := eb.realElapsedTime(e)

	if realElapsedTime < 0 {
		frame := b.MegamanSprites.IdleAnimation.Frames[0]
		return draw.ImageWithFrame(b.MegamanSprites.Image, frame)
	}

	rootNode := &draw.OptionsNode{}

	if realElapsedTime < 5 {
		megamanBusterAnimTime := int(realElapsedTime)
		if megamanBusterAnimTime >= len(b.MegamanSprites.BusterAnimation.Frames) {
			megamanBusterAnimTime = len(b.MegamanSprites.BusterAnimation.Frames) - 1
		}
		rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.BusterAnimation.Frames[megamanBusterAnimTime]))
	} else {
		rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.BusterEndAnimation.Frames[0]))
	}

	busterFrames := b.BusterSprites.Animations[0]
	busterAnimTime := int(realElapsedTime)
	if busterAnimTime >= len(busterFrames.Frames) {
		busterAnimTime = len(busterFrames.Frames) - 1
	}
	busterFrame := busterFrames.Frames[busterAnimTime]
	rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.BusterSprites.Image, busterFrame))

	if !eb.isJammed {
		muzzleFlashAnimTime := int(realElapsedTime) - 1
		if muzzleFlashAnimTime > 0 && muzzleFlashAnimTime < len(b.MuzzleFlashSprites.Animations[0].Frames) {
			muzzleFlashNode := &draw.OptionsNode{}
			muzzleFlashFrame := b.MuzzleFlashSprites.Animations[0].Frames[muzzleFlashAnimTime]
			// TODO: Figure out how to draw the muzzle flash.
			muzzleFlashNode.Children = append(muzzleFlashNode.Children, draw.ImageWithFrame(b.MuzzleFlashSprites.Image, muzzleFlashFrame))
			rootNode.Children = append(rootNode.Children, muzzleFlashNode)
		}
	}

	return rootNode
}

func (eb *BusterEntityBehavior) Interrupts(e *Entity) EntityBehaviorInterrupts {
	realElapsedTime := eb.realElapsedTime(e)
	return EntityBehaviorInterrupts{
		OnMove: realElapsedTime >= 5,
	}
}

type busterShotEntityBehavior struct {
	isPowerShot bool
}

func (eb *busterShotEntityBehavior) Clone() EntityBehavior {
	return &busterShotEntityBehavior{
		eb.isPowerShot,
	}
}

func (eb *busterShotEntityBehavior) Appearance(e *Entity, b *bundle.Bundle) draw.Node {
	return nil
}

func (eb *busterShotEntityBehavior) Interrupts(e *Entity) EntityBehaviorInterrupts {
	return EntityBehaviorInterrupts{}
}

func (eb *busterShotEntityBehavior) Step(e *Entity, sh *StepHandle) {
}

type distanceMetric func(src TilePos, dest TilePos) int

func dxForward(isFlipped bool) int {
	if isFlipped {
		return -1
	}
	return 1
}

func isInFrontOf(x int, targetX int, isFlipped bool) bool {
	if isFlipped {
		return targetX < x
	}
	return targetX > x
}

func horizontalDistance(src TilePos, dest TilePos) int {
	x1, _ := src.XY()
	x2, _ := dest.XY()
	if x1 > x2 {
		return x1 - x2
	}
	return x2 - x1
}

func findNearestEntity(s *State, myEntityID int, pos TilePos, isAlliedWithAnswerer bool, isFlipped bool, distance distanceMetric) (int, int) {
	x, _ := pos.XY()

	bestDist := tileCols

	var targetID int
	for candID, cand := range s.entities {
		if candID == myEntityID || cand.isAlliedWithAnswerer == isAlliedWithAnswerer {
			continue
		}

		candX, _ := cand.futureTilePos.XY()

		if !isInFrontOf(x, candX, isFlipped) {
			continue
		}

		if d := distance(pos, cand.futureTilePos); d >= 0 && d < bestDist {
			targetID = candID
			bestDist = d
		}
	}

	return targetID, bestDist
}
