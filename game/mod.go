package game

import (
	"context"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yumland/ctxwebrtc"
	"github.com/yumland/ringbuf"
	"github.com/yumland/yumbattle/packets"
	"golang.org/x/sync/errgroup"
)

const renderWidth = 240
const renderHeight = 160

type Game struct {
	rootGeoM ebiten.GeoM
	dc       *ctxwebrtc.DataChannel

	delayRingbuf   *ringbuf.RingBuf[time.Duration]
	delayRingbufMu sync.RWMutex
}

func New(dc *ctxwebrtc.DataChannel) *Game {
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("yumbattle")
	const defaultScale = 4
	ebiten.SetWindowSize(renderWidth*defaultScale, renderHeight*defaultScale)
	g := &Game{
		dc:           dc,
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
