package state

import (
	"math/rand"
	"sort"

	"github.com/yumland/clone"
	"github.com/yumland/syncrand"
	"github.com/yumland/yumbattle/input"
)

type State struct {
	elapsedTicks int

	randSource *syncrand.Source

	field    Field
	entities map[int]*Entity
}

func New(randSource *syncrand.Source) State {
	field := newField()

	return State{
		randSource: randSource,

		field:    field,
		entities: make(map[int]*Entity),
	}
}

func (s *State) ElapsedTicks() int {
	return s.elapsedTicks
}

func (s State) Clone() State {
	return State{
		s.elapsedTicks,
		s.randSource.Clone(),
		s.field.Clone(), clone.Map(s.entities),
	}
}

func (s *State) Apply(offererIntent input.Intent, answererIntent input.Intent) {
	intents := []struct {
		isOfferer bool
		intent    input.Intent
	}{
		{true, offererIntent},
		{false, answererIntent},
	}
	rand.New(s.randSource).Shuffle(len(intents), func(i, j int) {
		intents[i], intents[j] = intents[j], intents[i]
	})
	for _, wrapped := range intents {
		intent := wrapped.intent
		if wrapped.isOfferer {
			_ = intent
		} else {
			_ = intent
		}
	}
}

func (s *State) Step() {
	s.elapsedTicks++

	s.field.Step()

	// Step entities in a random order.
	entities := make([]struct {
		id     int
		entity *Entity
	}, 0, len(s.entities))
	for id, entity := range s.entities {
		entities = append(entities, struct {
			id     int
			entity *Entity
		}{id, entity})
	}
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].id < entities[j].id
	})
	rand.New(s.randSource).Shuffle(len(entities), func(i, j int) {
		entities[i], entities[j] = entities[j], entities[i]
	})
	for _, wrapped := range entities {
		wrapped.entity.Step()
	}
}
