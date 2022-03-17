package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type AreaGrab struct {
	Owner state.EntityID
}

func (tb *AreaGrab) Clone() state.TimestopBehavior {
	return &AreaGrab{}
}

func (tb *AreaGrab) Step(t *state.Timestop, s *state.State) {
	owner := s.Entities[tb.Owner]

	if t.BehaviorElapsedTime == 0 {
		xStart := 1
		xEnd := state.TileCols - 2
		xStep := 1

		if owner.IsAlliedWithAnswerer {
			xStart, xEnd = xEnd, xStart
			xStep = -1
		}

		x := xStart
		for ; x != xEnd; x += xStep {
			for y := 1; y < 4; y++ {
				pos := state.TilePosXY(x, y)
				t := s.Field.Tiles[pos]
				if t.IsAlliedWithAnswerer != owner.IsAlliedWithAnswerer {
					goto found
				}
			}
		}
	found:
		for y := 1; y < 4; y++ {
			s.AttachEntity(&state.Entity{
				TilePos: state.TilePosXY(x, y),

				RunsInTimestop: true,

				IsAlliedWithAnswerer: owner.IsAlliedWithAnswerer,

				Traits: state.EntityTraits{
					CanStepOnHoleLikeTiles: true,
					IgnoresTileEffects:     true,
					CannotFlinch:           true,
					IgnoresTileOwnership:   true,
					CannotSlide:            true,
					Intangible:             true,
				},

				BehaviorState: state.EntityBehaviorState{
					Behavior: &areaGrabBall{owner.ID()},
				},
			})
		}
	} else if t.BehaviorElapsedTime == 40 {
		// TODO: Probably not 40!
		t.IsPendingDestruction = true
	}
}

type areaGrabBall struct {
	Owner state.EntityID
}

func (eb *areaGrabBall) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *areaGrabBall) Clone() state.EntityBehavior {
	return &areaGrabBall{eb.Owner}
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
	if e.BehaviorState.ElapsedTime == 0 {
		s.AttachSound(&state.Sound{Type: bundle.SoundTypeAreaGrabStart})
	} else if e.BehaviorState.ElapsedTime == 31 {
		s.AttachSound(&state.Sound{Type: bundle.SoundTypeAreaGrabEnd})
		x, _ := e.TilePos.XY()
		if x == 1 || x == state.TileCols-2 {
			return
		}

		tile := s.Field.Tiles[e.TilePos]

		var h state.Hit
		h.Flinch = true
		h.AddDamage(state.Damage{Base: 10})
		h.RemovesFullSynchro = true
		if !s.ApplyHit(s.Entities[eb.Owner], e.TilePos, h) {
			tile.IsAlliedWithAnswerer = e.IsAlliedWithAnswerer
			s.Field.ColumnInfo[x].AllySwapTimeLeft = 1800
		}
	} else if e.BehaviorState.ElapsedTime == 30+15 {
		e.IsPendingDestruction = true
		return
	}
}
