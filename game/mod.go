package game

import "github.com/hajimehoshi/ebiten/v2"

const RenderWidth = 240
const RenderHeight = 160

type Game struct {
	rootGeoM ebiten.GeoM
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	scaleFactor := outsideWidth / RenderWidth
	if s := outsideHeight / RenderHeight; s < scaleFactor {
		scaleFactor = s
	}

	insideWidth := RenderWidth * scaleFactor
	insideHeight := RenderHeight * scaleFactor

	g.rootGeoM = ebiten.GeoM{}
	g.rootGeoM.Scale(float64(scaleFactor), float64(scaleFactor))
	g.rootGeoM.Translate(float64(outsideWidth-insideWidth)/2, float64(outsideHeight-insideHeight)/2)

	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
}

func (g *Game) Update() error {
	return nil
}
