package step

import (
	"math/rand"

	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/state"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func resolveOne(e *state.Entity, s *state.State) {
	if e.Traits.CannotFlinch || e.Emotion == state.EmotionAngry {
		// TODO: Double check if this
		e.HitResolution.Flinch = false
	}

	if e.Traits.CannotSlide {
		e.HitResolution.ForcedMovement = state.ForcedMovement{}
	}

	if e.HitResolution.RemovesFullSynchro && e.Emotion == state.EmotionFullSynchro {
		e.Emotion = state.EmotionNormal
	}

	// Set anger, if required.
	if e.HitResolution.Damage >= 300 {
		e.Emotion = state.EmotionAngry
	}

	// TODO: Process poison damage.

	// Process hit damage.
	// TODO: Should this be in ApplyHit?
	if e.HitResolution.Damage > 0 {
		e.PerTickState.WasHit = true
		s.AttachSound(&state.Sound{
			Type: bundle.SoundTypeOuch,
		})
	}

	mustLeave1HP := e.HP > 1 && e.Traits.FatalHitLeaves1HP
	e.HP -= e.HitResolution.Damage
	if e.HP < 0 {
		e.HP = 0
	}
	if mustLeave1HP {
		e.HP = 1
	}
	e.HitResolution.Damage = 0

	// TODO: Pop bubble, if required.

	if e.ForcedMovementState.ForcedMovement.Type != state.ForcedMovementTypeNone {
		// TODO: Is this even in the right place?
		e.ForcedMovementState.ElapsedTime++
	}

	if !s.IsInTimeStop {
		if e.DragLockoutTimeLeft > 0 {
			e.DragLockoutTimeLeft--
		}

		if e.HitResolution.ForcedMovement.Type.IsDrag() || e.DragLockoutTimeLeft > 0 {
			if e.HitResolution.FlashTime != 0 && state.BehaviorIs[*behaviors.Paralyzed](e.BehaviorState.Behavior) {
				e.SetBehaviorImmediate(&behaviors.Idle{}, s)
			}

			if state.BehaviorIs[*behaviors.Idle](e.BehaviorState.Behavior) {
				if e.HitResolution.Flinch {
					e.SetBehaviorImmediate(&behaviors.Flinch{}, s)
				}
			}

			e.HitResolution.Flinch = false

			e.ForcedMovementState = state.ForcedMovementState{ForcedMovement: e.HitResolution.ForcedMovement}
			e.HitResolution.ForcedMovement = state.ForcedMovement{}
			resolveSlideOrDrag(e, s)
		} else {
			if !e.ForcedMovementState.ForcedMovement.Type.IsDrag() {
				if e.HitResolution.ForcedMovement.Type == state.ForcedMovementTypeSlide {
					// HACK: Allow immediate application of slide if the last slide is ending.
					if e.ForcedMovementState.ForcedMovement.Type == state.ForcedMovementTypeNone || e.ForcedMovementState.ElapsedTime == 4 {
						e.ForcedMovementState = state.ForcedMovementState{ForcedMovement: e.HitResolution.ForcedMovement}
					}
					resolveSlideOrDrag(e, s)
				} else {
					resolveSlideOrDrag(e, s)
				}
				e.HitResolution.ForcedMovement = state.ForcedMovement{}

				if e.HitResolution.Flinch {
					if state.BehaviorIs[*behaviors.Paralyzed](e.BehaviorState.Behavior) && e.HitResolution.FlashTime == 0 {
						e.HitResolution.Flinch = false
					}
				}
				if e.HitResolution.Flinch {
					e.SetBehaviorImmediate(&behaviors.Flinch{}, s)
				}
				e.HitResolution.Flinch = false

				// Process flashing.
				if e.HitResolution.FlashTime > 0 {
					e.Flashing = state.Flashing{TimeLeft: e.HitResolution.FlashTime}
					e.HitResolution.FlashTime = 0
				}
				if e.Flashing.TimeLeft > 0 {
					e.Flashing.TimeLeft--
				} else {
					e.RemoveFlashing(s)
				}

				// Process paralyzed.
				if e.HitResolution.ParalyzeTime > 0 {
					e.SetBehaviorImmediate(&behaviors.Paralyzed{Duration: e.HitResolution.ParalyzeTime}, s)
					e.HitResolution.ConfuseTime = 0
					e.HitResolution.ParalyzeTime = 0
				}

				// Process frozen.
				if e.HitResolution.FreezeTime > 0 {
					e.SetBehaviorImmediate(&behaviors.Frozen{Duration: e.HitResolution.FreezeTime}, s)
					e.HitResolution.BubbleTime = 0
					e.HitResolution.ConfuseTime = 0
					e.HitResolution.FreezeTime = 0
				}

				// Process bubbled.
				if e.HitResolution.BubbleTime > 0 {
					e.SetBehaviorImmediate(&behaviors.Bubbled{Duration: e.HitResolution.BubbleTime}, s)
					e.ConfusedTimeLeft = 0
					e.HitResolution.ConfuseTime = 0
					e.HitResolution.BubbleTime = 0
				}

				// Process confused.
				if e.HitResolution.ConfuseTime > 0 {
					e.ConfusedTimeLeft = e.HitResolution.ConfuseTime
					// TODO: Double check if this is correct.
					if state.BehaviorIs[*behaviors.Paralyzed](e.BehaviorState.Behavior) ||
						state.BehaviorIs[*behaviors.Frozen](e.BehaviorState.Behavior) ||
						state.BehaviorIs[*behaviors.Bubbled](e.BehaviorState.Behavior) {
						e.SetBehaviorImmediate(&behaviors.Idle{}, s)
					}
					e.HitResolution.FreezeTime = 0
					e.HitResolution.BubbleTime = 0
					e.HitResolution.ParalyzeTime = 0
					e.HitResolution.ConfuseTime = 0
				}
				if e.ConfusedTimeLeft > 0 {
					e.ConfusedTimeLeft--
				}

				// Process immobilized.
				if e.HitResolution.ImmobilizeTime > 0 {
					e.ImmobilizedTimeLeft = e.HitResolution.ImmobilizeTime
					e.HitResolution.ImmobilizeTime = 0
				}
				if e.ImmobilizedTimeLeft > 0 {
					e.ImmobilizedTimeLeft--
				}

				// Process blinded.
				if e.HitResolution.BlindTime > 0 {
					e.BlindedTimeLeft = e.HitResolution.BlindTime
					e.HitResolution.BlindTime = 0
				}
				if e.BlindedTimeLeft > 0 {
					e.BlindedTimeLeft--
				}

				// Process invincible.
				if e.InvincibleTimeLeft > 0 {
					e.InvincibleTimeLeft--
				}
			} else {
				resolveSlideOrDrag(e, s)
			}
		}
	}
}

func resolveSlideOrDrag(e *state.Entity, s *state.State) {
	if e.ForcedMovementState.ForcedMovement.Direction != state.DirectionNone {
		if e.ForcedMovementState.ElapsedTime == 0 {
			x, y := e.TilePos.XY()
			dx, dy := e.ForcedMovementState.ForcedMovement.Direction.XY()

			if !e.StartMove(state.TilePosXY(x+dx, y+dy), s) {
				if e.ForcedMovementState.ForcedMovement.Type.IsDrag() {
					e.DragLockoutTimeLeft = 20
				}
				e.ForcedMovementState = state.ForcedMovementState{}
			}
		} else if e.ForcedMovementState.ElapsedTime == 2 {
			e.FinishMove(s)
		} else if e.ForcedMovementState.ElapsedTime == 4 {
			e.ForcedMovementState = state.ForcedMovementState{}
		}
	}
}

func Step(s *state.State, b *bundle.Bundle) {
	s.ElapsedTime++

	if s.CounterPlaqueTimeLeft > 0 {
		s.CounterPlaqueTimeLeft--
	}

	for _, e := range s.Entities {
		e.PerTickState = state.EntityPerTickState{}
	}

	for _, snd := range s.Sounds {
		bbuf := b.Sounds[snd.Type]
		if state.TicksToSampleOffset(bbuf.Format().SampleRate, snd.ElapsedTime) >= bbuf.Len() {
			delete(s.Sounds, snd.ID())
			continue
		}
		snd.Step()
	}

	for _, d := range s.Decorations {
		if int(d.ElapsedTime) >= len(b.DecorationSprites[d.Type].Animation.Frames) {
			delete(s.Decorations, d.ID())
			continue
		}

		if !s.IsInTimeStop || d.RunsInTimestop {
			d.Step()
		}
	}

	if !s.IsInTimeStop {
		s.Field.Step(s)
	}

	// Step all entities in a random order.
	pending := maps.Values(s.Entities)
	slices.SortFunc(pending, func(a *state.Entity, b *state.Entity) bool {
		return a.ID() < b.ID()
	})
	rand.New(s.RandSource).Shuffle(len(pending), func(i, j int) {
		pending[i], pending[j] = pending[j], pending[i]
	})
	for _, e := range pending {
		if e.IsPendingDestruction {
			delete(s.Entities, e.ID())
			continue
		}

		if !e.ForcedMovementState.ForcedMovement.Type.IsDrag() && (!s.IsInTimeStop || e.RunsInTimestop) {
			e.Step(s)
			e.LastIntent = e.Intent
		}
	}

	// Resolve any hits.
	pending = maps.Values(s.Entities)
	slices.SortFunc(pending, func(a *state.Entity, b *state.Entity) bool {
		return a.ID() < b.ID()
	})
	rand.New(s.RandSource).Shuffle(len(pending), func(i, j int) {
		pending[i], pending[j] = pending[j], pending[i]
	})
	for _, e := range pending {
		resolveOne(e, s)

		if e.HP == 0 {
			// Do something special.
		}

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

}
