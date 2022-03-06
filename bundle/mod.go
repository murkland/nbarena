package bundle

import (
	"context"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/murkland/moreio"
	"github.com/murkland/nbarena/loader"
	"github.com/murkland/pngsheet"
	"github.com/zachomedia/go-bdf"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Sheet struct {
	Info  *pngsheet.Info
	Image image.Image
}

func loadSheet(ctx context.Context, f moreio.File) (*Sheet, error) {
	defer f.Close()

	img, info, err := pngsheet.Load(f)
	if err != nil {
		return nil, err
	}

	return &Sheet{info, img}, nil
}

type Battletiles struct {
	Info          *pngsheet.Info
	OffererTiles  *ebiten.Image
	AnswererTiles *ebiten.Image
}

func loadBattletiles(ctx context.Context, f moreio.File) (*Battletiles, error) {
	sheet, err := loadSheet(ctx, f)
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
	FlinchAnimation              *pngsheet.Animation
	StuckAnimation               *pngsheet.Animation
	TeleportEndAnimation         *pngsheet.Animation
	TeleportStartAnimation       *pngsheet.Animation
	SlashAnimation               *pngsheet.Animation
	ThrowAnimation               *pngsheet.Animation
	BraceAnimation               *pngsheet.Animation
	CannonAnimation              *pngsheet.Animation
	FireAndSlideAnimation        *pngsheet.Animation
	BusterEndAnimation           *pngsheet.Animation
	BusterAnimation              *pngsheet.Animation
	FlourishAnimation            *pngsheet.Animation
	GattlingAnimation            *pngsheet.Animation
	TwoHandedSlashStartAnimation *pngsheet.Animation
	TwoHandedSlashAnimation      *pngsheet.Animation
}

func makeSpriteLoader[T any](f func(sheet *Sheet) T) func(ctx context.Context, file moreio.File) (T, error) {
	return func(ctx context.Context, file moreio.File) (T, error) {
		sheet, err := loadSheet(ctx, file)
		if err != nil {
			return *new(T), err
		}
		return f(sheet), nil
	}
}

func loadCharacterSprite(ctx context.Context, f moreio.File) (*CharacterSprites, error) {
	return makeSpriteLoader(func(sheet *Sheet) *CharacterSprites {
		return &CharacterSprites{
			Image: ebiten.NewImageFromImage(sheet.Image.(*image.Paletted)),

			IdleAnimation:                sheet.Info.Animations[0],
			FlinchAnimation:              sheet.Info.Animations[1],
			StuckAnimation:               sheet.Info.Animations[2],
			TeleportEndAnimation:         sheet.Info.Animations[3],
			TeleportStartAnimation:       sheet.Info.Animations[4],
			SlashAnimation:               sheet.Info.Animations[5],
			ThrowAnimation:               sheet.Info.Animations[6],
			BraceAnimation:               sheet.Info.Animations[7],
			CannonAnimation:              sheet.Info.Animations[8],
			FireAndSlideAnimation:        sheet.Info.Animations[9],
			BusterEndAnimation:           sheet.Info.Animations[10],
			BusterAnimation:              sheet.Info.Animations[11],
			FlourishAnimation:            sheet.Info.Animations[12],
			GattlingAnimation:            sheet.Info.Animations[13],
			TwoHandedSlashStartAnimation: sheet.Info.Animations[18],
			TwoHandedSlashAnimation:      sheet.Info.Animations[19],
		}
	})(ctx, f)
}

func makeFontFaceLoader(size int) func(ctx context.Context, f moreio.File) (font.Face, error) {
	return func(ctx context.Context, f moreio.File) (font.Face, error) {
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

type CannonSprites struct {
	CannonImage   *ebiten.Image
	HiCannonImage *ebiten.Image
	MCannonImage  *ebiten.Image

	Animation *pngsheet.Animation
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
	CannonSprites      *CannonSprites
	ChipIconSprites    *Sprites
	FontBold           font.Face
	EnemyHPFont        font.Face
}

func sheetToSprites(sheet *Sheet) *Sprites {
	return &Sprites{
		Image:      ebiten.NewImageFromImage(sheet.Image.(*image.Paletted)),
		Animations: sheet.Info.Animations,
	}
}

func Load(ctx context.Context, loaderCallback loader.Callback) (*Bundle, error) {
	b := &Bundle{}

	l, ctx := loader.New(ctx, loaderCallback)
	loader.Add(ctx, l, "assets/battletiles.png", &b.Battletiles, loadBattletiles)
	loader.Add(ctx, l, "assets/sprites/0000.png", &b.MegamanSprites, loadCharacterSprite)
	loader.Add(ctx, l, "assets/sprites/0274.png", &b.ChargingSprites, makeSpriteLoader(func(sheet *Sheet) *ChargingSprites {
		return &ChargingSprites{
			Image: ebiten.NewImageFromImage(sheet.Image.(*image.Paletted)),

			ChargingAnimation: sheet.Info.Animations[1],
			ChargedAnimation:  sheet.Info.Animations[2],
		}
	}))
	loader.Add(ctx, l, "assets/sprites/0069.png", &b.SwordSprites, makeSpriteLoader(func(sheet *Sheet) *SwordSprites {
		return &SwordSprites{
			Image: ebiten.NewImageFromImage(sheet.Image.(*image.Paletted)),

			BaseAnimation: sheet.Info.Animations[0],
		}
	}))
	loader.Add(ctx, l, "assets/sprites/0070.png", &b.CannonSprites, makeSpriteLoader(func(sheet *Sheet) *CannonSprites {
		img := sheet.Image.(*image.Paletted)
		palette := append(img.Palette, sheet.Info.SuggestedPalettes["extra"]...)
		img.Palette = palette[0 : 0+16]

		cannonImage := ebiten.NewImageFromImage(img)

		img.Palette = palette[16 : 16+16]
		hiCannonImage := ebiten.NewImageFromImage(img)

		img.Palette = palette[32 : 32+16]
		mCannonImage := ebiten.NewImageFromImage(img)

		return &CannonSprites{
			CannonImage:   cannonImage,
			HiCannonImage: hiCannonImage,
			MCannonImage:  mCannonImage,

			Animation: sheet.Info.Animations[0],
		}
	}))
	loader.Add(ctx, l, "assets/sprites/0072.png", &b.BusterSprites, makeSpriteLoader(func(sheet *Sheet) *BusterSprites {
		return &BusterSprites{
			Image: ebiten.NewImageFromImage(sheet.Image.(*image.Paletted)),

			BaseAnimation: sheet.Info.Animations[0],
		}
	}))
	loader.Add(ctx, l, "assets/sprites/0075.png", &b.MuzzleFlashSprites, makeSpriteLoader(sheetToSprites))
	loader.Add(ctx, l, "assets/sprites/0089.png", &b.SlashSprites, makeSpriteLoader(func(sheet *Sheet) *SlashSprites {
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
	loader.Add(ctx, l, "assets/chipicons.png", &b.ChipIconSprites, makeSpriteLoader(sheetToSprites))
	loader.Add(ctx, l, "assets/fonts/FontBold.ttf", &b.FontBold, makeFontFaceLoader(16))
	loader.Add(ctx, l, "assets/fonts/enemyhp.bdf", &b.EnemyHPFont, func(ctx context.Context, f moreio.File) (font.Face, error) {
		defer f.Close()

		buf, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}

		font, err := bdf.Parse(buf)
		if err != nil {
			return nil, err
		}

		log.Printf("%#v", font)
		return font.NewFace(), nil
	})

	if err := l.Load(); err != nil {
		return nil, err
	}

	return b, nil
}
