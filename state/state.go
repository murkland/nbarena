package state

import (
	"github.com/murkland/clone"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/syncrand"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Ticks int

type State struct {
	ElapsedTime Ticks

	RandSource *syncrand.Source

	Field *Field

	Entities     map[int]*Entity
	nextEntityID int
}

func New(randSource *syncrand.Source) State {
	field := newField()
	return State{
		RandSource: randSource,

		Field: field,

		Entities:     map[int]*Entity{},
		nextEntityID: 1,
	}
}

func (s *State) AddEntity(e *Entity) int {
	e.id = s.nextEntityID
	e.elapsedTime = -1
	e.behaviorElapsedTime = -1
	s.Entities[e.id] = e
	s.nextEntityID++
	return e.id
}

func (s *State) RemoveEntity(id int) {
	delete(s.Entities, id)
}

func (s State) Clone() State {
	return State{
		s.ElapsedTime,
		s.RandSource.Clone(),
		s.Field.Clone(),
		clone.Map(s.Entities), s.nextEntityID,
	}
}

func (s *State) Flip() {
	s.Field.Flip()
	for _, entity := range s.Entities {
		entity.Flip()
	}
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
		entities := maps.Values(s.Entities)
		slices.SortFunc(entities, func(a *Entity, b *Entity) bool {
			x1, y1 := a.TilePos.XY()
			x2, y2 := b.TilePos.XY()
			if y1 != y2 {
				return y1 < y2
			}
			if x1 != x2 {
				return x1 < x2
			}
			return a.ID() < b.ID()
		})
		for _, entity := range entities {
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
