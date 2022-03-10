package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/murkland/nbarena/state"
)

func PressedKeysToIntent(keys []ebiten.Key) state.Intent {
	var intent state.Intent
	for _, k := range keys {
		switch k {
		case ebiten.KeyUp:
			intent.Direction |= state.DirectionUp
		case ebiten.KeyDown:
			intent.Direction |= state.DirectionDown
		case ebiten.KeyLeft:
			intent.Direction |= state.DirectionLeft
		case ebiten.KeyRight:
			intent.Direction |= state.DirectionRight
		case ebiten.KeyZ:
			intent.UseChip = true
		case ebiten.KeyA, ebiten.KeyS:
			intent.EndTurn = true
		case ebiten.KeyX:
			intent.ChargeBasicWeapon = true
		}
	}
	return intent
}

func CurrentIntent() state.Intent {
	return PressedKeysToIntent(inpututil.PressedKeys())
}
