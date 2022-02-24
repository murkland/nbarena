package bundle

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yumland/pngsheet"
)

type Sheet struct {
	Info  *pngsheet.Info
	Image image.Image
}

func loadSheet(filename string) (*Sheet, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("%w while loading %s", err, filename)
	}
	defer f.Close()

	img, info, err := pngsheet.Load(f)
	if err != nil {
		return nil, fmt.Errorf("%w while loading %s", err, filename)
	}

	return &Sheet{info, img}, nil
}

type Battletiles struct {
	Info          *pngsheet.Info
	OffererTiles  *ebiten.Image
	AnswererTiles *ebiten.Image
}

type Megaman struct {
	Info        *pngsheet.Info
	BaseSprites *ebiten.Image
}

func loadBattleTiles() (*Battletiles, error) {
	sheet, err := loadSheet("assets/battletiles.png")
	if err != nil {
		return nil, err
	}

	img := sheet.Image.(*image.Paletted)
	offererImg := ebiten.NewImageFromImage(img)

	img.Palette = sheet.Info.SuggestedPalettes["alt"]
	answererImg := ebiten.NewImageFromImage(img)

	return &Battletiles{sheet.Info, offererImg, answererImg}, nil
}

func loadMegaman() (*Megaman, error) {
	sheet, err := loadSheet("assets/sprites/0000.png")
	if err != nil {
		return nil, err
	}

	img := sheet.Image.(*image.Paletted)

	return &Megaman{sheet.Info, ebiten.NewImageFromImage(img)}, nil
}

type Bundle struct {
	Battletiles *Battletiles
	Megaman     *Megaman
}

func Load() (*Bundle, error) {
	b := &Bundle{}

	{
		battletiles, err := loadBattleTiles()
		if err != nil {
			return nil, err
		}
		b.Battletiles = battletiles
	}

	{
		megaman, err := loadMegaman()
		if err != nil {
			return nil, err
		}
		b.Megaman = megaman
	}

	return b, nil
}
