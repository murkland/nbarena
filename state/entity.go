package state

import (
	"flag"
	"image"
	"image/color"
	"strconv"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/murkland/clone"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
)

var (
	debugDrawEntityMarker = flag.Bool("debug_draw_entity_markers", false, "draw entity markers")
)

type EntityTraits struct {
	CanStepOnHoleLikeTiles bool
	IgnoresTileEffects     bool
	CannotFlinch           bool
	FatalHitLeaves1HP      bool
	IgnoresTileOwnership   bool
}

type EntityPerTickState struct {
	IsStepped               bool
	IsPendingDeletion       bool
	DoubleDamageWasConsumed bool
	Hit                     Hit
}

type Entity struct {
	id int

	elapsedTime Ticks

	behaviorElapsedTime Ticks
	behavior            EntityBehavior
	lastInterrupts      EntityBehaviorInterrupts

	TilePos       TilePos
	FutureTilePos TilePos

	IsAlliedWithAnswerer bool

	IsFlipped bool

	IsDeleted bool

	HP        int
	DisplayHP int

	Traits EntityTraits

	CanUseChip             bool
	ChipUseLockoutTimeLeft Ticks

	ChargingElapsedTime Ticks
	PowerShotChargeTime Ticks

	ParalyzedTimeLeft   Ticks
	ConfusedTimeLeft    Ticks
	BlindedTimeLeft     Ticks
	ImmobilizedTimeLeft Ticks
	FlashingTimeLeft    Ticks
	InvincibleTimeLeft  Ticks
	FrozenTimeLeft      Ticks
	BubbledTimeLeft     Ticks

	IsAngry        bool
	IsFullSynchro  bool
	IsBeingDragged bool
	IsSliding      bool
	IsCounterable  bool

	PerTickState EntityPerTickState
}

func (e *Entity) ID() int {
	return e.id
}

func (e *Entity) LastInterrupts() EntityBehaviorInterrupts {
	return e.lastInterrupts
}

func (e *Entity) Clone() *Entity {
	return &Entity{
		e.id,
		e.elapsedTime,
		e.behaviorElapsedTime, e.behavior.Clone(), e.lastInterrupts,
		e.TilePos, e.FutureTilePos,
		e.IsAlliedWithAnswerer,
		e.IsFlipped,
		e.IsDeleted,
		e.HP, e.DisplayHP,
		e.Traits,
		e.CanUseChip, e.ChipUseLockoutTimeLeft,
		e.ChargingElapsedTime, e.PowerShotChargeTime,
		e.ParalyzedTimeLeft, e.ConfusedTimeLeft, e.BlindedTimeLeft, e.ImmobilizedTimeLeft, e.FlashingTimeLeft, e.InvincibleTimeLeft, e.FrozenTimeLeft, e.BubbledTimeLeft,
		e.IsAngry, e.IsFullSynchro, e.IsBeingDragged, e.IsSliding, e.IsCounterable,
		e.PerTickState,
	}
}

func (e *Entity) Behavior() EntityBehavior {
	return e.behavior
}

func (e *Entity) SetBehavior(behavior EntityBehavior) {
	e.behaviorElapsedTime = 0
	e.behavior = behavior
}

func (e *Entity) BehaviorElapsedTime() Ticks {
	return e.behaviorElapsedTime
}

func (e *Entity) StartMove(tilePos TilePos, field *Field) bool {
	if tilePos < 0 {
		return false
	}

	x, y := tilePos.XY()
	if x < 0 || x >= TileCols || y < 0 || y >= TileRows {
		return false
	}

	tile := &field.Tiles[tilePos]
	if tilePos == e.TilePos ||
		(!e.Traits.IgnoresTileOwnership && e.IsAlliedWithAnswerer != tile.IsAlliedWithAnswerer) ||
		!tile.CanEnter(e) {
		return false
	}

	e.FutureTilePos = tilePos
	return true
}

func (e *Entity) FinishMove() {
	e.TilePos = e.FutureTilePos
}

var debugEntityMarkerImage *ebiten.Image
var debugEntityMarkerImageOnce sync.Once

func (e *Entity) Appearance(b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	x, y := e.TilePos.XY()
	rootNode.Opts.GeoM.Translate(float64((x-1)*TileRenderedWidth+TileRenderedWidth/2), float64((y-1)*TileRenderedHeight+TileRenderedHeight/2))

	characterNode := &draw.OptionsNode{}
	if e.IsFlipped {
		characterNode.Opts.GeoM.Scale(-1, 1)
	}
	if e.FrozenTimeLeft > 0 {
		// TODO: Render ice.
		characterNode.Opts.ColorM.Translate(float64(0xa5)/float64(0xff), float64(0xa5)/float64(0xff), float64(0xff)/float64(0xff), 0.0)
	}
	if e.ParalyzedTimeLeft > 0 && (e.elapsedTime/2)%2 == 1 {
		characterNode.Opts.ColorM.Translate(1.0, 1.0, 0.0, 0.0)
	}
	if e.FlashingTimeLeft > 0 && (e.elapsedTime/2)%2 == 0 {
		characterNode.Opts.ColorM.Translate(0.0, 0.0, 0.0, -1.0)
	}
	if e.PerTickState.Hit.TotalDamage > 0 {
		characterNode.Opts.ColorM.Translate(1.0, 1.0, 1.0, 0.0)
	}
	characterNode.Children = append(characterNode.Children, e.behavior.Appearance(e, b))

	if e.ChargingElapsedTime >= 10 {
		chargingNode := &draw.OptionsNode{}
		characterNode.Children = append(characterNode.Children, chargingNode)

		frames := b.ChargingSprites.ChargingAnimation.Frames
		if e.ChargingElapsedTime >= e.PowerShotChargeTime {
			frames = b.ChargingSprites.ChargedAnimation.Frames
		}
		frame := frames[int(e.ChargingElapsedTime)%len(frames)]
		chargingNode.Children = append(chargingNode.Children, draw.ImageWithFrame(b.ChargingSprites.Image, frame))
	}

	rootNode.Children = append(rootNode.Children, characterNode)

	if *debugDrawEntityMarker {
		debugEntityMarkerImageOnce.Do(func() {
			debugEntityMarkerImage = ebiten.NewImage(5, 5)
			for x := 0; x < 5; x++ {
				debugEntityMarkerImage.Set(x, 2, color.RGBA{0, 255, 0, 255})
			}
			for y := 0; y < 5; y++ {
				debugEntityMarkerImage.Set(2, y, color.RGBA{0, 255, 0, 255})
			}
		})
		rootNode.Children = append(rootNode.Children, draw.ImageWithOrigin(debugEntityMarkerImage, image.Point{2, 2}))
	}

	if e.HP > 0 && e.IsAlliedWithAnswerer {
		hpNode := &draw.OptionsNode{}
		rootNode.Children = append(rootNode.Children, hpNode)

		// Render HP.
		hpText := strconv.Itoa(int(e.DisplayHP))
		rect := text.BoundString(b.FontBold, hpText)
		hpNode.Opts.GeoM.Translate(float64(-rect.Max.X/2), float64(rect.Dy()/2))

		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				strokeNode := &draw.OptionsNode{}
				hpNode.Children = append(hpNode.Children, strokeNode)
				strokeNode.Opts.GeoM.Translate(float64(dx), float64(dy))
				strokeNode.Opts.ColorM.Scale(float64(0x31)/float64(0xFF), float64(0x39)/float64(0xFF), float64(0x52)/float64(0xFF), 1.0)
				strokeNode.Children = append(strokeNode.Children, &draw.TextNode{Text: hpText, Face: b.FontBold})
			}

			fillNode := &draw.OptionsNode{}
			hpNode.Children = append(hpNode.Children, fillNode)
			if e.DisplayHP > e.HP {
				fillNode.Opts.ColorM.Scale(float64(0xFF)/float64(0xFF), float64(0x84)/float64(0xFF), float64(0x5A)/float64(0xFF), 1.0)
			} else if e.DisplayHP < e.HP {
				fillNode.Opts.ColorM.Scale(float64(0x73)/float64(0xFF), float64(0xFF)/float64(0xFF), float64(0x4A)/float64(0xFF), 1.0)
			}
			fillNode.Children = append(fillNode.Children, &draw.TextNode{Text: hpText, Face: b.FontBold})
		}
	}

	return rootNode
}

func (e *Entity) Step(s *State) {
	e.lastInterrupts = e.behavior.Interrupts(e)

	if e.ChipUseLockoutTimeLeft > 0 {
		e.ChipUseLockoutTimeLeft--
	}

	e.elapsedTime++
	// Tick timers.
	// TODO: Verify this behavior is correct.
	e.behaviorElapsedTime++
	e.behavior.Step(e, s)
}

func (e *Entity) MakeDamageAndConsume(base int) Damage {
	dmg := Damage{
		Base: base,

		DoubleDamage: e.IsAngry || e.IsFullSynchro,
	}
	e.IsAngry = false
	e.IsFullSynchro = false
	if dmg.DoubleDamage {
		e.PerTickState.DoubleDamageWasConsumed = true
	}
	return dmg
}

type EntityBehaviorInterrupts struct {
	OnMove    bool
	OnChipUse bool
	OnCharge  bool
}

type EntityBehavior interface {
	clone.Cloner[EntityBehavior]
	Appearance(e *Entity, b *bundle.Bundle) draw.Node
	Step(e *Entity, s *State)
	Interrupts(e *Entity) EntityBehaviorInterrupts
}
