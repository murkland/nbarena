package styledtext

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/murkland/nbarena/draw"
	"golang.org/x/image/font"
)

type Anchor int

const (
	AnchorLeft   Anchor = 0b0000
	AnchorCenter Anchor = 0b0001
	AnchorRight  Anchor = 0b0010
	AnchorTop    Anchor = 0b1000
	AnchorMiddle Anchor = 0b0100
	AnchorBottom Anchor = 0b0000
)

type Border int

const (
	BorderNone         = 0b000000000
	BorderLeftTop      = 0b000000001
	BorderCenterTop    = 0b000000010
	BorderRightTop     = 0b000000100
	BorderLeftMiddle   = 0b000001000
	BorderRightMiddle  = 0b000100000
	BorderLeftBottom   = 0b001000000
	BorderCenterBottom = 0b010000000
	BorderRightBottom  = 0b100000000
	BorderAll          = 0b111111111
)

type Span struct {
	Text       string
	Background *ebiten.Image
}

func MakeNode(spans []Span, anchor Anchor, face font.Face, border Border, borderColor color.RGBA) draw.Node {
	if len(spans) == 0 {
		return nil
	}

	spanBounds := make([]image.Rectangle, len(spans))
	for i, span := range spans {
		spanBounds[i] = text.BoundString(face, span.Text)
	}

	width := 0
	height := 0

	for _, bounds := range spanBounds {
		width += bounds.Max.X
		if height < bounds.Dy() {
			height = bounds.Dy()
		}
	}

	textOnlyImage := ebiten.NewImage(width, height)
	{
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(0), float64(-spanBounds[0].Min.Y))
		for i, span := range spans {
			bounds := spanBounds[i]
			text.DrawWithOptions(textOnlyImage, span.Text, face, opts)
			opts.GeoM.Translate(float64(bounds.Dx()), float64(0))
		}
	}

	finalImage := ebiten.NewImage(textOnlyImage.Bounds().Dx(), textOnlyImage.Bounds().Dy())
	{
		var advanceX int
		for i, span := range spans {
			bounds := spanBounds[i]
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Scale(float64(bounds.Dx()), 1.0)
			opts.GeoM.Translate(float64(advanceX), 0)
			finalImage.DrawImage(span.Background, opts)
			advanceX += bounds.Dx()
		}
	}

	{
		opts := &ebiten.DrawImageOptions{}
		opts.CompositeMode = ebiten.CompositeModeDestinationIn
		finalImage.DrawImage(textOnlyImage, opts)
	}

	origin := spanBounds[0].Min
	if anchor&AnchorCenter != 0 {
		origin.X = -width / 2
	} else if anchor&AnchorRight != 0 {
		origin.X = -width
	}

	if anchor&AnchorMiddle != 0 {
		origin.Y = origin.Y / 2
	} else if anchor&AnchorTop != 0 {
		origin.Y = 0
	}

	rootNode := &draw.OptionsNode{}
	rootNode.Opts.GeoM.Translate(float64(origin.X), float64(origin.Y))

	textNode := &draw.OptionsNode{}
	rootNode.Children = append(rootNode.Children, textNode)

	borderNode := &draw.OptionsNode{}
	textNode.Children = append(textNode.Children, borderNode)
	borderNode.Opts.ColorM.Translate(1.0, 1.0, 1.0, 0.0)
	borderNode.Opts.ColorM.Scale(float64(borderColor.R)/float64(0xff), float64(borderColor.G)/float64(0xff), float64(borderColor.B)/float64(0xff), float64(borderColor.A)/float64(0xff))

	if border&BorderLeftTop != 0 {
		edgeNode := &draw.OptionsNode{}
		borderNode.Children = append(borderNode.Children, edgeNode)
		edgeNode.Opts.GeoM.Translate(-1, -1)
		edgeNode.Children = append(edgeNode.Children, draw.ImageWithOrigin(finalImage, image.Point{}))
	}

	if border&BorderCenterTop != 0 {
		edgeNode := &draw.OptionsNode{}
		borderNode.Children = append(borderNode.Children, edgeNode)
		edgeNode.Opts.GeoM.Translate(0, -1)
		edgeNode.Children = append(edgeNode.Children, draw.ImageWithOrigin(finalImage, image.Point{}))
	}

	if border&BorderRightTop != 0 {
		edgeNode := &draw.OptionsNode{}
		borderNode.Children = append(borderNode.Children, edgeNode)
		edgeNode.Opts.GeoM.Translate(1, -1)
		edgeNode.Children = append(edgeNode.Children, draw.ImageWithOrigin(finalImage, image.Point{}))
	}

	if border&BorderLeftMiddle != 0 {
		edgeNode := &draw.OptionsNode{}
		borderNode.Children = append(borderNode.Children, edgeNode)
		edgeNode.Opts.GeoM.Translate(-1, 0)
		edgeNode.Children = append(edgeNode.Children, draw.ImageWithOrigin(finalImage, image.Point{}))
	}

	if border&BorderRightMiddle != 0 {
		edgeNode := &draw.OptionsNode{}
		borderNode.Children = append(borderNode.Children, edgeNode)
		edgeNode.Opts.GeoM.Translate(1, 0)
		edgeNode.Children = append(edgeNode.Children, draw.ImageWithOrigin(finalImage, image.Point{}))
	}

	if border&BorderLeftBottom != 0 {
		edgeNode := &draw.OptionsNode{}
		borderNode.Children = append(borderNode.Children, edgeNode)
		edgeNode.Opts.GeoM.Translate(-1, 1)
		edgeNode.Children = append(edgeNode.Children, draw.ImageWithOrigin(finalImage, image.Point{}))
	}

	if border&BorderCenterBottom != 0 {
		edgeNode := &draw.OptionsNode{}
		borderNode.Children = append(borderNode.Children, edgeNode)
		edgeNode.Opts.GeoM.Translate(0, 1)
		edgeNode.Children = append(edgeNode.Children, draw.ImageWithOrigin(finalImage, image.Point{}))
	}

	if border&BorderRightBottom != 0 {
		edgeNode := &draw.OptionsNode{}
		borderNode.Children = append(borderNode.Children, edgeNode)
		edgeNode.Opts.GeoM.Translate(1, 1)
		edgeNode.Children = append(edgeNode.Children, draw.ImageWithOrigin(finalImage, image.Point{}))
	}

	textNode.Children = append(textNode.Children, &draw.ImageNode{Image: finalImage})

	return rootNode
}
