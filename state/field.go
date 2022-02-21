package state

import (
	"github.com/yumland/clone"
	"github.com/yumland/yumbattle/draw"
)

type ColumnInfo struct {
	ownerSwapTimeLeft int
}

func (c ColumnInfo) Clone() ColumnInfo {
	return ColumnInfo{c.ownerSwapTimeLeft}
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
		if f.columnInfo[j].ownerSwapTimeLeft > 0 {
			f.columnInfo[j].ownerSwapTimeLeft--
			if f.columnInfo[j].ownerSwapTimeLeft <= 0 {
				for i := 0; i < tileCols; i++ {
					t := &f.tiles[int(TilePosXY(i, j))]
					// TODO: Check if tile is occupied: if occupied, do not switch owner.
					t.isOwnedByAnswerer = !t.isOwnedByAnswerer
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
