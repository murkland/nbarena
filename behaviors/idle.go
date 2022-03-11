package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type Idle struct {
	ChargingElapsedTime state.Ticks
}

func (eb *Idle) Flip() {
}

func (eb *Idle) Clone() state.EntityBehavior {
	return &Idle{eb.ChargingElapsedTime}
}

func (eb *Idle) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *Idle) Step(e *state.Entity, s *state.State) {
	if e.Intent.UseChip && e.LastIntent.UseChip != e.Intent.UseChip && e.ChipUseLockoutTimeLeft == 0 {
		// TODO: Add 1 frame delay between losing chip and switching to chip behavior.
		e.UseChip(s)
		return
	}

	if e.Intent.ChargeBasicWeapon {
		eb.ChargingElapsedTime++
	}

	if eb.ChargingElapsedTime > 0 && !e.Intent.ChargeBasicWeapon {
		// Release buster shot.
		e.NextBehavior = &Buster{BaseDamage: 1, IsPowerShot: eb.ChargingElapsedTime >= e.PowerShotChargeTime}
		eb.ChargingElapsedTime = 0
	}

	dir := e.Intent.Direction
	if e.ConfusedTimeLeft > 0 {
		dir = dir.FlipH().FlipV()
	}

	x, y := e.TilePos.XY()
	dx, dy := dir.XY()

	if e.StartMove(state.TilePosXY(x+dx, y+dy), s) {
		e.NextBehavior = &Teleport{ChargingElapsedTime: eb.ChargingElapsedTime}
	}
}

func (eb *Idle) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}

	rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.IdleAnimation.Frames[int(e.BehaviorState.ElapsedTime)%len(b.MegamanSprites.IdleAnimation.Frames)]))

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
