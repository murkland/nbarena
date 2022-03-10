package state

import (
	"github.com/murkland/clone"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
)

type ColumnInfo struct {
	allySwapTimeLeft int
}

func (c *ColumnInfo) Clone() *ColumnInfo {
	return &ColumnInfo{c.allySwapTimeLeft}
}

type Field struct {
	Tiles      []*Tile
	ColumnInfo []*ColumnInfo
}

func (f *Field) Clone() *Field {
	return &Field{clone.Slice(f.Tiles), clone.Slice(f.ColumnInfo)}
}

func newField() *Field {
	tiles := make([]*Tile, TileCols*TileRows)
	for j := 0; j < TileRows; j++ {
		for i := 0; i < TileCols; i++ {
			t := &Tile{}
			if i >= 1 && i < TileCols-1 && j >= 1 && j < TileRows-1 {
				t.SetBehavior(&NormalTileBehavior{})
			}
			t.IsAlliedWithAnswerer = i >= TileCols/2
			t.ShouldBeAlliedWithAnswerer = t.IsAlliedWithAnswerer
			tiles[j*TileCols+i] = t
		}
	}

	columnInfos := make([]*ColumnInfo, TileCols)
	for i := 0; i < TileCols; i++ {
		columnInfos[i] = &ColumnInfo{}
	}

	return &Field{tiles, columnInfos}
}

func (f *Field) Flip() {
	for j := 0; j < TileRows; j++ {
		for i := 0; i < TileCols/2; i++ {
			pos := TilePosXY(i, j)
			newPos := pos.Flipped()
			f.Tiles[pos], f.Tiles[newPos] = f.Tiles[newPos], f.Tiles[pos]
		}
	}

	for i := range f.Tiles {
		f.Tiles[i].Flip()
	}

	for i := 0; i < len(f.ColumnInfo)/2; i++ {
		j := len(f.ColumnInfo) - i - 1
		f.ColumnInfo[i], f.ColumnInfo[j] = f.ColumnInfo[j], f.ColumnInfo[i]
	}
}

func (f *Field) Step(s *State) {
	for i, ci := range f.ColumnInfo {
		if ci.allySwapTimeLeft > 0 {
			ci.allySwapTimeLeft--
			if ci.allySwapTimeLeft <= 0 {
				for j := 0; j < TileRows; j++ {
					t := f.Tiles[int(TilePosXY(i, j))]
					t.ShouldBeAlliedWithAnswerer = !t.ShouldBeAlliedWithAnswerer
				}
				ci.allySwapTimeLeft = 0
			}
		}
	}

	columnOccupiers := make([]bool, len(f.ColumnInfo))
	for _, entity := range s.Entities {
		if entity.Traits.ExtendsTileLifetime {
			{
				x, _ := entity.TilePos.XY()
				columnOccupiers[x] = true
			}
			{
				x, _ := entity.FutureTilePos.XY()
				columnOccupiers[x] = true
			}
		}
	}

	for i := range f.Tiles {
		x, _ := TilePos(i).XY()

		t := f.Tiles[i]
		if t.ShouldBeAlliedWithAnswerer != t.IsAlliedWithAnswerer && !columnOccupiers[x] {
			t.IsAlliedWithAnswerer = t.ShouldBeAlliedWithAnswerer
		}
		t.Step()
	}
}

const (
	TileRenderedWidth  = 40
	TileRenderedHeight = 24
)

func (f *Field) Appearance(b *bundle.Bundle) draw.Node {
	optsNode := &draw.OptionsNode{}
	for i, tile := range f.Tiles {
		x, y := TilePos(i).XY()
		node := tile.Appearance(y, b)
		if node == nil {
			continue
		}

		childNode := &draw.OptionsNode{}
		childNode.Opts.GeoM.Translate(float64((x-1)*TileRenderedWidth), float64((y-1)*TileRenderedHeight))
		childNode.Children = append(childNode.Children, node)
		optsNode.Children = append(optsNode.Children, childNode)
	}

	for x := 1; x < 7; x++ {
		childNode := &draw.OptionsNode{}
		childNode.Opts.GeoM.Translate(float64((x-1)*TileRenderedWidth), float64((4-1)*TileRenderedHeight))
		frame := b.Battletiles.Info.Animations[len(b.Battletiles.Info.Animations)-1].Frames[0]
		tiles := b.Battletiles.OffererTiles
		if f.Tiles[TilePosXY(x, 3)].IsAlliedWithAnswerer {
			tiles = b.Battletiles.AnswererTiles
		}
		childNode.Children = append(childNode.Children, draw.ImageWithFrame(tiles, frame))
		optsNode.Children = append(optsNode.Children, childNode)
	}

	return optsNode
}
