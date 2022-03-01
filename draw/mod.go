package draw

import (
	"flag"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/yumland/pngsheet"
	"golang.org/x/image/font"
)

var (
	debugDrawImageNodeOutlines = flag.Bool("debug_draw_image_node_outlines", false, "draw image node outlines")
)

type Node interface {
	Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions)
}

type OptionsNode struct {
	Opts     ebiten.DrawImageOptions
	Children []Node
}

func (n *OptionsNode) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	o := n.Opts
	o.GeoM.Concat(opts.GeoM)
	o.ColorM.Concat(opts.ColorM)

	for _, c := range n.Children {
		if c == nil {
			continue
		}
		c.Draw(screen, &o)
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

func (n *ImageNode) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(n.Image, opts)
	if *debugDrawImageNodeOutlines {
		opts2 := *opts
		opts2.ColorM.Reset()
		screen.DrawImage(makeDebugOutline(n.Image.Bounds()), &opts2)
	}
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

func (n *TextNode) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	o := *opts
	bounds := text.BoundString(n.Face, n.Text)
	o.GeoM = ebiten.GeoM{}
	o.GeoM.Translate(0, float64(-bounds.Min.Y))
	o.GeoM.Concat(opts.GeoM)
	text.DrawWithOptions(screen, n.Text, n.Face, &o)
}
