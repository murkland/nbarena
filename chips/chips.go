package chips

import (
	"image"

	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/state"
)

var Chips = []*state.Chip{
	{
		Index:      0,
		Name:       "Cannon",
		BaseDamage: 40,
		OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
			e.NextBehavior = &behaviors.Cannon{Style: behaviors.CannonStyleCannon, Damage: damage}
		},
	},
	{
		Index:      1,
		Name:       "HiCannon",
		BaseDamage: 100,
		OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
			e.NextBehavior = &behaviors.Cannon{Style: behaviors.CannonStyleHiCannon, Damage: damage}
		},
	},
	{
		Index:      2,
		Name:       "M-Cannon",
		BaseDamage: 180,
		OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
			e.NextBehavior = &behaviors.Cannon{Style: behaviors.CannonStyleMCannon, Damage: damage}
		},
	},
	{
		Index:      3,
		Name:       "AirShot",
		BaseDamage: 20,
		OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
			e.NextBehavior = &behaviors.AirShot{Damage: damage}
		},
	},
	{
		Index:      4,
		Name:       "Vulcan1",
		BaseDamage: 10,
		OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
			e.NextBehavior = &behaviors.Vulcan{Shots: 3, Damage: damage, ExplosionDecorationType: bundle.DecorationTypeVulcanExplosion}
		},
	},
	{
		Index:      5,
		Name:       "Vulcan2",
		BaseDamage: 15,
		OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
			e.NextBehavior = &behaviors.Vulcan{Shots: 4, Damage: damage, ExplosionDecorationType: bundle.DecorationTypeVulcanExplosion}
		},
	},
	{
		Index:      6,
		Name:       "Vulcan3",
		BaseDamage: 20,
		OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
			e.NextBehavior = &behaviors.Vulcan{Shots: 5, Damage: damage, ExplosionDecorationType: bundle.DecorationTypeVulcanExplosion}
		},
	},
	{
		Index:      7,
		Name:       "SuprVulc",
		BaseDamage: 20,
		OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
			e.NextBehavior = &behaviors.Vulcan{Shots: 10, Damage: damage, ExplosionDecorationType: bundle.DecorationTypeSuperVulcanExplosion}
		},
	},
	{
		Index:      73,
		Name:       "WideBlde",
		BaseDamage: 150,
		OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
			e.NextBehavior = &behaviors.Sword{Damage: damage, Style: behaviors.SwordStyleBlade, Range: behaviors.SwordRangeWide}
		},
	},
	{
		Index:      79,
		Name:       "WindRack",
		BaseDamage: 140,
		OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
			e.NextBehavior = &behaviors.WindRack{Damage: damage}
		},
	},
	{
		Index:      160,
		Name:       "Recov300",
		BaseDamage: 0,
		OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
			e.HP += 300
			if e.HP > e.MaxHP {
				e.HP = e.MaxHP
			}

			// TODO: Trigger antirecovery.

			s.AttachSound(&state.Sound{
				Type: bundle.SoundTypeRecov,
			})
			s.AttachDecoration(&state.Decoration{
				Type:      bundle.DecorationTypeRecov,
				TilePos:   e.TilePos,
				Offset:    image.Point{0, 0},
				IsFlipped: e.IsFlipped,
			})

			e.ChipUseLockoutTimeLeft = 30
		},
	},
	{
		Index: 162,
		Name:  "AreaGrab",
		OnUse: func(s *state.State, e *state.Entity, damage state.Damage) {
			e.NextBehavior = &behaviors.AreaGrab{}
		},
	},
}
