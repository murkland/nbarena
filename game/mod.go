package game

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yumland/ctxwebrtc"
	"github.com/yumland/ringbuf"
	"github.com/yumland/syncrand"
	"github.com/yumland/yumbattle/input"
	"github.com/yumland/yumbattle/packets"
	"github.com/yumland/yumbattle/state"
	"golang.org/x/sync/errgroup"
)

const renderWidth = 240
const renderHeight = 160

const maxPendingIntents = 60

type clientState struct {
	isOfferer      bool
	committedState *state.State
	dirtyState     *state.State

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

		var tickIntents []input.Intent
		if cs.isOfferer {
			tickIntents = []input.Intent{ourIntent, theirIntent}
		} else {
			tickIntents = []input.Intent{theirIntent, ourIntent}
		}

		if err := cs.committedState.Update(tickIntents); err != nil {
			return err
		}
	}

	cs.dirtyState = cs.committedState.Clone()
	for _, intent := range ourIntents[n:] {
		if err := cs.dirtyState.Update([]input.Intent{intent}); err != nil {
			return err
		}
	}

	return nil
}

type Game struct {
	rootGeoM ebiten.GeoM
	dc       *ctxwebrtc.DataChannel

	cs   *clientState
	csMu sync.Mutex

	delayRingbuf   *ringbuf.RingBuf[time.Duration]
	delayRingbufMu sync.RWMutex
}

func New(dc *ctxwebrtc.DataChannel, rng *syncrand.Source, isOfferer bool) *Game {
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("yumbattle")
	const defaultScale = 4
	ebiten.SetWindowSize(renderWidth*defaultScale, renderHeight*defaultScale)

	s := &state.State{Rng: rng}

	g := &Game{
		dc: dc,
		cs: &clientState{
			isOfferer:       isOfferer,
			committedState:  s,
			dirtyState:      s.Clone(),
			incomingIntents: ringbuf.New[input.Intent](maxPendingIntents),
			outgoingIntents: ringbuf.New[input.Intent](maxPendingIntents),
		},
		delayRingbuf: ringbuf.New[time.Duration](10),
	}
	return g
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

				nextTick := g.cs.committedState.ElapsedTicks + uint32(g.cs.incomingIntents.Used()) + 1
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
	scaleFactor := outsideWidth / renderWidth
	if s := outsideHeight / renderHeight; s < scaleFactor {
		scaleFactor = s
	}

	insideWidth := renderWidth * scaleFactor
	insideHeight := renderHeight * scaleFactor

	g.rootGeoM = ebiten.GeoM{}
	g.rootGeoM.Scale(float64(scaleFactor), float64(scaleFactor))
	g.rootGeoM.Translate(float64(outsideWidth-insideWidth)/2, float64(outsideHeight-insideHeight)/2)

	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
}

func (g *Game) Update() error {
	g.csMu.Lock()
	defer g.csMu.Unlock()

	if g.cs.outgoingIntents.Free() == 0 {
		return nil
	}

	intent := input.CurrentIntent()
	forTick := g.cs.dirtyState.ElapsedTicks + 1

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
