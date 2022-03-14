package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/state/query"
)

type AreaGrab struct {
}

func (eb *AreaGrab) Flip() {
}

func (eb *AreaGrab) Clone() state.EntityBehavior {
	return &AreaGrab{}
}

func (eb *AreaGrab) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *AreaGrab) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == 0 {
		s.IsInTimeStop = true
		e.RunsInTimestop = true

		xStart := 1
		xEnd := state.TileCols - 2
		xStep := 1

		if e.IsAlliedWithAnswerer {
			xStart, xEnd = xEnd, xStart
			xStep = -1
		}

		x := xStart
		for ; x != xEnd; x += xStep {
			for y := 1; y < 4; y++ {
				pos := state.TilePosXY(x, y)
				t := s.Field.Tiles[pos]
				if t.IsAlliedWithAnswerer != e.IsAlliedWithAnswerer {
					goto found
				}
			}
		}
	found:

		for y := 1; y < 4; y++ {
			s.AddEntity(&state.Entity{
				TilePos: state.TilePosXY(x, y),

				IsFlipped:            e.IsFlipped,
				IsAlliedWithAnswerer: e.IsAlliedWithAnswerer,

				Traits: state.EntityTraits{
					CanStepOnHoleLikeTiles: true,
					IgnoresTileEffects:     true,
					CannotFlinch:           true,
					IgnoresTileOwnership:   true,
					CannotSlide:            true,
					Intangible:             true,
				},

				BehaviorState: state.EntityBehaviorState{
					Behavior: &areaGrabBall{},
				},
			})
		}
	} else if e.BehaviorState.ElapsedTime == 40 {
		// TODO: Probably not 40!
		e.NextBehavior = &Idle{}
	}
}

func (eb *AreaGrab) Cleanup(e *state.Entity, s *state.State) {
	s.IsInTimeStop = false
	e.RunsInTimestop = false
}

func (eb *AreaGrab) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	return draw.ImageWithAnimation(b.MegamanSprites.Image, b.MegamanSprites.IdleAnimation, int(e.ElapsedTime()))
}

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

func (eb *areaGrabBall) Cleanup(e *state.Entity, s *state.State) {
}

func (eb *areaGrabBall) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	ballNode := &draw.OptionsNode{Layer: 7}

	if e.BehaviorState.ElapsedTime < 32 {
		frames := b.AreaGrabSprites.Animations[0].Frames
		ballNode.Opts.GeoM.Translate(0, float64(-9*(31-e.BehaviorState.ElapsedTime-1)))
		ballNode.Children = append(ballNode.Children, draw.ImageWithFrame(b.AreaGrabSprites.Image, frames[int(e.BehaviorState.ElapsedTime)%len(frames)]))
	} else {
		frames := b.AreaGrabSprites.Animations[1].Frames
		ballNode.Children = append(ballNode.Children, draw.ImageWithFrame(b.AreaGrabSprites.Image, frames[int(e.BehaviorState.ElapsedTime)-31]))
	}

	return ballNode
}

func (eb *areaGrabBall) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == 31 {
		x, y := e.TilePos.XY()

		if x == 1 || x == state.TileCols-2 {
			return
		}

		tile := s.Field.Tiles[e.TilePos]

		if tile.Reserver == 0 {
			tile.IsAlliedWithAnswerer = e.IsAlliedWithAnswerer
			s.Field.ColumnInfo[x].AllySwapTimeLeft = 1800
		} else {
			for _, entity := range query.HittableEntitiesAt(s, e, state.TilePosXY(x, y)) {
				var h state.Hit
				h.Traits.Flinch = true
				h.AddDamage(state.Damage{Base: 10})
				entity.Hit.Merge(h)
			}
		}
	} else if e.BehaviorState.ElapsedTime == 30+15 {
		e.IsPendingDestruction = true
		return
	}
}
