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

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/keegancsmith/nth"
	"github.com/murkland/ctxwebrtc"
	"github.com/murkland/nbarena/behaviors"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/chips"
	"github.com/murkland/nbarena/draw"
	"github.com/murkland/nbarena/input"
	"github.com/murkland/nbarena/packets"
	"github.com/murkland/nbarena/state"
	"github.com/murkland/nbarena/step"
	"github.com/murkland/ringbuf"
	"github.com/murkland/syncrand"
	"golang.org/x/exp/constraints"
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

	OffererEntityID  int
	AnswererEntityID int

	committedState state.State
	dirtyState     state.State

	lastIncomingIntent input.Intent

	incomingIntents *ringbuf.RingBuf[input.Intent]
	outgoingIntents *ringbuf.RingBuf[input.Intent]
}

func (cs *clientState) SelfEntityID() int {
	if cs.isAnswerer {
		return cs.AnswererEntityID
	}
	return cs.OffererEntityID
}

func (cs *clientState) fastForward() error {
	n := cs.outgoingIntents.Used()
	if cs.incomingIntents.Used() < n {
		n = cs.incomingIntents.Used()
	}

	ourIntents := make([]input.Intent, cs.outgoingIntents.Used())
	if err := cs.outgoingIntents.Peek(ourIntents, 0); err != nil {
		return err
	}
	if err := cs.outgoingIntents.Advance(n); err != nil {
		return err
	}

	theirIntents := make([]input.Intent, n)
	if err := cs.incomingIntents.Peek(theirIntents, 0); err != nil {
		return err
	}
	if err := cs.incomingIntents.Advance(n); err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		ourIntent := ourIntents[i]
		theirIntent := theirIntents[i]

		var offererIntent input.Intent
		var answererIntent input.Intent
		if cs.isAnswerer {
			offererIntent = theirIntent
			answererIntent = ourIntent
		} else {
			offererIntent = ourIntent
			answererIntent = theirIntent
		}

		cs.committedState.Entities[cs.OffererEntityID].Intent = offererIntent
		cs.committedState.Entities[cs.AnswererEntityID].Intent = answererIntent
		step.Step(&cs.committedState)
	}

	cs.dirtyState = cs.committedState.Clone()
	for _, intent := range ourIntents[n:] {
		var offererIntent input.Intent
		var answererIntent input.Intent
		if cs.isAnswerer {
			offererIntent = cs.lastIncomingIntent
			offererIntent.Direction = input.DirectionNone
			answererIntent = intent
		} else {
			offererIntent = intent
			answererIntent = cs.lastIncomingIntent
			answererIntent.Direction = input.DirectionNone
		}

		cs.dirtyState.Entities[cs.OffererEntityID].Intent = offererIntent
		cs.dirtyState.Entities[cs.AnswererEntityID].Intent = answererIntent
		step.Step(&cs.dirtyState)
	}

	return nil
}

type Game struct {
	dc *ctxwebrtc.DataChannel

	compositor *draw.Compositor

	cs   *clientState
	csMu sync.Mutex

	bundle *bundle.Bundle

	paused bool

	inputFrameDelay int

	delayRingbuf   *ringbuf.RingBuf[time.Duration]
	delayRingbufMu sync.RWMutex
}

func New(b *bundle.Bundle, dc *ctxwebrtc.DataChannel, rng *syncrand.Source, isAnswerer bool, delaysWindowSize int, inputFrameDelay int) *Game {
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("nbarena")
	const defaultScale = 4
	ebiten.SetWindowSize(sceneWidth*defaultScale, sceneHeight*defaultScale)

	s := state.New(rng)
	var offererEntityID int
	{
		e := &state.Entity{
			HP:        1000,
			DisplayHP: 1000,

			Chips: []state.Chip{chips.Chips[0], chips.Chips[1], chips.Chips[2], chips.Chips[3], chips.Chips[4]},

			Traits: state.EntityTraits{
				ExtendsTileLifetime: true,
			},

			PowerShotChargeTime: state.Ticks(50),

			TilePos:       state.TilePosXY(2, 2),
			FutureTilePos: state.TilePosXY(2, 2),
		}
		e.SetBehavior(&behaviors.Idle{})
		offererEntityID = s.AddEntity(e)
	}

	var answererEntityID int
	{
		e := &state.Entity{
			HP:        1000,
			DisplayHP: 1000,

			Chips: []state.Chip{chips.Chips[0], chips.Chips[1], chips.Chips[2], chips.Chips[3], chips.Chips[4]},

			Traits: state.EntityTraits{
				ExtendsTileLifetime: true,
			},

			PowerShotChargeTime: state.Ticks(50),

			IsFlipped:            true,
			IsAlliedWithAnswerer: true,

			TilePos:       state.TilePosXY(5, 2),
			FutureTilePos: state.TilePosXY(5, 2),
		}
		e.SetBehavior(&behaviors.Idle{})
		answererEntityID = s.AddEntity(e)
	}

	g := &Game{
		bundle: b,
		dc:     dc,
		cs: &clientState{
			OffererEntityID:  offererEntityID,
			AnswererEntityID: answererEntityID,

			isAnswerer: isAnswerer,

			committedState: s,
			dirtyState:     s.Clone(),

			incomingIntents: ringbuf.New[input.Intent](maxPendingIntents),
			outgoingIntents: ringbuf.New[input.Intent](maxPendingIntents),
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

				if err := g.cs.incomingIntents.Push([]input.Intent{p.Intent}); err != nil {
					return err
				}

				if err := g.cs.fastForward(); err != nil {
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

func makeHPGradient(color0 color.Color, color1 color.Color, color2 color.Color) *ebiten.Image {
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
	hpNeutralGradient = makeHPGradient(color.RGBA{0xCE, 0xE7, 0xFF, 0xFF}, color.RGBA{0xE7, 0xE7, 0xFF, 0xFF}, color.RGBA{0xF7, 0xF7, 0xF7, 0xFF})
	hpLossGradient    = makeHPGradient(color.RGBA{0xFF, 0xA5, 0x21, 0xFF}, color.RGBA{0xFF, 0xC6, 0x63, 0xFF}, color.RGBA{0xFF, 0xEF, 0xAD, 0xFF})
	hpGainGradient    = makeHPGradient(color.RGBA{0x39, 0xFF, 0x94, 0xFF}, color.RGBA{0x84, 0xFF, 0xC6, 0xFF}, color.RGBA{0xD7, 0xFF, 0xF7, 0xFF})
)

func (g *Game) Draw(screen *ebiten.Image) {
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

	rootNode := &draw.OptionsNode{}
	sceneNode := &draw.OptionsNode{}
	sceneNode.Opts.GeoM.Scale(float64(k), float64(k))
	sceneNode.Children = append(sceneNode.Children, state.Appearance(g.bundle))
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

func (g *Game) uiAppearance() draw.Node {
	rootNode := &draw.OptionsNode{}
	{
		self := g.cs.dirtyState.Entities[g.cs.SelfEntityID()]

		hpText := strconv.Itoa(int(self.DisplayHP))
		rect := text.BoundString(g.bundle.TallFont, hpText)

		textOnlyImage := ebiten.NewImage(rect.Max.X, rect.Dy())
		{
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(0), float64(-rect.Min.Y))
			text.DrawWithOptions(textOnlyImage, hpText, g.bundle.TallFont, opts)
		}

		gradientImage := hpNeutralGradient
		if self.DisplayHP > self.HP {
			gradientImage = hpLossGradient
		} else if self.DisplayHP < self.HP {
			gradientImage = hpGainGradient
		}

		hpTextImage := ebiten.NewImage(textOnlyImage.Bounds().Dx(), textOnlyImage.Bounds().Dy())
		{
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Scale(float64(hpTextImage.Bounds().Dx()), 1.0)
			hpTextImage.DrawImage(gradientImage, opts)
		}

		{
			opts := &ebiten.DrawImageOptions{}
			opts.CompositeMode = ebiten.CompositeModeDestinationIn
			hpTextImage.DrawImage(textOnlyImage, opts)
		}

		hpPlaqueNode := &draw.OptionsNode{}
		hpPlaqueNode.Opts.GeoM.Translate(float64(-hpTextImage.Bounds().Dx()+39), float64(3))
		hpPlaqueNode.Children = append(hpPlaqueNode.Children, &draw.ImageNode{Image: hpTextImage})
		rootNode.Children = append(rootNode.Children, hpPlaqueNode)
	}

	// TODO: Render chip. Must be not in chip use lockout.
	self := g.cs.dirtyState.Entities[g.cs.SelfEntityID()]
	if self.ChipUseLockoutTimeLeft == 0 && len(self.Chips) > 0 {
		chip := self.Chips[len(self.Chips)-1]

		rect := text.BoundString(g.bundle.TallFont, chip.Name)

		chipTextNode := &draw.OptionsNode{}
		chipTextNode.Opts.GeoM.Translate(1, float64(sceneHeight-12))
		rootNode.Children = append(rootNode.Children, chipTextNode)

		chipTextBgNode := &draw.OptionsNode{}
		chipTextNode.Children = append(chipTextNode.Children, chipTextBgNode)
		chipTextBgNode.Opts.ColorM.Translate(-1.0, -1.0, -1.0, 0.0)
		chipTextBgNode.Opts.GeoM.Translate(float64(1), float64(1))
		chipTextBgNode.Children = append(chipTextBgNode.Children, &draw.TextNode{Text: chip.Name, Face: g.bundle.TallFont})

		chipTextFgNode := &draw.OptionsNode{}
		chipTextNode.Children = append(chipTextNode.Children, chipTextFgNode)
		chipTextFgNode.Opts.ColorM.Translate(1.0, 1.0, 1.0, 0.0)
		chipTextFgNode.Children = append(chipTextFgNode.Children, &draw.TextNode{Text: chip.Name, Face: g.bundle.TallFont})

		if chip.Damage > 0 {
			chipDamageBgNode := &draw.OptionsNode{}
			chipTextNode.Children = append(chipTextNode.Children, chipDamageBgNode)
			chipDamageBgNode.Opts.ColorM.Translate(-1.0, -1.0, -1.0, 0.0)
			chipDamageBgNode.Opts.GeoM.Translate(float64(rect.Max.X+2), 0)
			chipDamageBgNode.Opts.GeoM.Translate(float64(1), float64(1))
			chipDamageBgNode.Children = append(chipDamageBgNode.Children, &draw.TextNode{Text: strconv.Itoa(chip.Damage), Face: g.bundle.TallFont})

			chipDamageFgNode := &draw.OptionsNode{}
			chipTextNode.Children = append(chipTextNode.Children, chipDamageFgNode)
			chipDamageFgNode.Opts.GeoM.Translate(float64(rect.Max.X+2), 0)
			chipDamageFgNode.Opts.ColorM.Translate(1.0, 1.0, 1.0, 0.0)
			chipDamageFgNode.Children = append(chipDamageFgNode.Children, &draw.TextNode{Text: strconv.Itoa(chip.Damage), Face: g.bundle.TallFont})
		}
	}

	return rootNode
}

func (g *Game) Update() error {
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

	if err := g.cs.outgoingIntents.Push([]input.Intent{intent}); err != nil {
		return err
	}

	if err := g.cs.fastForward(); err != nil {
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
