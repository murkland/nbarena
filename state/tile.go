package state

import (
	"github.com/yumland/clone"
	"github.com/yumland/yumbattle/draw"
)

type Tile struct {
	behavior             TileBehavior
	isAlliedWithAnswerer bool
}

func (t Tile) Clone() Tile {
	return Tile{clone.Interface[TileBehavior](t.behavior), t.isAlliedWithAnswerer}
}

func (t *Tile) SetBehavior(b TileBehavior) {
	t.behavior = b
}

func (t *Tile) IsAlliedWithAnswerer() bool {
	return t.isAlliedWithAnswerer
}

func (t *Tile) Step() {
	if t.behavior == nil {
		return
	}

	t.behavior.Step(t)
}

const tileRows = 5
const tileCols = 8

type TilePos int

func TilePosXY(x int, y int) TilePos {
	return TilePos(y*tileCols + x)
}

func (p TilePos) XY() (int, int) {
	return int(p) / tileCols, int(p) % tileCols
}

type TileBehavior interface {
	clone.Cloner[TileBehavior]
	Appearance(t *Tile) draw.Node
	CanStepOn(e *Entity) bool
	OnEnter(t *Tile, e *Entity)
	OnLeave(t *Tile, e *Entity)
	Step(t *Tile)
}

type HoleTileBehavior struct {
}

func (tb *HoleTileBehavior) Clone() TileBehavior {
	return &HoleTileBehavior{}
}

func (tb *HoleTileBehavior) Appearance(t *Tile) draw.Node {
	return nil
}

func (tb *HoleTileBehavior) CanStepOn(e *Entity) bool {
	return e.CanStepOnHoleLikeTiles()
}
func (tb *HoleTileBehavior) OnEnter(t *Tile, e *Entity) {}
func (tb *HoleTileBehavior) OnLeave(t *Tile, e *Entity) {}
func (tb *HoleTileBehavior) Step(t *Tile)               {}

type BrokenTileBehavior struct {
	returnToNormalTimeLeft int
}

func (tb *BrokenTileBehavior) Clone() TileBehavior {
	return &BrokenTileBehavior{tb.returnToNormalTimeLeft}
}

func (tb *BrokenTileBehavior) Appearance(t *Tile) draw.Node {
	return nil
}

func (tb *BrokenTileBehavior) CanStepOn(e *Entity) bool {
	return e.CanStepOnHoleLikeTiles()
}
func (tb *BrokenTileBehavior) OnEnter(t *Tile, e *Entity) {}
func (tb *BrokenTileBehavior) OnLeave(t *Tile, e *Entity) {}

func (tb *BrokenTileBehavior) Step(t *Tile) {
	if tb.returnToNormalTimeLeft > 0 {
		tb.returnToNormalTimeLeft--
		if tb.returnToNormalTimeLeft <= 0 {
			t.SetBehavior(&NormalTileBehavior{})
		}
	}
}

type NormalTileBehavior struct {
}

func (tb *NormalTileBehavior) Clone() TileBehavior {
	return &NormalTileBehavior{}
}

func (tb *NormalTileBehavior) Appearance(t *Tile) draw.Node {
	return nil
}

func (tb *NormalTileBehavior) CanStepOn(e *Entity) bool {
	return true
}
func (tb *NormalTileBehavior) OnEnter(t *Tile, e *Entity) {}
func (tb *NormalTileBehavior) OnLeave(t *Tile, e *Entity) {}
func (tb *NormalTileBehavior) Step(t *Tile)               {}

type CrackedTile struct {
}

func (tb *CrackedTile) Clone() TileBehavior {
	return &CrackedTile{}
}

func (tb *CrackedTile) Appearance(t *Tile) draw.Node {
	return nil
}

func (tb *CrackedTile) CanStepOn(e *Entity) bool {
	return true
}
func (tb *CrackedTile) OnEnter(t *Tile, e *Entity) {
}
func (tb *CrackedTile) OnLeave(t *Tile, e *Entity) {
	if e.IgnoresTileEffects() {
		return
	}
	t.SetBehavior(&BrokenTileBehavior{})
}
func (tb *CrackedTile) Step(t *Tile) {}
