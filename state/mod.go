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

type Ticks int

const (
	OffererEntityID  = 1
	AnswererEntityID = 2
)

type State struct {
	elapsedTime Ticks

	randSource *syncrand.Source

	field Field

	entities     map[int]*Entity
	nextEntityID int
}

func New(randSource *syncrand.Source) State {
	field := newField()
	entities := map[int]*Entity{
		OffererEntityID: {
			behavior: &IdleEntityBehavior{},

			powerShotChargeTime: Ticks(50),

			tilePos:       TilePosXY(2, 2),
			futureTilePos: TilePosXY(2, 2),
		},
		AnswererEntityID: {
			isFlipped:            true,
			isAlliedWithAnswerer: true,

			powerShotChargeTime: Ticks(50),

			behavior:      &IdleEntityBehavior{},
			tilePos:       TilePosXY(5, 2),
			futureTilePos: TilePosXY(5, 2),
		},
	}

	return State{
		randSource: randSource,

		field: field,

		entities:     entities,
		nextEntityID: AnswererEntityID + 1,
	}
}

func (s *State) AddEntity(e *Entity) int {
	id := s.nextEntityID
	s.entities[id] = e
	s.nextEntityID++
	return id
}

func (s *State) RemoveEntity(id int) {
	delete(s.entities, id)
}

func (s *State) ElapsedTime() Ticks {
	return s.elapsedTime
}

func (s State) Clone() State {
	return State{
		s.elapsedTime,
		s.randSource.Clone(),
		s.field.Clone(),
		clone.Map(s.entities), s.nextEntityID,
	}
}

func (s *State) applyPlayerIntent(e *Entity, intent input.Intent, isOfferer bool) {
	interrupts := e.behavior.Interrupts(e)
	if intent.ChargeBasicWeapon && (interrupts.OnCharge || e.chargingElapsedTime > 0) {
		e.chargingElapsedTime++
	}

	if interrupts.OnCharge && !intent.ChargeBasicWeapon && e.chargingElapsedTime > 0 {
		// Release buster shot.
		e.SetBehavior(&BusterEntityBehavior{IsPowerShot: e.chargingElapsedTime >= e.powerShotChargeTime})
		e.chargingElapsedTime = 0
	}

	interrupts = e.behavior.Interrupts(e)
	if interrupts.OnMove {
		dir := intent.Direction
		if e.confusedTimeLeft > 0 {
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
		if tilePos != e.tilePos && e.isAlliedWithAnswerer == tile.isAlliedWithAnswerer && tile.CanEnter(e) {
			e.futureTilePos = tilePos
			e.SetBehavior(&MoveEntityBehavior{})
		}
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

type entityAndID struct {
	ID     int
	Entity *Entity
}

type StepHandle struct {
	state *State
	sq    *updateStack
}

func (sh *StepHandle) SpawnEntity(e *Entity) int {
	id := sh.state.AddEntity(e)
	sh.sq.Push(id, e)
	return id
}

func (sh *StepHandle) RemoveEntity(id int) {
	sh.state.RemoveEntity(id)
	sh.sq.Remove(id)
}

type updateStack struct {
	pending []entityAndID
}

func (sq *updateStack) HasMore() bool {
	return len(sq.pending) > 0
}

func (sq *updateStack) Remove(id int) {
	pending := make([]entityAndID, 0, cap(sq.pending))
	for _, eid := range sq.pending {
		if eid.ID == id {
			continue
		}
		pending = append(pending, eid)
	}
	sq.pending = pending
}

func (sq *updateStack) Push(id int, entity *Entity) {
	sq.pending = append(sq.pending, entityAndID{id, entity})
}

func (sq *updateStack) Pop() (int, *Entity) {
	eid := &sq.pending[len(sq.pending)-1]
	entity := eid.Entity
	eid.Entity = nil
	sq.pending = sq.pending[: len(sq.pending)-1 : cap(sq.pending)]
	return eid.ID, entity
}

func (s *State) Step() {
	s.elapsedTime++

	// Step entities in a random order.
	pending := make([]entityAndID, 0, len(s.entities))
	for id, entity := range s.entities {
		pending = append(pending, entityAndID{id, entity})
	}
	sort.Slice(pending, func(i, j int) bool {
		return pending[i].ID < pending[j].ID
	})
	rand.New(s.randSource).Shuffle(len(pending), func(i, j int) {
		pending[i], pending[j] = pending[j], pending[i]
	})

	sq := &updateStack{pending}
	for sq.HasMore() {
		_, entity := sq.Pop()
		sh := &StepHandle{s, sq}
		entity.Step(sh)
	}

	s.field.Step()
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

func (s *State) Entity(id int) *Entity {
	return s.entities[id]
}
