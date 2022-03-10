package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/input"
	"github.com/murkland/nbarena/state"
)

type Idle struct {
}

func (eb *Idle) Clone() state.EntityBehavior {
	return &Idle{}
}

func (eb *Idle) Step(e *state.Entity, s *state.State) {
	if e.Intent.UseChip && e.LastIntent.UseChip != e.Intent.UseChip && e.ChipUseLockoutTimeLeft == 0 {
		e.UseChip(s)
		return
	}

	if e.Intent.ChargeBasicWeapon {
		e.ChargingElapsedTime++
	}

	if e.ChargingElapsedTime > 0 && !e.Intent.ChargeBasicWeapon {
		// Release buster shot.
		e.SetBehavior(&Buster{BaseDamage: 1, IsPowerShot: e.ChargingElapsedTime >= e.PowerShotChargeTime})
		e.ChargingElapsedTime = 0
	}

	dir := e.Intent.Direction
	if e.ConfusedTimeLeft > 0 {
		dir = dir.FlipH().FlipV()
	}

	x, y := e.TilePos.XY()
	if dir&input.DirectionLeft != 0 {
		x--
	}
	if dir&input.DirectionRight != 0 {
		x++
	}
	if dir&input.DirectionUp != 0 {
		y--
	}
	if dir&input.DirectionDown != 0 {
		y++
	}

	if e.StartMove(state.TilePosXY(x, y), s.Field) {
		e.SetBehavior(&Teleport{})
	}
}

func (eb *Idle) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	frames := b.MegamanSprites.IdleAnimation.Frames
	frame := frames[int(e.BehaviorElapsedTime())%len(frames)]
	return draw.ImageWithFrame(b.MegamanSprites.Image, frame)
}
