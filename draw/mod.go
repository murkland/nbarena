package draw

import (
	"flag"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/murkland/pngsheet"
	"golang.org/x/image/font"
)

var (
	debugDrawImageNodeOutlines = flag.Bool("debug_draw_image_node_outlines", false, "draw image node outlines")
)

type Compositor struct {
	currentLayer *ebiten.Image
	layers       []*ebiten.Image
}

func NewCompositor(rect image.Rectangle, numLayers int) *Compositor {
	layers := make([]*ebiten.Image, numLayers)
	for i := 0; i < numLayers; i++ {
		layers[i] = ebiten.NewImage(rect.Dx(), rect.Dy())
	}
	return &Compositor{layers[0], layers}
}

func (c *Compositor) Bounds() image.Rectangle {
	return c.layers[0].Bounds()
}

func (c *Compositor) Clear() {
	for _, layer := range c.layers {
		layer.Clear()
	}
}

func (c *Compositor) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	for _, layer := range c.layers {
		screen.DrawImage(layer, opts)
	}
}

type Node interface {
	Draw(compositor *Compositor, opts *ebiten.DrawImageOptions)
}

type OptionsNode struct {
	Opts     ebiten.DrawImageOptions
	Layer    int
	Children []Node
}

func (n *OptionsNode) Draw(compositor *Compositor, opts *ebiten.DrawImageOptions) {
	o := n.Opts
	o.GeoM.Concat(opts.GeoM)
	o.ColorM.Concat(opts.ColorM)

	layer := compositor.currentLayer
	if n.Layer != 0 {
		layer = compositor.layers[n.Layer-1]
	}

	for _, c := range n.Children {
		if c == nil {
			continue
		}
		c.Draw(&Compositor{layer, compositor.layers}, &o)
	}
}

type ImageNode struct {
	Image *ebiten.Image
}

func makeDebugOutline(bounds image.Rectangle) *ebiten.Image {
	img := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	img.Fill(color.RGBA{255, 0, 0, 255})
	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			img.Set(x, y, color.Transparent)
		}
	}
	return img
}

func (n *ImageNode) Draw(compositor *Compositor, opts *ebiten.DrawImageOptions) {
	compositor.currentLayer.DrawImage(n.Image, opts)
	if *debugDrawImageNodeOutlines && n.Image.Bounds().Dx() > 0 && n.Image.Bounds().Dy() > 0 {
		opts2 := *opts
		opts2.ColorM.Reset()
		compositor.currentLayer.DrawImage(makeDebugOutline(n.Image.Bounds()), &opts2)
	}
}

func ImageWithAnimation(img *ebiten.Image, animation *pngsheet.Animation, t int) Node {
	if animation.IsLooping {
		t = t % len(animation.Frames)
	}
	frame := animation.Frames[t]
	return ImageWithOrigin(img.SubImage(frame.Rect).(*ebiten.Image), frame.Origin)
}

func ImageWithFrame(img *ebiten.Image, frame *pngsheet.Frame) Node {
	return ImageWithOrigin(img.SubImage(frame.Rect).(*ebiten.Image), frame.Origin)
}

func ImageWithOrigin(img *ebiten.Image, origin image.Point) Node {
	node := &OptionsNode{Children: []Node{
		&ImageNode{
			Image: img,
		},
	}}
	node.Opts.GeoM.Translate(float64(-origin.X), float64(-origin.Y))
	return node
}

type TextNode struct {
	Face font.Face
	Text string
}

func (n *TextNode) Draw(compositor *Compositor, opts *ebiten.DrawImageOptions) {
	o := *opts
	bounds := text.BoundString(n.Face, n.Text)
	o.GeoM = ebiten.GeoM{}
	o.GeoM.Translate(0, float64(-bounds.Min.Y))
	o.GeoM.Concat(opts.GeoM)
	text.DrawWithOptions(compositor.currentLayer, n.Text, n.Face, &o)
}
