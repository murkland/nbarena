package state

import (
	"github.com/yumland/clone"
	"github.com/yumland/yumbattle/bundle"
	"github.com/yumland/yumbattle/draw"
)

type ColumnInfo struct {
	allySwapTimeLeft int
}

func (c ColumnInfo) Clone() ColumnInfo {
	return ColumnInfo{c.allySwapTimeLeft}
}

type Field struct {
	tiles      []Tile
	columnInfo []ColumnInfo
}

func (f Field) Clone() Field {
	return Field{clone.Slice(f.tiles), clone.Slice(f.columnInfo)}
}

func newField() Field {
	tiles := make([]Tile, tileCols*tileRows)
	for j := 0; j < 5; j++ {
		for i := 0; i < 8; i++ {
			t := &tiles[int(TilePosXY(i, j))]
			if i >= 1 && i < 7 && j >= 1 && j < 4 {
				t.behavior = &NormalTileBehavior{}
			}
			t.isAlliedWithAnswerer = i >= 4
		}
	}
	return Field{tiles, make([]ColumnInfo, tileCols)}
}

func (f *Field) Step() {
	for j := range f.columnInfo {
		if f.columnInfo[j].allySwapTimeLeft > 0 {
			f.columnInfo[j].allySwapTimeLeft--
			if f.columnInfo[j].allySwapTimeLeft <= 0 {
				for i := 0; i < tileCols; i++ {
					t := &f.tiles[int(TilePosXY(i, j))]
					// TODO: Check if tile is occupied: if occupied, do not switch ally.
					t.isAlliedWithAnswerer = !t.isAlliedWithAnswerer
				}
			}
		}
	}

	for i := range f.tiles {
		f.tiles[i].Step()
	}
}

const (
	tileRenderedWidth  = 40
	tileRenderedHeight = 24
)

func (f *Field) Appearance(b *bundle.Bundle) draw.Node {
	optsNode := draw.OptionsNode{}
	for i, tile := range f.tiles {
		x, y := TilePos(i).XY()
		node := tile.Appearance(y, b)
		if node == nil {
			continue
		}

		childNode := draw.OptionsNode{}
		childNode.Opts.GeoM.Translate(float64((x-1)*tileRenderedWidth), float64((y-1)*tileRenderedHeight))
		childNode.Children = append(childNode.Children, node)
		optsNode.Children = append(optsNode.Children, childNode)
	}
	return optsNode
}
