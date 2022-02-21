package state

import (
	"math/rand"

	"github.com/yumland/clone"
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
	elapsedTicks int

	randSource *syncrand.Source

	field    Field
	entities map[int]Entity

	OffererPlayer  PlayerState
	AnswererPlayer PlayerState
}

func New(randSource *syncrand.Source) State {
	field := newField()

	return State{
		randSource: randSource,

		field:    field,
		entities: make(map[int]Entity),
	}
}

func (s *State) ElapsedTicks() int {
	return s.elapsedTicks
}

func (s State) Clone() State {
	return State{
		s.elapsedTicks,
		s.randSource.Clone(),
		s.field, clone.Map(s.entities),
		s.OffererPlayer, s.AnswererPlayer,
	}
}

func (s *State) Apply(offererIntent input.Intent, answererIntent input.Intent) {
	wrappedIntents := []struct {
		isOfferer bool
		intent    input.Intent
	}{
		{true, offererIntent},
		{false, answererIntent},
	}
	rand.New(s.randSource).Shuffle(len(wrappedIntents), func(i, j int) {
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
	s.elapsedTicks++
	// TODO: Step everything in a random order.
	s.OffererPlayer.Step()
	s.AnswererPlayer.Step()
	s.field.Step()
}
