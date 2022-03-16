package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/state/query"
)

type AirShot struct {
	Damage state.Damage
}

func (eb *AirShot) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{
		CanBeCountered: true,
	}
}

func (eb *AirShot) Clone() state.EntityBehavior {
	return &AirShot{
		eb.Damage,
	}
}

func (eb *AirShot) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == 6 {
		x, y := e.TilePos.XY()
		dx := query.DXForward(e.IsFlipped)
		s.AttachEntity(MakeShotEntity(e, state.TilePosXY(x+dx, y), &Shot{
			Damage: eb.Damage,
			Hit: state.Hit{
				Element:        state.ElementWind,
				ForcedMovement: state.ForcedMovement{Type: state.ForcedMovementTypeSmallDrag, Direction: e.Facing()},
				CanCounter:     true,
			},
			ExplosionDecorationType: bundle.DecorationTypeCannonExplosion,
		}))
	} else if e.BehaviorState.ElapsedTime == 21-1 {
		e.NextBehavior = &Idle{}
	}
}

func (eb *AirShot) Cleanup(e *state.Entity, s *state.State) {
}

func (eb *AirShot) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	rootNode.Children = append(rootNode.Children, draw.ImageWithAnimation(b.MegamanSprites.Image, b.MegamanSprites.RecoilShotAnimation, int(e.BehaviorState.ElapsedTime)))

	airShooterNode := &draw.OptionsNode{Layer: 6}
	airShooterNode.Opts.GeoM.Translate(float64(16), float64(-24))
	rootNode.Children = append(rootNode.Children, airShooterNode)
	airShooterNode.Children = append(airShooterNode.Children, draw.ImageWithAnimation(b.AirShooterSprites.Image, b.AirShooterSprites.Animations[0], int(e.BehaviorState.ElapsedTime)))
	return rootNode
}
