package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/pngsheet"
)

const teleportEndlagTicks = 6

type Teleport struct {
	ChargingElapsedTime state.Ticks
	useChip             bool
}

func (eb *Teleport) Flip() {
}

func (eb *Teleport) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *Teleport) Clone() state.EntityBehavior {
	return &Teleport{eb.ChargingElapsedTime, eb.useChip}
}

func (eb *Teleport) Step(e *state.Entity, s *state.State) {
	if e.Intent.UseChip && e.LastIntent.UseChip != e.Intent.UseChip {
		eb.useChip = true
	}

	if e.Intent.ChargeBasicWeapon {
		eb.ChargingElapsedTime++
	}

	if e.BehaviorState.ElapsedTime == 3 {
		e.FinishMove(s)
	}

	if e.BehaviorState.ElapsedTime == 6+teleportEndlagTicks-1 {
		e.NextBehavior = &Idle{eb.ChargingElapsedTime}
		if eb.useChip && e.ChipUseLockoutTimeLeft == 0 {
			UseNextChip(e)
		}
	}
}

func (eb *Teleport) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}

	var frame *pngsheet.Frame
	if e.BehaviorState.ElapsedTime < 3 {
		frame = b.MegamanSprites.TeleportStartAnimation.Frames[e.BehaviorState.ElapsedTime]
	} else if e.BehaviorState.ElapsedTime < 6 {
		frame = b.MegamanSprites.TeleportEndAnimation.Frames[e.BehaviorState.ElapsedTime-3]
	} else {
		frame = b.MegamanSprites.TeleportEndAnimation.Frames[len(b.MegamanSprites.TeleportEndAnimation.Frames)-1]
	}
	rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.MegamanSprites.Image, frame))

	if eb.ChargingElapsedTime >= 10 {
		chargingNode := &draw.OptionsNode{}
		rootNode.Children = append(rootNode.Children, chargingNode)

		frames := b.ChargingSprites.ChargingAnimation.Frames
		if eb.ChargingElapsedTime >= e.PowerShotChargeTime {
			frames = b.ChargingSprites.ChargedAnimation.Frames
		}
		frame := frames[int(eb.ChargingElapsedTime)%len(frames)]
		chargingNode.Children = append(chargingNode.Children, draw.ImageWithFrame(b.ChargingSprites.Image, frame))
	}

	return rootNode
}
