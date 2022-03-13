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

type Damage struct {
	Base int

	ParalyzeTime Ticks
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

type HitTraits struct {
	FlashTime      Ticks
	ParalyzeTime   Ticks
	ConfuseTime    Ticks
	BlindTime      Ticks
	ImmobilizeTime Ticks
	FreezeTime     Ticks
	BubbleTime     Ticks

	Drag           DragType
	SlideDirection Direction

	SecondaryElementSword bool
	GuardPiercing         bool
	RemovesFlashing       bool
	Flinch                bool
}

type Hit struct {
	Traits      HitTraits
	TotalDamage int
}

func (h *Hit) AddDamage(d Damage) {
	v := d.Base + d.AttackPlus
	if d.DoubleDamage {
		v *= 2
	}
	h.TotalDamage += v
	if d.ParalyzeTime > 0 {
		h.Traits.ParalyzeTime = d.ParalyzeTime
	}
}

func (h *Hit) Merge(h2 Hit) {
	h.TotalDamage += h2.TotalDamage

	// TODO: Verify this is correct behavior.
	if h2.Traits.ParalyzeTime > h.Traits.ParalyzeTime {
		h.Traits.ParalyzeTime = h2.Traits.ParalyzeTime
	}
	if h2.Traits.ConfuseTime > h.Traits.ConfuseTime {
		h.Traits.ConfuseTime = h2.Traits.ConfuseTime
	}
	if h2.Traits.BlindTime > h.Traits.BlindTime {
		h.Traits.BlindTime = h2.Traits.BlindTime
	}
	if h2.Traits.ImmobilizeTime > h.Traits.ImmobilizeTime {
		h.Traits.ImmobilizeTime = h2.Traits.ImmobilizeTime
	}
	if h2.Traits.FreezeTime > h.Traits.FreezeTime {
		h.Traits.FreezeTime = h2.Traits.FreezeTime
	}
	if h2.Traits.BubbleTime > h.Traits.BubbleTime {
		h.Traits.BubbleTime = h2.Traits.BubbleTime
	}
	if h2.Traits.FlashTime > h.Traits.FlashTime {
		h.Traits.FlashTime = h2.Traits.FlashTime
	}
	if h2.Traits.Flinch {
		h.Traits.Flinch = true
	}
	if h2.Traits.Drag != DragTypeNone {
		h.Traits.Drag = h2.Traits.Drag
	}
	if h2.Traits.SlideDirection != DirectionNone {
		h.Traits.SlideDirection = h2.Traits.SlideDirection
	}
}

func MaybeApplyCounter(target *Entity, owner *Entity, h *Hit) {
	// From Alyrsc#7506:
	// I was mostly sure that it's frames 2-16 of an action.
	// I gathered that by frame stepping P2 while P1 had FullSynchro. The timing of the blue flashes was somewhat inconsistent, possibly because it's based on a global clock or counter, but those were the earliest and latest frames I saw.
	// TODO: Check the code for this.
	if target.BehaviorState.Behavior.Traits(target).CanBeCountered && target.BehaviorState.ElapsedTime < 15 {
		owner.IsFullSynchro = true
		h.Traits.FlashTime = 0
		h.Traits.ParalyzeTime = DefaultParalyzeTime
	}
}
