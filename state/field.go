package state

import "github.com/yumland/clone"

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
