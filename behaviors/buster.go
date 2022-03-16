package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/state/query"
)

type Buster struct {
	Speed        int
	BaseDamage   int
	IsPowerShot  bool
	isJammed     bool
	cooldownTime state.Ticks
}

func (eb *Buster) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *Buster) realElapsedTime(e *state.Entity) state.Ticks {
	t := e.BehaviorState.ElapsedTime
	if eb.IsPowerShot {
		t -= 5
	}
	return t
}

func (eb *Buster) Clone() state.EntityBehavior {
	return &Buster{
		eb.Speed,
		eb.BaseDamage,
		eb.IsPowerShot,
		eb.isJammed,
		eb.cooldownTime,
	}
}

// Buster cooldown time:
var busterCooldownDurations = [][]state.Ticks{
	// d = 1, 2, 3, 4, 5, 6
	{5, 9, 13, 17, 21, 25}, // Lv1
	{4, 8, 11, 15, 18, 21}, // Lv2
	{4, 7, 10, 13, 16, 18}, // Lv3
	{3, 5, 7, 9, 11, 13},   // Lv4
	{3, 4, 5, 6, 7, 8},     // Lv5
}

func (eb *Buster) Step(e *state.Entity, s *state.State) {
	realElapsedTime := eb.realElapsedTime(e)

	if realElapsedTime == 5+eb.cooldownTime-1 {
		e.NextBehavior = &Idle{}
		return
	}

	if realElapsedTime == 1 {
		_, d := query.FindNearestEntity(s, e.ID(), e.TilePos, e.IsAlliedWithAnswerer, e.IsFlipped, query.HorizontalDistance)
		eb.cooldownTime = busterCooldownDurations[eb.Speed][d]

		x, y := e.TilePos.XY()
		dx, _ := e.Facing().XY()

		damage := eb.BaseDamage
		if eb.IsPowerShot {
			damage *= 10
		}
		decorationType := bundle.DecorationTypeBusterExplosion
		if eb.IsPowerShot {
			decorationType = bundle.DecorationTypeBusterPowerShotExplosion
		}
		s.AttachSound(&state.Sound{
			Type: bundle.SoundTypeBuster,
		})
		s.AttachEntity(MakeShotEntity(e, state.TilePosXY(x+dx, y), &Shot{
			Damage: state.Damage{Base: damage},
			Hit: state.Hit{
				Element: state.ElementNull,
			},
			ExplosionDecorationType: decorationType,
		}))
	}

	if e.Intent.Direction != state.DirectionNone && realElapsedTime >= 5 {
		dir := e.Intent.Direction
		if e.ConfusedTimeLeft > 0 {
			dir = dir.FlipH().FlipV()
		}

		x, y := e.TilePos.XY()
		dx, dy := dir.XY()

		if e.StartMove(state.TilePosXY(x+dx, y+dy), s) {
			e.NextBehavior = &Teleport{}
		}
	}

}

func (eb *Buster) Cleanup(e *state.Entity, s *state.State) {
}

func (eb *Buster) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	realElapsedTime := eb.realElapsedTime(e)

	if realElapsedTime < 0 {
		return draw.ImageWithAnimation(b.MegamanSprites.Image, b.MegamanSprites.IdleAnimation, int(e.ElapsedTime))
	}

	rootNode := &draw.OptionsNode{}
	rootNode.Children = append(rootNode.Children, draw.ImageWithAnimation(b.MegamanSprites.Image, b.MegamanSprites.BusterAnimation, int(realElapsedTime)))
	rootNode.Children = append(rootNode.Children, draw.ImageWithAnimation(b.BusterSprites.Image, b.BusterSprites.BaseAnimation, int(realElapsedTime)))

	if !eb.isJammed {
		muzzleFlashAnimTime := int(realElapsedTime) - 1
		if muzzleFlashAnimTime > 0 && muzzleFlashAnimTime < len(b.MuzzleFlashSprites.Animations[0].Frames) {
			muzzleFlashNode := &draw.OptionsNode{Layer: 7}
			muzzleFlashNode.Opts.GeoM.Translate(float64(state.TileRenderedWidth), float64(-26))
			muzzleFlashNode.Children = append(muzzleFlashNode.Children, draw.ImageWithAnimation(b.MuzzleFlashSprites.Image, b.MuzzleFlashSprites.Animations[0], muzzleFlashAnimTime))
			rootNode.Children = append(rootNode.Children, muzzleFlashNode)
		}
	}

	return rootNode
}
