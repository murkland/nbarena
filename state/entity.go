package state

import (
	"flag"
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yumland/clone"
	"github.com/yumland/yumbattle/bundle"
	"github.com/yumland/yumbattle/draw"
)

var (
	debugDrawEntityMarker = flag.Bool("debug_draw_entity_markers", false, "draw entity markers")
)

type EntityTraits struct {
	CanStepOnHoleLikeTiles bool
	IgnoresTileEffects     bool
	CannotFlinch           bool
	FatalHitLeaves1HP      bool
	IgnoresTileOwnership   bool
}

type Entity struct {
	id int

	elapsedTime Ticks

	behaviorElapsedTime Ticks
	behavior            EntityBehavior

	TilePos       TilePos
	FutureTilePos TilePos

	IsAlliedWithAnswerer bool

	IsFlipped bool

	isDeleted bool

	HP        int
	DisplayHP int

	Traits EntityTraits

	ChargingElapsedTime Ticks
	PowerShotChargeTime Ticks

	ParalyzedTimeLeft   Ticks
	ConfusedTimeLeft    Ticks
	BlindedTimeLeft     Ticks
	ImmobilizedTimeLeft Ticks
	FlashingTimeLeft    Ticks
	InvincibleTimeLeft  Ticks
	FrozenTimeLeft      Ticks
	BubbledTimeLeft     Ticks

	IsAngry        bool
	IsBeingDragged bool
	IsSliding      bool

	CurrentHit Hit
}

func (e *Entity) ID() int {
	return e.id
}

func (e *Entity) Interrupts() EntityBehaviorInterrupts {
	return e.behavior.Interrupts(e)
}

func (e *Entity) Clone() *Entity {
	return &Entity{
		e.id,
		e.elapsedTime,
		e.behaviorElapsedTime, e.behavior.Clone(),
		e.TilePos, e.FutureTilePos,
		e.IsAlliedWithAnswerer,
		e.IsFlipped,
		e.isDeleted,
		e.HP, e.DisplayHP,
		e.Traits,
		e.ChargingElapsedTime, e.PowerShotChargeTime,
		e.ParalyzedTimeLeft, e.ConfusedTimeLeft, e.BlindedTimeLeft, e.ImmobilizedTimeLeft, e.FlashingTimeLeft, e.InvincibleTimeLeft, e.FrozenTimeLeft, e.BubbledTimeLeft,
		e.IsAngry, e.IsBeingDragged, e.IsSliding,
		e.CurrentHit,
	}
}

func (e *Entity) SetBehavior(behavior EntityBehavior) {
	e.behaviorElapsedTime = 0
	e.behavior = behavior
}

func (e *Entity) BehaviorElapsedTime() Ticks {
	return e.behaviorElapsedTime
}

func (e *Entity) StartMove(tilePos TilePos, field *Field) bool {
	x, y := tilePos.XY()
	if x < 0 || x >= TileCols || y < 0 || y >= TileRows {
		return false
	}

	tile := &field.Tiles[tilePos]
	if tilePos == e.TilePos ||
		(!e.Traits.IgnoresTileOwnership && e.IsAlliedWithAnswerer != tile.IsAlliedWithAnswerer) ||
		!tile.CanEnter(e) {
		return false
	}

	e.FutureTilePos = tilePos
	return true
}

func (e *Entity) FinishMove() {
	e.TilePos = e.FutureTilePos
}

var debugEntityMarkerImage *ebiten.Image
var debugEntityMarkerImageOnce sync.Once

func (e *Entity) Appearance(b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	x, y := e.TilePos.XY()
	rootNode.Opts.GeoM.Translate(float64((x-1)*tileRenderedWidth+tileRenderedWidth/2), float64((y-1)*tileRenderedHeight+tileRenderedHeight/2))

	characterNode := &draw.OptionsNode{}
	if e.IsFlipped {
		characterNode.Opts.GeoM.Scale(-1, 1)
	}
	if e.FrozenTimeLeft > 0 {
		// TODO: Render ice.
		characterNode.Opts.ColorM.Translate(float64(0xa5)/float64(0xff), float64(0xa5)/float64(0xff), float64(0xff)/float64(0xff), 0.0)
	}
	if e.ParalyzedTimeLeft > 0 && (e.elapsedTime/2)%2 == 1 {
		characterNode.Opts.ColorM.Translate(1.0, 1.0, 0.0, 0.0)
	}
	if e.FlashingTimeLeft > 0 && (e.elapsedTime/2)%2 == 0 {
		characterNode.Opts.ColorM.Translate(0.0, 0.0, 0.0, -1.0)
	}
	characterNode.Children = append(characterNode.Children, e.behavior.Appearance(e, b))

	if e.ChargingElapsedTime >= 10 {
		chargingNode := &draw.OptionsNode{}
		characterNode.Children = append(characterNode.Children, chargingNode)

		frames := b.ChargingSprites.ChargingAnimation.Frames
		if e.ChargingElapsedTime >= e.PowerShotChargeTime {
			frames = b.ChargingSprites.ChargedAnimation.Frames
		}
		frame := frames[int(e.ChargingElapsedTime)%len(frames)]
		chargingNode.Children = append(chargingNode.Children, draw.ImageWithFrame(b.ChargingSprites.Image, frame))
	}

	rootNode.Children = append(rootNode.Children, characterNode)

	if *debugDrawEntityMarker {
		debugEntityMarkerImageOnce.Do(func() {
			debugEntityMarkerImage = ebiten.NewImage(5, 5)
			for x := 0; x < 5; x++ {
				debugEntityMarkerImage.Set(x, 2, color.RGBA{0, 255, 0, 255})
			}
			for y := 0; y < 5; y++ {
				debugEntityMarkerImage.Set(2, y, color.RGBA{0, 255, 0, 255})
			}
		})
		rootNode.Children = append(rootNode.Children, draw.ImageWithOrigin(debugEntityMarkerImage, image.Point{2, 2}))
	}

	return rootNode
}

func (e *Entity) Step(sh *StepHandle) {
	e.elapsedTime++

	// Set anger, if required.
	if e.CurrentHit.TotalDamage >= 300 {
		e.IsAngry = true
	}

	// TODO: Process poison damage.

	// Process hit damage.
	mustLeave1HP := e.HP > 1 && e.Traits.FatalHitLeaves1HP
	e.HP -= e.CurrentHit.TotalDamage
	if e.HP < 0 {
		e.HP = 0
	}
	if mustLeave1HP {
		e.HP = 1
	}
	e.CurrentHit.TotalDamage = 0

	// Tick timers.
	// TODO: Verify this behavior is correct.
	e.behaviorElapsedTime++
	e.behavior.Step(e, sh)

	if !e.CurrentHit.Drag {
		if !e.IsBeingDragged /* && !e.isInTimestop */ {
			// Process flashing.
			if e.CurrentHit.FlashTime > 0 {
				e.FlashingTimeLeft = e.CurrentHit.FlashTime
				e.CurrentHit.FlashTime = 0
			}
			if e.FlashingTimeLeft > 0 {
				e.FlashingTimeLeft--
			}

			// Process paralyzed.
			if e.CurrentHit.ParalyzeTime > 0 {
				e.ParalyzedTimeLeft = e.CurrentHit.ParalyzeTime
				e.CurrentHit.ConfuseTime = 0
				e.CurrentHit.ParalyzeTime = 0
			}
			if e.ParalyzedTimeLeft > 0 {
				e.ParalyzedTimeLeft--
				e.FrozenTimeLeft = 0
				e.BubbledTimeLeft = 0
				e.ConfusedTimeLeft = 0
			}

			// Process frozen.
			if e.CurrentHit.FreezeTime > 0 {
				e.FrozenTimeLeft = e.CurrentHit.FreezeTime
				e.ParalyzedTimeLeft = 0
				e.CurrentHit.BubbleTime = 0
				e.CurrentHit.ConfuseTime = 0
				e.CurrentHit.FreezeTime = 0
			}
			if e.FrozenTimeLeft > 0 {
				e.FrozenTimeLeft--
				e.BubbledTimeLeft = 0
				e.ConfusedTimeLeft = 0
			}

			// Process bubbled.
			if e.CurrentHit.BubbleTime > 0 {
				e.BubbledTimeLeft = e.CurrentHit.BubbleTime
				e.ConfusedTimeLeft = 0
				e.ParalyzedTimeLeft = 0
				e.FrozenTimeLeft = 0
				e.CurrentHit.ConfuseTime = 0
				e.CurrentHit.BubbleTime = 0
			}
			if e.BubbledTimeLeft > 0 {
				e.BubbledTimeLeft--
				e.ConfusedTimeLeft = 0
			}

			// Process confused.
			if e.CurrentHit.ConfuseTime > 0 {
				e.ConfusedTimeLeft = e.CurrentHit.ConfuseTime
				e.ParalyzedTimeLeft = 0
				e.FrozenTimeLeft = 0
				e.BubbledTimeLeft = 0
				e.CurrentHit.FreezeTime = 0
				e.CurrentHit.BubbleTime = 0
				e.CurrentHit.ParalyzeTime = 0
				e.CurrentHit.ConfuseTime = 0
			}
			if e.ConfusedTimeLeft > 0 {
				e.ConfusedTimeLeft--
			}

			// Process immobilized.
			if e.CurrentHit.ImmobilizeTime > 0 {
				e.ImmobilizedTimeLeft = e.CurrentHit.ImmobilizeTime
				e.CurrentHit.ImmobilizeTime = 0
			}
			if e.ImmobilizedTimeLeft > 0 {
				e.ImmobilizedTimeLeft--
			}

			// Process blinded.
			if e.CurrentHit.BlindTime > 0 {
				e.BlindedTimeLeft = e.CurrentHit.BlindTime
				e.CurrentHit.BlindTime = 0
			}
			if e.BlindedTimeLeft > 0 {
				e.BlindedTimeLeft--
			}

			// Process invincible.
			if e.InvincibleTimeLeft > 0 {
				e.InvincibleTimeLeft--
			}
		} else {
			// TODO: Interrupt player.
		}
	} else {
		e.CurrentHit.Drag = false

		e.FrozenTimeLeft = 0
		e.BubbledTimeLeft = 0
		e.ParalyzedTimeLeft = 0
		e.CurrentHit.BubbleTime = 0
		e.CurrentHit.FreezeTime = 0

		if false {
			e.ParalyzedTimeLeft = 0
		}

		// TODO: Interrupt player.
	}

	// Update UI.
	if e.DisplayHP != 0 && e.DisplayHP != e.HP {
		if e.HP < e.DisplayHP {
			e.DisplayHP -= ((e.DisplayHP-e.HP)>>3 + 4)
			if e.DisplayHP < e.HP {
				e.DisplayHP = e.HP
			}
		} else {
			e.DisplayHP += ((e.HP-e.DisplayHP)>>3 + 4)
			if e.DisplayHP > e.HP {
				e.DisplayHP = e.HP
			}
		}
	}
}

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
