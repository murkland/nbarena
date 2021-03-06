package query

import "github.com/murkland/nbarena/state"

type DistanceMetric func(src state.TilePos, dest state.TilePos) int

func IsInFrontOf(x int, targetX int, isFlipped bool) bool {
	if isFlipped {
		return targetX < x
	}
	return targetX > x
}

func HorizontalDistance(src state.TilePos, dest state.TilePos) int {
	x1, _ := src.XY()
	x2, _ := dest.XY()
	if x1 > x2 {
		return x1 - x2
	}
	return x2 - x1
}

func FindNearestEntity(s *state.State, myEntityID state.EntityID, pos state.TilePos, isAlliedWithAnswerer bool, isFlipped bool, distance DistanceMetric) (state.EntityID, int) {
	x, _ := pos.XY()

	bestDist := state.TileCols

	var targetID state.EntityID
	for _, cand := range s.Entities {
		if cand.ID() == myEntityID || cand.IsAlliedWithAnswerer == isAlliedWithAnswerer {
			continue
		}

		candX, _ := cand.FutureTilePos.XY()

		if !IsInFrontOf(x, candX, isFlipped) {
			continue
		}

		if d := distance(pos, cand.FutureTilePos); d >= 0 && d < bestDist {
			targetID = cand.ID()
			bestDist = d
		}
	}

	return targetID, bestDist
}
