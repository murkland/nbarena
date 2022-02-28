package bundle

import (
	"context"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yumland/moreio"
	"github.com/yumland/pngsheet"
	"github.com/yumland/yumbattle/loader"
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
	Sprites *ebiten.Image

	IdleAnimation                *pngsheet.Animation
	FlinchEndAnimation           *pngsheet.Animation
	FlinchingAnimation           *pngsheet.Animation
	MoveEndAnimation             *pngsheet.Animation
	MoveStartAnimation           *pngsheet.Animation
	SlashAnimation               *pngsheet.Animation
	ThrowAnimation               *pngsheet.Animation
	BraceEndAnimation            *pngsheet.Animation
	CannonAnimation              *pngsheet.Animation
	FireAndSlideAnimation        *pngsheet.Animation
	BusterEndAnimation           *pngsheet.Animation
	BusterAnimation              *pngsheet.Animation
	FlourishAnimation            *pngsheet.Animation
	GattlingAnimation            *pngsheet.Animation
	TwoHandedSlashStartAnimation *pngsheet.Animation
	TwoHandedSlashAnimation      *pngsheet.Animation
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

	return &Megaman{
		Sprites: ebiten.NewImageFromImage(img),

		IdleAnimation:                sheet.Info.Animations[0],
		FlinchEndAnimation:           sheet.Info.Animations[1],
		FlinchingAnimation:           sheet.Info.Animations[2],
		MoveEndAnimation:             sheet.Info.Animations[3],
		MoveStartAnimation:           sheet.Info.Animations[4],
		SlashAnimation:               sheet.Info.Animations[5],
		ThrowAnimation:               sheet.Info.Animations[6],
		BraceEndAnimation:            sheet.Info.Animations[7],
		CannonAnimation:              sheet.Info.Animations[8],
		FireAndSlideAnimation:        sheet.Info.Animations[9],
		BusterEndAnimation:           sheet.Info.Animations[10],
		BusterAnimation:              sheet.Info.Animations[11],
		FlourishAnimation:            sheet.Info.Animations[12],
		GattlingAnimation:            sheet.Info.Animations[13],
		TwoHandedSlashStartAnimation: sheet.Info.Animations[18],
		TwoHandedSlashAnimation:      sheet.Info.Animations[19],
	}, nil
}

type Bundle struct {
	Battletiles *Battletiles
	Megaman     *Megaman
}

func Load(ctx context.Context) (*Bundle, error) {
	b := &Bundle{}

	l, ctx := loader.New(ctx)
	loader.Add(ctx, l, &b.Battletiles, loadBattleTiles)
	loader.Add(ctx, l, &b.Megaman, loadMegaman)

	if err := l.Load(); err != nil {
		return nil, err
	}

	return b, nil
}
