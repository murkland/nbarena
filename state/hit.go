package state

// This is maybe 1 frame shorter than expected?
const DefaultFlashTime Ticks = 119

// Apparently this is 1 frame shorter than expected - BN6 will remove paralyze if timeLeft - 1 == 0, but we only remove it if timeLeft = 0.
const DefaultParalyzeTime Ticks = 149

type Damage struct {
	Base int

	ParalyzeTime Ticks
	Uninstall    bool
	Skull        bool
	DoubleDamage bool
	AttackPlus   int
}

type HitTraits struct {
	FlashTime      Ticks
	ParalyzeTime   Ticks
	ConfuseTime    Ticks
	BlindTime      Ticks
	ImmobilizeTime Ticks
	FreezeTime     Ticks
	BubbleTime     Ticks

	Slide Slide

	Drag                  bool
	SecondaryElementSword bool
	GuardPiercing         bool
	RemovesFlashing       bool
	Counters              bool
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
	if h2.Traits.Counters {
		h.Traits.Counters = true
	}
	if h2.Traits.Drag {
		h.Traits.Drag = true
	}
	if h2.Traits.Slide.Direction != DirectionNone {
		h.Traits.Slide = h2.Traits.Slide
	}
}
