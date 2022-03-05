package step

import (
	"math/rand"

	"github.com/yumland/nbarena/behaviors"
	"github.com/yumland/nbarena/state"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
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

func Step(s *state.State) {
	s.ElapsedTime++

	// Mark all entities as pending step.
	for _, e := range s.Entities {
		e.IsPendingStep = true
	}

	// Step all entities in a random order.
	for {
		pending := make([]*state.Entity, 0, len(s.Entities))
		for _, e := range s.Entities {
			if !e.IsPendingStep {
				continue
			}
			pending = append(pending, e)
		}
		if len(pending) == 0 {
			break
		}

		slices.SortFunc(pending, func(a *state.Entity, b *state.Entity) bool {
			return a.ID() < b.ID()
		})
		rand.New(s.RandSource).Shuffle(len(pending), func(i, j int) {
			pending[i], pending[j] = pending[j], pending[i]
		})
		for _, e := range pending {
			e.Step(s)
			e.IsPendingStep = false
		}
	}

	// Resolve any hits.
	for _, e := range maps.Values(s.Entities) {
		resolveHit(e, e.CurrentHit)
		e.CurrentHit = state.Hit{}

		// Update UI.
		if e.DisplayHP != 0 && e.DisplayHP != e.HP {
			if e.HP == 0 {
				e.DisplayHP = 0
			} else {
				if e.HP < e.DisplayHP {
					e.DisplayHP -= ((e.DisplayHP-e.HP)>>3 + 4)
					if e.DisplayHP < e.HP {
						e.DisplayHP = e.HP
					}
				} else {
					e.DisplayHP += ((e.HP-e.DisplayHP)>>3 + 4)
					if e.DisplayHP > e.HP {
						e.DisplayHP = e.HP
					}
				}
			}
		}
	}

	// Delete any entities pending deletion.
	for k, e := range s.Entities {
		if e.IsPendingDeletion {
			delete(s.Entities, k)
		}
	}

	s.Field.Step(s)
}
