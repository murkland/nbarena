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

	IsInTimeStop bool

	Entities     map[EntityID]*Entity
	nextEntityID EntityID

	Decorations      map[DecorationID]*Decoration
	nextDecorationID DecorationID
}

func New(randSource *syncrand.Source) *State {
	field := newField()
	return &State{
		RandSource: randSource,

		Field: field,

		Entities:     map[EntityID]*Entity{},
		nextEntityID: 1,

		Decorations:      map[DecorationID]*Decoration{},
		nextDecorationID: 1,
	}
}

func (s *State) AttachEntity(e *Entity) {
	e.id = s.nextEntityID
	if s.IsInTimeStop {
		e.RunsInTimestop = true
	}
	s.Entities[e.id] = e
	e.BehaviorState.Behavior.Step(e, s)
	s.nextEntityID++
}

func (s *State) RemoveEntity(id EntityID) {
	delete(s.Entities, id)
}

func (s *State) AddDecoration(d *Decoration) {
	d.id = s.nextDecorationID
	s.Decorations[d.id] = d
	s.nextDecorationID++
}

func (s *State) RemoveDecoration(id DecorationID) {
	delete(s.Decorations, id)
}
func (s *State) Clone() *State {
	return &State{
		s.ElapsedTime,
		s.RandSource.Clone(),
		s.Field.Clone(),
		s.IsInTimeStop,
		clone.Map(s.Entities), s.nextEntityID,
		clone.Map(s.Decorations), s.nextDecorationID,
	}
}

func (s *State) Flip() {
	s.Field.Flip()
	for _, entity := range s.Entities {
		entity.Flip()
	}
	for _, decoration := range s.Decorations {
		decoration.Flip()
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
			return a.ID() > b.ID()
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
	{
		decorationsNode := &draw.OptionsNode{}
		decorations := maps.Values(s.Decorations)
		slices.SortFunc(decorations, func(a *Decoration, b *Decoration) bool {
			x1, y1 := a.TilePos.XY()
			x2, y2 := b.TilePos.XY()
			if y1 != y2 {
				return y1 < y2
			}
			if x1 != x2 {
				return x1 < x2
			}
			return a.ID() > b.ID()
		})
		for _, decoration := range decorations {
			node := decoration.Appearance(b)
			if node == nil {
				continue
			}
			decorationsNode.Children = append(decorationsNode.Children, node)
		}
		rootNode.Children = append(rootNode.Children, decorationsNode)
	}
	return rootNode
}

func (s *State) EntitiesAt(pos TilePos) []*Entity {
	var entities []*Entity
	for _, e := range s.Entities {
		if e.TilePos != pos {
			continue
		}
		entities = append(entities, e)
	}
	slices.SortFunc(entities, func(a, b *Entity) bool {
		return a.ID() > b.ID()
	})
	return entities
}

func (s *State) ApplyHit(owner *Entity, pos TilePos, h Hit) bool {
	for _, target := range s.EntitiesAt(pos) {
		if target.IsAlliedWithAnswerer == owner.IsAlliedWithAnswerer {
			continue
		}

		if target.Traits.Intangible {
			continue
		}

		if h.RemovesFlashing {
			target.FlashingTimeLeft = 0
		}

		if target.FlashingTimeLeft > 0 {
			continue
		}

		if h.CanCounter && target.BehaviorState.Behavior.Traits(target).CanBeCountered && target.BehaviorState.ElapsedTime < 15 {
			// From Alyrsc#7506:
			// I was mostly sure that it's frames 2-16 of an action.
			// I gathered that by frame stepping P2 while P1 had FullSynchro. The timing of the blue flashes was somewhat inconsistent, possibly because it's based on a global clock or counter, but those were the earliest and latest frames I saw.
			// TODO: Check the code for this.
			owner.Emotion = EmotionFullSynchro
			h.FlashTime = 0
			h.ParalyzeTime = DefaultParalyzeTime
		}

		target.ApplyHit(h)
		return true
	}

	return false
}
