package state

import (
	"math/rand"
	"sort"

	"github.com/yumland/clone"
	"github.com/yumland/syncrand"
	"github.com/yumland/yumbattle/bundle"
	"github.com/yumland/yumbattle/draw"
	"github.com/yumland/yumbattle/input"
)

const (
	OffererEntityID  = 1
	AnswererEntityID = 2
)

type State struct {
	elapsedTicks int

	randSource *syncrand.Source

	field    Field
	entities map[int]*Entity
}

func New(randSource *syncrand.Source) State {
	field := newField()
	entities := map[int]*Entity{
		OffererEntityID: {
			tilePos:       TilePosXY(2, 2),
			futureTilePos: TilePosXY(2, 2),
		},
		AnswererEntityID: {
			isFlipped:            true,
			isAlliedWithAnswerer: true,
			tilePos:              TilePosXY(5, 2),
			futureTilePos:        TilePosXY(5, 2),
		},
	}

	return State{
		randSource: randSource,

		field:    field,
		entities: entities,
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

func (s *State) applyPlayerIntent(e *Entity, intent input.Intent, isOfferer bool) {
	dir := intent.Direction
	if e.confusedFramesLeft > 0 {
		dir = dir.FlipH().FlipV()
	}

	x, y := e.tilePos.XY()
	if dir&input.DirectionLeft != 0 {
		x--
	}
	if dir&input.DirectionRight != 0 {
		x++
	}
	if dir&input.DirectionUp != 0 {
		y--
	}
	if dir&input.DirectionDown != 0 {
		y++
	}

	tilePos := TilePosXY(x, y)
	tile := &s.field.tiles[tilePos]
	if e.isAlliedWithAnswerer == tile.isAlliedWithAnswerer && tile.CanEnter(e) {
		e.futureTilePos = tilePos
		// TODO: Switch behavior.
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
		var entity *Entity
		if wrapped.isOfferer {
			entity = s.entities[OffererEntityID]
		} else {
			entity = s.entities[AnswererEntityID]
		}
		s.applyPlayerIntent(entity, wrapped.intent, wrapped.isOfferer)
	}
}

func (s *State) Step() {
	s.elapsedTicks++

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

	s.field.Step()
}

const (
	fieldOffsetTopFull = 87
	fieldOffsetTop     = 72
)

func (s *State) Appearance(b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	rootNode.Opts.GeoM.Translate(0, fieldOffsetTopFull)
	{
		tilesNode := &draw.OptionsNode{}
		tilesNode.Children = append(tilesNode.Children, s.field.Appearance(b))
		rootNode.Children = append(rootNode.Children, tilesNode)
	}
	{
		entitiesNode := &draw.OptionsNode{}
		for _, entity := range s.entities {
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
