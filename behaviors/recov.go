package behaviors

import (
	"image"

	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type Recov struct {
	HP int
}

func (eb *Recov) Flip() {
}

func (eb *Recov) Clone() state.EntityBehavior {
	return &Recov{}
}

func (eb *Recov) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *Recov) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == 0 {
		e.HP += eb.HP
		if e.HP > e.MaxHP {
			e.HP = e.MaxHP
		}

		// TODO: Trigger antirecovery.

		s.AddDecoration(&state.Decoration{
			Type:      bundle.DecorationTypeRecov,
			TilePos:   e.TilePos,
			Offset:    image.Point{0, 0},
			IsFlipped: e.IsFlipped,
		})

		e.NextBehavior = &Idle{}
	}
}

func (eb *Recov) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	return draw.ImageWithAnimation(b.MegamanSprites.Image, b.MegamanSprites.IdleAnimation, int(e.ElapsedTime()))
}
