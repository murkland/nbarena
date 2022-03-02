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
	Drag bool
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
	h.ParalyzeTime = h2.ParalyzeTime
	h.ConfuseTime = h2.ConfuseTime
	h.BlindTime = h2.BlindTime
	h.ImmobilizeTime = h2.ImmobilizeTime
	h.FreezeTime = h2.FreezeTime
	h.BubbleTime = h2.BubbleTime
}
