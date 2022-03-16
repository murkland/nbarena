package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/murkland/clone"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
)

type TileBehaviorState struct {
	Behavior    TileBehavior
	ElapsedTime Ticks
}

func (s TileBehaviorState) Clone() TileBehaviorState {
	return TileBehaviorState{clone.Interface[TileBehavior](s.Behavior), s.ElapsedTime}
}

type Tile struct {
	BehaviorState TileBehaviorState

	IsFlipped     bool
	IsHighlighted bool

	Reserver EntityID

	IsAlliedWithAnswerer bool
}

func (t *Tile) Clone() *Tile {
	return &Tile{
		t.BehaviorState.Clone(),
		t.IsFlipped, t.IsHighlighted,
		t.Reserver,
		t.IsAlliedWithAnswerer,
	}
}

func (t *Tile) Flip() {
	t.IsFlipped = !t.IsFlipped
	t.IsAlliedWithAnswerer = !t.IsAlliedWithAnswerer
}

func (t *Tile) CanEnter(e *Entity) bool {
	if t.BehaviorState.Behavior == nil {
		return false
	}

	return t.BehaviorState.Behavior.CanEnter(t, e)
}

func (t *Tile) ReplaceBehavior(b TileBehavior) {
	t.BehaviorState.ElapsedTime = 0
	t.BehaviorState.Behavior = b
	t.BehaviorState.Behavior.Step(t)
}

func (t *Tile) ElapsedTime() Ticks {
	if t.BehaviorState.Behavior == nil {
		return 0
	}

	return t.BehaviorState.ElapsedTime
}

func (t *Tile) Step() {
	if t.BehaviorState.Behavior == nil {
		return
	}

	t.BehaviorState.ElapsedTime++
	t.BehaviorState.Behavior.Step(t)
}

func (t *Tile) OnEnter(e *Entity, s *State) {
	if t.BehaviorState.Behavior == nil {
		return
	}

	t.BehaviorState.Behavior.OnEnter(t, e, s)
}
func (t *Tile) OnLeave(e *Entity, s *State) {
	if t.BehaviorState.Behavior == nil {
		return
	}

	t.BehaviorState.Behavior.OnLeave(t, e, s)
}

func (t *Tile) Appearance(y int, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	if t.BehaviorState.Behavior == nil {
		return nil
	}
	tiles := b.Battletiles.OffererTiles
	if t.IsAlliedWithAnswerer {
		tiles = b.Battletiles.AnswererTiles
	}
	rootNode.Children = append(rootNode.Children, t.BehaviorState.Behavior.Appearance(t, y, b, tiles))
	if t.IsHighlighted {
		rootNode.Opts.ColorM.Translate(1.0, 1.0, 1.0, 1.0)
		rootNode.Opts.ColorM.Scale(1.0, 1.0, 0.0, 1.0)
	}
	return rootNode
}

const TileRows = 5
const TileCols = 8

type TilePos int

func TilePosXY(x int, y int) TilePos {
	if x < 0 || x >= TileCols || y < 0 || y >= TileRows {
		return -1
	}
	return TilePos(y*TileCols + x)
}

func (p TilePos) XY() (int, int) {
	return int(p) % TileCols, int(p) / TileCols
}

func (p TilePos) Flipped() TilePos {
	x, y := p.XY()
	return TilePosXY(TileCols-x-1, y)
}

type TileBehavior interface {
	clone.Cloner[TileBehavior]
	Appearance(t *Tile, y int, b *bundle.Bundle, tiles *ebiten.Image) draw.Node
	CanEnter(t *Tile, e *Entity) bool
	OnEnter(t *Tile, e *Entity, s *State)
	OnLeave(t *Tile, e *Entity, s *State)
	Step(t *Tile)
}

type HoleTileBehavior struct {
}

func (tb *HoleTileBehavior) Clone() TileBehavior {
	return &HoleTileBehavior{}
}

func (tb *HoleTileBehavior) Appearance(t *Tile, y int, b *bundle.Bundle, tiles *ebiten.Image) draw.Node {
	return nil
}

func (tb *HoleTileBehavior) CanEnter(t *Tile, e *Entity) bool {
	return e.Traits.CanStepOnHoleLikeTiles
}
func (tb *HoleTileBehavior) OnEnter(t *Tile, e *Entity, s *State) {}
func (tb *HoleTileBehavior) OnLeave(t *Tile, e *Entity, s *State) {}
func (tb *HoleTileBehavior) Step(t *Tile)                         {}

type BrokenTileBehavior struct {
	returnToNormalTimeLeft int
}

func (tb *BrokenTileBehavior) Clone() TileBehavior {
	return &BrokenTileBehavior{tb.returnToNormalTimeLeft}
}

func (tb *BrokenTileBehavior) Appearance(t *Tile, y int, b *bundle.Bundle, tiles *ebiten.Image) draw.Node {
	frame := b.Battletiles.Info.Animations[1*3+(y-1)].Frames[0]
	return draw.ImageWithFrame(tiles, frame)
}

func (tb *BrokenTileBehavior) CanEnter(t *Tile, e *Entity) bool {
	return e.Traits.CanStepOnHoleLikeTiles
}
func (tb *BrokenTileBehavior) OnEnter(t *Tile, e *Entity, s *State) {}
func (tb *BrokenTileBehavior) OnLeave(t *Tile, e *Entity, s *State) {}

func (tb *BrokenTileBehavior) Step(t *Tile) {
	if tb.returnToNormalTimeLeft > 0 {
		tb.returnToNormalTimeLeft--
		if tb.returnToNormalTimeLeft <= 0 {
			t.ReplaceBehavior(&NormalTileBehavior{})
		}
	}
}

type NormalTileBehavior struct {
}

func (tb *NormalTileBehavior) Clone() TileBehavior {
	return &NormalTileBehavior{}
}

func (tb *NormalTileBehavior) Appearance(t *Tile, y int, b *bundle.Bundle, tiles *ebiten.Image) draw.Node {
	frame := b.Battletiles.Info.Animations[2*3+(y-1)].Frames[0]
	return draw.ImageWithFrame(tiles, frame)
}

func (tb *NormalTileBehavior) CanEnter(t *Tile, e *Entity) bool {
	return true
}
func (tb *NormalTileBehavior) OnEnter(t *Tile, e *Entity, s *State) {}
func (tb *NormalTileBehavior) OnLeave(t *Tile, e *Entity, s *State) {}
func (tb *NormalTileBehavior) Step(t *Tile)                         {}

type CrackedTileBehavior struct {
}

func (tb *CrackedTileBehavior) Clone() TileBehavior {
	return &CrackedTileBehavior{}
}

func (tb *CrackedTileBehavior) Appearance(t *Tile, y int, b *bundle.Bundle, tiles *ebiten.Image) draw.Node {
	frame := b.Battletiles.Info.Animations[3*3+(y-1)].Frames[0]
	return draw.ImageWithFrame(tiles, frame)
}

func (tb *CrackedTileBehavior) CanEnter(t *Tile, e *Entity) bool {
	return true
}
func (tb *CrackedTileBehavior) OnEnter(t *Tile, e *Entity, s *State) {
}
func (tb *CrackedTileBehavior) OnLeave(t *Tile, e *Entity, s *State) {
	if e.Traits.IgnoresTileEffects {
		return
	}
	// TODO: Play cracking sound.
	// TODO: Add returnToNormalTimeLeft
	t.ReplaceBehavior(&BrokenTileBehavior{})
}
func (tb *CrackedTileBehavior) Step(t *Tile) {}
