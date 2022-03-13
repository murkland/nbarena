package state

import (
	"image"

	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
)

type DecorationID int

type Decoration struct {
	id DecorationID

	Delay       Ticks
	ElapsedTime Ticks

	RunsInTimestop bool

	IsFlipped bool

	Type bundle.DecorationType

	TilePos TilePos
	Offset  image.Point
}

func (d *Decoration) ID() DecorationID {
	return d.id
}

func (d *Decoration) Flip() {
	d.TilePos = d.TilePos.Flipped()
	d.Offset.X = -d.Offset.X
	d.IsFlipped = !d.IsFlipped
}

func (d *Decoration) Clone() *Decoration {
	return &Decoration{
		d.id,
		d.Delay, d.ElapsedTime,
		d.RunsInTimestop,
		d.IsFlipped,
		d.Type,
		d.TilePos, d.Offset,
	}
}

func (d *Decoration) Step() {
	d.ElapsedTime++
}

func (d *Decoration) Appearance(b *bundle.Bundle) draw.Node {
	if d.ElapsedTime < 0 {
		return nil
	}

	rootNode := &draw.OptionsNode{}
	x, y := d.TilePos.XY()

	rootNode.Opts.GeoM.Translate(
		float64((x-1)*TileRenderedWidth+TileRenderedWidth/2+d.Offset.X),
		float64((y-1)*TileRenderedHeight+TileRenderedHeight/2+d.Offset.Y),
	)

	spriteNode := &draw.OptionsNode{}
	rootNode.Children = append(rootNode.Children, spriteNode)

	sprite := b.DecorationSprites[d.Type]
	spriteNode.Children = append(spriteNode.Children, draw.ImageWithAnimation(sprite.Image, sprite.Animation, int(d.ElapsedTime)))
	if d.IsFlipped {
		spriteNode.Opts.GeoM.Scale(-1, 1)
	}

	return rootNode
}
