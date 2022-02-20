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
	ElapsedTicks uint32

	RandSource *syncrand.Source

	Tiles    []*Tile
	Entities map[int]*Entity

	OffererPlayer  PlayerState
	AnswererPlayer PlayerState
}

func New(randSource *syncrand.Source) *State {
	tiles := EmptyTiles()

	return &State{RandSource: randSource, Tiles: tiles, Entities: make(map[int]*Entity)}
}

func (s *State) Clone() *State {
	return &State{
		s.ElapsedTicks,
		s.RandSource.Clone(),
		clone.Slice(s.Tiles), clone.Map(s.Entities),
		s.OffererPlayer, s.AnswererPlayer,
	}
}

func (s *State) Apply(offererIntent input.Intent, answererIntent input.Intent) {
	type wrappedIntent struct {
		isOfferer bool
		intent    input.Intent
	}

	wrappedIntents := []wrappedIntent{{true, offererIntent}, {false, answererIntent}}
	rand.New(s.RandSource).Shuffle(len(wrappedIntents), func(i, j int) {
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
