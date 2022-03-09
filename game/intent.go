package game

import (
	"math/rand"

	"github.com/murkland/nbarena/input"
	"github.com/murkland/nbarena/state"
)

func applyPlayerIntents(s *state.State, offererEntityID int, offererIntent input.Intent, answererEntityID int, answererIntent input.Intent) {
	intents := []struct {
		isOfferer bool
		intent    input.Intent
	}{
		{true, offererIntent},
		{false, answererIntent},
	}
	rand.New(s.RandSource).Shuffle(len(intents), func(i, j int) {
		intents[i], intents[j] = intents[j], intents[i]
	})
	for _, wrapped := range intents {
		var entity *state.Entity
		if wrapped.isOfferer {
			entity = s.Entities[offererEntityID]
		} else {
			entity = s.Entities[answererEntityID]
		}
		entity.Behavior().ApplyIntent(entity, s, wrapped.intent)
		entity.LastIntent = wrapped.intent
	}
}
