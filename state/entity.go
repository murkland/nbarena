package state

import (
	"github.com/yumland/clone"
	"github.com/yumland/yumbattle/draw"
)

type Hit struct {
	Damage int

	ParalyzeFrames   int
	ConfuseFrames    int
	BlindFrames      int
	ImmobilizeFrames int
	FreezeFrames     int
	BubbleFrames     int
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
	appearance draw.Node

	tilePos       TilePos
	futureTilePos TilePos

	isAlliedWithAnswerer bool

	isFlipped bool

	hp        int
	displayHP *int

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

	isBeingDragged bool
	isSliding      bool
}

func (e *Entity) Clone() *Entity {
	return &Entity{
		e.appearance, // Appearances are not cloned: they are considered immutable enough.
		e.tilePos, e.futureTilePos,
		e.isAlliedWithAnswerer,
		e.isFlipped,
		e.hp, clone.Shallow(e.displayHP),
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
		e.isBeingDragged,
		e.isSliding,
	}
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

func (e *Entity) Appearance() draw.Node {
	return e.appearance
}

func (e *Entity) Step() {
	// TODO: Handle action.

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
	if !e.isBeingDragged /* && !e.isInTimestop */ {
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
	}

	// Update UI.
	if e.displayHP != nil {
		dhp := *e.displayHP - e.hp
		var newDisplayHP int
		if dhp < 0 {
			newDisplayHP := *e.displayHP - (dhp >> 3) + 2
			if newDisplayHP < e.hp {
				newDisplayHP = e.hp
			}
		} else {
			newDisplayHP := *e.displayHP + (dhp >> 3) + 2
			if newDisplayHP > e.hp {
				newDisplayHP = e.hp
			}
		}
		*e.displayHP = newDisplayHP
	}
}
