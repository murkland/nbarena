package game

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/keegancsmith/nth"
	"github.com/yumland/ctxwebrtc"
	"github.com/yumland/ringbuf"
	"github.com/yumland/syncrand"
	"github.com/yumland/yumbattle/bundle"
	"github.com/yumland/yumbattle/draw"
	"github.com/yumland/yumbattle/input"
	"github.com/yumland/yumbattle/packets"
	"github.com/yumland/yumbattle/state"
	"golang.org/x/exp/constraints"
	"golang.org/x/sync/errgroup"
)

const sceneWidth = 240
const sceneHeight = 160

const maxPendingIntents = 60

type clientState struct {
	isAnswerer bool

	committedState state.State
	dirtyState     state.State

	incomingIntents *ringbuf.RingBuf[input.Intent]
	outgoingIntents *ringbuf.RingBuf[input.Intent]
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
		var answerwerIntent input.Intent
		if cs.isAnswerer {
			offererIntent = theirIntent
			answerwerIntent = ourIntent
		} else {
			offererIntent = ourIntent
			answerwerIntent = theirIntent
		}

		cs.committedState.Step()
		cs.committedState.Apply(offererIntent, answerwerIntent)
	}

	cs.dirtyState = cs.committedState.Clone()
	for _, intent := range ourIntents[n:] {
		var offererIntent input.Intent
		var answerwerIntent input.Intent
		if cs.isAnswerer {
			answerwerIntent = intent
		} else {
			offererIntent = intent
		}

		cs.dirtyState.Step()
		cs.dirtyState.Apply(offererIntent, answerwerIntent)
	}

	return nil
}

type Game struct {
	sceneGeoM ebiten.GeoM
	dc        *ctxwebrtc.DataChannel

	cs   *clientState
	csMu sync.Mutex

	bundle *bundle.Bundle

	delayRingbuf   *ringbuf.RingBuf[time.Duration]
	delayRingbufMu sync.RWMutex
}

func New(b *bundle.Bundle, dc *ctxwebrtc.DataChannel, rng *syncrand.Source, isAnswerer bool, delaysWindowSize int) *Game {
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("yumbattle")
	const defaultScale = 4
	ebiten.SetWindowSize(sceneWidth*defaultScale, sceneHeight*defaultScale)

	s := state.New(rng)

	g := &Game{
		bundle: b,
		dc:     dc,
		cs: &clientState{
			isAnswerer: isAnswerer,

			committedState: s,
			dirtyState:     s.Clone(),

			incomingIntents: ringbuf.New[input.Intent](maxPendingIntents),
			outgoingIntents: ringbuf.New[input.Intent](maxPendingIntents),
		},
		delayRingbuf: ringbuf.New[time.Duration](delaysWindowSize),
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

				nextTick := uint32(g.cs.committedState.ElapsedTicks() + g.cs.incomingIntents.Used() + 1)
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

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	scaleFactor := outsideWidth / sceneWidth
	if s := outsideHeight / sceneHeight; s < scaleFactor {
		scaleFactor = s
	}

	insideWidth := sceneWidth * scaleFactor
	insideHeight := sceneHeight * scaleFactor

	g.sceneGeoM = ebiten.GeoM{}
	g.sceneGeoM.Scale(float64(scaleFactor), float64(scaleFactor))
	g.sceneGeoM.Translate(float64(outsideWidth-insideWidth)/2, float64(outsideHeight-insideHeight)/2)

	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.csMu.Lock()
	defer g.csMu.Unlock()

	rootNode := &draw.OptionsNode{}
	sceneNode := &draw.OptionsNode{Opts: ebiten.DrawImageOptions{GeoM: g.sceneGeoM}}
	sceneNode.Children = append(sceneNode.Children, g.cs.dirtyState.Appearance(g.bundle))
	rootNode.Children = append(rootNode.Children, sceneNode)
	rootNode.Children = append(rootNode.Children, g.makeDebugDrawNode())
	rootNode.Draw(screen, &ebiten.DrawImageOptions{})
}

func (g *Game) Update() error {
	g.csMu.Lock()
	defer g.csMu.Unlock()

	if g.cs.outgoingIntents.Used() >= int(g.medianDelay()*time.Duration(ebiten.MaxTPS())/2/time.Second+1) {
		// Pause until we have enough space.
		return nil
	}

	intent := input.CurrentIntent()
	forTick := uint32(g.cs.dirtyState.ElapsedTicks() + 1)

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
