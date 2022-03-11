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

func (eb *Cannon) Flip() {
}

func (eb *Cannon) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{
		CanBeCountered: true,
	}
}

func (eb *Cannon) Clone() state.EntityBehavior {
	return &Cannon{
		eb.Style,
		eb.Damage,
	}
}

func (eb *Cannon) Step(e *state.Entity, s *state.State) {
	// TODO: Counter timing.

	if e.BehaviorState.ElapsedTime == 16 {
		x, y := e.TilePos.XY()
		dx := query.DXForward(e.IsFlipped)
		s.AddEntity(MakeShot(e, state.TilePosXY(x+dx, y), e.MakeDamageAndConsume(eb.Damage), state.HitTraits{
			Flinch:    true,
			Counters:  true,
			FlashTime: state.DefaultFlashTime,
		}))
	} else if e.BehaviorState.ElapsedTime == 33-1 {
		e.NextBehavior = &Idle{}
	}
}

func (eb *Cannon) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	if e.BehaviorState.ElapsedTime >= 29 {
		return draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.BraceAnimation.Frames[int(e.BehaviorState.ElapsedTime-29)])
	}

	rootNode := &draw.OptionsNode{}
	rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.CannonAnimation.Frames[e.BehaviorState.ElapsedTime]))

	cannonNode := &draw.OptionsNode{Layer: 6}
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
	cannonNode.Children = append(cannonNode.Children, draw.ImageWithFrame(img, b.CannonSprites.Animation.Frames[e.BehaviorState.ElapsedTime]))
	return rootNode
}
