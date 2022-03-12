package step

import (
	"math/rand"

	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/state"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func resolveHit(e *state.Entity, s *state.State) {
	if e.Traits.CannotFlinch || e.IsAngry {
		// TODO: Double check if this
		e.Hit.Traits.Flinch = false
	}

	// Set anger, if required.
	if e.Hit.TotalDamage >= 300 {
		e.IsAngry = true
	}

	if e.Hit.Traits.RemovesFlashing {
		e.FlashingTimeLeft = 0
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

	// TODO: Pop bubble, if required.

	if e.SlideState.Direction != state.DirectionNone {
		// TODO: Is this even in the right place?
		e.SlideState.ElapsedTime++
	}

	if e.Hit.Traits.Drag == state.DragTypeNone {
		if !state.BehaviorIs[*behaviors.Dragged](e.BehaviorState.Behavior) && !s.IsInTimeStop {
			if e.SlideState.Direction == state.DirectionNone {
				if e.Hit.Traits.SlideDirection != state.DirectionNone {
					e.SlideState = state.SlideState{Direction: e.Hit.Traits.SlideDirection, ElapsedTime: 0}
				}
				resolveSlide(e, s)
			} else {
				resolveSlide(e, s)
			}
			e.Hit.Traits.SlideDirection = state.DirectionNone

			if e.Hit.Traits.Flinch {
				if state.BehaviorIs[*behaviors.Paralyzed](e.BehaviorState.Behavior) && e.Hit.Traits.FlashTime == 0 {
					e.Hit.Traits.Flinch = false
				}

				if e.Hit.Traits.Flinch {
					// TODO: This should probably not be here...
					e.FinishMove(s)
					e.BehaviorState = state.EntityBehaviorState{Behavior: &behaviors.Flinch{}}
				}
			}
			e.Hit.Traits.Flinch = false

			// Process flashing.
			if e.Hit.Traits.FlashTime > 0 {
				e.FlashingTimeLeft = e.Hit.Traits.FlashTime
				e.Hit.Traits.FlashTime = 0
			}
			if e.FlashingTimeLeft > 0 {
				e.FlashingTimeLeft--
			}

			// Process paralyzed.
			if e.Hit.Traits.ParalyzeTime > 0 {
				e.FinishMove(s)
				e.BehaviorState = state.EntityBehaviorState{Behavior: &behaviors.Paralyzed{Duration: e.Hit.Traits.ParalyzeTime}}
				e.Hit.Traits.ConfuseTime = 0
				e.Hit.Traits.ParalyzeTime = 0
			}

			// Process frozen.
			if e.Hit.Traits.FreezeTime > 0 {
				e.FinishMove(s)
				e.BehaviorState = state.EntityBehaviorState{Behavior: &behaviors.Frozen{Duration: e.Hit.Traits.FreezeTime}}
				e.Hit.Traits.BubbleTime = 0
				e.Hit.Traits.ConfuseTime = 0
				e.Hit.Traits.FreezeTime = 0
			}

			// Process bubbled.
			if e.Hit.Traits.BubbleTime > 0 {
				e.FinishMove(s)
				e.BehaviorState = state.EntityBehaviorState{Behavior: &behaviors.Bubbled{Duration: e.Hit.Traits.BubbleTime}}
				e.ConfusedTimeLeft = 0
				e.Hit.Traits.ConfuseTime = 0
				e.Hit.Traits.BubbleTime = 0
			}

			// Process confused.
			if e.Hit.Traits.ConfuseTime > 0 {
				e.ConfusedTimeLeft = e.Hit.Traits.ConfuseTime
				// TODO: Double check if this is correct.
				if state.BehaviorIs[*behaviors.Paralyzed](e.BehaviorState.Behavior) || state.BehaviorIs[*behaviors.Frozen](e.BehaviorState.Behavior) || state.BehaviorIs[*behaviors.Bubbled](e.BehaviorState.Behavior) {
					e.BehaviorState = state.EntityBehaviorState{Behavior: &behaviors.Idle{}}
				}
				e.Hit.Traits.FreezeTime = 0
				e.Hit.Traits.BubbleTime = 0
				e.Hit.Traits.ParalyzeTime = 0
				e.Hit.Traits.ConfuseTime = 0
			}
			if e.ConfusedTimeLeft > 0 {
				e.ConfusedTimeLeft--
			}

			// Process immobilized.
			if e.Hit.Traits.ImmobilizeTime > 0 {
				e.ImmobilizedTimeLeft = e.Hit.Traits.ImmobilizeTime
				e.Hit.Traits.ImmobilizeTime = 0
			}
			if e.ImmobilizedTimeLeft > 0 {
				e.ImmobilizedTimeLeft--
			}

			// Process blinded.
			if e.Hit.Traits.BlindTime > 0 {
				e.BlindedTimeLeft = e.Hit.Traits.BlindTime
				e.Hit.Traits.BlindTime = 0
			}
			if e.BlindedTimeLeft > 0 {
				e.BlindedTimeLeft--
			}

			// Process invincible.
			if e.InvincibleTimeLeft > 0 {
				e.InvincibleTimeLeft--
			}
		}
	} else {
		var postDragParalyzeTime state.Ticks
		if e.Hit.Traits.FlashTime == 0 {
			// Only add post drag paralysis if we're not going to be flashing afterwards.
			if paralyzed, ok := e.BehaviorState.Behavior.(*behaviors.Paralyzed); ok {
				postDragParalyzeTime = paralyzed.Duration - e.BehaviorState.ElapsedTime
			}
		}
		e.FinishMove(s)
		e.BehaviorState = state.EntityBehaviorState{Behavior: &behaviors.Dragged{PostDragParalyzeTime: postDragParalyzeTime, Type: e.Hit.Traits.Drag}}
		e.SlideState = state.SlideState{Direction: e.Hit.Traits.SlideDirection, ElapsedTime: 0}
		e.Hit.Traits.Drag = state.DragTypeNone
		e.Hit.Traits.SlideDirection = state.DirectionNone
	}

	if state.BehaviorIs[*behaviors.Dragged](e.BehaviorState.Behavior) && !s.IsInTimeStop {
		// Resolve drag-based slide.
		resolveSlide(e, s)
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

		if !s.IsInTimeStop || e.BehaviorState.Behavior.Traits(e).RunsInTimestop {
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
		resolveHit(e, s)

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
