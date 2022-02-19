package draw

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Node interface {
	Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions)
}

type OptionsNode struct {
	Opts     *ebiten.DrawImageOptions
	Children []Node
}

func (n *OptionsNode) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	o := *n.Opts
	o.GeoM.Concat(opts.GeoM)
	o.ColorM.Concat(opts.ColorM)

	for _, c := range n.Children {
		c.Draw(screen, &o)
	}
}

type ImageNode struct {
	Image *ebiten.Image
}

func (n *ImageNode) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(n.Image, opts)
}

type TextNode struct {
	Face font.Face
	Text string
}

func (n *TextNode) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	o := *opts
	bounds := text.BoundString(n.Face, n.Text)
	o.GeoM = ebiten.GeoM{}
	o.GeoM.Translate(0, float64(-bounds.Min.Y))
	o.GeoM.Concat(opts.GeoM)
	text.DrawWithOptions(screen, n.Text, n.Face, &o)
}
