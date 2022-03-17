package state

// 12725 -> 12743 (text appears) -> 12750 (text completes) -> 12793 (text starts disappearing) -> 12801 (text disappears) -> 12803 (action start) -> 12844 (Action end) -> 12882 (tf end)

type Timestop struct {
}

func (t *Timestop) Clone() *Timestop {
	return &Timestop{}
}
