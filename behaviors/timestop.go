package behaviors

import "github.com/murkland/nbarena/state"

type Timestop struct {
	returnBehaviorState state.EntityBehaviorState
}

func (eb *Timestop) Step(e *state.Entity, s *state.State) {
	s.IsInTimeStop = true
}
