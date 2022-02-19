package state

import "github.com/yumland/yumbattle/draw"

type Entity struct {
	Appearance draw.Node

	TilePos TilePos

	HP        int
	DisplayHP *int

	CanStepOnHoleLikeTiles bool
	IgnoresTileEffects     bool
}
