package behaviors

import (
	"image"
	"math/rand"

	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/state/query"
)

type Vulcan struct {
	Shots  int
	Damage state.Damage
}

func (eb *Vulcan) Flip() {
}

func (eb *Vulcan) Clone() state.EntityBehavior {
	return &Vulcan{
		eb.Shots,
		eb.Damage,
	}
}

func (eb *Vulcan) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{
		CanBeCountered: true,
	}
}

func (eb *Vulcan) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == state.Ticks(2+11*eb.Shots)-1 {
		e.NextBehavior = &Idle{}
		return
	}

	if (e.BehaviorState.ElapsedTime-2)%11 == 0 {
		x, y := e.TilePos.XY()
		dx := query.DXForward(e.IsFlipped)
		s.AddEntity(&state.Entity{
			TilePos: state.TilePosXY(x+dx, y),

			IsFlipped:            e.IsFlipped,
			IsAlliedWithAnswerer: e.IsAlliedWithAnswerer,

			Traits: state.EntityTraits{
				CanStepOnHoleLikeTiles: true,
				IgnoresTileEffects:     true,
				CannotFlinch:           true,
				IgnoresTileOwnership:   true,
			},

			BehaviorState: state.EntityBehaviorState{
				Behavior: &vulcanShot{e.ID(), eb.Damage},
			},
		})
	}
}

func (eb *Vulcan) Cleanup(e *state.Entity, s *state.State) {
}

func (eb *Vulcan) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	var megamanImageNode draw.Node

	if e.BehaviorState.ElapsedTime < 2 {
		megamanImageNode = draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.HoldInFrontAnimation.Frames[int(e.BehaviorState.ElapsedTime)%len(b.MegamanSprites.HoldInFrontAnimation.Frames)])
	} else {
		megamanImageNode = draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.GattlingAnimation.Frames[int(e.BehaviorState.ElapsedTime-2)%len(b.MegamanSprites.GattlingAnimation.Frames)])
	}
	rootNode.Children = append(rootNode.Children, megamanImageNode)

	vulcanNode := &draw.OptionsNode{Layer: 6}
	rootNode.Children = append(rootNode.Children, vulcanNode)
	vulcanNode.Opts.GeoM.Translate(float64(24), float64(-24))
	var vulcanImageNode draw.Node
	if e.BehaviorState.ElapsedTime < 2 {
		vulcanImageNode = draw.ImageWithFrame(b.VulcanSprites.Image, b.VulcanSprites.Animations[0].Frames[e.BehaviorState.ElapsedTime])
	} else {
		vulcanFrames := b.VulcanSprites.Animations[1].Frames
		vulcanImageNode = draw.ImageWithFrame(b.VulcanSprites.Image, vulcanFrames[int(e.BehaviorState.ElapsedTime-2)%len(vulcanFrames)])
	}
	vulcanNode.Children = append(vulcanNode.Children, vulcanImageNode)

	return rootNode
}

type vulcanShot struct {
	Owner  state.EntityID
	Damage state.Damage
}

func (eb *vulcanShot) Flip() {
}

func (eb *vulcanShot) Clone() state.EntityBehavior {
	return &vulcanShot{
		eb.Owner,
		eb.Damage,
	}
}

func (eb *vulcanShot) Traits(e *state.Entity) state.EntityBehaviorTraits {
	return state.EntityBehaviorTraits{}
}

func (eb *vulcanShot) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	return nil
}

func (eb *vulcanShot) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == 0 {
		return
	}

	for _, target := range query.EntitiesAt(s, e.TilePos) {
		if target.IsAlliedWithAnswerer == e.IsAlliedWithAnswerer || target.FlashingTimeLeft > 0 {
			continue
		}

		var h state.Hit
		h.Traits = state.HitTraits{
			Flinch: true,
		}
		h.AddDamage(eb.Damage)
		h.Traits.Element = state.ElementNull
		state.MaybeApplyCounter(target, s.Entities[eb.Owner], &h)
		target.Hit.Merge(h)

		rand := rand.New(s.RandSource)

		xOff := rand.Intn(state.TileRenderedWidth / 4)
		yOff := -rand.Intn(state.TileRenderedHeight)

		s.AddDecoration(&state.Decoration{
			Type:      bundle.DecorationTypeBusterExplosion,
			TilePos:   e.TilePos,
			Offset:    image.Point{xOff + rand.Intn(2) - 4, yOff + rand.Intn(2) - 4},
			IsFlipped: e.IsFlipped,
		})

		s.AddDecoration(&state.Decoration{
			ElapsedTime: -1,
			Type:        bundle.DecorationTypeBusterExplosion,
			TilePos:     e.TilePos,
			Offset:      image.Point{xOff + rand.Intn(2) - 4, yOff + rand.Intn(2) - 4},
			IsFlipped:   e.IsFlipped,
		})

		s.AddDecoration(&state.Decoration{
			Type:      bundle.DecorationTypePiercingExplosion,
			TilePos:   e.TilePos,
			Offset:    image.Point{xOff + rand.Intn(2) - 4, yOff + rand.Intn(2) - 4},
			IsFlipped: e.IsFlipped,
		})

		e.IsPendingDestruction = true
		return
	}

	x, y := e.TilePos.XY()
	x += query.DXForward(e.IsFlipped)
	if !e.MoveDirectly(state.TilePosXY(x, y)) {
		e.IsPendingDestruction = true
		return
	}
}

func (eb *vulcanShot) Cleanup(e *state.Entity, s *state.State) {
}
