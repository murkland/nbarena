package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func PressedKeysToIntent(keys []ebiten.Key) Intent {
	var intent Intent
	for _, k := range keys {
		switch k {
		case ebiten.KeyUp:
			intent.Direction |= DirectionUp
		case ebiten.KeyDown:
			intent.Direction |= DirectionDown
		case ebiten.KeyLeft:
			intent.Direction |= DirectionLeft
		case ebiten.KeyRight:
			intent.Direction |= DirectionRight
		case ebiten.KeyZ:
			if intent.Action == ActionNone {
				// TODO: Need to figure out the actual intent.
				intent.Action = ActionUseChip
			}
		case ebiten.KeyA, ebiten.KeyS:
			intent.Action = ActionEndTurn
		case ebiten.KeyX:
			intent.ChargeBasicWeapon = true
		}
	}
	return intent
}

func CurrentIntent() Intent {
	return PressedKeysToIntent(inpututil.PressedKeys())
}
