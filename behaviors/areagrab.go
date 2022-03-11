package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/state/query"
)

type areaGrabBall struct {
}

func (eb *areaGrabBall) Flip() {
}

func (eb *areaGrabBall) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *areaGrabBall) Clone() state.EntityBehavior {
	return &areaGrabBall{}
}

func (eb *areaGrabBall) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	ballNode := &draw.OptionsNode{Layer: 7}

	if e.BehaviorState.ElapsedTime < 31 {
		frames := b.AreaGrabSprites.Animations[0].Frames
		ballNode.Children = append(ballNode.Children, draw.ImageWithFrame(b.AreaGrabSprites.Image, frames[int(e.BehaviorState.ElapsedTime)%len(frames)]))
	} else {
		frames := b.AreaGrabSprites.Animations[1].Frames
		ballNode.Children = append(ballNode.Children, draw.ImageWithFrame(b.AreaGrabSprites.Image, frames[int(e.BehaviorState.ElapsedTime)-31]))
	}

	return ballNode
}

func (eb *areaGrabBall) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == 30 {
		x, y := e.TilePos.XY()

		if x == 1 || x == state.TileCols-2 {
			return
		}

		tile := s.Field.Tiles[e.TilePos]

		if tile.Reserver == 0 {
			tile.IsAlliedWithAnswerer = e.IsAlliedWithAnswerer
			s.Field.ColumnInfo[x].AllySwapTimeLeft = 1800
		} else {
			entities := query.EntitiesAt(s, state.TilePosXY(x, y))
			_ = entities
			// TODO: Damage.
		}
	} else if e.BehaviorState.ElapsedTime == 30+15 {
		e.PerTickState.IsPendingDeletion = true
		return
	}
}
