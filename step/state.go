package step

import (
	"math/rand"

	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/state"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func resolveHit(e *state.Entity, s *state.State) {
	if e.Hit.RemovesFlashing {
		e.FlashingTimeLeft = 0
	}

	if e.FlashingTimeLeft > 0 {
		e.Hit = state.Hit{}
	}

	// Set anger, if required.
	if e.Hit.TotalDamage >= 300 {
		e.IsAngry = true
	}

	// TODO: Process poison damage.

	// Process hit damage.
	if e.Hit.TotalDamage > 0 {
		e.PerTickState.WasHit = true
	}

	mustLeave1HP := e.HP > 1 && e.Traits.FatalHitLeaves1HP
	e.HP -= e.Hit.TotalDamage
	if e.HP < 0 {
		e.HP = 0
	}
	if mustLeave1HP {
		e.HP = 1
	}
	e.Hit.TotalDamage = 0

	if e.Hit.Flinch && !e.Traits.CannotFlinch {
		e.FinishMove(s)
		e.SetBehavior(&behaviors.Flinch{}, s)
	}
	e.Hit.Flinch = false

	if e.IsCounterable && e.Hit.Counters {
		e.FlashingTimeLeft = 0
		e.FinishMove(s)
		e.SetBehavior(&behaviors.Paralyzed{Duration: 150}, s)
	}
	e.Hit.Counters = false

	if e.SlideState.Slide.Direction != state.DirectionNone {
		// TODO: Is this even in the right place?
		e.SlideState.ElapsedTime++
	}

	if !e.Hit.Drag {
		if !state.BehaviorIs[*behaviors.Dragged](e.BehaviorState.Behavior) && !s.IsInTimeStop {
			if e.Hit.Slide.Direction != state.DirectionNone {
				e.SlideState.Slide = e.Hit.Slide
				e.SlideState.ElapsedTime = 0
				e.Hit.Slide = state.Slide{}
			}

			// Process flashing.
			if e.Hit.FlashTime > 0 {
				e.FlashingTimeLeft = e.Hit.FlashTime
				e.Hit.FlashTime = 0
			}
			if e.FlashingTimeLeft > 0 {
				e.FlashingTimeLeft--
			}

			// Process paralyzed.
			if e.Hit.ParalyzeTime > 0 {
				e.FinishMove(s)
				e.SetBehavior(&behaviors.Paralyzed{Duration: e.Hit.ParalyzeTime}, s)
				e.Hit.ConfuseTime = 0
				e.Hit.ParalyzeTime = 0
			}

			// Process frozen.
			if e.Hit.FreezeTime > 0 {
				e.FinishMove(s)
				e.SetBehavior(&behaviors.Frozen{Duration: e.Hit.FreezeTime}, s)
				e.Hit.BubbleTime = 0
				e.Hit.ConfuseTime = 0
				e.Hit.FreezeTime = 0
			}

			// Process bubbled.
			if e.Hit.BubbleTime > 0 {
				e.FinishMove(s)
				e.SetBehavior(&behaviors.Bubbled{Duration: e.Hit.BubbleTime}, s)
				e.ConfusedTimeLeft = 0
				e.Hit.ConfuseTime = 0
				e.Hit.BubbleTime = 0
			}

			// Process confused.
			if e.Hit.ConfuseTime > 0 {
				e.ConfusedTimeLeft = e.Hit.ConfuseTime
				// TODO: Double check if this is correct.
				if state.BehaviorIs[*behaviors.Paralyzed](e.BehaviorState.Behavior) || state.BehaviorIs[*behaviors.Frozen](e.BehaviorState.Behavior) || state.BehaviorIs[*behaviors.Bubbled](e.BehaviorState.Behavior) {
					e.SetBehavior(&behaviors.Idle{}, s)
				}
				e.Hit.FreezeTime = 0
				e.Hit.BubbleTime = 0
				e.Hit.ParalyzeTime = 0
				e.Hit.ConfuseTime = 0
			}
			if e.ConfusedTimeLeft > 0 {
				e.ConfusedTimeLeft--
			}

			// Process immobilized.
			if e.Hit.ImmobilizeTime > 0 {
				e.ImmobilizedTimeLeft = e.Hit.ImmobilizeTime
				e.Hit.ImmobilizeTime = 0
			}
			if e.ImmobilizedTimeLeft > 0 {
				e.ImmobilizedTimeLeft--
			}

			// Process blinded.
			if e.Hit.BlindTime > 0 {
				e.BlindedTimeLeft = e.Hit.BlindTime
				e.Hit.BlindTime = 0
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
		var postDragParalyzeTime state.Ticks
		if paralyzed, ok := e.BehaviorState.Behavior.(*behaviors.Paralyzed); ok {
			postDragParalyzeTime = paralyzed.Duration - e.BehaviorState.ElapsedTime
		}
		e.FinishMove(s)
		e.SetBehavior(&behaviors.Dragged{PostDragParalyzeTime: postDragParalyzeTime}, s)
		e.SlideState.Slide = e.Hit.Slide
		e.SlideState.ElapsedTime = 0
		e.Hit.Drag = false
		e.Hit.Slide = state.Slide{}
	}
}

func resolveSlide(e *state.Entity, s *state.State) {
	if e.SlideState.Slide.Direction != state.DirectionNone {
		if e.SlideState.ElapsedTime%4 == 0 {
			x, y := e.TilePos.XY()
			dx, dy := e.SlideState.Slide.Direction.XY()

			if !e.StartMove(state.TilePosXY(x+dx, y+dy), s) {
				e.SlideState = state.SlideState{}
			}
		} else if e.SlideState.ElapsedTime%4 == 2 {
			e.FinishMove(s)
			if !e.SlideState.Slide.IsBig {
				e.SlideState = state.SlideState{}
			}
		}
	}
}

func Step(s *state.State) {
	s.ElapsedTime++

	for _, e := range s.Entities {
		e.PerTickState = state.EntityPerTickState{}
	}

	// Step all entities in a random order.
	pending := make([]*state.Entity, 0, len(s.Entities))
	for _, e := range s.Entities {
		pending = append(pending, e)
	}

	slices.SortFunc(pending, func(a *state.Entity, b *state.Entity) bool {
		return a.ID() < b.ID()
	})
	rand.New(s.RandSource).Shuffle(len(pending), func(i, j int) {
		pending[i], pending[j] = pending[j], pending[i]
	})

	for _, e := range pending {
		if !s.IsInTimeStop || state.BehaviorIs[state.TimestopMaskedEntityBehavior](e.BehaviorState.Behavior) {
			e.Step(s)
			e.LastIntent = e.Intent
		}
	}

	// Resolve any hits.
	for _, e := range maps.Values(s.Entities) {
		resolveHit(e, s)
		resolveSlide(e, s)

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
		if e.PerTickState.IsPendingDeletion {
			delete(s.Entities, k)
		}
	}

	if !s.IsInTimeStop {
		s.Field.Step(s)
	}
}
