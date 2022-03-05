package state

const DefaultFlashTime Ticks = 120

type Damage struct {
	Base int

	ParalyzeTime Ticks
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

	// ???
	RemovesFlashing bool
	Drag            bool
	Flinch          bool
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
}

func (h *Hit) Merge(h2 Hit) {
	h.TotalDamage += h2.TotalDamage

	// TODO: Verify this is correct behavior.
	if h2.ParalyzeTime > h.ParalyzeTime {
		h.ParalyzeTime = h2.ParalyzeTime
	}
	if h2.ConfuseTime > h.ConfuseTime {
		h.ConfuseTime = h2.ConfuseTime
	}
	if h2.BlindTime > h.BlindTime {
		h.BlindTime = h2.BlindTime
	}
	if h2.ImmobilizeTime > h.ImmobilizeTime {
		h.ImmobilizeTime = h2.ImmobilizeTime
	}
	if h2.FreezeTime > h.FreezeTime {
		h.FreezeTime = h2.FreezeTime
	}
	if h2.BubbleTime > h.BubbleTime {
		h.BubbleTime = h2.BubbleTime
	}
	if h2.FlashTime > h.FlashTime {
		h.FlashTime = h2.FlashTime
	}
	if h2.Flinch {
		h.Flinch = true
	}
}
