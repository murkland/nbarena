package behaviors

import (
	"image"

	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
)

type SwordRange int

const (
	SwordRangeShort    SwordRange = 0
	SwordRangeWide     SwordRange = 1
	SwordRangeLong     SwordRange = 2
	SwordRangeVeryLong SwordRange = 3
)

type SwordStyle int

const (
	SwordStyleSword SwordStyle = 0
	SwordStyleBlade SwordStyle = 1
)

func swordSlashDecorationType(s SwordStyle, r SwordRange) bundle.DecorationType {
	switch s {
	case SwordStyleSword:
		switch r {
		case SwordRangeShort:
			return bundle.DecorationTypeNullShortSwordSlash
		case SwordRangeWide:
			return bundle.DecorationTypeNullWideSwordSlash
		case SwordRangeLong:
			return bundle.DecorationTypeNullLongSwordSlash
		case SwordRangeVeryLong:
			return bundle.DecorationTypeNullVeryLongSwordSlash
		}
	case SwordStyleBlade:
		switch r {
		case SwordRangeShort:
			return bundle.DecorationTypeNullShortBladeSlash
		case SwordRangeWide:
			return bundle.DecorationTypeNullWideBladeSlash
		case SwordRangeLong:
			return bundle.DecorationTypeNullLongBladeSlash
		case SwordRangeVeryLong:
			return bundle.DecorationTypeNullVeryLongBladeSlash
		}
	}
	return bundle.DecorationTypeNone
}

type Sword struct {
	Range  SwordRange
	Style  SwordStyle
	Damage state.Damage
}

func (eb *Sword) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{
		CanBeCountered: true,
	}
}

func (eb *Sword) Cleanup(e *state.Entity, s *state.State) {
}

func (eb *Sword) Clone() state.EntityBehavior {
	return &Sword{
		eb.Range,
		eb.Style,
		eb.Damage,
	}
}

func swordTargetEntities(s *state.State, e *state.Entity, r SwordRange) []state.TilePos {
	x, y := e.TilePos.XY()
	dx, _ := e.Facing().XY()
	var positions []state.TilePos
	positions = append(positions, state.TilePosXY(x+dx, y))

	switch r {
	case SwordRangeWide:
		positions = append(positions, state.TilePosXY(x+dx, y+1), state.TilePosXY(x+dx, y-1))
	case SwordRangeLong:
		positions = append(positions, state.TilePosXY(x+2*dx, y))
	case SwordRangeVeryLong:
		positions = append(positions, state.TilePosXY(x+2*dx, y), state.TilePosXY(x+3*dx, y))
	}

	return positions
}

func (eb *Sword) Step(e *state.Entity, s *state.State) {
	// Only hits while the slash is coming out.
	if e.BehaviorState.ElapsedTime == 9 {
		s.AttachDecoration(&state.Decoration{
			Type:      swordSlashDecorationType(eb.Style, eb.Range),
			TilePos:   e.TilePos,
			Offset:    image.Point{state.TileRenderedWidth, -16},
			IsFlipped: e.IsFlipped,
		})
		s.AttachSound(&state.Sound{
			Type: bundle.SoundTypeSwordSlash,
		})

		for _, pos := range swordTargetEntities(s, e, eb.Range) {
			var h state.Hit
			h.Flinch = true
			h.FlashTime = state.DefaultFlashTime
			h.Element = state.ElementSword
			h.SecondaryElementSword = true
			h.CanCounter = true
			h.RemovesFullSynchro = true
			h.AddDamage(eb.Damage)
			s.ApplyHit(e, pos, h)
		}
	} else if e.BehaviorState.ElapsedTime == 21-1 {
		e.NextBehavior = &Idle{}
	}
}

func (eb *Sword) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	rootNode.Children = append(rootNode.Children, draw.ImageWithAnimation(b.MegamanSprites.Image, b.MegamanSprites.SlashAnimation, int(e.BehaviorState.ElapsedTime)))

	swordNode := &draw.OptionsNode{Layer: 6}
	rootNode.Children = append(rootNode.Children, swordNode)
	swordNode.Children = append(swordNode.Children, draw.ImageWithAnimation(b.SwordSprites.Image, b.SwordSprites.BaseAnimation, int(e.BehaviorState.ElapsedTime)))

	return rootNode
}
