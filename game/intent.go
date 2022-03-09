package game

import (
	"math/rand"

	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/input"
	"github.com/murkland/nbarena/state"
)

func applyPlayerIntent(s *state.State, e *state.Entity, intent input.Intent, isOfferer bool) {
	if e.LastIntent.UseChip != intent.UseChip && intent.UseChip && e.ChipUseLockoutTimeLeft == 0 {
		if e.PerTickState.Interrupts.WithChipUse != state.WithChipUseInterruptTypeIgnore {
			e.ChipUseQueued = true
		}
	}

	if intent.ChargeBasicWeapon && (e.PerTickState.Interrupts.WithCharge || e.ChargingElapsedTime > 0) {
		e.ChargingElapsedTime++
	}

	if !intent.ChargeBasicWeapon && e.PerTickState.Interrupts.WithCharge && e.ChargingElapsedTime > 0 {
		// Release buster shot.
		e.SetBehavior(&behaviors.Buster{BaseDamage: 1, IsPowerShot: e.ChargingElapsedTime >= e.PowerShotChargeTime})
		e.ChargingElapsedTime = 0
	}

	if e.PerTickState.Interrupts.WithMove {
		dir := intent.Direction
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

		if e.StartMove(state.TilePosXY(x, y), &s.Field) {
			e.SetBehavior(&behaviors.Teleport{})
		}
	}
}

func applyPlayerIntents(s *state.State, offererEntityID int, offererIntent input.Intent, answererEntityID int, answererIntent input.Intent) {
	intents := []struct {
		isOfferer bool
		intent    input.Intent
	}{
		{true, offererIntent},
		{false, answererIntent},
	}
	rand.New(s.RandSource).Shuffle(len(intents), func(i, j int) {
		intents[i], intents[j] = intents[j], intents[i]
	})
	for _, wrapped := range intents {
		var entity *state.Entity
		if wrapped.isOfferer {
			entity = s.Entities[offererEntityID]
		} else {
			entity = s.Entities[answererEntityID]
		}
		applyPlayerIntent(s, entity, wrapped.intent, wrapped.isOfferer)
		entity.LastIntent = wrapped.intent
	}
}
