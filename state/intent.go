package state

type Intent struct {
	Direction         Direction
	UseChip           bool
	Confirm           bool
	CutIn             bool
	EndTurn           bool
	ChargeBasicWeapon bool
}
