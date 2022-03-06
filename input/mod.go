package input

type Direction uint8

const (
	DirectionNone  Direction = 0
	DirectionUp    Direction = 0b0001
	DirectionDown  Direction = 0b0010
	DirectionLeft  Direction = 0b0100
	DirectionRight Direction = 0b1000
)

func (d Direction) FlipH() Direction {
	d2 := d & ^(DirectionLeft | DirectionRight)
	if d&DirectionLeft == DirectionLeft {
		d2 |= DirectionRight
	}
	if d&DirectionRight == DirectionRight {
		d2 |= DirectionLeft
	}
	return d2
}

func (d Direction) FlipV() Direction {
	d2 := d & ^(DirectionUp | DirectionDown)
	if d&DirectionUp == DirectionUp {
		d2 |= DirectionDown
	}
	if d&DirectionDown == DirectionDown {
		d2 |= DirectionUp
	}
	return d2
}

type Action uint8

type Intent struct {
	Direction         Direction
	UseChip           bool
	Confirm           bool
	CutIn             bool
	EndTurn           bool
	ChargeBasicWeapon bool
}
