package state

import (
	"github.com/yumland/clone"
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
	return Field{make([]Tile, tileCols*tileRows), make([]ColumnInfo, tileCols)}
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

func (f *Field) DrawNode() draw.Node {
	return draw.OptionsNode{}
}
