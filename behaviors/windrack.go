package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/state/query"
)

type WindRack struct {
	Damage int
}

func (eb *WindRack) Flip() {
}

func (eb *WindRack) Clone() state.EntityBehavior {
	return &WindRack{
		eb.Damage,
	}
}

func (eb *WindRack) Step(e *state.Entity, s *state.State) {
	// Only hits while the slash is coming out.
	if e.BehaviorState.ElapsedTime == 9 {
		x, y := e.TilePos.XY()
		dx := query.DXForward(e.IsFlipped)

		var entities []*state.Entity
		entities = append(entities, query.EntitiesAt(s, state.TilePosXY(x+dx, y))...)
		entities = append(entities, query.EntitiesAt(s, state.TilePosXY(x+dx, y+1))...)
		entities = append(entities, query.EntitiesAt(s, state.TilePosXY(x+dx, y-1))...)

		for _, target := range entities {
			if target.FlashingTimeLeft == 0 {
				var h state.Hit
				h.Counters = true
				h.Drag = true
				h.Slide.Direction = e.Facing()
				h.Slide.IsBig = true
				h.AddDamage(e.MakeDamageAndConsume(eb.Damage))
				target.Hit.Merge(h)
			}
		}

		for i := 1; i <= 3; i++ {
			shot := &state.Entity{
				TilePos: state.TilePosXY(x+dx, i),

				IsFlipped:            e.IsFlipped,
				IsAlliedWithAnswerer: e.IsAlliedWithAnswerer,

				Traits: state.EntityTraits{
					CanStepOnHoleLikeTiles: true,
					IgnoresTileEffects:     true,
					CannotFlinch:           true,
					IgnoresTileOwnership:   true,
				},
			}
			shot.SetBehavior(&windRackGust{e.Facing()}, s)
			s.AddEntity(shot)
		}
	} else if e.BehaviorState.ElapsedTime == 27 {
		e.SetBehavior(&Idle{}, s)
	}
}

func (eb *WindRack) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	megamanFrameIdx := int(e.BehaviorState.ElapsedTime)
	if megamanFrameIdx >= len(b.MegamanSprites.SlashAnimation.Frames) {
		megamanFrameIdx = len(b.MegamanSprites.SlashAnimation.Frames) - 1
	}
	rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.SlashAnimation.Frames[megamanFrameIdx]))

	swordNode := &draw.OptionsNode{Layer: 6}
	rootNode.Children = append(rootNode.Children, swordNode)
	rackFrameIdx := int(e.BehaviorState.ElapsedTime)
	if rackFrameIdx >= len(b.WindRackSprites.Animations[0].Frames) {
		rackFrameIdx = len(b.WindRackSprites.Animations[0].Frames) - 1
	}
	swordNode.Children = append(swordNode.Children, draw.ImageWithFrame(b.WindRackSprites.Image, b.WindRackSprites.Animations[0].Frames[rackFrameIdx]))

	slashNode := &draw.OptionsNode{Layer: 7}
	slashNode.Opts.GeoM.Translate(0, float64(-16))
	rootNode.Children = append(rootNode.Children, slashNode)
	slashFrameIdx := int(e.BehaviorState.ElapsedTime)
	if slashFrameIdx >= len(b.WindSlashSprites.Animations[0].Frames) {
		slashFrameIdx = len(b.WindSlashSprites.Animations[0].Frames) - 1
	}
	slashNode.Children = append(slashNode.Children, draw.ImageWithFrame(b.WindSlashSprites.Image, b.WindSlashSprites.Animations[0].Frames[slashFrameIdx]))

	return rootNode
}

type windRackGust struct {
	direction state.Direction
}

func (eb *windRackGust) Flip() {
}

func (eb *windRackGust) Clone() state.EntityBehavior {
	return &windRackGust{
		eb.direction,
	}
}

func (eb *windRackGust) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	return nil
}

func (eb *windRackGust) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime%2 == 1 {
		x, y := e.TilePos.XY()
		x += query.DXForward(e.IsFlipped)
		if !e.MoveDirectly(state.TilePosXY(x, y)) {
			e.PerTickState.IsPendingDeletion = true
			return
		}
	}

	for _, target := range query.EntitiesAt(s, e.TilePos) {
		if target.IsAlliedWithAnswerer == e.IsAlliedWithAnswerer {
			continue
		}

		var h state.Hit
		h.Slide.Direction = eb.direction
		h.Slide.IsBig = true
		target.Hit.Merge(h)

		// Gust doesn't stop on hit.
		return
	}
}
