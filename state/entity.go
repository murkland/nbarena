package state

import (
	"github.com/yumland/clone"
	"github.com/yumland/yumbattle/draw"
)

type Hit struct {
	ParalyzeFrames   int
	ConfuseFrames    int
	BlindFrames      int
	ImmobilizeFrames int
	FreezeFrames     int
	BubbleFrames     int
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
		h,
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

func (e *Entity) Hit(h Hit) {
	if h.ParalyzeFrames > 0 {
		e.paralyzedFramesLeft = h.ParalyzeFrames
		e.frozenFramesLeft = 0
		e.bubbledFramesLeft = 0
		e.confusedFramesLeft = 0
		h.ConfuseFrames = 0
	}
	h.ParalyzeFrames = 0

	if h.FreezeFrames > 0 {
		e.frozenFramesLeft = h.FreezeFrames
		e.bubbledFramesLeft = 0
		e.confusedFramesLeft = 0
		e.paralyzedFramesLeft = 0
		h.BubbleFrames = 0
		h.ConfuseFrames = 0
	}
	h.FreezeFrames = 0

	if h.BubbleFrames > 0 {
		e.bubbledFramesLeft = h.BubbleFrames
		e.confusedFramesLeft = 0
		e.paralyzedFramesLeft = 0
		e.frozenFramesLeft = 0
		e.confusedFramesLeft = 0
		h.ConfuseFrames = 0
	}
	h.BubbleFrames = 0

	if h.ConfuseFrames > 0 {
		e.confusedFramesLeft = h.ConfuseFrames
		e.paralyzedFramesLeft = 0
		e.frozenFramesLeft = 0
		e.bubbledFramesLeft = 0
		h.FreezeFrames = 0
		h.BubbleFrames = 0
		h.ParalyzeFrames = 0
	}
	h.ConfuseFrames = 0

	if h.ImmobilizeFrames > 0 {
		e.immobilizedFramesLeft = h.ImmobilizeFrames
	}
	h.ImmobilizeFrames = 0

	if h.BlindFrames > 0 {
		e.blindedFramesLeft = h.BlindFrames
	}
	h.BlindFrames = 0
}

func (e *Entity) Step() {
	// TODO: Handle action.

	// Tick timers.
	if !e.isBeingDragged /* && !e.isInTimestop */ {
		if e.paralyzedFramesLeft > 0 {
			e.paralyzedFramesLeft--
		}

		if e.frozenFramesLeft > 0 {
			e.frozenFramesLeft--
		}

		if e.bubbledFramesLeft > 0 {
			e.bubbledFramesLeft--
		}

		if e.confusedFramesLeft > 0 {
			e.confusedFramesLeft--
		}

		if e.immobilizedFramesLeft > 0 {
			e.immobilizedFramesLeft--
		}

		if e.blindedFramesLeft > 0 {
			e.blindedFramesLeft--
		}

		if e.invincibleFramesLeft > 0 {
			e.invincibleFramesLeft--
		}
	}
}
