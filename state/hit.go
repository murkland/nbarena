package state

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
	Flinch bool
	Drag   bool
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
	if h2.ParalyzeTime > 0 {
		h.ParalyzeTime = h2.ParalyzeTime
	}
	if h2.ConfuseTime > 0 {
		h.ConfuseTime = h2.ConfuseTime
	}
	if h2.BlindTime > 0 {
		h.BlindTime = h2.BlindTime
	}
	if h2.ImmobilizeTime > 0 {
		h.ImmobilizeTime = h2.ImmobilizeTime
	}
	if h2.FreezeTime > 0 {
		h.FreezeTime = h2.FreezeTime
	}
	if h2.BubbleTime > 0 {
		h.BubbleTime = h2.BubbleTime
	}
	if h2.FlashTime > 0 {
		h.FlashTime = h2.FlashTime
	}
}
