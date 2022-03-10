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
	if e.BehaviorElapsedTime() == 9 {
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

		// TODO: Spawn gusts as well.
	} else if e.BehaviorElapsedTime() == 23 {
		e.SetBehavior(&Idle{}, s)
	}
}

func (eb *WindRack) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	megamanFrameIdx := int(e.BehaviorElapsedTime())
	if megamanFrameIdx >= len(b.MegamanSprites.SlashAnimation.Frames) {
		megamanFrameIdx = len(b.MegamanSprites.SlashAnimation.Frames) - 1
	}
	rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.SlashAnimation.Frames[megamanFrameIdx]))

	swordNode := &draw.OptionsNode{Layer: 6}
	rootNode.Children = append(rootNode.Children, swordNode)
	rackFrameIdx := int(e.BehaviorElapsedTime())
	if rackFrameIdx >= len(b.WindRackSprites.Animations[0].Frames) {
		rackFrameIdx = len(b.WindRackSprites.Animations[0].Frames) - 1
	}
	swordNode.Children = append(swordNode.Children, draw.ImageWithFrame(b.WindRackSprites.Image, b.WindRackSprites.Animations[0].Frames[rackFrameIdx]))

	slashNode := &draw.OptionsNode{Layer: 7}
	slashNode.Opts.GeoM.Translate(0, float64(-16))
	rootNode.Children = append(rootNode.Children, slashNode)
	slashFrameIdx := int(e.BehaviorElapsedTime())
	if slashFrameIdx >= len(b.WindSlashSprites.Animations[0].Frames) {
		slashFrameIdx = len(b.WindSlashSprites.Animations[0].Frames) - 1
	}
	slashNode.Children = append(slashNode.Children, draw.ImageWithFrame(b.WindSlashSprites.Image, b.WindSlashSprites.Animations[0].Frames[slashFrameIdx]))

	return rootNode
}
