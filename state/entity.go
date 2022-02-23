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

	currentHit *Hit

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

	// Tick timers.
	if e.currentHit != nil {
		if !e.isBeingDragged /* && !e.isInTimestop */ {
			// Process paralyzed.
			if e.paralyzedFramesLeft > 0 {
				if e.currentHit.ParalyzeFrames > 0 {
					e.paralyzedFramesLeft = e.currentHit.ParalyzeFrames
					e.currentHit.ConfuseFrames = 0
				}
				e.paralyzedFramesLeft--
				e.frozenFramesLeft = 0
				e.bubbledFramesLeft = 0
				e.confusedFramesLeft = 0
			}
			e.currentHit.ParalyzeFrames = 0

			// Process frozen.
			if e.frozenFramesLeft > 0 {
				if e.currentHit.FreezeFrames > 0 {
					e.frozenFramesLeft = e.currentHit.FreezeFrames
					e.currentHit.BubbleFrames = 0
					e.currentHit.ConfuseFrames = 0
					e.paralyzedFramesLeft = 0
				}
				e.frozenFramesLeft--
				e.bubbledFramesLeft = 0
				e.confusedFramesLeft = 0
			}
			e.currentHit.FreezeFrames = 0

			// Process bubbled.
			if e.bubbledFramesLeft > 0 {
				if e.currentHit.BubbleFrames > 0 {
					e.bubbledFramesLeft = e.currentHit.BubbleFrames
					e.currentHit.ConfuseFrames = 0
					e.confusedFramesLeft = 0
					e.paralyzedFramesLeft = 0
					e.frozenFramesLeft = 0
				}
				e.bubbledFramesLeft--
				e.confusedFramesLeft = 0
			}
			e.currentHit.BubbleFrames = 0

			// Process confused.
			if e.confusedFramesLeft > 0 {
				if e.currentHit.ConfuseFrames > 0 {
					e.confusedFramesLeft = e.currentHit.ConfuseFrames
					e.currentHit.FreezeFrames = 0
					e.currentHit.BubbleFrames = 0
					e.currentHit.ParalyzeFrames = 0
					e.paralyzedFramesLeft = 0
					e.frozenFramesLeft = 0
					e.bubbledFramesLeft = 0
				}
				e.confusedFramesLeft--
			}
			e.currentHit.ConfuseFrames = 0

			// Process immobilized.
			if e.immobilizedFramesLeft > 0 {
				if e.currentHit.ImmobilizeFrames > 0 {
					e.immobilizedFramesLeft = e.currentHit.ImmobilizeFrames
				}
				e.immobilizedFramesLeft--
			}
			e.currentHit.ImmobilizeFrames = 0

			// Process blinded.
			if e.blindedFramesLeft > 0 {
				if e.currentHit.BlindFrames > 0 {
					e.blindedFramesLeft = e.currentHit.BlindFrames
				}
				e.blindedFramesLeft--
			}
			e.currentHit.BlindFrames = 0

			// Process invincible.
			if e.invincibleFramesLeft > 0 {
				e.invincibleFramesLeft--
			}
		}
		e.currentHit = nil
	}
}
