package state

import (
	"github.com/yumland/clone"
	"github.com/yumland/yumbattle/draw"
)

type Entity struct {
	appearance draw.Node

	tilePos       TilePos
	futureTilePos TilePos

	isOwnedByAnswerer bool

	isFlipped bool

	hp        int
	displayHP *int

	canStepOnHoleLikeTiles bool
	ignoresTileEffects     bool
	cannotFlinch           bool
	fatalHitLeaves1HP      bool

	paralyzedFramesLeft   uint16
	confusedFramesLeft    uint16
	blindedFramesLeft     uint16
	immobilizedFramesLeft uint16
	flashingFramesLeft    uint16
	invincibleFramesLeft  uint16
	frozenFramesLeft      uint16
	bubbledFramesLeft     uint16

	isBeingDragged bool
}

func (e Entity) Clone() Entity {
	return Entity{
		e.appearance, // Appearances are not cloned: they are considered immutable enough.
		e.tilePos, e.futureTilePos,
		e.isOwnedByAnswerer,
		e.isFlipped,
		e.hp, clone.Shallow(e.displayHP),
		e.canStepOnHoleLikeTiles, e.ignoresTileEffects, e.cannotFlinch, e.fatalHitLeaves1HP,
		e.paralyzedFramesLeft, e.confusedFramesLeft, e.blindedFramesLeft, e.immobilizedFramesLeft, e.flashingFramesLeft, e.invincibleFramesLeft, e.frozenFramesLeft, e.bubbledFramesLeft,
		e.isBeingDragged,
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
