package bundle

import (
	"context"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yumland/moreio"
	"github.com/yumland/pngsheet"
	"golang.org/x/sync/errgroup"
)

type Sheet struct {
	Info  *pngsheet.Info
	Image image.Image
}

func loadSheet(ctx context.Context, filename string) (*Sheet, error) {
	f, err := moreio.Open(ctx, filename)
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

func loadBattleTiles(ctx context.Context) (*Battletiles, error) {
	sheet, err := loadSheet(ctx, "assets/battletiles.png")
	if err != nil {
		return nil, err
	}

	img := sheet.Image.(*image.Paletted)
	offererImg := ebiten.NewImageFromImage(img)

	img.Palette = sheet.Info.SuggestedPalettes["alt"]
	answererImg := ebiten.NewImageFromImage(img)

	return &Battletiles{sheet.Info, offererImg, answererImg}, nil
}

func loadMegaman(ctx context.Context) (*Megaman, error) {
	sheet, err := loadSheet(ctx, "assets/sprites/0000.png")
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

func Load(ctx context.Context) (*Bundle, error) {
	b := &Bundle{}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		battletiles, err := loadBattleTiles(ctx)
		if err != nil {
			return err
		}
		b.Battletiles = battletiles
		return nil
	})

	g.Go(func() error {
		megaman, err := loadMegaman(ctx)
		if err != nil {
			return nil
		}
		b.Megaman = megaman
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return b, nil
}
