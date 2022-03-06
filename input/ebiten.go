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
		case ebiten.KeyX:
			intent.ChargeBasicWeapon = true
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		// TODO: Need to figure out the actual intent.
		intent.Action = ActionUseChip
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		intent.Action = ActionEndTurn
	}

	return intent
}

func CurrentIntent() Intent {
	return PressedKeysToIntent(inpututil.PressedKeys())
}
