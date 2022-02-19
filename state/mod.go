package state

import (
	"math/rand"

	"github.com/yumland/syncrand"
	"github.com/yumland/yumbattle/input"
)

type PlayerState struct {
	X int
	Y int
}

func (ps *PlayerState) Apply(intent input.Intent) {
	if intent.Direction&input.DirectionLeft == input.DirectionLeft {
		ps.X--
	}
	if intent.Direction&input.DirectionRight == input.DirectionRight {
		ps.X++
	}
	if intent.Direction&input.DirectionUp == input.DirectionUp {
		ps.Y--
	}
	if intent.Direction&input.DirectionDown == input.DirectionDown {
		ps.Y++
	}
}

func (ps *PlayerState) Step() {

}

type State struct {
	Rng          *syncrand.Source
	ElapsedTicks uint32

	OffererPlayer  PlayerState
	AnswererPlayer PlayerState
}

func New(rng *syncrand.Source) *State {
	return &State{Rng: rng}
}

func (s *State) Clone() *State {
	return &State{s.Rng.Clone(), s.ElapsedTicks, s.OffererPlayer, s.AnswererPlayer}
}

func (s *State) Apply(offererIntent input.Intent, answererIntent input.Intent) {
	type wrappedIntent struct {
		isOfferer bool
		intent    input.Intent
	}

	wrappedIntents := []wrappedIntent{{true, offererIntent}, {false, answererIntent}}
	rand.New(s.Rng).Shuffle(len(wrappedIntents), func(i, j int) {
		wrappedIntents[i], wrappedIntents[j] = wrappedIntents[j], wrappedIntents[i]
	})

	for _, wrapped := range wrappedIntents {
		intent := wrapped.intent
		if wrapped.isOfferer {
			s.OffererPlayer.Apply(intent)
		} else {
			s.AnswererPlayer.Apply(intent)
		}
	}
}

func (s *State) Step() {
	// TODO: Step everything in a random order.
	s.OffererPlayer.Step()
	s.AnswererPlayer.Step()
	s.ElapsedTicks++
}
