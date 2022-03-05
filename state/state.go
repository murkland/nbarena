package state

import (
	"github.com/yumland/clone"
	"github.com/yumland/nbarena/bundle"
	"github.com/yumland/nbarena/draw"
	"github.com/yumland/syncrand"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Ticks int

type State struct {
	ElapsedTime Ticks

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
		nextEntityID: 1,
	}
}

func (s *State) AddEntity(e *Entity) int {
	e.id = s.nextEntityID
	e.elapsedTime = -1
	e.behaviorElapsedTime = -1
	e.IsPendingStep = true
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
