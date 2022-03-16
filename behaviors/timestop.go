package behaviors

import "github.com/murkland/nbarena/state"

type Timestop struct {
	activeBehavior      state.EntityBehaviorState
	returnBehaviorState state.EntityBehaviorState
}

func (eb *Timestop) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == 0 {
		s.IsInTimeStop = true
		e.RunsInTimestop = true
	}
}

func (eb *Timestop) Cleanup(e *state.Entity, s *state.State) {
	s.IsInTimeStop = false
	e.RunsInTimestop = false
}

// 12725 -> 12743 (text appears) -> 12750 (text completes) -> 12793 (text starts disappearing) -> 12801 (text disappears) -> 12803 (action start) -> 12844 (Action end) -> 12882 (tf end)
