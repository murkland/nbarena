package game

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/sanity-io/litter"
	"github.com/yumland/yumbattle/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	mplusNormalFont font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) makeDebugDrawNode() draw.Node {
	colorM := ebiten.ColorM{}
	colorM.Translate(1.0, 1.0, 1.0, 0.0)

	geoM := ebiten.GeoM{}
	geoM.Translate(12, 12)

	delay := g.medianDelay()
	return &draw.OptionsNode{
		Opts: ebiten.DrawImageOptions{
			GeoM:   geoM,
			ColorM: colorM,
		},
		Children: []draw.Node{
			&draw.TextNode{Face: mplusNormalFont, Text: fmt.Sprintf("delay: %6.2fms\n%s", float64(delay)/float64(time.Millisecond), litter.Options{
				HidePrivateFields: false,
				FieldFilter: func(sf reflect.StructField, v reflect.Value) bool {
					if sf.Type.Implements(reflect.TypeOf((*rand.Source)(nil)).Elem()) {
						return false
					}
					return true
				},
			}.Sdump(g.cs.dirtyState))},
		},
	}
}
