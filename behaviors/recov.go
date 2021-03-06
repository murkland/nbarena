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

func (eb *Recov) Clone() state.EntityBehavior {
	return &Recov{eb.HP}
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

		s.AttachSound(&state.Sound{
			Type: bundle.SoundTypeRecov,
		})
		s.AttachDecoration(&state.Decoration{
			Type:      bundle.DecorationTypeRecov,
			TilePos:   e.TilePos,
			Offset:    image.Point{0, 0},
			IsFlipped: e.IsFlipped,
		})

		e.ChipUseLockoutTimeLeft = 30

		e.SetBehaviorImmediate(&Idle{}, s)
	}
}

func (eb *Recov) Cleanup(e *state.Entity, s *state.State) {
}

func (eb *Recov) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	return nil
}
