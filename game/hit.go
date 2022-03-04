package game

import (
	"github.com/yumland/yumbattle/behaviors"
	"github.com/yumland/yumbattle/state"
)

func resolveHit(e *state.Entity, hit state.Hit) {
	// Set anger, if required.
	if hit.TotalDamage >= 300 {
		e.IsAngry = true
	}

	// TODO: Process poison damage.

	// Process hit damage.
	mustLeave1HP := e.HP > 1 && e.Traits.FatalHitLeaves1HP
	e.HP -= hit.TotalDamage
	if e.HP < 0 {
		e.HP = 0
	}
	if mustLeave1HP {
		e.HP = 1
	}
	hit.TotalDamage = 0

	if !hit.Drag {
		if !e.IsBeingDragged /* && !e.isInTimestop */ {
			// Process flashing.
			if hit.FlashTime > 0 {
				e.FlashingTimeLeft = hit.FlashTime
				hit.FlashTime = 0
			}
			if e.FlashingTimeLeft > 0 {
				e.FlashingTimeLeft--
			}

			// Process paralyzed.
			if hit.ParalyzeTime > 0 {
				e.ParalyzedTimeLeft = hit.ParalyzeTime
				hit.ConfuseTime = 0
				hit.ParalyzeTime = 0
			}
			if e.ParalyzedTimeLeft > 0 {
				e.ParalyzedTimeLeft--
				e.FrozenTimeLeft = 0
				e.BubbledTimeLeft = 0
				e.ConfusedTimeLeft = 0
			}

			// Process frozen.
			if hit.FreezeTime > 0 {
				e.FrozenTimeLeft = hit.FreezeTime
				e.ParalyzedTimeLeft = 0
				hit.BubbleTime = 0
				hit.ConfuseTime = 0
				hit.FreezeTime = 0
			}
			if e.FrozenTimeLeft > 0 {
				e.FrozenTimeLeft--
				e.BubbledTimeLeft = 0
				e.ConfusedTimeLeft = 0
			}

			// Process bubbled.
			if hit.BubbleTime > 0 {
				e.BubbledTimeLeft = hit.BubbleTime
				e.ConfusedTimeLeft = 0
				e.ParalyzedTimeLeft = 0
				e.FrozenTimeLeft = 0
				hit.ConfuseTime = 0
				hit.BubbleTime = 0
			}
			if e.BubbledTimeLeft > 0 {
				e.BubbledTimeLeft--
				e.ConfusedTimeLeft = 0
			}

			// Process confused.
			if hit.ConfuseTime > 0 {
				e.ConfusedTimeLeft = hit.ConfuseTime
				e.ParalyzedTimeLeft = 0
				e.FrozenTimeLeft = 0
				e.BubbledTimeLeft = 0
				hit.FreezeTime = 0
				hit.BubbleTime = 0
				hit.ParalyzeTime = 0
				hit.ConfuseTime = 0
			}
			if e.ConfusedTimeLeft > 0 {
				e.ConfusedTimeLeft--
			}

			// Process immobilized.
			if hit.ImmobilizeTime > 0 {
				e.ImmobilizedTimeLeft = hit.ImmobilizeTime
				hit.ImmobilizeTime = 0
			}
			if e.ImmobilizedTimeLeft > 0 {
				e.ImmobilizedTimeLeft--
			}

			// Process blinded.
			if hit.BlindTime > 0 {
				e.BlindedTimeLeft = hit.BlindTime
				hit.BlindTime = 0
			}
			if e.BlindedTimeLeft > 0 {
				e.BlindedTimeLeft--
			}

			// Process invincible.
			if e.InvincibleTimeLeft > 0 {
				e.InvincibleTimeLeft--
			}
		} else {
			// TODO: Interrupt player.
		}
	} else {
		hit.Drag = false

		e.FrozenTimeLeft = 0
		e.BubbledTimeLeft = 0
		e.ParalyzedTimeLeft = 0
		hit.BubbleTime = 0
		hit.FreezeTime = 0

		if false {
			e.ParalyzedTimeLeft = 0
		}

		// TODO: Interrupt player.
	}

	if hit.Flinch && !e.Traits.CannotFlinch {
		e.SetBehavior(&behaviors.Flinch{})
	}
	hit.Flinch = false
}
