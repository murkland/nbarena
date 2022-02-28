package state

import (
	"github.com/yumland/yumbattle/bundle"
	"github.com/yumland/yumbattle/draw"
)

type Hit struct {
	Damage int

	FlashTime      Ticks
	ParalyzeTime   Ticks
	ConfuseTime    Ticks
	BlindTime      Ticks
	ImmobilizeTime Ticks
	FreezeTime     Ticks
	BubbleTime     Ticks

	// ???
	Drag bool
}

func (h *Hit) Merge(h2 Hit) {
	h.Damage += h2.Damage

	// TODO: Verify this is correct behavior.
	h.ParalyzeTime = h2.ParalyzeTime
	h.ConfuseTime = h2.ConfuseTime
	h.BlindTime = h2.BlindTime
	h.ImmobilizeTime = h2.ImmobilizeTime
	h.FreezeTime = h2.FreezeTime
	h.BubbleTime = h2.BubbleTime
}

type Entity struct {
	elapsedTime Ticks

	behaviorElapsedTime Ticks
	behavior            EntityBehavior

	tilePos       TilePos
	futureTilePos TilePos

	isAlliedWithAnswerer bool

	isFlipped bool

	isDeleted bool

	hp        int
	displayHP int

	canStepOnHoleLikeTiles bool
	ignoresTileEffects     bool
	cannotFlinch           bool
	fatalHitLeaves1HP      bool

	chargingElapsedTime Ticks
	powerShotChargeTime Ticks

	paralyzedTimeLeft   Ticks
	confusedTimeLeft    Ticks
	blindedTimeLeft     Ticks
	immobilizedTimeLeft Ticks
	flashingTimeLeft    Ticks
	invincibleTimeLeft  Ticks
	frozenTimeLeft      Ticks
	bubbledTimeLeft     Ticks

	currentHit Hit

	isAngry        bool
	isBeingDragged bool
	isSliding      bool
}

func (e *Entity) Clone() *Entity {
	return &Entity{
		e.elapsedTime,
		e.behaviorElapsedTime, e.behavior.Clone(),
		e.tilePos, e.futureTilePos,
		e.isAlliedWithAnswerer,
		e.isFlipped,
		e.isDeleted,
		e.hp, e.displayHP,
		e.canStepOnHoleLikeTiles, e.ignoresTileEffects, e.cannotFlinch, e.fatalHitLeaves1HP,
		e.chargingElapsedTime, e.powerShotChargeTime,
		e.paralyzedTimeLeft,
		e.confusedTimeLeft,
		e.blindedTimeLeft,
		e.immobilizedTimeLeft,
		e.flashingTimeLeft,
		e.invincibleTimeLeft,
		e.frozenTimeLeft,
		e.bubbledTimeLeft,
		e.currentHit,
		e.isAngry,
		e.isBeingDragged,
		e.isSliding,
	}
}

func (e *Entity) SetBehavior(behavior EntityBehavior) {
	e.behaviorElapsedTime = 0
	e.behavior = behavior
}

func (e *Entity) TilePos() TilePos {
	return e.tilePos
}

func (e *Entity) SetTilePos(tilePos TilePos) {
	e.tilePos = tilePos
}

func (e *Entity) HP() int {
	return e.hp
}

func (e *Entity) SetHP(hp int) {
	e.hp = hp
}

func (e *Entity) CanStepOnHoleLikeTiles() bool {
	return e.canStepOnHoleLikeTiles
}

func (e *Entity) IgnoresTileEffects() bool {
	return e.ignoresTileEffects
}

func (e *Entity) Appearance(b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	x, y := e.tilePos.XY()
	rootNode.Opts.GeoM.Translate(float64((x-1)*tileRenderedWidth+tileRenderedWidth/2), float64((y-1)*tileRenderedHeight+tileRenderedHeight/2))

	characterNode := &draw.OptionsNode{}
	if e.isFlipped {
		characterNode.Opts.GeoM.Scale(-1, 1)
	}
	if e.frozenTimeLeft > 0 {
		// TODO: Render ice.
		characterNode.Opts.ColorM.Translate(float64(0xa5)/float64(0xff), float64(0xa5)/float64(0xff), float64(0xff)/float64(0xff), 0.0)
	}
	if e.paralyzedTimeLeft > 0 && (e.elapsedTime/2)%2 == 1 {
		characterNode.Opts.ColorM.Translate(1.0, 1.0, 0.0, 0.0)
	}
	if e.flashingTimeLeft > 0 && (e.elapsedTime/2)%2 == 0 {
		characterNode.Opts.ColorM.Translate(0.0, 0.0, 0.0, -1.0)
	}
	characterNode.Children = append(characterNode.Children, e.behavior.Appearance(e, b))

	if e.chargingElapsedTime >= 10 {
		chargingNode := &draw.OptionsNode{}
		characterNode.Children = append(characterNode.Children, chargingNode)

		frames := b.ChargingSprites.ChargingAnimation.Frames
		if e.chargingElapsedTime >= e.powerShotChargeTime {
			frames = b.ChargingSprites.ChargedAnimation.Frames
		}
		frame := frames[int(e.chargingElapsedTime)%len(frames)]
		chargingNode.Children = append(chargingNode.Children, draw.ImageWithFrame(b.ChargingSprites.Image, frame))
	}

	rootNode.Children = append(rootNode.Children, characterNode)

	return rootNode
}

func (e *Entity) Step() {
	e.elapsedTime++

	// Set anger, if required.
	if e.currentHit.Damage >= 300 {
		e.isAngry = true
	}

	// TODO: Process poison damage.

	// Process hit damage.
	mustLeave1HP := e.hp > 1 && e.fatalHitLeaves1HP
	e.hp -= e.currentHit.Damage
	if e.hp < 0 {
		e.hp = 0
	}
	if mustLeave1HP {
		e.hp = 1
	}
	e.currentHit.Damage = 0

	// Tick timers.
	// TODO: Verify this behavior is correct.
	e.behaviorElapsedTime++
	e.behavior.Step(e)

	if !e.currentHit.Drag {
		if !e.isBeingDragged /* && !e.isInTimestop */ {
			// Process flashing.
			if e.currentHit.FlashTime > 0 {
				e.flashingTimeLeft = e.currentHit.FlashTime
				e.currentHit.FlashTime = 0
			}
			if e.flashingTimeLeft > 0 {
				e.flashingTimeLeft--
			}

			// Process paralyzed.
			if e.currentHit.ParalyzeTime > 0 {
				e.paralyzedTimeLeft = e.currentHit.ParalyzeTime
				e.currentHit.ConfuseTime = 0
				e.currentHit.ParalyzeTime = 0
			}
			if e.paralyzedTimeLeft > 0 {
				e.paralyzedTimeLeft--
				e.frozenTimeLeft = 0
				e.bubbledTimeLeft = 0
				e.confusedTimeLeft = 0
			}

			// Process frozen.
			if e.currentHit.FreezeTime > 0 {
				e.frozenTimeLeft = e.currentHit.FreezeTime
				e.paralyzedTimeLeft = 0
				e.currentHit.BubbleTime = 0
				e.currentHit.ConfuseTime = 0
				e.currentHit.FreezeTime = 0
			}
			if e.frozenTimeLeft > 0 {
				e.frozenTimeLeft--
				e.bubbledTimeLeft = 0
				e.confusedTimeLeft = 0
			}

			// Process bubbled.
			if e.currentHit.BubbleTime > 0 {
				e.bubbledTimeLeft = e.currentHit.BubbleTime
				e.confusedTimeLeft = 0
				e.paralyzedTimeLeft = 0
				e.frozenTimeLeft = 0
				e.currentHit.ConfuseTime = 0
				e.currentHit.BubbleTime = 0
			}
			if e.bubbledTimeLeft > 0 {
				e.bubbledTimeLeft--
				e.confusedTimeLeft = 0
			}

			// Process confused.
			if e.currentHit.ConfuseTime > 0 {
				e.confusedTimeLeft = e.currentHit.ConfuseTime
				e.paralyzedTimeLeft = 0
				e.frozenTimeLeft = 0
				e.bubbledTimeLeft = 0
				e.currentHit.FreezeTime = 0
				e.currentHit.BubbleTime = 0
				e.currentHit.ParalyzeTime = 0
				e.currentHit.ConfuseTime = 0
			}
			if e.confusedTimeLeft > 0 {
				e.confusedTimeLeft--
			}

			// Process immobilized.
			if e.currentHit.ImmobilizeTime > 0 {
				e.immobilizedTimeLeft = e.currentHit.ImmobilizeTime
				e.currentHit.ImmobilizeTime = 0
			}
			if e.immobilizedTimeLeft > 0 {
				e.immobilizedTimeLeft--
			}

			// Process blinded.
			if e.currentHit.BlindTime > 0 {
				e.blindedTimeLeft = e.currentHit.BlindTime
				e.currentHit.BlindTime = 0
			}
			if e.blindedTimeLeft > 0 {
				e.blindedTimeLeft--
			}

			// Process invincible.
			if e.invincibleTimeLeft > 0 {
				e.invincibleTimeLeft--
			}
		} else {
			// TODO: Interrupt player.
		}
	} else {
		e.currentHit.Drag = false

		e.frozenTimeLeft = 0
		e.bubbledTimeLeft = 0
		e.paralyzedTimeLeft = 0
		e.currentHit.BubbleTime = 0
		e.currentHit.FreezeTime = 0

		if false {
			e.paralyzedTimeLeft = 0
		}

		// TODO: Interrupt player.
	}

	// Update UI.
	if e.displayHP != 0 {
		var newDisplayHP int
		dhp := e.displayHP - e.hp
		if dhp < 0 {
			newDisplayHP := e.displayHP - (-dhp >> 3) + 2
			if newDisplayHP < e.hp {
				newDisplayHP = e.hp
			}
		} else {
			newDisplayHP := e.displayHP + (dhp >> 3) + 2
			if newDisplayHP > e.hp {
				newDisplayHP = e.hp
			}
		}
		e.displayHP = newDisplayHP
	}
}
