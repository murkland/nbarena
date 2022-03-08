package state

import (
	"github.com/murkland/clone"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
)

type ColumnInfo struct {
	allySwapTimeLeft int
}

func (c ColumnInfo) Clone() ColumnInfo {
	return ColumnInfo{c.allySwapTimeLeft}
}

type Field struct {
	Tiles      []Tile
	ColumnInfo []ColumnInfo
}

func (f Field) Clone() Field {
	return Field{clone.Slice(f.Tiles), clone.Slice(f.ColumnInfo)}
}

func newField() Field {
	tiles := make([]Tile, TileCols*TileRows)
	for j := 0; j < TileRows; j++ {
		for i := 0; i < TileCols; i++ {
			t := &tiles[int(TilePosXY(i, j))]
			if i >= 1 && i < TileCols-1 && j >= 1 && j < TileRows-1 {
				t.behavior = &NormalTileBehavior{}
			}
			t.IsAlliedWithAnswerer = i >= TileCols/2
			t.ShouldBeAlliedWithAnswerer = t.IsAlliedWithAnswerer
		}
	}
	return Field{tiles, make([]ColumnInfo, TileCols)}
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
	for j := range f.ColumnInfo {
		if f.ColumnInfo[j].allySwapTimeLeft > 0 {
			f.ColumnInfo[j].allySwapTimeLeft--
			if f.ColumnInfo[j].allySwapTimeLeft <= 0 {
				for i := 0; i < TileCols; i++ {
					t := &f.Tiles[int(TilePosXY(i, j))]
					t.ShouldBeAlliedWithAnswerer = !t.ShouldBeAlliedWithAnswerer
				}
			}
		}
	}

	for i := range f.Tiles {
		f.Tiles[i].Step()
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
		if x > 3 {
			tiles = b.Battletiles.AnswererTiles
		}
		childNode.Children = append(childNode.Children, draw.ImageWithFrame(tiles, frame))
		optsNode.Children = append(optsNode.Children, childNode)
	}

	return optsNode
}
