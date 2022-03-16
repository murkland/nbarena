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

	TilePos TilePos

	IsFlipped     bool
	IsHighlighted bool

	Reserver EntityID

	IsAlliedWithAnswerer bool
}

func (t *Tile) Clone() *Tile {
	return &Tile{
		t.BehaviorState.Clone(),
		t.TilePos,
		t.IsFlipped, t.IsHighlighted,
		t.Reserver,
		t.IsAlliedWithAnswerer,
	}
}

func (t *Tile) Flip() {
	t.IsFlipped = !t.IsFlipped
	t.IsAlliedWithAnswerer = !t.IsAlliedWithAnswerer
	if t.BehaviorState.Behavior != nil {
		t.BehaviorState.Behavior.Flip()
	}
}

func (t *Tile) CanEnter(e *Entity) bool {
	if t.BehaviorState.Behavior == nil {
		return false
	}

	return t.BehaviorState.Behavior.CanEnter(t, e)
}

func (t *Tile) ReplaceBehavior(b TileBehavior, s *State) {
	t.BehaviorState.ElapsedTime = 0
	t.BehaviorState.Behavior = b
	t.BehaviorState.Behavior.Step(t, s)
}

func (t *Tile) ElapsedTime() Ticks {
	if t.BehaviorState.Behavior == nil {
		return 0
	}

	return t.BehaviorState.ElapsedTime
}

func (t *Tile) Step(s *State) {
	if t.BehaviorState.Behavior == nil {
		return
	}

	t.BehaviorState.ElapsedTime++
	t.BehaviorState.Behavior.Step(t, s)
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
	OnLeave(t *Tile, e *Entity, s *State)
	Flip()
	Step(t *Tile, s *State)
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
func (tb *HoleTileBehavior) OnLeave(t *Tile, e *Entity, s *State) {}
func (tb *HoleTileBehavior) Flip()                                {}
func (tb *HoleTileBehavior) Step(t *Tile, s *State)               {}

type BrokenTileBehavior struct {
	returnToNormalTimeLeft int
}

func (tb *BrokenTileBehavior) Clone() TileBehavior {
	return &BrokenTileBehavior{tb.returnToNormalTimeLeft}
}

func (tb *BrokenTileBehavior) Appearance(t *Tile, y int, b *bundle.Bundle, tiles *ebiten.Image) draw.Node {
	return draw.ImageWithAnimation(tiles, b.Battletiles.Info.Animations[1*3+(y-1)], int(t.BehaviorState.ElapsedTime))
}

func (tb *BrokenTileBehavior) CanEnter(t *Tile, e *Entity) bool {
	return e.Traits.CanStepOnHoleLikeTiles
}
func (tb *BrokenTileBehavior) OnLeave(t *Tile, e *Entity, s *State) {}
func (tb *BrokenTileBehavior) Flip()                                {}

func (tb *BrokenTileBehavior) Step(t *Tile, s *State) {
	if tb.returnToNormalTimeLeft > 0 {
		tb.returnToNormalTimeLeft--
		if tb.returnToNormalTimeLeft <= 0 {
			t.ReplaceBehavior(&NormalTileBehavior{}, s)
		}
	}
}

type NormalTileBehavior struct {
}

func (tb *NormalTileBehavior) Clone() TileBehavior {
	return &NormalTileBehavior{}
}

func (tb *NormalTileBehavior) Appearance(t *Tile, y int, b *bundle.Bundle, tiles *ebiten.Image) draw.Node {
	return draw.ImageWithAnimation(tiles, b.Battletiles.Info.Animations[2*3+(y-1)], int(t.BehaviorState.ElapsedTime))
}

func (tb *NormalTileBehavior) CanEnter(t *Tile, e *Entity) bool {
	return true
}
func (tb *NormalTileBehavior) OnLeave(t *Tile, e *Entity, s *State) {}
func (tb *NormalTileBehavior) Flip()                                {}
func (tb *NormalTileBehavior) Step(t *Tile, s *State)               {}

type CrackedTileBehavior struct {
}

func (tb *CrackedTileBehavior) Clone() TileBehavior {
	return &CrackedTileBehavior{}
}

func (tb *CrackedTileBehavior) Appearance(t *Tile, y int, b *bundle.Bundle, tiles *ebiten.Image) draw.Node {
	return draw.ImageWithAnimation(tiles, b.Battletiles.Info.Animations[3*3+(y-1)], int(t.BehaviorState.ElapsedTime))
}

func (tb *CrackedTileBehavior) CanEnter(t *Tile, e *Entity) bool {
	return true
}

func (tb *CrackedTileBehavior) OnLeave(t *Tile, e *Entity, s *State) {
	if e.Traits.IgnoresTileEffects {
		return
	}
	// TODO: Play cracking sound.
	// TODO: Add returnToNormalTimeLeft
	t.ReplaceBehavior(&BrokenTileBehavior{}, s)
}

func (tb *CrackedTileBehavior) Flip() {}

func (tb *CrackedTileBehavior) Step(t *Tile, s *State) {}

type RoadTileBehavior struct {
	Direction Direction
}

func (tb *RoadTileBehavior) Clone() TileBehavior {
	return &RoadTileBehavior{tb.Direction}
}

func (tb *RoadTileBehavior) Appearance(t *Tile, y int, b *bundle.Bundle, tiles *ebiten.Image) draw.Node {
	var offset int
	switch tb.Direction {
	case DirectionUp:
		offset = 0
	case DirectionDown:
		offset = 1
	case DirectionLeft:
		offset = 2
	case DirectionRight:
		offset = 3
	}
	return draw.ImageWithAnimation(tiles, b.Battletiles.Info.Animations[(9+offset)*3+(y-1)], int(t.BehaviorState.ElapsedTime))
}

func (tb *RoadTileBehavior) Flip() {
	tb.Direction = tb.Direction.FlipH()
}

func (tb *RoadTileBehavior) CanEnter(t *Tile, e *Entity) bool {
	return true
}
func (tb *RoadTileBehavior) OnLeave(t *Tile, e *Entity, s *State) {
}
func (tb *RoadTileBehavior) Step(t *Tile, s *State) {
	for _, e := range s.Entities {
		if e.TilePos != t.TilePos {
			continue
		}

		if e.Traits.IgnoresTileEffects {
			return
		}

		if e.ForcedMovementState.ForcedMovement.Type == ForcedMovementTypeNone {
			var h Hit
			h.ForcedMovement = ForcedMovement{Type: ForcedMovementTypeSlide, Direction: tb.Direction}
			e.ApplyHit(h)
			// TODO: Play conveyor noise.
		}
	}
}

type IceTileBehavior struct {
	direction Direction
}

func (tb *IceTileBehavior) Clone() TileBehavior {
	return &IceTileBehavior{tb.direction}
}

func (tb *IceTileBehavior) Appearance(t *Tile, y int, b *bundle.Bundle, tiles *ebiten.Image) draw.Node {
	return draw.ImageWithAnimation(tiles, b.Battletiles.Info.Animations[7*3+(y-1)], int(t.BehaviorState.ElapsedTime))
}

func (tb *IceTileBehavior) Flip() {
}

func (tb *IceTileBehavior) CanEnter(t *Tile, e *Entity) bool {
	return true
}
func (tb *IceTileBehavior) OnLeave(t *Tile, e *Entity, s *State) {
}
func (tb *IceTileBehavior) Step(t *Tile, s *State) {
	if tb.direction == DirectionNone {
		for _, e := range s.Entities {
			if e.FutureTilePos != t.TilePos || e.TilePos == e.FutureTilePos {
				continue
			}

			ex, ey := e.TilePos.XY()
			x, y := e.FutureTilePos.XY()

			tb.direction = DirectionDXDY(x-ex, y-ey)
			break
		}
	}

	for _, e := range s.Entities {
		if e.TilePos != t.TilePos {
			continue
		}

		if e.Traits.IgnoresTileEffects {
			return
		}

		if e.ForcedMovementState.ForcedMovement.Type == ForcedMovementTypeNone && tb.direction != DirectionNone {
			var h Hit
			h.ForcedMovement = ForcedMovement{Type: ForcedMovementTypeSlide, Direction: tb.direction}
			e.ApplyHit(h)
			tb.direction = DirectionNone
		}
	}
}
