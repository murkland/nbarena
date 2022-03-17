package state

import "github.com/murkland/clone"

// 12725 -> 12743 (text appears) -> 12750 (text completes) -> 12793 (text starts disappearing) -> 12801 (text disappears) -> 12803 (action start) -> 12844 (Action end) -> 12882 (tf end)

type Timestop struct {
	Parent *Timestop

	Owner EntityID

	Behavior            TimestopBehavior
	BehaviorElapsedTime Ticks

	IsPendingDestruction bool
}

func (t *Timestop) Step(s *State) {
	t.BehaviorElapsedTime++
	t.Behavior.Step(t, s)
}

func (t *Timestop) Clone() *Timestop {
	return &Timestop{
		clone.ValuePointer(t.Parent),
		t.Owner,
		t.Behavior.Clone(),
		t.BehaviorElapsedTime,
		t.IsPendingDestruction,
	}
}

type TimestopBehavior interface {
	clone.Cloner[TimestopBehavior]
	Step(t *Timestop, s *State)
}
