package state

import (
	"flag"
	"image"
	"image/color"
	"strconv"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/murkland/clone"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/draw"
	"golang.org/x/exp/slices"
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
	CannotSlide            bool
	Intangible             bool
}

type EntityPerTickState struct {
	WasHit bool
}

type ForcedMovementType int

const (
	ForcedMovementTypeNone      ForcedMovementType = 0
	ForcedMovementTypeSlide     ForcedMovementType = 1
	ForcedMovementTypeSmallDrag ForcedMovementType = 2
	ForcedMovementTypeBigDrag   ForcedMovementType = 2
)

func (t ForcedMovementType) IsDrag() bool {
	return t == ForcedMovementTypeSmallDrag || t == ForcedMovementTypeBigDrag
}

type ForcedMovement struct {
	Type      ForcedMovementType
	Direction Direction
}

type EntityBehaviorTraits struct {
	CanBeCountered bool
}

type EntityBehaviorState struct {
	Behavior    EntityBehavior
	ElapsedTime Ticks
}

func (s EntityBehaviorState) Clone() EntityBehaviorState {
	return EntityBehaviorState{s.Behavior.Clone(), s.ElapsedTime}
}

type EntityID uint64

type ChipPlaque struct {
	ElapsedTime  Ticks
	Chip         *Chip
	DoubleDamage bool
	AttackPlus   int
}

type Emotion int

const (
	EmotionNormal      Emotion = 0
	EmotionFullSynchro Emotion = 1
	EmotionAngry       Emotion = 2
)

type HitResolution struct {
	Damage int

	Flinch         bool
	FlashTime      Ticks
	ParalyzeTime   Ticks
	ConfuseTime    Ticks
	BlindTime      Ticks
	ImmobilizeTime Ticks
	FreezeTime     Ticks
	BubbleTime     Ticks

	ForcedMovement ForcedMovement
}

type Flashing struct {
	TimeLeft Ticks
	IsInvis  bool
}

type ForcedMovementState struct {
	ForcedMovement ForcedMovement
	ElapsedTime    Ticks
}

type Entity struct {
	id EntityID

	ElapsedTime Ticks

	RunsInTimestop bool

	BehaviorState        EntityBehaviorState
	NextBehavior         EntityBehavior
	IsPendingDestruction bool

	Intent     Intent
	LastIntent Intent

	TilePos       TilePos
	FutureTilePos TilePos

	ForcedMovementState ForcedMovementState

	IsAlliedWithAnswerer bool

	IsFlipped bool

	IsDead bool

	Element Element

	HP        int
	MaxHP     int
	DisplayHP int

	Traits EntityTraits

	PowerShotChargeTime Ticks

	ConfusedTimeLeft    Ticks
	BlindedTimeLeft     Ticks
	ImmobilizedTimeLeft Ticks
	Flashing            Flashing
	InvincibleTimeLeft  Ticks

	Emotion Emotion

	HitResolution HitResolution
	PerTickState  EntityPerTickState

	Chips         []*Chip
	ChipUseQueued bool

	DragLockoutTimeLeft    Ticks
	ChipUseLockoutTimeLeft Ticks

	ChipPlaque ChipPlaque
}

func (e *Entity) ID() EntityID {
	return e.id
}

func (e *Entity) Flip() {
	e.IsAlliedWithAnswerer = !e.IsAlliedWithAnswerer
	e.IsFlipped = !e.IsFlipped
	e.TilePos = e.TilePos.Flipped()
	e.FutureTilePos = e.FutureTilePos.Flipped()
	e.ForcedMovementState.ForcedMovement.Direction = e.ForcedMovementState.ForcedMovement.Direction.FlipH()
}

func (e *Entity) Clone() *Entity {
	return &Entity{
		e.id,
		e.ElapsedTime,
		e.RunsInTimestop,
		e.BehaviorState.Clone(), clone.Interface[EntityBehavior](e.NextBehavior), e.IsPendingDestruction,
		e.Intent, e.LastIntent,
		e.TilePos, e.FutureTilePos,
		e.ForcedMovementState,
		e.IsAlliedWithAnswerer,
		e.IsFlipped,
		e.IsDead,
		e.Element,
		e.HP, e.MaxHP, e.DisplayHP,
		e.Traits,
		e.PowerShotChargeTime,
		e.ConfusedTimeLeft, e.BlindedTimeLeft, e.ImmobilizedTimeLeft, e.Flashing, e.InvincibleTimeLeft,
		e.Emotion,
		e.HitResolution, e.PerTickState,
		slices.Clone(e.Chips), e.ChipUseQueued,
		e.DragLockoutTimeLeft, e.ChipUseLockoutTimeLeft,
		e.ChipPlaque,
	}
}

func (e *Entity) Facing() Direction {
	if e.IsFlipped {
		return DirectionLeft
	}
	return DirectionRight
}

func (e *Entity) DoubleDamage() bool {
	return e.Emotion == EmotionAngry || e.Emotion == EmotionFullSynchro
}

func (e *Entity) UseChip(s *State) bool {
	if len(e.Chips) == 0 {
		return false
	}
	chip := e.Chips[len(e.Chips)-1]
	e.Chips = e.Chips[:len(e.Chips)-1]

	dmg := Damage{
		Base: chip.BaseDamage,

		DoubleDamage: e.DoubleDamage(),
	}
	e.Emotion = EmotionNormal
	if dmg.DoubleDamage {
		s.AttachSound(&Sound{
			Type: bundle.SoundTypeDoubleDamageConsumed,
		})
	}

	chip.OnUse(s, e, dmg)
	e.ChipPlaque = ChipPlaque{Chip: chip, DoubleDamage: dmg.DoubleDamage}
	return true
}

func (e *Entity) MoveDirectly(tilePos TilePos) bool {
	if tilePos < 0 {
		return false
	}

	x, y := tilePos.XY()
	if x < 0 || x >= TileCols || y < 0 || y >= TileRows {
		return false
	}

	e.TilePos = tilePos
	return true
}

func (e *Entity) StartMove(tilePos TilePos, s *State) bool {
	if tilePos < 0 {
		return false
	}

	x, y := tilePos.XY()
	if x < 0 || x >= TileCols || y < 0 || y >= TileRows {
		return false
	}

	tile := s.Field.Tiles[tilePos]
	if tilePos == e.TilePos ||
		(!e.Traits.IgnoresTileOwnership && e.IsAlliedWithAnswerer != tile.IsAlliedWithAnswerer) ||
		(tile.Reserver != 0 && tile.Reserver != e.id) ||
		!tile.CanEnter(e) {
		return false
	}

	// TODO: Figure out when to trigger onleave/onenter callbacks
	tile.Reserver = e.ID()
	e.FutureTilePos = tilePos
	return true
}

func (e *Entity) FinishMove(s *State) {
	// TODO: Trigger on leave?
	s.Field.Tiles[e.TilePos].Reserver = 0
	e.TilePos = e.FutureTilePos
	s.Field.Tiles[e.TilePos].Reserver = e.ID()
}

// SetBehaviorImmediate sets the entity's behavior immediately to the next state and steps once. You probably don't want to call this: you should probably use NextBehavior instead.
func (e *Entity) SetBehaviorImmediate(behavior EntityBehavior, s *State) {
	if e.ForcedMovementState.ForcedMovement.Type != ForcedMovementTypeSlide || e.ForcedMovementState.ElapsedTime > 0 {
		e.FinishMove(s)
	}
	e.BehaviorState.Behavior.Cleanup(e, s)
	e.BehaviorState = EntityBehaviorState{Behavior: behavior}
	e.NextBehavior = nil
	e.BehaviorState.Behavior.Step(e, s)
}

var debugEntityMarkerImage *ebiten.Image
var debugEntityMarkerImageOnce sync.Once

func BehaviorIs[T EntityBehavior](behavior EntityBehavior) bool {
	_, ok := behavior.(T)
	return ok
}

func (e *Entity) Appearance(b *bundle.Bundle) draw.Node {
	rootNode := &draw.OptionsNode{}
	x, y := e.TilePos.XY()

	dx, dy := e.ForcedMovementState.ForcedMovement.Direction.XY()
	offset := (int(e.ForcedMovementState.ElapsedTime)+2+4)%4 - 2
	dx *= offset
	dy *= offset

	rootNode.Opts.GeoM.Translate(
		float64((x-1)*TileRenderedWidth+TileRenderedWidth/2+dx*TileRenderedWidth/4),
		float64((y-1)*TileRenderedHeight+TileRenderedHeight/2+dy*TileRenderedHeight/4),
	)

	rootCharacterNode := &draw.OptionsNode{}
	rootNode.Children = append(rootNode.Children, rootCharacterNode)

	if e.IsFlipped {
		rootCharacterNode.Opts.GeoM.Scale(-1, 1)
	}

	characterNode := &draw.OptionsNode{}
	rootCharacterNode.Children = append(rootCharacterNode.Children, characterNode)

	characterNode.Children = append(characterNode.Children, e.BehaviorState.Behavior.Appearance(e, b))

	if e.Flashing.TimeLeft > 0 && (e.ElapsedTime/2)%2 == 0 {
		characterNode.Opts.ColorM.Translate(0.0, 0.0, 0.0, -1.0)
	}
	if e.PerTickState.WasHit {
		characterNode.Opts.ColorM.Translate(1.0, 1.0, 1.0, 0.0)
	}
	if e.Emotion == EmotionFullSynchro {
		characterNode.Opts.ColorM.Translate(float64(0x29)/float64(0xff), float64(0x29)/float64(0xff), float64(0x29)/float64(0xff), 0.0)

		fullSynchroNode := &draw.OptionsNode{Layer: 8}
		fullSynchroNode.Children = append(fullSynchroNode.Children, draw.ImageWithAnimation(b.FullSynchroSprites.Image, b.FullSynchroSprites.Animations[0], int(e.ElapsedTime)))
		rootCharacterNode.Children = append(rootCharacterNode.Children, fullSynchroNode)
	} else if e.Emotion == EmotionAngry {
		characterNode.Opts.ColorM.Translate(float64(0x80)/float64(0xff), float64(0)/float64(0xff), float64(0)/float64(0xff), 0.0)
	}

	if *debugDrawEntityMarker {
		debugEntityMarkerImageOnce.Do(func() {
			debugEntityMarkerImage = ebiten.NewImage(5, 5)
			for x := 0; x < 5; x++ {
				debugEntityMarkerImage.Set(x, 2, color.RGBA{255, 255, 255, 255})
			}
			for y := 0; y < 5; y++ {
				debugEntityMarkerImage.Set(2, y, color.RGBA{255, 255, 255, 255})
			}
		})
		debugEntityMarkerNode := &draw.OptionsNode{}
		debugEntityMarkerNode.Children = append(debugEntityMarkerNode.Children, draw.ImageWithOrigin(debugEntityMarkerImage, image.Point{2, 2}))
		debugEntityMarkerNode.Opts.ColorM.Scale(1.0, 0.0, 1.0, 0.5)
		rootNode.Children = append(rootNode.Children, debugEntityMarkerNode)
	}

	if e.IsAlliedWithAnswerer {
		if e.DisplayHP != 0 {
			hpNode := &draw.OptionsNode{}
			rootNode.Children = append(rootNode.Children, hpNode)

			// Render HP.
			hpText := strconv.Itoa(int(e.DisplayHP))
			hpNode.Opts.GeoM.Translate(float64(0), float64(4))

			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					strokeNode := &draw.OptionsNode{}
					hpNode.Children = append(hpNode.Children, strokeNode)
					strokeNode.Opts.GeoM.Translate(float64(dx), float64(dy))
					strokeNode.Opts.ColorM.Scale(float64(0x31)/float64(0xFF), float64(0x39)/float64(0xFF), float64(0x52)/float64(0xFF), 1.0)
					strokeNode.Children = append(strokeNode.Children, &draw.TextNode{Text: hpText, Face: b.TinyNumFont, Anchor: draw.TextAnchorCenter | draw.TextAnchorBottom})
				}

				fillNode := &draw.OptionsNode{}
				hpNode.Children = append(hpNode.Children, fillNode)
				if e.DisplayHP > e.HP {
					fillNode.Opts.ColorM.Scale(float64(0xFF)/float64(0xFF), float64(0x84)/float64(0xFF), float64(0x5A)/float64(0xFF), 1.0)
				} else if e.DisplayHP < e.HP {
					fillNode.Opts.ColorM.Scale(float64(0x73)/float64(0xFF), float64(0xFF)/float64(0xFF), float64(0x4A)/float64(0xFF), 1.0)
				}
				fillNode.Children = append(fillNode.Children, &draw.TextNode{Text: hpText, Face: b.TinyNumFont, Anchor: draw.TextAnchorCenter | draw.TextAnchorBottom})
			}
		}
	} else {
		chipsNode := &draw.OptionsNode{}
		chipsNode.Opts.GeoM.Translate(0, float64(-56))
		rootNode.Children = append(rootNode.Children, chipsNode)

		for i, chip := range e.Chips {
			chipNode := &draw.OptionsNode{Layer: 8}
			j := len(e.Chips) - i - 1
			chipNode.Opts.GeoM.Translate(float64(-j*2), float64(-j*2))
			chipsNode.Children = append(chipsNode.Children, chipNode)

			chipNode.Children = append(chipNode.Children, draw.ImageWithFrame(b.ChipIconSprites.Image, b.ChipIconSprites.Animations[chip.Index].Frames[0]))
		}
	}

	return rootNode
}

func (e *Entity) ApplyHit(h Hit) {
	if h.Element.IsSuperEffectiveAgainst(e.Element) {
		h.TotalDamage *= 2
	}

	e.HitResolution.Damage += h.TotalDamage

	// TODO: Verify this is correct behavior.
	if h.ParalyzeTime > e.HitResolution.ParalyzeTime {
		e.HitResolution.ParalyzeTime = h.ParalyzeTime
	}
	if h.ConfuseTime > e.HitResolution.ConfuseTime {
		e.HitResolution.ConfuseTime = h.ConfuseTime
	}
	if h.BlindTime > e.HitResolution.BlindTime {
		e.HitResolution.BlindTime = h.BlindTime
	}
	if h.ImmobilizeTime > e.HitResolution.ImmobilizeTime {
		e.HitResolution.ImmobilizeTime = h.ImmobilizeTime
	}
	if h.FreezeTime > e.HitResolution.FreezeTime {
		e.HitResolution.FreezeTime = h.FreezeTime
	}
	if h.BubbleTime > e.HitResolution.BubbleTime {
		e.HitResolution.BubbleTime = h.BubbleTime
	}
	if h.FlashTime > e.HitResolution.FlashTime {
		e.HitResolution.FlashTime = h.FlashTime
	}
	if h.Flinch {
		e.HitResolution.Flinch = true
	}
	if h.ForcedMovement.Type != ForcedMovementTypeNone && (e.HitResolution.ForcedMovement.Type == ForcedMovementTypeNone || h.ForcedMovement.Type.IsDrag()) {
		e.HitResolution.ForcedMovement = h.ForcedMovement
	}
}

func (e *Entity) RemoveFlashing(s *State) {
	if e.Flashing.IsInvis {
		// TODO: Play un-invis sound effect.
	}
	e.Flashing = Flashing{}
}

func (e *Entity) Step(s *State) {
	if e.ChipUseLockoutTimeLeft > 0 {
		e.ChipUseLockoutTimeLeft--
	}

	if e.ChipPlaque.Chip != nil {
		e.ChipPlaque.ElapsedTime++
		if e.ChipPlaque.ElapsedTime >= 60 {
			e.ChipPlaque = ChipPlaque{}
		}
	}

	e.ElapsedTime++
	// Tick timers.
	// TODO: Verify this behavior is correct.
	e.BehaviorState.ElapsedTime++
	if e.NextBehavior != nil {
		e.BehaviorState.Behavior.Cleanup(e, s)
		e.BehaviorState = EntityBehaviorState{e.NextBehavior, 0}
	}
	e.NextBehavior = nil
	e.BehaviorState.Behavior.Step(e, s)
}

type EntityBehavior interface {
	clone.Cloner[EntityBehavior]
	Appearance(e *Entity, b *bundle.Bundle) draw.Node
	Traits(e *Entity) EntityBehaviorTraits
	Step(e *Entity, s *State)
	Cleanup(e *Entity, s *State)
}
