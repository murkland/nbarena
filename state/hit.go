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

type Hit struct {
	TotalDamage int

	FlashTime      Ticks
	ParalyzeTime   Ticks
	ConfuseTime    Ticks
	BlindTime      Ticks
	ImmobilizeTime Ticks
	FreezeTime     Ticks
	BubbleTime     Ticks
	Flinch         bool

	ForcedMovement ForcedMovement

	Element                 Element
	MustParalyzeImmediately bool
	CanCounter              bool
	SecondaryElementSword   bool
	GuardPiercing           bool
	RemovesFlashing         bool
}

func (h *Hit) AddDamage(d Damage) {
	v := d.Base + d.AttackPlus
	if d.DoubleDamage {
		v *= 2
	}
	h.TotalDamage += v
	if d.ParalyzeTime > 0 {
		if d.ParalyzeTime > h.ParalyzeTime {
			h.ParalyzeTime = d.ParalyzeTime
		}
		h.MustParalyzeImmediately = true
	}
	if d.Flinch {
		h.Flinch = true
	}
}
