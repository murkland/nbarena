package bundle

import (
	"context"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yumland/moreio"
	"github.com/yumland/nbarena/loader"
	"github.com/yumland/pngsheet"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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

type CharacterSprites struct {
	Image *ebiten.Image

	IdleAnimation                *pngsheet.Animation
	FlinchEndAnimation           *pngsheet.Animation
	FlinchAnimation              *pngsheet.Animation
	TeleportEndAnimation         *pngsheet.Animation
	TeleportStartAnimation       *pngsheet.Animation
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

func makeSpriteLoader[T any](path string, f func(sheet *Sheet) T) func(ctx context.Context) (T, error) {
	return func(ctx context.Context) (T, error) {
		sheet, err := loadSheet(ctx, path)
		if err != nil {
			return *new(T), err
		}
		return f(sheet), nil
	}
}

func makeCharacterSpriteLoader(path string) func(ctx context.Context) (*CharacterSprites, error) {
	return makeSpriteLoader(path, func(sheet *Sheet) *CharacterSprites {
		return &CharacterSprites{
			Image: ebiten.NewImageFromImage(sheet.Image.(*image.Paletted)),

			IdleAnimation:                sheet.Info.Animations[0],
			FlinchEndAnimation:           sheet.Info.Animations[1],
			FlinchAnimation:              sheet.Info.Animations[2],
			TeleportEndAnimation:         sheet.Info.Animations[3],
			TeleportStartAnimation:       sheet.Info.Animations[4],
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
		}
	})
}

func makeFontFaceLoader(path string, size int) func(ctx context.Context) (font.Face, error) {
	return func(ctx context.Context) (font.Face, error) {
		f, err := moreio.Open(ctx, path)
		if err != nil {
			return nil, err
		}

		fnt, err := opentype.ParseReaderAt(f)
		if err != nil {
			return nil, err
		}

		return opentype.NewFace(fnt, &opentype.FaceOptions{
			Size:    16,
			DPI:     72,
			Hinting: font.HintingNone,
		})
	}
}

type ChargingSprites struct {
	Image *ebiten.Image

	ChargingAnimation *pngsheet.Animation
	ChargedAnimation  *pngsheet.Animation
}

type SlashSprites struct {
	SwordImage *ebiten.Image
	BladeImage *ebiten.Image

	WideAnimation     *pngsheet.Animation
	LongAnimation     *pngsheet.Animation
	ShortAnimation    *pngsheet.Animation
	VeryLongAnimation *pngsheet.Animation
}

type SwordSprites struct {
	Image *ebiten.Image

	BaseAnimation *pngsheet.Animation
}

type BusterSprites struct {
	Image *ebiten.Image

	BaseAnimation *pngsheet.Animation
}

type Sprites struct {
	Image      *ebiten.Image
	Animations []*pngsheet.Animation
}

type Bundle struct {
	Battletiles        *Battletiles
	MegamanSprites     *CharacterSprites
	ChargingSprites    *ChargingSprites
	BusterSprites      *BusterSprites
	MuzzleFlashSprites *Sprites
	SwordSprites       *SwordSprites
	SlashSprites       *SlashSprites
	FontBold           font.Face
}

func sheetToSprites(sheet *Sheet) *Sprites {
	return &Sprites{
		Image:      ebiten.NewImageFromImage(sheet.Image.(*image.Paletted)),
		Animations: sheet.Info.Animations,
	}
}

func Load(ctx context.Context) (*Bundle, error) {
	b := &Bundle{}

	l, ctx := loader.New(ctx)
	loader.Add(ctx, l, &b.Battletiles, loadBattleTiles)
	loader.Add(ctx, l, &b.MegamanSprites, makeCharacterSpriteLoader("assets/sprites/0000.png"))
	loader.Add(ctx, l, &b.ChargingSprites, makeSpriteLoader("assets/sprites/0274.png", func(sheet *Sheet) *ChargingSprites {
		return &ChargingSprites{
			Image: ebiten.NewImageFromImage(sheet.Image.(*image.Paletted)),

			ChargingAnimation: sheet.Info.Animations[1],
			ChargedAnimation:  sheet.Info.Animations[2],
		}
	}))
	loader.Add(ctx, l, &b.SwordSprites, makeSpriteLoader("assets/sprites/0069.png", func(sheet *Sheet) *SwordSprites {
		return &SwordSprites{
			Image: ebiten.NewImageFromImage(sheet.Image.(*image.Paletted)),

			BaseAnimation: sheet.Info.Animations[0],
		}
	}))
	loader.Add(ctx, l, &b.BusterSprites, makeSpriteLoader("assets/sprites/0072.png", func(sheet *Sheet) *BusterSprites {
		return &BusterSprites{
			Image: ebiten.NewImageFromImage(sheet.Image.(*image.Paletted)),

			BaseAnimation: sheet.Info.Animations[0],
		}
	}))
	loader.Add(ctx, l, &b.MuzzleFlashSprites, makeSpriteLoader("assets/sprites/0075.png", sheetToSprites))
	loader.Add(ctx, l, &b.SlashSprites, makeSpriteLoader("assets/sprites/0089.png", func(sheet *Sheet) *SlashSprites {
		img := sheet.Image.(*image.Paletted)
		palette := append(img.Palette, sheet.Info.SuggestedPalettes["extra"]...)
		img.Palette = palette[0 : 0+16]

		swordImage := ebiten.NewImageFromImage(img)

		img.Palette = palette[16 : 16+16]
		bladeImage := ebiten.NewImageFromImage(img)

		return &SlashSprites{
			SwordImage: swordImage,
			BladeImage: bladeImage,

			WideAnimation:     sheet.Info.Animations[0],
			LongAnimation:     sheet.Info.Animations[1],
			ShortAnimation:    sheet.Info.Animations[2],
			VeryLongAnimation: sheet.Info.Animations[3],
		}
	}))
	loader.Add(ctx, l, &b.FontBold, makeFontFaceLoader("assets/fonts/FontBold.ttf", 16))

	if err := l.Load(); err != nil {
		return nil, err
	}

	return b, nil
}
