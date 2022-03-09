package behaviors

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/input"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/state/query"
	"github.com/murkland/pngsheet"
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

func slashAnimation(b *bundle.Bundle, r SwordRange) *pngsheet.Animation {
	switch r {
	case SwordRangeShort:
		return b.SlashSprites.ShortAnimation
	case SwordRangeWide:
		return b.SlashSprites.WideAnimation
	case SwordRangeLong:
		return b.SlashSprites.LongAnimation
	case SwordRangeVeryLong:
		return b.SlashSprites.VeryLongAnimation
	}
	return nil
}

type Sword struct {
	Range  SwordRange
	Style  SwordStyle
	Damage int
}

func (eb *Sword) Clone() state.EntityBehavior {
	return &Sword{
		eb.Range,
		eb.Style,
		eb.Damage,
	}
}
func (eb *Sword) ApplyIntent(e *state.Entity, s *state.State, intent input.Intent) {
}

func swordTargetEntities(s *state.State, e *state.Entity, r SwordRange) []*state.Entity {
	x, y := e.TilePos.XY()
	dx := query.DXForward(e.IsFlipped)
	var entities []*state.Entity
	entities = append(entities, query.EntitiesAt(s, state.TilePosXY(x+dx, y))...)

	switch r {
	case SwordRangeWide:
		entities = append(entities, query.EntitiesAt(s, state.TilePosXY(x+dx, y+1))...)
		entities = append(entities, query.EntitiesAt(s, state.TilePosXY(x+dx, y-1))...)
	case SwordRangeLong:
		entities = append(entities, query.EntitiesAt(s, state.TilePosXY(x+2*dx, y))...)
	case SwordRangeVeryLong:
		entities = append(entities, query.EntitiesAt(s, state.TilePosXY(x+2*dx, y))...)
		entities = append(entities, query.EntitiesAt(s, state.TilePosXY(x+3*dx, y))...)
	}

	return entities
}

func (eb *Sword) Step(e *state.Entity, s *state.State) {
	// Only hits while the slash is coming out.
	if e.BehaviorElapsedTime() == 9 {
		for _, target := range swordTargetEntities(s, e, eb.Range) {
			if target.FlashingTimeLeft == 0 {
				var h state.Hit
				h.Flinch = true
				h.Counters = true
				h.FlashTime = state.DefaultFlashTime
				h.AddDamage(e.MakeDamageAndConsume(eb.Damage))
				target.PerTickState.Hit.Merge(h)
			}
		}
	} else if e.BehaviorElapsedTime() == 21 {
		e.SetBehavior(&Idle{})
	}
}

func (eb *Sword) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	rootNode.Children = append(rootNode.Children, draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.SlashAnimation.Frames[e.BehaviorElapsedTime()]))

	swordNode := &draw.OptionsNode{Layer: 8}
	rootNode.Children = append(rootNode.Children, swordNode)
	swordNode.Children = append(swordNode.Children, draw.ImageWithFrame(b.SwordSprites.Image, b.SwordSprites.BaseAnimation.Frames[e.BehaviorElapsedTime()]))

	if e.BehaviorElapsedTime() >= 9 && e.BehaviorElapsedTime() < 19 {
		slashNode := &draw.OptionsNode{Layer: 7}
		slashNode.Opts.GeoM.Translate(float64(state.TileRenderedWidth), float64(-16))
		rootNode.Children = append(rootNode.Children, slashNode)

		slashAnim := slashAnimation(b, eb.Range)
		var img *ebiten.Image
		switch eb.Style {
		case SwordStyleSword:
			img = b.SlashSprites.SwordImage
		case SwordStyleBlade:
			img = b.SlashSprites.BladeImage
		}
		slashNode.Children = append(slashNode.Children, draw.ImageWithFrame(img, slashAnim.Frames[e.BehaviorElapsedTime()-9]))
	}

	return rootNode
}
