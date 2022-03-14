package state

// This is maybe 1 frame shorter than expected?
const DefaultFlashTime Ticks = 119

// Apparently this is 1 frame shorter than expected - BN6 will remove paralyze if timeLeft - 1 == 0, but we only remove it if timeLeft = 0.
const DefaultParalyzeTime Ticks = 149

type Element int

const (
	ElementNull   Element = 0
	ElementFire   Element = 1
	ElementAqua   Element = 2
	ElementElec   Element = 3
	ElementWood   Element = 4
	ElementSword  Element = 5
	ElementWind   Element = 6
	ElementCursor Element = 7
	ElementBreak  Element = 8
)

func (e Element) IsSuperEffectiveAgainst(e2 Element) bool {
	return (e == ElementFire && e2 == ElementWood) ||
		(e == ElementAqua && e2 == ElementFire) ||
		(e == ElementElec && e2 == ElementAqua) ||
		(e == ElementWood && e2 == ElementElec) ||
		(e == ElementSword && e2 == ElementWind) ||
		(e == ElementWind && e2 == ElementCursor) ||
		(e == ElementCursor && e2 == ElementBreak) ||
		(e == ElementBreak && e2 == ElementSword)
}

type Damage struct {
	Base int

	ParalyzeTime Ticks
	Flinch       bool
	Uninstall    bool
	Skull        bool
	DoubleDamage bool
	AttackPlus   int
}

type DragType int

const (
	DragTypeNone  DragType = 0
	DragTypeSmall DragType = 1
	DragTypeBig   DragType = 2
)

type Hit struct {
	TotalDamage int

	FlashTime      Ticks
	ParalyzeTime   Ticks
	ConfuseTime    Ticks
	BlindTime      Ticks
	ImmobilizeTime Ticks
	FreezeTime     Ticks
	BubbleTime     Ticks

	Drag           DragType
	SlideDirection Direction

	Element               Element
	SecondaryElementSword bool
	GuardPiercing         bool
	RemovesFlashing       bool
	Flinch                bool
}

func (h *Hit) AddDamage(d Damage) {
	v := d.Base + d.AttackPlus
	if d.DoubleDamage {
		v *= 2
	}
	h.TotalDamage += v
	if d.ParalyzeTime > 0 {
		h.ParalyzeTime = d.ParalyzeTime
	}
	if d.Flinch {
		h.Flinch = true
	}
}

func MaybeApplyCounter(target *Entity, owner *Entity, h *Hit) {
	// From Alyrsc#7506:
	// I was mostly sure that it's frames 2-16 of an action.
	// I gathered that by frame stepping P2 while P1 had FullSynchro. The timing of the blue flashes was somewhat inconsistent, possibly because it's based on a global clock or counter, but those were the earliest and latest frames I saw.
	// TODO: Check the code for this.
	if target.BehaviorState.Behavior.Traits(target).CanBeCountered && target.BehaviorState.ElapsedTime < 15 {
		owner.Emotion = EmotionFullSynchro
		h.FlashTime = 0
		h.ParalyzeTime = DefaultParalyzeTime
	}
}
