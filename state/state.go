package state

import (
	"math/rand"
	"sort"

	"github.com/yumland/clone"
	"github.com/yumland/syncrand"
	"github.com/yumland/yumbattle/bundle"
	"github.com/yumland/yumbattle/draw"
)

type Ticks int

type State struct {
	elapsedTime Ticks

	RandSource *syncrand.Source

	Field Field

	Entities     map[int]*Entity
	nextEntityID int
}

func New(randSource *syncrand.Source) State {
	return State{
		RandSource: randSource,

		Field: newField(),

		Entities:     map[int]*Entity{},
		nextEntityID: 0,
	}
}

func (s *State) AddEntity(e *Entity) int {
	e.id = s.nextEntityID
	s.Entities[e.id] = e
	s.nextEntityID++
	return e.id
}

func (s *State) RemoveEntity(id int) {
	delete(s.Entities, id)
}

func (s *State) ElapsedTime() Ticks {
	return s.elapsedTime
}

func (s State) Clone() State {
	return State{
		s.elapsedTime,
		s.RandSource.Clone(),
		s.Field.Clone(),
		clone.Map(s.Entities), s.nextEntityID,
	}
}

type entityAndID struct {
	ID     int
	Entity *Entity
}

type StepHandle struct {
	State *State
	sq    *updateStack
}

func (sh *StepHandle) SpawnEntity(e *Entity) {
	e.behaviorElapsedTime = -1
	e.elapsedTime = -1
	sh.State.AddEntity(e)
	sh.sq.Push(e)
}

func (sh *StepHandle) RemoveEntity(id int) {
	sh.State.RemoveEntity(id)
	sh.sq.Remove(id)
}

type updateStack struct {
	pending []*Entity
}

func (sq *updateStack) HasMore() bool {
	return len(sq.pending) > 0
}

func (sq *updateStack) Remove(id int) {
	pending := make([]*Entity, 0, cap(sq.pending))
	for _, entity := range sq.pending {
		if entity.id == id {
			continue
		}
		pending = append(pending, entity)
	}
	sq.pending = pending
}

func (sq *updateStack) Push(entity *Entity) {
	sq.pending = append(sq.pending, entity)
}

func (sq *updateStack) Pop() *Entity {
	slot := &sq.pending[len(sq.pending)-1]
	entity := *slot
	*slot = nil
	sq.pending = sq.pending[: len(sq.pending)-1 : cap(sq.pending)]
	return entity
}

func (s *State) Step() {
	s.elapsedTime++

	// Step Entities in a random order.
	pending := make([]*Entity, 0, len(s.Entities))
	for _, entity := range s.Entities {
		pending = append(pending, entity)
	}
	sort.Slice(pending, func(i, j int) bool {
		return pending[i].id < pending[j].id
	})
	rand.New(s.RandSource).Shuffle(len(pending), func(i, j int) {
		pending[i], pending[j] = pending[j], pending[i]
	})

	sq := &updateStack{pending}
	for sq.HasMore() {
		entity := sq.Pop()
		sh := &StepHandle{s, sq}
		entity.Step(sh)
	}

	s.Field.Step()
}

const (
	fieldOffsetTopFull = 87
	fieldOffsetTop     = 72
)

func (s *State) Appearance(b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	rootNode.Opts.GeoM.Translate(0, fieldOffsetTop)
	{
		tilesNode := &draw.OptionsNode{}
		tilesNode.Children = append(tilesNode.Children, s.Field.Appearance(b))
		rootNode.Children = append(rootNode.Children, tilesNode)
	}
	{
		entitiesNode := &draw.OptionsNode{}
		for _, entity := range s.Entities {
			node := entity.Appearance(b)
			if node == nil {
				continue
			}
			entitiesNode.Children = append(entitiesNode.Children, node)
		}
		rootNode.Children = append(rootNode.Children, entitiesNode)
	}
	return rootNode
}
