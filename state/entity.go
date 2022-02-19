package state

import (
	"github.com/yumland/clone"
	"github.com/yumland/yumbattle/draw"
)

type Entity struct {
	appearance draw.Node

	tilePos TilePos

	hp        int
	displayHP *int

	canStepOnHoleLikeTiles bool
	ignoresTileEffects     bool
}

func (e *Entity) Clone() *Entity {
	return &Entity{
		e.appearance, // Appearances are not cloned: they are considered immutable enough.
		e.tilePos,
		e.hp, clone.Shallow(e.displayHP),
		e.canStepOnHoleLikeTiles, e.ignoresTileEffects,
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
