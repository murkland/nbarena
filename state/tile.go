package state

import (
	"github.com/yumland/clone"
	"github.com/yumland/yumbattle/draw"
)

type Tile struct {
	behavior TileBehavior
	team     Team
}

func (t *Tile) Clone() *Tile {
	return &Tile{t.behavior.Clone(), t.team}
}

func (t *Tile) SetBehavior(b TileBehavior) {
	t.behavior = b
}

func (t *Tile) Team() Team {
	return t.team
}

const tileRows = 5
const tileCols = 8

type TilePos int

func EmptyTiles() []*Tile {
	return make([]*Tile, tileRows*tileCols)
}

func TilePosXY(x int, y int) TilePos {
	return TilePos(y*tileCols + x)
}

func (p TilePos) XY() (int, int) {
	return int(p) / tileCols, int(p) % tileCols
}

type TileBehavior interface {
	clone.Interface[TileBehavior]
	Appearance(t *Tile) draw.Node
	CanStepOn(e *Entity) bool
	OnEnter(t *Tile, e *Entity)
	OnLeave(t *Tile, e *Entity)
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
	return e.CanStepOnHoleLikeTiles
}
func (tb *HoleTileBehavior) OnEnter(t *Tile, e *Entity) {}
func (tb *HoleTileBehavior) OnLeave(t *Tile, e *Entity) {}

type BrokenTileBehavior struct {
}

func (tb *BrokenTileBehavior) Clone() TileBehavior {
	return &BrokenTileBehavior{}
}

func (tb *BrokenTileBehavior) Appearance(t *Tile) draw.Node {
	return nil
}

func (tb *BrokenTileBehavior) CanStepOn(e *Entity) bool {
	return e.CanStepOnHoleLikeTiles
}
func (tb *BrokenTileBehavior) OnEnter(t *Tile, e *Entity) {}
func (tb *BrokenTileBehavior) OnLeave(t *Tile, e *Entity) {}

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
	if e.IgnoresTileEffects {
		return
	}
	t.SetBehavior(&BrokenTileBehavior{})
}
