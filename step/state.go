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
		e.HitResolution.SlideDirection = state.DirectionNone
		e.HitResolution.Drag = state.DragTypeNone
	}

	// Set anger, if required.
	if e.HitResolution.Damage >= 300 {
		e.Emotion = state.EmotionAngry
	}

	// TODO: Process poison damage.

	// Process hit damage.
	if e.HitResolution.Damage > 0 {
		e.PerTickState.WasHit = true
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

	if e.SlideState.Direction != state.DirectionNone {
		// TODO: Is this even in the right place?
		e.SlideState.ElapsedTime++
	}

	if e.HitResolution.Drag != state.DragTypeNone {
		var postDragParalyzeTime state.Ticks
		if e.HitResolution.FlashTime == 0 {
			// Only add post drag paralysis if we're not going to be flashing afterwards.
			if paralyzed, ok := e.BehaviorState.Behavior.(*behaviors.Paralyzed); ok {
				postDragParalyzeTime = paralyzed.Duration - e.BehaviorState.ElapsedTime
			}
		}
		e.FinishMove(s)
		isBig := false
		if e.HitResolution.Drag == state.DragTypeBig {
			isBig = true
		}
		e.SlideState = state.SlideState{Direction: e.HitResolution.SlideDirection}
		e.SetBehaviorImmediate(&behaviors.Dragged{PostDragParalyzeTime: postDragParalyzeTime, IsBig: isBig}, s)
		e.HitResolution.Drag = state.DragTypeNone
		e.HitResolution.SlideDirection = state.DirectionNone
	} else {
		if !state.BehaviorIs[*behaviors.Dragged](e.BehaviorState.Behavior) && !s.IsInTimeStop {
			if e.HitResolution.SlideDirection != state.DirectionNone {
				// HACK: Allow immediate application of slide if the last slide is ending.
				if e.SlideState.Direction == state.DirectionNone || e.SlideState.ElapsedTime == 4 {
					e.SlideState = state.SlideState{Direction: e.HitResolution.SlideDirection, ElapsedTime: 0}
				}
				resolveSlide(e, s)
			} else {
				resolveSlide(e, s)
			}
			e.HitResolution.SlideDirection = state.DirectionNone

			if e.HitResolution.Flinch {
				if state.BehaviorIs[*behaviors.Paralyzed](e.BehaviorState.Behavior) && e.HitResolution.FlashTime == 0 {
					e.HitResolution.Flinch = false
				}

				if e.HitResolution.Flinch {
					// TODO: This should probably not be here...
					e.FinishMove(s)
					e.SetBehaviorImmediate(&behaviors.Flinch{}, s)
				}
			}
			e.HitResolution.Flinch = false

			// Process flashing.
			if e.HitResolution.FlashTime > 0 {
				e.FlashingTimeLeft = e.HitResolution.FlashTime
				e.HitResolution.FlashTime = 0
			}
			if e.FlashingTimeLeft > 0 {
				e.FlashingTimeLeft--
			}

			// Process paralyzed.
			if e.HitResolution.ParalyzeTime > 0 {
				e.FinishMove(s)
				e.SetBehaviorImmediate(&behaviors.Paralyzed{Duration: e.HitResolution.ParalyzeTime}, s)
				e.HitResolution.ConfuseTime = 0
				e.HitResolution.ParalyzeTime = 0
			}

			// Process frozen.
			if e.HitResolution.FreezeTime > 0 {
				e.FinishMove(s)
				e.SetBehaviorImmediate(&behaviors.Frozen{Duration: e.HitResolution.FreezeTime}, s)
				e.HitResolution.BubbleTime = 0
				e.HitResolution.ConfuseTime = 0
				e.HitResolution.FreezeTime = 0
			}

			// Process bubbled.
			if e.HitResolution.BubbleTime > 0 {
				e.FinishMove(s)
				e.SetBehaviorImmediate(&behaviors.Bubbled{Duration: e.HitResolution.BubbleTime}, s)
				e.ConfusedTimeLeft = 0
				e.HitResolution.ConfuseTime = 0
				e.HitResolution.BubbleTime = 0
			}

			// Process confused.
			if e.HitResolution.ConfuseTime > 0 {
				e.ConfusedTimeLeft = e.HitResolution.ConfuseTime
				// TODO: Double check if this is correct.
				if state.BehaviorIs[*behaviors.Paralyzed](e.BehaviorState.Behavior) || state.BehaviorIs[*behaviors.Frozen](e.BehaviorState.Behavior) || state.BehaviorIs[*behaviors.Bubbled](e.BehaviorState.Behavior) {
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
		}
	}
}

func resolveSlide(e *state.Entity, s *state.State) {
	if e.SlideState.Direction != state.DirectionNone {
		if e.SlideState.ElapsedTime == 0 {
			x, y := e.TilePos.XY()
			dx, dy := e.SlideState.Direction.XY()

			if !e.StartMove(state.TilePosXY(x+dx, y+dy), s) {
				e.SlideState = state.SlideState{}
			}
		} else if e.SlideState.ElapsedTime == 2 {
			e.FinishMove(s)
		} else if e.SlideState.ElapsedTime == 4 {
			e.SlideState = state.SlideState{}
		}
	}
}

func Step(s *state.State, b *bundle.Bundle) {
	s.ElapsedTime++

	for _, e := range s.Entities {
		e.PerTickState = state.EntityPerTickState{}
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

		if !s.IsInTimeStop || e.RunsInTimestop {
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

	if !s.IsInTimeStop {
		s.Field.Step(s)
	}
}
