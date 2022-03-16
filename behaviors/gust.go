package behaviors

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type GustStyle int

const (
	GustStyleNone GustStyle = 0
	GustStyleWind GustStyle = 1
	GustStyleFan  GustStyle = 2
)

type Gust struct {
	Owner state.EntityID
	Style GustStyle
}

func (eb *Gust) Clone() state.EntityBehavior {
	return &Gust{
		eb.Owner,
		eb.Style,
	}
}

func (eb *Gust) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *Gust) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	if e.IsPendingDestruction {
		return nil
	}

	var gustImg *ebiten.Image = nil
	switch eb.Style {
	case GustStyleWind:
		gustImg = b.GustSprites.WindImage
	case GustStyleFan:
		gustImg = b.GustSprites.FanImage
	}

	dx := (int(e.BehaviorState.ElapsedTime)-1+4)%4 - 2

	rootNode := &draw.OptionsNode{}
	rootNode.Opts.GeoM.Translate(float64(dx*state.TileRenderedWidth/4), 0)
	rootNode.Children = append(rootNode.Children, draw.ImageWithAnimation(gustImg, b.GustSprites.Animation, int(e.BehaviorState.ElapsedTime)))
	return rootNode
}

func (eb *Gust) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime%4 == 1 {
		x, y := e.TilePos.XY()
		dx, _ := e.Facing().XY()
		if !e.MoveDirectly(state.TilePosXY(x+dx, y)) {
			e.IsPendingDestruction = true
			return
		}
	}

	var h state.Hit
	h.Element = state.ElementWind
	h.ForcedMovement = state.ForcedMovement{Type: state.ForcedMovementTypeSlide, Direction: e.Facing()}
	s.ApplyHit(s.Entities[eb.Owner], e.TilePos, h)
}

func (eb *Gust) Cleanup(e *state.Entity, s *state.State) {
}
