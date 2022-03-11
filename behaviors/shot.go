package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/state/query"
)

type Shot struct {
	Damage    state.Damage
	HitTraits state.HitTraits
}

func (eb *Shot) Flip() {
}

func (eb *Shot) Clone() state.EntityBehavior {
	return &Shot{
		eb.Damage,
		eb.HitTraits,
	}
}

func (eb *Shot) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	return nil
}

func (eb *Shot) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime%2 == 1 {
		x, y := e.TilePos.XY()
		x += query.DXForward(e.IsFlipped)
		if !e.MoveDirectly(state.TilePosXY(x, y)) {
			e.PerTickState.IsPendingDeletion = true
			return
		}
	}

	for _, target := range query.EntitiesAt(s, e.TilePos) {
		if target.IsAlliedWithAnswerer == e.IsAlliedWithAnswerer {
			continue
		}

		var h state.Hit
		h.Traits = eb.HitTraits
		h.AddDamage(eb.Damage)
		target.Hit.Merge(h)

		e.PerTickState.IsPendingDeletion = true
		return
	}
}

func MakeShot(owner *state.Entity, pos state.TilePos, damage state.Damage, hitTraits state.HitTraits) *state.Entity {
	return &state.Entity{
		TilePos: pos,

		IsFlipped:            owner.IsFlipped,
		IsAlliedWithAnswerer: owner.IsAlliedWithAnswerer,

		Traits: state.EntityTraits{
			CanStepOnHoleLikeTiles: true,
			IgnoresTileEffects:     true,
			CannotFlinch:           true,
			IgnoresTileOwnership:   true,
		},

		BehaviorState: state.EntityBehaviorState{
			Behavior: &Shot{
				Damage:    damage,
				HitTraits: hitTraits,
			},
		},
	}
}
