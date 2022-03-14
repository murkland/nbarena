package game

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/color"
	"strconv"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/keegancsmith/nth"
	"github.com/murkland/ctxwebrtc"
	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/chips"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/draw/styledtext"
	"github.com/murkland/nbarena/input"
	"github.com/murkland/nbarena/packets"
	"github.com/murkland/nbarena/sound"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/step"
	"github.com/murkland/ringbuf"
	"github.com/murkland/syncrand"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
)

var (
	debugSpewEntityState = flag.Bool("debug_spew_entity_state", false, "spew entity state")
)

const sceneWidth = 240
const sceneHeight = 160

const maxPendingIntents = 60

type clientState struct {
	isAnswerer bool

	OffererEntityID  state.EntityID
	AnswererEntityID state.EntityID

	committedState *state.State
	dirtyState     *state.State

	incomingIntents *ringbuf.RingBuf[state.Intent]
	outgoingIntents *ringbuf.RingBuf[state.Intent]
}

func (cs *clientState) SelfEntityID() state.EntityID {
	if cs.isAnswerer {
		return cs.AnswererEntityID
	}
	return cs.OffererEntityID
}

func (cs *clientState) OpponentEntityID() state.EntityID {
	if cs.isAnswerer {
		return cs.OffererEntityID
	}
	return cs.AnswererEntityID
}

func (cs *clientState) fastForward(b *bundle.Bundle) error {
	n := cs.outgoingIntents.Used()
	if cs.incomingIntents.Used() < n {
		n = cs.incomingIntents.Used()
	}

	ourIntents := make([]state.Intent, cs.outgoingIntents.Used())
	if err := cs.outgoingIntents.Peek(ourIntents, 0); err != nil {
		return err
	}
	if err := cs.outgoingIntents.Advance(n); err != nil {
		return err
	}

	theirIntents := make([]state.Intent, n)
	if err := cs.incomingIntents.Peek(theirIntents, 0); err != nil {
		return err
	}
	if err := cs.incomingIntents.Advance(n); err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		ourIntent := ourIntents[i]
		theirIntent := theirIntents[i]

		var offererIntent state.Intent
		var answererIntent state.Intent
		if cs.isAnswerer {
			offererIntent = theirIntent
			answererIntent = ourIntent
		} else {
			offererIntent = ourIntent
			answererIntent = theirIntent
		}

		cs.committedState.Entities[cs.OffererEntityID].Intent = offererIntent
		cs.committedState.Entities[cs.AnswererEntityID].Intent = answererIntent
		step.Step(cs.committedState, b)
	}

	cs.dirtyState = cs.committedState.Clone()
	for _, intent := range ourIntents[n:] {
		var offererIntent state.Intent
		var answererIntent state.Intent
		if cs.isAnswerer {
			offererIntent = cs.committedState.Entities[cs.OffererEntityID].LastIntent
			offererIntent.Direction = state.DirectionNone
			answererIntent = intent
		} else {
			offererIntent = intent
			answererIntent = cs.committedState.Entities[cs.AnswererEntityID].LastIntent
			answererIntent.Direction = state.DirectionNone
		}

		cs.dirtyState.Entities[cs.OffererEntityID].Intent = offererIntent
		cs.dirtyState.Entities[cs.AnswererEntityID].Intent = answererIntent
		step.Step(cs.dirtyState, b)
	}

	return nil
}

type Game struct {
	dc *ctxwebrtc.DataChannel

	compositor *draw.Compositor

	mixer          *beep.Mixer
	volume         *effects.Volume
	soundScheduler sound.Scheduler

	cs   *clientState
	csMu sync.Mutex

	bundle *bundle.Bundle

	paused bool

	inputFrameDelay int

	delayRingbuf   *ringbuf.RingBuf[time.Duration]
	delayRingbufMu sync.RWMutex
}

var sampleRate = beep.SampleRate(48000)

func New(b *bundle.Bundle, dc *ctxwebrtc.DataChannel, rng *syncrand.Source, isAnswerer bool, delaysWindowSize int, inputFrameDelay int) *Game {
	speaker.Init(sampleRate, 128)
	mixer := &beep.Mixer{}
	volume := &effects.Volume{Streamer: mixer}
	speaker.Play(volume)

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("nbarena")
	const defaultScale = 4
	ebiten.SetWindowSize(sceneWidth*defaultScale, sceneHeight*defaultScale)

	mixer.Add(b.BattleBGM.Streamer())

	s := state.New(rng)
	var offererEntityID state.EntityID
	{
		e := &state.Entity{
			HP:        1000,
			MaxHP:     1000,
			DisplayHP: 1000,

			Chips: slices.Clone(chips.Chips),

			PowerShotChargeTime: state.Ticks(50),

			TilePos:       state.TilePosXY(2, 2),
			FutureTilePos: state.TilePosXY(2, 2),

			BehaviorState: state.EntityBehaviorState{
				Behavior: &behaviors.Idle{},
			},
		}
		s.AttachEntity(e)
		offererEntityID = e.ID()
		s.Field.Tiles[e.TilePos].Reserver = e.ID()
	}

	var answererEntityID state.EntityID
	{
		e := &state.Entity{
			HP:        1000,
			MaxHP:     1000,
			DisplayHP: 1000,

			Chips: slices.Clone(chips.Chips),

			PowerShotChargeTime: state.Ticks(50),

			IsFlipped:            true,
			IsAlliedWithAnswerer: true,

			TilePos:       state.TilePosXY(5, 2),
			FutureTilePos: state.TilePosXY(5, 2),

			BehaviorState: state.EntityBehaviorState{
				Behavior: &behaviors.Idle{},
			},
		}
		s.AttachEntity(e)
		answererEntityID = e.ID()
		s.Field.Tiles[e.TilePos].Reserver = e.ID()
	}

	g := &Game{
		bundle:         b,
		dc:             dc,
		volume:         volume,
		mixer:          mixer,
		soundScheduler: sound.NewScheduler(sampleRate, mixer),
		cs: &clientState{
			OffererEntityID:  offererEntityID,
			AnswererEntityID: answererEntityID,

			isAnswerer: isAnswerer,

			committedState: s,
			dirtyState:     s.Clone(),

			incomingIntents: ringbuf.New[state.Intent](maxPendingIntents),
			outgoingIntents: ringbuf.New[state.Intent](maxPendingIntents),
		},
		inputFrameDelay: inputFrameDelay,
		delayRingbuf:    ringbuf.New[time.Duration](delaysWindowSize),
	}
	return g
}

type orderableSlice[T constraints.Ordered] []T

func (s orderableSlice[T]) Len() int {
	return len(s)
}

func (s orderableSlice[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s orderableSlice[T]) Less(i, j int) bool {
	return s[i] < s[j]
}

func (g *Game) medianDelay() time.Duration {
	g.delayRingbufMu.RLock()
	defer g.delayRingbufMu.RUnlock()

	if g.delayRingbuf.Used() == 0 {
		return 0
	}

	delays := make([]time.Duration, g.delayRingbuf.Used())
	if err := g.delayRingbuf.Peek(delays, 0); err != nil {
		panic(err)
	}

	i := len(delays) / 2
	nth.Element(orderableSlice[time.Duration](delays), i)
	return delays[i]
}

func (g *Game) sendPings(ctx context.Context) error {
	for {
		now := time.Now()
		if err := packets.Send(ctx, g.dc, packets.Ping{
			ID: uint64(now.UnixMicro()),
		}); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
		}
	}
}

func (g *Game) handleConn(ctx context.Context) error {
	for {
		packet, err := packets.Recv(ctx, g.dc)
		if err != nil {
			return err
		}

		switch p := packet.(type) {
		case packets.Ping:
			if err := packets.Send(ctx, g.dc, packets.Pong{ID: p.ID}); err != nil {
				return err
			}
		case packets.Pong:
			if err := (func() error {
				g.delayRingbufMu.Lock()
				defer g.delayRingbufMu.Unlock()

				if g.delayRingbuf.Free() == 0 {
					g.delayRingbuf.Advance(1)
				}

				delay := time.Now().Sub(time.UnixMicro(int64(p.ID)))
				if err := g.delayRingbuf.Push([]time.Duration{delay}); err != nil {
					return err
				}
				return nil
			})(); err != nil {
				return err
			}
		case packets.Intent:
			if err := (func() error {
				g.csMu.Lock()
				defer g.csMu.Unlock()

				nextTick := uint32(int(g.cs.committedState.ElapsedTime) + g.cs.incomingIntents.Used() + 1)
				if p.ForTick != nextTick {
					return fmt.Errorf("expected intent for %d but it was for %d", nextTick, p.ForTick)
				}

				if err := g.cs.incomingIntents.Push([]state.Intent{p.Intent}); err != nil {
					return err
				}

				if err := g.cs.fastForward(g.bundle); err != nil {
					return err
				}

				return nil
			})(); err != nil {
				return err
			}

		}
	}
}

func scaleFactor(bounds image.Rectangle) int {
	k := bounds.Dx() / sceneWidth
	if s := bounds.Dy() / sceneHeight; s < k {
		k = s
	}
	return k
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func makeTextGradient(color0 color.Color, color1 color.Color, color2 color.Color) *ebiten.Image {
	gradientImage := ebiten.NewImage(1, 10)
	gradientImage.Set(0, 0, color0)
	gradientImage.Set(0, 1, color0)
	gradientImage.Set(0, 2, color0)
	gradientImage.Set(0, 3, color1)
	gradientImage.Set(0, 4, color2)
	gradientImage.Set(0, 5, color2)
	gradientImage.Set(0, 6, color1)
	gradientImage.Set(0, 7, color0)
	gradientImage.Set(0, 8, color0)
	gradientImage.Set(0, 9, color0)
	return gradientImage
}

var (
	whiteTextGradient      = makeTextGradient(color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
	chipDamageTextGradient = makeTextGradient(color.RGBA{0xFF, 0xA5, 0x00, 0xFF}, color.RGBA{0xFF, 0xDE, 0x00, 0xFF}, color.RGBA{0xFF, 0xDE, 0x00, 0xFF})
	hpNeutralTextGradient  = makeTextGradient(color.RGBA{0xCE, 0xE7, 0xFF, 0xFF}, color.RGBA{0xE7, 0xEF, 0xFF, 0xFF}, color.RGBA{0xF7, 0xF7, 0xF7, 0xFF})
	hpLossTextGradient     = makeTextGradient(color.RGBA{0xFF, 0xA5, 0x21, 0xFF}, color.RGBA{0xFF, 0xC6, 0x63, 0xFF}, color.RGBA{0xFF, 0xEF, 0xAD, 0xFF})
	hpGainTextGradient     = makeTextGradient(color.RGBA{0x39, 0xFF, 0x94, 0xFF}, color.RGBA{0x84, 0xFF, 0xC6, 0xFF}, color.RGBA{0xD7, 0xFF, 0xF7, 0xFF})
)

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xCE, 0xFF, 0xC6, 0xFF})
	k := scaleFactor(screen.Bounds())
	compositorBounds := image.Rect(0, 0, sceneWidth*k, sceneHeight*k)

	if g.compositor == nil || g.compositor.Bounds() != compositorBounds {
		g.compositor = draw.NewCompositor(compositorBounds, 9)
	}

	g.csMu.Lock()
	defer g.csMu.Unlock()

	state := g.cs.dirtyState
	if g.cs.isAnswerer {
		state = g.cs.dirtyState.Clone()
		state.Flip()
	}

	g.soundScheduler.Dispatch(g.bundle, state.Sounds)

	rootNode := &draw.OptionsNode{}
	sceneNode := &draw.OptionsNode{}
	sceneNode.Opts.GeoM.Scale(float64(k), float64(k))
	sceneNode.Children = append(sceneNode.Children, state.Appearance(g.bundle))
	if state.IsInTimeStop {
		timestopOverlay := &draw.OptionsNode{}
		overlay := ebiten.NewImage(sceneWidth, sceneHeight)
		overlay.Fill(color.Black)
		timestopOverlay.Opts.ColorM.Scale(1.0, 1.0, 1.0, 0.25)
		timestopOverlay.Children = append(timestopOverlay.Children, &draw.ImageNode{Image: overlay})
		sceneNode.Children = append(sceneNode.Children, timestopOverlay)
	}
	sceneNode.Children = append(sceneNode.Children, g.uiAppearance())
	rootNode.Children = append(rootNode.Children, sceneNode)
	if *debugSpewEntityState {
		rootNode.Children = append(rootNode.Children, g.makeDebugDrawNode())
	}

	g.compositor.Clear()
	rootNode.Draw(g.compositor, &ebiten.DrawImageOptions{})

	var opts ebiten.DrawImageOptions
	opts.GeoM.Translate(float64((screen.Bounds().Dx()-g.compositor.Bounds().Dx())/2), float64((screen.Bounds().Dy()-g.compositor.Bounds().Dy())/2))
	g.compositor.Draw(screen, &opts)
}

func chipPlaqueApperance(b *bundle.Bundle, chip *state.Chip, attackPlus int, doubleDamage bool, anchor styledtext.Anchor) draw.Node {
	spans := []styledtext.Span{{Text: chip.Name, Background: whiteTextGradient}}
	if chip.BaseDamage > 0 {
		s := strconv.Itoa(chip.BaseDamage)
		if attackPlus > 0 {
			s += "+" + strconv.Itoa(attackPlus)
		}
		spans = append(spans, styledtext.Span{Text: s, Background: chipDamageTextGradient})
	}
	if doubleDamage {
		spans = append(spans, styledtext.Span{Text: "×2", Background: whiteTextGradient})
	}
	return styledtext.MakeNode(spans, anchor, b.TallFont, styledtext.BorderRightBottom, color.RGBA{0x00, 0x00, 0x00, 0xff})
}

var (
	hpBoxImage = func() *ebiten.Image {
		img := ebiten.NewImage(44, 16)
		img.Fill(color.RGBA{0xf8, 0xff, 0xff, 0xff})
		img.SubImage(image.Rect(2, 1, 42, 15)).(*ebiten.Image).Fill(color.RGBA{0x39, 0x52, 0x6b, 0xff})
		return img
	}()
)

func (g *Game) uiAppearance() draw.Node {
	rootNode := &draw.OptionsNode{Layer: 9}
	{
		self := g.cs.dirtyState.Entities[g.cs.SelfEntityID()]

		gradientImage := hpNeutralTextGradient
		if self.DisplayHP > self.HP {
			gradientImage = hpLossTextGradient
		} else if self.DisplayHP < self.HP {
			gradientImage = hpGainTextGradient
		}

		hpPlaqueNode := &draw.OptionsNode{}
		hpPlaqueNode.Opts.GeoM.Translate(float64(2), float64(0))
		rootNode.Children = append(rootNode.Children, hpPlaqueNode)

		hpPlaqueBgNode := &draw.OptionsNode{}
		hpPlaqueNode.Children = append(hpPlaqueNode.Children, hpPlaqueBgNode)
		hpPlaqueBgNode.Children = append(hpPlaqueBgNode.Children, &draw.ImageNode{Image: hpBoxImage})

		hpPlaqueTextNode := &draw.OptionsNode{}
		hpPlaqueTextNode.Opts.GeoM.Translate(float64(38), float64(3))
		hpPlaqueNode.Children = append(hpPlaqueNode.Children, hpPlaqueTextNode)
		hpPlaqueTextNode.Children = append(hpPlaqueTextNode.Children, styledtext.MakeNode([]styledtext.Span{{Text: strconv.Itoa(self.DisplayHP), Background: gradientImage}}, styledtext.AnchorRight|styledtext.AnchorTop, g.bundle.TallFont, styledtext.BorderNone, color.RGBA{}))
	}

	{
		opponent := g.cs.dirtyState.Entities[g.cs.OpponentEntityID()]
		if opponent.ChipPlaque.Chip != nil {
			chipPlaqueNode := &draw.OptionsNode{}
			rootNode.Children = append(rootNode.Children, chipPlaqueNode)
			chipPlaqueNode.Opts.GeoM.Translate(float64(sceneWidth-16), 36)
			chipPlaqueNode.Children = append(chipPlaqueNode.Children, chipPlaqueApperance(g.bundle, opponent.ChipPlaque.Chip, opponent.ChipPlaque.AttackPlus, opponent.ChipPlaque.DoubleDamage, styledtext.AnchorRight|styledtext.AnchorTop))
		}
	}

	// TODO: Render chip. Must be not in chip use lockout.
	self := g.cs.dirtyState.Entities[g.cs.SelfEntityID()]
	if self.ChipUseLockoutTimeLeft == 0 && len(self.Chips) > 0 {
		chip := self.Chips[len(self.Chips)-1]

		chipTextNode := &draw.OptionsNode{}
		rootNode.Children = append(rootNode.Children, chipTextNode)
		chipTextNode.Opts.GeoM.Translate(1, float64(sceneHeight-12))
		chipTextNode.Children = append(chipTextNode.Children, chipPlaqueApperance(g.bundle, chip, 0, self.DoubleDamage(), styledtext.AnchorLeft|styledtext.AnchorTop))
	}

	return rootNode
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.volume.Silent = !g.volume.Silent
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.paused = !g.paused
	}

	if g.paused && !inpututil.IsKeyJustPressed(ebiten.KeyPeriod) {
		return nil
	}

	g.csMu.Lock()
	defer g.csMu.Unlock()

	highWaterMark := int(g.medianDelay()*time.Duration(ebiten.MaxTPS())/2/time.Second+1) - g.inputFrameDelay
	if highWaterMark < 1 {
		highWaterMark = 1
	}

	if g.cs.outgoingIntents.Used() >= highWaterMark {
		// Pause until we have enough space.
		return nil
	}

	intent := input.CurrentIntent()
	if g.cs.isAnswerer {
		intent.Direction = intent.Direction.FlipH()
	}

	forTick := uint32(g.cs.dirtyState.ElapsedTime + 1)

	ctx := context.Background()

	if err := packets.Send(ctx, g.dc, packets.Intent{ForTick: forTick, Intent: intent}); err != nil {
		return err
	}

	if err := g.cs.outgoingIntents.Push([]state.Intent{intent}); err != nil {
		return err
	}

	if err := g.cs.fastForward(g.bundle); err != nil {
		return err
	}

	return nil
}

func (g *Game) RunBackgroundTasks(ctx context.Context) error {
	errg, ctx := errgroup.WithContext(ctx)

	errg.Go(func() error {
		return g.handleConn(ctx)
	})

	errg.Go(func() error {
		return g.sendPings(ctx)
	})

	return errg.Wait()
}
