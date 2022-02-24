package state

import (
	"github.com/yumland/clone"
	"github.com/yumland/yumbattle/bundle"
	"github.com/yumland/yumbattle/draw"
)

type Hit struct {
	Damage int

	FlashFrames      int
	ParalyzeFrames   int
	ConfuseFrames    int
	BlindFrames      int
	ImmobilizeFrames int
	FreezeFrames     int
	BubbleFrames     int

	// ???
	Drag bool
}

func (h *Hit) Merge(h2 Hit) {
	h.Damage += h2.Damage

	// TODO: Verify this is correct behavior.
	h.ParalyzeFrames = h2.ParalyzeFrames
	h.ConfuseFrames = h2.ConfuseFrames
	h.BlindFrames = h2.BlindFrames
	h.ImmobilizeFrames = h2.ImmobilizeFrames
	h.FreezeFrames = h2.FreezeFrames
	h.BubbleFrames = h2.BubbleFrames
}

type Entity struct {
	behaviorElapsed int
	behavior        EntityBehavior

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

	paralyzedFramesLeft   int
	confusedFramesLeft    int
	blindedFramesLeft     int
	immobilizedFramesLeft int
	flashingFramesLeft    int
	invincibleFramesLeft  int
	frozenFramesLeft      int
	bubbledFramesLeft     int

	currentHit Hit

	isAngry        bool
	isBeingDragged bool
	isSliding      bool
}

func (e *Entity) Clone() *Entity {
	return &Entity{
		e.behaviorElapsed, clone.Interface[EntityBehavior](e.behavior),
		e.tilePos, e.futureTilePos,
		e.isAlliedWithAnswerer,
		e.isFlipped,
		e.isDeleted,
		e.hp, e.displayHP,
		e.canStepOnHoleLikeTiles, e.ignoresTileEffects, e.cannotFlinch, e.fatalHitLeaves1HP,
		e.paralyzedFramesLeft,
		e.confusedFramesLeft,
		e.blindedFramesLeft,
		e.immobilizedFramesLeft,
		e.flashingFramesLeft,
		e.invincibleFramesLeft,
		e.frozenFramesLeft,
		e.bubbledFramesLeft,
		e.currentHit,
		e.isAngry,
		e.isBeingDragged,
		e.isSliding,
	}
}

func (e *Entity) SetBehavior(behavior EntityBehavior) {
	e.behaviorElapsed = 0
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
	characterNode.Children = append(characterNode.Children, e.behavior.Appearance(e, b))
	rootNode.Children = append(rootNode.Children, characterNode)

	return rootNode
}

func (e *Entity) Step() {
	// TODO: Handle action.

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
	e.behaviorElapsed++
	e.behavior.Step(e)

	if !e.currentHit.Drag {
		if !e.isBeingDragged /* && !e.isInTimestop */ {
			// Process flashing.
			if e.currentHit.FlashFrames > 0 {
				e.flashingFramesLeft = e.currentHit.FlashFrames
				e.currentHit.FlashFrames = 0
			}
			if e.flashingFramesLeft > 0 {
				e.flashingFramesLeft--
			}

			// Process paralyzed.
			if e.currentHit.ParalyzeFrames > 0 {
				e.paralyzedFramesLeft = e.currentHit.ParalyzeFrames
				e.currentHit.ConfuseFrames = 0
				e.currentHit.ParalyzeFrames = 0
			}
			if e.paralyzedFramesLeft > 0 {
				e.paralyzedFramesLeft--
				e.frozenFramesLeft = 0
				e.bubbledFramesLeft = 0
				e.confusedFramesLeft = 0
			}

			// Process frozen.
			if e.currentHit.FreezeFrames > 0 {
				e.frozenFramesLeft = e.currentHit.FreezeFrames
				e.paralyzedFramesLeft = 0
				e.currentHit.BubbleFrames = 0
				e.currentHit.ConfuseFrames = 0
				e.currentHit.FreezeFrames = 0
			}
			if e.frozenFramesLeft > 0 {
				e.frozenFramesLeft--
				e.bubbledFramesLeft = 0
				e.confusedFramesLeft = 0
			}

			// Process bubbled.
			if e.currentHit.BubbleFrames > 0 {
				e.bubbledFramesLeft = e.currentHit.BubbleFrames
				e.confusedFramesLeft = 0
				e.paralyzedFramesLeft = 0
				e.frozenFramesLeft = 0
				e.currentHit.ConfuseFrames = 0
				e.currentHit.BubbleFrames = 0
			}
			if e.bubbledFramesLeft > 0 {
				e.bubbledFramesLeft--
				e.confusedFramesLeft = 0
			}

			// Process confused.
			if e.currentHit.ConfuseFrames > 0 {
				e.confusedFramesLeft = e.currentHit.ConfuseFrames
				e.paralyzedFramesLeft = 0
				e.frozenFramesLeft = 0
				e.bubbledFramesLeft = 0
				e.currentHit.FreezeFrames = 0
				e.currentHit.BubbleFrames = 0
				e.currentHit.ParalyzeFrames = 0
				e.currentHit.ConfuseFrames = 0
			}
			if e.confusedFramesLeft > 0 {
				e.confusedFramesLeft--
			}

			// Process immobilized.
			if e.currentHit.ImmobilizeFrames > 0 {
				e.immobilizedFramesLeft = e.currentHit.ImmobilizeFrames
				e.currentHit.ImmobilizeFrames = 0
			}
			if e.immobilizedFramesLeft > 0 {
				e.immobilizedFramesLeft--
			}

			// Process blinded.
			if e.currentHit.BlindFrames > 0 {
				e.blindedFramesLeft = e.currentHit.BlindFrames
				e.currentHit.BlindFrames = 0
			}
			if e.blindedFramesLeft > 0 {
				e.blindedFramesLeft--
			}

			// Process invincible.
			if e.invincibleFramesLeft > 0 {
				e.invincibleFramesLeft--
			}
		} else {
			// TODO: Interrupt player.
		}
	} else {
		e.currentHit.Drag = false

		e.frozenFramesLeft = 0
		e.bubbledFramesLeft = 0
		e.paralyzedFramesLeft = 0
		e.currentHit.BubbleFrames = 0
		e.currentHit.FreezeFrames = 0

		if false {
			e.paralyzedFramesLeft = 0
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
