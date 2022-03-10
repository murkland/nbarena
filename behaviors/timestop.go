package behaviors

import "github.com/murkland/nbarena/state"

type Timestop struct {
	returnBehavior            state.EntityBehavior
	returnBehaviorElapsedTime state.Ticks
}

func (eb *Timestop) Step(e *state.Entity, s *state.State) {
	s.IsInTimeStop = true
}
