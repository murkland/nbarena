package state

import (
	"math/rand"

	"github.com/yumland/syncrand"
	"github.com/yumland/yumbattle/input"
)

type State struct {
	Rng          *syncrand.Source
	ElapsedTicks uint32
}

func (s *State) Clone() *State {
	return &State{s.Rng.Clone(), s.ElapsedTicks}
}

func (s *State) Update(intents []input.Intent) error {
	rand.New(s.Rng).Shuffle(len(intents), func(i, j int) {
		intents[i], intents[j] = intents[j], intents[i]
	})
	s.ElapsedTicks++
	return nil
}
