package chips

import (
	"github.com/murkland/nbarena/behaviors/timestopbehaviors"
	"github.com/murkland/nbarena/state"
)

var AreaGrab = &state.Chip{
	Index: 162,
	Name:  "AreaGrab",
	OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
		s.StartTimestop(&timestopbehaviors.AreaGrab{Owner: e.ID()})
	},
}
