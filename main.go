package main

import (
	"context"
	"encoding/hex"
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pion/webrtc/v3"
	"github.com/yumland/clone"
	"github.com/yumland/ctxwebrtc"
	signorclient "github.com/yumland/signor/client"
	"github.com/yumland/yumbattle/game"
	"github.com/yumland/yumbattle/netsyncrand"
)

var (
	connectAddr = flag.String("connect_addr", "http://localhost:12345", "address to connect to")
	answer      = flag.Bool("answer", false, "if true, answers a session instead of offers")
	sessionID   = flag.String("session_id", "test-session", "session to join to")
)

func main() {
	flag.Parse()

	signorClient := signorclient.New(*connectAddr)
	ctx := context.Background()

	peerConn, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Fatalf("failed to create RTC peer connection: %s", err)
	}

	rtcDc, err := peerConn.CreateDataChannel("game", &webrtc.DataChannelInit{
		ID:         clone.P(uint16(1)),
		Negotiated: clone.P(true),
		Ordered:    clone.P(true),
	})
	if err != nil {
		log.Fatalf("failed to create RTC peer connection: %s", err)
	}

	dc := ctxwebrtc.WrapDataChannel(rtcDc)

	isOfferer := !*answer
	if isOfferer {
		if err := signorClient.Offer(ctx, []byte(*sessionID), peerConn); err != nil {
			log.Fatalf("failed to offer: %s", err)
		}
	} else {
		if err := signorClient.Answer(ctx, []byte(*sessionID), peerConn); err != nil {
			log.Fatalf("failed to offer: %s", err)
		}
	}

	randSource, seed, err := netsyncrand.Negotiate(ctx, dc)
	if err != nil {
		log.Fatalf("failed to negotiate randSource: %s", err)
	}

	log.Printf("negotiated rng, seed: %s", hex.EncodeToString(seed))

	g := game.New(dc, randSource, isOfferer)
	go func() {
		if err := g.RunBackgroundTasks(ctx); err != nil {
			log.Fatalf("error running background tasks: %s", err)
		}
	}()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("failed to run game: %s", err)
	}
}
