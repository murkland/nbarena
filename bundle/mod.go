package bundle

import (
	"context"
	"image"
	_ "image/png"
	"io/ioutil"

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

	IdleAnimation          *pngsheet.Animation
	FlinchAnimation        *pngsheet.Animation // e.g. flinch, pauses on first frame on drag!
	StuckAnimation         *pngsheet.Animation // e.g. paralyzed, bubbled, frozen
	TeleportEndAnimation   *pngsheet.Animation
	TeleportStartAnimation *pngsheet.Animation
	SlashAnimation         *pngsheet.Animation // e.g. Sword
	ThrowAnimation         *pngsheet.Animation // e.g. MiniBomb
	BraceAnimation         *pngsheet.Animation // e.g. end of Cannon
	CannonAnimation        *pngsheet.Animation // e.g. Cannon
	RecoilShotAnimation    *pngsheet.Animation // e.g. AirShot
	HoldInFrontAnimation   *pngsheet.Animation // e.g. DolThndr, RskyHony, TankCan, Tornado
	BusterAnimation        *pngsheet.Animation
	FlourishAnimation      *pngsheet.Animation // e.g. BublStar
	GattlingAnimation      *pngsheet.Animation // e.g. Vulcan, MachGun
	TwoHandedAnimation     *pngsheet.Animation // e.g. AquaWhirl
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

			IdleAnimation:          sheet.Info.Animations[0],
			FlinchAnimation:        sheet.Info.Animations[1],
			StuckAnimation:         sheet.Info.Animations[2],
			TeleportEndAnimation:   sheet.Info.Animations[3],
			TeleportStartAnimation: sheet.Info.Animations[4],
			SlashAnimation:         sheet.Info.Animations[5],
			ThrowAnimation:         sheet.Info.Animations[6],
			BraceAnimation:         sheet.Info.Animations[7],
			CannonAnimation:        sheet.Info.Animations[8],
			RecoilShotAnimation:    sheet.Info.Animations[9],
			HoldInFrontAnimation:   sheet.Info.Animations[10],
			FlourishAnimation:      sheet.Info.Animations[12],
			GattlingAnimation:      sheet.Info.Animations[13],
			BusterAnimation:        sheet.Info.Animations[14],
			TwoHandedAnimation:     sheet.Info.Animations[18],
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

type SlashDecorationSprites struct {
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

type Sprite struct {
	Image     *ebiten.Image
	Animation *pngsheet.Animation
}

type DecorationType int

const (
	DecorationTypeNone            DecorationType = 0
	DecorationTypeCannonExplosion DecorationType = iota
	DecorationTypeBusterPowerShotExplosion
	DecorationTypeBusterExplosion
	DecorationTypePiercingExplosion
	DecorationTypeUninstallExplosion
	DecorationTypeChipDeleteExplosion
	DecorationTypeShieldHitExplosion
	DecorationTypeNullShortSwordSlash
	DecorationTypeNullWideSwordSlash
	DecorationTypeNullLongSwordSlash
	DecorationTypeNullVeryLongSwordSlash
	DecorationTypeNullShortBladeSlash
	DecorationTypeNullWideBladeSlash
	DecorationTypeNullLongBladeSlash
	DecorationTypeNullVeryLongBladeSlash
	DecorationTypeWindSlash
)

type Bundle struct {
	Battletiles *Battletiles

	MegamanSprites     *CharacterSprites
	SwordSprites       *SwordSprites
	CannonSprites      *CannonSprites
	AirShooterSprites  *Sprites
	BusterSprites      *BusterSprites
	MuzzleFlashSprites *Sprites
	AreaGrabSprites    *Sprites
	VulcanSprites      *Sprites
	WindRackSprites    *Sprites
	FullSynchroSprites *Sprites
	IcedSprites        *Sprites

	DecorationSprites map[DecorationType]*Sprite

	ChargingSprites *ChargingSprites

	ChipIconSprites *Sprites

	TallFont    font.Face
	Tall2Font   font.Face
	TinyNumFont font.Face
}

func loadBDF(ctx context.Context, f moreio.File) (font.Face, error) {
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	font, err := bdf.Parse(buf)
	if err != nil {
		return nil, err
	}

	return font.NewFace(), nil
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
	loader.Add(ctx, l, "assets/sprites/0088.png", &b.AreaGrabSprites, makeSpriteLoader(sheetToSprites))
	loader.Add(ctx, l, "assets/sprites/0093.png", &b.AirShooterSprites, makeSpriteLoader(sheetToSprites))
	loader.Add(ctx, l, "assets/sprites/0098.png", &b.VulcanSprites, makeSpriteLoader(sheetToSprites))
	loader.Add(ctx, l, "assets/sprites/0108.png", &b.WindRackSprites, makeSpriteLoader(sheetToSprites))
	loader.Add(ctx, l, "assets/sprites/0288.png", &b.FullSynchroSprites, makeSpriteLoader(sheetToSprites))
	loader.Add(ctx, l, "assets/sprites/0294.png", &b.IcedSprites, makeSpriteLoader(sheetToSprites))

	var slashDecorationSprites *SlashDecorationSprites
	loader.Add(ctx, l, "assets/sprites/0089.png", &slashDecorationSprites, makeSpriteLoader(func(sheet *Sheet) *SlashDecorationSprites {
		img := sheet.Image.(*image.Paletted)
		palette := append(img.Palette, sheet.Info.SuggestedPalettes["extra"]...)
		img.Palette = palette[0 : 0+16]

		swordImage := ebiten.NewImageFromImage(img)

		img.Palette = palette[16 : 16+16]
		bladeImage := ebiten.NewImageFromImage(img)

		return &SlashDecorationSprites{
			SwordImage: swordImage,
			BladeImage: bladeImage,

			WideAnimation:     sheet.Info.Animations[0],
			LongAnimation:     sheet.Info.Animations[1],
			ShortAnimation:    sheet.Info.Animations[2],
			VeryLongAnimation: sheet.Info.Animations[3],
		}
	}))

	var windSlashDecorationSprites *Sprites
	loader.Add(ctx, l, "assets/sprites/0109.png", &windSlashDecorationSprites, makeSpriteLoader(sheetToSprites))

	var cannonExplosionDecorationSprites *Sprites
	loader.Add(ctx, l, "assets/sprites/0267.png", &cannonExplosionDecorationSprites, makeSpriteLoader(sheetToSprites))

	var chargeShotExplosionDecorationSprites *Sprites
	loader.Add(ctx, l, "assets/sprites/0270.png", &chargeShotExplosionDecorationSprites, makeSpriteLoader(sheetToSprites))

	var explosionDecorationSprites *Sprites
	loader.Add(ctx, l, "assets/sprites/0271.png", &explosionDecorationSprites, makeSpriteLoader(sheetToSprites))

	var shieldHitExplosionDecorationSprites *Sprites
	loader.Add(ctx, l, "assets/sprites/0272.png", &shieldHitExplosionDecorationSprites, makeSpriteLoader(sheetToSprites))

	var chipDeleteExplosionDecorationSprites *Sprites
	loader.Add(ctx, l, "assets/sprites/0278.png", &chipDeleteExplosionDecorationSprites, makeSpriteLoader(sheetToSprites))

	var piercingExplosionDecorationSprites *Sprites
	loader.Add(ctx, l, "assets/sprites/0281.png", &piercingExplosionDecorationSprites, makeSpriteLoader(sheetToSprites))

	var uninstallExplosionDecorationSprites *Sprites
	loader.Add(ctx, l, "assets/sprites/0290.png", &uninstallExplosionDecorationSprites, makeSpriteLoader(sheetToSprites))

	loader.Add(ctx, l, "assets/sprites/0274.png", &b.ChargingSprites, makeSpriteLoader(func(sheet *Sheet) *ChargingSprites {
		return &ChargingSprites{
			Image: ebiten.NewImageFromImage(sheet.Image.(*image.Paletted)),

			ChargingAnimation: sheet.Info.Animations[1],
			ChargedAnimation:  sheet.Info.Animations[2],
		}
	}))

	loader.Add(ctx, l, "assets/chipicons.png", &b.ChipIconSprites, makeSpriteLoader(sheetToSprites))

	loader.Add(ctx, l, "assets/fonts/tall.bdf", &b.TallFont, loadBDF)
	loader.Add(ctx, l, "assets/fonts/tall2.bdf", &b.Tall2Font, loadBDF)
	loader.Add(ctx, l, "assets/fonts/tinynum.bdf", &b.TinyNumFont, loadBDF)

	if err := l.Load(); err != nil {
		return nil, err
	}

	b.DecorationSprites = map[DecorationType]*Sprite{
		DecorationTypeCannonExplosion:          {cannonExplosionDecorationSprites.Image, cannonExplosionDecorationSprites.Animations[0]},
		DecorationTypeBusterPowerShotExplosion: {chargeShotExplosionDecorationSprites.Image, chargeShotExplosionDecorationSprites.Animations[0]},
		DecorationTypeBusterExplosion:          {explosionDecorationSprites.Image, explosionDecorationSprites.Animations[0]},
		DecorationTypePiercingExplosion:        {piercingExplosionDecorationSprites.Image, piercingExplosionDecorationSprites.Animations[0]},
		DecorationTypeUninstallExplosion:       {uninstallExplosionDecorationSprites.Image, uninstallExplosionDecorationSprites.Animations[0]},
		DecorationTypeChipDeleteExplosion:      {chipDeleteExplosionDecorationSprites.Image, chipDeleteExplosionDecorationSprites.Animations[0]},
		DecorationTypeShieldHitExplosion:       {shieldHitExplosionDecorationSprites.Image, shieldHitExplosionDecorationSprites.Animations[0]},
		DecorationTypeNullShortSwordSlash:      {slashDecorationSprites.SwordImage, slashDecorationSprites.ShortAnimation},
		DecorationTypeNullWideSwordSlash:       {slashDecorationSprites.SwordImage, slashDecorationSprites.WideAnimation},
		DecorationTypeNullLongSwordSlash:       {slashDecorationSprites.SwordImage, slashDecorationSprites.LongAnimation},
		DecorationTypeNullVeryLongSwordSlash:   {slashDecorationSprites.SwordImage, slashDecorationSprites.VeryLongAnimation},
		DecorationTypeNullShortBladeSlash:      {slashDecorationSprites.BladeImage, slashDecorationSprites.ShortAnimation},
		DecorationTypeNullWideBladeSlash:       {slashDecorationSprites.BladeImage, slashDecorationSprites.WideAnimation},
		DecorationTypeNullLongBladeSlash:       {slashDecorationSprites.BladeImage, slashDecorationSprites.LongAnimation},
		DecorationTypeNullVeryLongBladeSlash:   {slashDecorationSprites.BladeImage, slashDecorationSprites.VeryLongAnimation},
		DecorationTypeWindSlash:                {windSlashDecorationSprites.Image, windSlashDecorationSprites.Animations[0]},
	}

	return b, nil
}
