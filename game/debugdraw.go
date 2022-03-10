package game

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/state"
	"github.com/sanity-io/litter"
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
	rootNode := &draw.OptionsNode{Layer: 9}

	rootNode.Opts.ColorM.Scale(0.0, 1.0, 0.0, 1.0)
	rootNode.Opts.GeoM.Translate(12, 12)

	var entity *state.Entity
	if !g.cs.isAnswerer {
		entity = g.cs.dirtyState.Entities[g.cs.OffererEntityID]
	} else {
		entity = g.cs.dirtyState.Entities[g.cs.AnswererEntityID]
	}

	delay := g.medianDelay()
	rootNode.Children = append(rootNode.Children, &draw.TextNode{
		Face: mplusNormalFont,
		Text: fmt.Sprintf("delay: %6.2fms\n%s", float64(delay)/float64(time.Millisecond), litter.Options{
			HidePrivateFields: false,
		}.Sdump(entity)),
	})

	return rootNode
}
