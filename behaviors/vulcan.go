package behaviors

import (
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/state/query"
)

type Vulcan struct {
	Shots  int
	Damage int
}

func (eb *Vulcan) Flip() {
}

func (eb *Vulcan) Clone() state.EntityBehavior {
	return &Vulcan{
		eb.Shots,
		eb.Damage,
	}
}

func (eb *Vulcan) Step(e *state.Entity, s *state.State) {
	if e.BehaviorState.ElapsedTime == state.Ticks(2+11*eb.Shots) {
		e.ReplaceBehavior(&Idle{}, s)
		return
	}

	if (e.BehaviorState.ElapsedTime-2)%11 == 0 {
		x, y := e.TilePos.XY()
		dx := query.DXForward(e.IsFlipped)
		s.AddEntity(MakeShot(e, state.TilePosXY(x+dx, y), e.MakeDamageAndConsume(eb.Damage), state.HitTraits{
			Flinch:   true,
			Counters: true,
		}))
	}
}

func (eb *Vulcan) Appearance(e *state.Entity, b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	var megamanImageNode draw.Node

	if e.BehaviorState.ElapsedTime < 2 {
		megamanImageNode = draw.ImageWithFrame(b.MegamanSprites.Image, b.MegamanSprites.GattlingAnimation.Frames[0])
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
