package chips

import (
	"github.com/murkland/nbarena/state"
)

var AreaGrab = &state.Chip{
	Index: 162,
	Name:  "AreaGrab",
	MakeBehavior: func(damage state.Damage) state.EntityBehavior {
		// s.StartTimestop(e, &timestopbehaviors.AreaGrab{Owner: e.ID()})
		return nil
	},
}
