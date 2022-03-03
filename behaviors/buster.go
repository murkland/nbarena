package behaviors

import (
	"github.com/yumland/yumbattle/bundle"
	"github.com/yumland/yumbattle/draw"
	"github.com/yumland/yumbattle/state"
)

type Buster struct {
	IsPowerShot  bool
	isJammed     bool
	AnimIndex    int
	cooldownTime state.Ticks
}

func (eb *Buster) realElapsedTime(e *state.Entity) state.Ticks {
	t := e.BehaviorElapsedTime()
	if eb.IsPowerShot {
		t -= 5
	}
	return t
}

func (eb *Buster) Clone() state.EntityBehavior {
	return &Buster{
		eb.IsPowerShot,
		eb.isJammed,
		eb.AnimIndex,
		eb.cooldownTime,
	}
}

// Buster cooldown time:
var busterCooldownDurations = [][]state.Ticks{
	// d = 1, 2, 3, 4, 5, 6
	{5, 9, 13, 17, 21, 25}, // Lv1
	{4, 8, 11, 15, 18, 21}, // Lv2
	{4, 7, 10, 13, 16, 18}, // Lv3
	{3, 5, 7, 9, 11, 13},   // Lv4
	{3, 4, 5, 6, 7, 8},     // Lv5
}

func (eb *Buster) Step(e *state.Entity, sh *state.StepHandle) {
	realElapsedTime := eb.realElapsedTime(e)

	if realElapsedTime == 5+eb.cooldownTime {
		e.SetBehavior(&Idle{})
	}

	if realElapsedTime == 1 {
		_, d := findNearestEntity(sh.State, e.ID(), e.TilePos, e.IsAlliedWithAnswerer, e.IsFlipped, horizontalDistance)
		eb.cooldownTime = busterCooldownDurations[0][d]

		x, y := e.TilePos.XY()
		if e.IsFlipped {
			x--
		} else {
			x++
		}

		e := &state.Entity{
			TilePos:       state.TilePosXY(x, y),
			FutureTilePos: state.TilePosXY(x, y),

			IsFlipped:            e.IsFlipped,
			IsAlliedWithAnswerer: e.IsAlliedWithAnswerer,

			Traits: state.EntityTraits{
				CanStepOnHoleLikeTiles: true,
				IgnoresTileEffects:     true,
				CannotFlinch:           true,
				IgnoresTileOwnership:   true,
			},
		}
		e.SetBehavior(&busterShot{
			isPowerShot: eb.IsPowerShot,
		})
		sh.SpawnEntity(e)
	}
}

func (eb *Buster) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
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

	busterFrames := b.BusterSprites.Animations[eb.AnimIndex]
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

func (eb *Buster) Interrupts(e *state.Entity) state.EntityBehaviorInterrupts {
	realElapsedTime := eb.realElapsedTime(e)
	return state.EntityBehaviorInterrupts{
		OnMove: realElapsedTime >= 5,
	}
}

type busterShot struct {
	isPowerShot bool
}

func (eb *busterShot) Clone() state.EntityBehavior {
	return &busterShot{
		eb.isPowerShot,
	}
}

func (eb *busterShot) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	return nil
}

func (eb *busterShot) Interrupts(e *state.Entity) state.EntityBehaviorInterrupts {
	return state.EntityBehaviorInterrupts{}
}

func (eb *busterShot) Step(e *state.Entity, sh *state.StepHandle) {
	if e.BehaviorElapsedTime()%2 == 1 {
		x, y := e.TilePos.XY()
		x += dxForward(e.IsFlipped)
		if !e.StartMove(state.TilePosXY(x, y), &sh.State.Field) {
			sh.RemoveEntity(e.ID())
			return
		}
	} else {
		e.FinishMove()

		for _, e2 := range entitiesAt(sh.State, e.TilePos) {
			if e2.IsAlliedWithAnswerer == e.IsAlliedWithAnswerer {
				continue
			}

			damage := 10
			if eb.isPowerShot {
				damage *= 10
			}
			var h state.Hit
			h.AddDamage(state.Damage{Base: damage})
			e2.AddHit(h)

			sh.RemoveEntity(e.ID())
			return
		}
	}
}

type distanceMetric func(src state.TilePos, dest state.TilePos) int

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

func horizontalDistance(src state.TilePos, dest state.TilePos) int {
	x1, _ := src.XY()
	x2, _ := dest.XY()
	if x1 > x2 {
		return x1 - x2
	}
	return x2 - x1
}

func findNearestEntity(s *state.State, myEntityID int, pos state.TilePos, isAlliedWithAnswerer bool, isFlipped bool, distance distanceMetric) (int, int) {
	x, _ := pos.XY()

	bestDist := state.TileCols

	var targetID int
	for _, cand := range s.Entities {
		if cand.ID() == myEntityID || cand.IsAlliedWithAnswerer == isAlliedWithAnswerer {
			continue
		}

		candX, _ := cand.FutureTilePos.XY()

		if !isInFrontOf(x, candX, isFlipped) {
			continue
		}

		if d := distance(pos, cand.FutureTilePos); d >= 0 && d < bestDist {
			targetID = cand.ID()
			bestDist = d
		}
	}

	return targetID, bestDist
}

func entitiesAt(s *state.State, pos state.TilePos) []*state.Entity {
	var entities []*state.Entity
	for _, e := range s.Entities {
		if e.TilePos != pos {
			continue
		}
		entities = append(entities, e)
	}
	return entities
}
