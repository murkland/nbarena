package state

type Direction uint8

const (
	DirectionNone  Direction = 0
	DirectionUp    Direction = 0b0001
	DirectionDown  Direction = 0b0010
	DirectionLeft  Direction = 0b0100
	DirectionRight Direction = 0b1000
)

func (d Direction) XY() (int, int) {
	x := 0
	y := 0
	if d&DirectionLeft != 0 {
		x--
	}
	if d&DirectionRight != 0 {
		x++
	}
	if d&DirectionUp != 0 {
		y--
	}
	if d&DirectionDown != 0 {
		y++
	}
	return x, y
}

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

func DirectionDXDY(dx int, dy int) Direction {
	var d Direction
	if dx < 0 {
		d |= DirectionLeft
	} else if dx > 0 {
		d |= DirectionRight
	}
	if dy < 0 {
		d |= DirectionUp
	} else if dy > 0 {
		d |= DirectionDown
	}
	return d
}
