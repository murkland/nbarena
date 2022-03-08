package behaviors

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/state/query"
)

type CannonStyle int

const (
	CannonStyleCannon   CannonStyle = 0
	CannonStyleHiCannon CannonStyle = 1
	CannonStyleMCannon  CannonStyle = 2
)

type Cannon struct {
	Style  CannonStyle
	Damage int
}

func (eb *Cannon) Clone() state.EntityBehavior {
	return &Cannon{
		eb.Style,
		eb.Damage,
	}
}

func (eb *Cannon) Step(e *state.Entity, s *state.State) {
	if e.BehaviorElapsedTime() == 16 {
		x, y := e.TilePos.XY()
		if e.IsFlipped {
			x--
		} else {
			x++
		}

		shot := &state.Entity{
			TilePos:       state.TilePosXY(x, y),
			FutureTilePos: state.TilePosXY(x, y),

			IsFlipped:            e.IsFlipped,
			IsAlliedWithAnswerer: e.IsAlliedWithAnswerer,

			Traits: state.EntityTraits{
				CanStepOnHoleLikeTiles: true,
				IgnoresTileEffects:     true,
				CannotFlinch:           true,
				IgnoresTileOwnership:   true,
			},
		}
		shot.SetBehavior(&cannonShot{e.MakeDamageAndConsume(eb.Damage)})
		s.AddEntity(shot)
	} else if e.BehaviorElapsedTime() == 33 {
		e.SetBehavior(&Idle{})
	}
}

func (eb *Cannon) Interrupts(e *state.Entity) state.EntityBehaviorInterrupts {
	return state.EntityBehaviorInterrupts{}
}

func (eb *Cannon) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	if e.BehaviorElapsedTime() < 29 {
		rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.CannonAnimation.Frames[e.BehaviorElapsedTime()]))

		cannonNode := &draw.OptionsNode{Layer: 8}
		cannonNode.Opts.GeoM.Translate(float64(16), float64(-24))
		rootNode.Children = append(rootNode.Children, cannonNode)
		var img *ebiten.Image
		switch eb.Style {
		case CannonStyleCannon:
			img = b.CannonSprites.CannonImage
		case CannonStyleHiCannon:
			img = b.CannonSprites.HiCannonImage
		case CannonStyleMCannon:
			img = b.CannonSprites.MCannonImage
		}
		cannonNode.Children = append(cannonNode.Children, draw.ImageWithFrame(img, b.CannonSprites.Animation.Frames[e.BehaviorElapsedTime()]))
	} else {
		return draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.BraceAnimation.Frames[int(e.BehaviorElapsedTime()-29)])
	}

	return rootNode
}

type cannonShot struct {
	damage state.Damage
}

func (eb *cannonShot) Clone() state.EntityBehavior {
	return &cannonShot{
		eb.damage,
	}
}

func (eb *cannonShot) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	return nil
}

func (eb *cannonShot) Interrupts(e *state.Entity) state.EntityBehaviorInterrupts {
	return state.EntityBehaviorInterrupts{}
}

func (eb *cannonShot) Step(e *state.Entity, s *state.State) {
	if e.BehaviorElapsedTime()%2 == 1 {
		x, y := e.TilePos.XY()
		x += query.DXForward(e.IsFlipped)
		if !e.StartMove(state.TilePosXY(x, y), &s.Field) {
			e.PerTickState.IsPendingDeletion = true
			return
		}
		e.FinishMove()
	}

	for _, target := range query.EntitiesAt(s, e.TilePos) {
		if target.IsAlliedWithAnswerer == e.IsAlliedWithAnswerer {
			continue
		}

		var h state.Hit
		h.Flinch = true
		h.Counters = true
		h.FlashTime = state.DefaultFlashTime
		h.AddDamage(eb.damage)
		target.PerTickState.Hit.Merge(h)

		e.PerTickState.IsPendingDeletion = true
		return
	}
}
