package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/pngsheet"
)

const teleportEndlagTicks = 8

type Teleport struct {
	useChip bool
}

func (eb *Teleport) Flip() {
}

func (eb *Teleport) Clone() state.EntityBehavior {
	return &Teleport{eb.useChip}
}

func (eb *Teleport) Step(e *state.Entity, s *state.State) {
	if e.Intent.UseChip && e.LastIntent.UseChip != e.Intent.UseChip {
		eb.useChip = true
	}

	if e.Intent.ChargeBasicWeapon {
		e.ChargingElapsedTime++
	}

	if e.BehaviorState.ElapsedTime == 3 {
		e.FinishMove(s)
	}

	if e.BehaviorState.ElapsedTime == 6+teleportEndlagTicks {
		e.ReplaceBehavior(&Idle{}, s)
		if eb.useChip && e.ChipUseLockoutTimeLeft == 0 {
			e.UseChip(s)
		}
	}
}

func (eb *Teleport) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	var frame *pngsheet.Frame
	if e.BehaviorState.ElapsedTime < 3 {
		frame = b.MegamanSprites.TeleportStartAnimation.Frames[e.BehaviorState.ElapsedTime]
	} else if e.BehaviorState.ElapsedTime < 6 {
		frame = b.MegamanSprites.TeleportEndAnimation.Frames[e.BehaviorState.ElapsedTime-3]
	} else {
		frame = b.MegamanSprites.TeleportEndAnimation.Frames[len(b.MegamanSprites.TeleportEndAnimation.Frames)-1]
	}
	return draw.ImageWithFrame(b.MegamanSprites.Image, frame)
}
