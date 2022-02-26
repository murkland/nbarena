package main

import (
	"context"
	"encoding/hex"
	"flag"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pion/webrtc/v3"
	"github.com/yumland/clone"
	"github.com/yumland/ctxwebrtc"
	"github.com/yumland/moreflag"
	signorclient "github.com/yumland/signor/client"
	"github.com/yumland/yumbattle/bundle"
	"github.com/yumland/yumbattle/game"
	"github.com/yumland/yumbattle/netsyncrand"
)

var (
	connectAddr      = flag.String("connect_addr", "http://localhost:12345", "address to connect to")
	answer           = flag.Bool("answer", false, "if true, answers a session instead of offers")
	sessionID        = flag.String("session_id", "test-session", "session to join to")
	stunServers      = flag.String("stun_servers", "stun:stun.l.google.com:19302,stun:stun1.l.google.com:19302,stun:stun2.l.google.com:19302,stun:stun3.l.google.com:19302,stun:stun4.l.google.com:19302", "stun servers")
	delaysWindowSize = flag.Int("delays_window_size", 5, "size of window for calculating delay")
)

func main() {
	moreflag.Parse()

	var iceServers []webrtc.ICEServer
	for _, url := range strings.Split(*stunServers, ",") {
		iceServers = append(iceServers, webrtc.ICEServer{URLs: []string{url}})
	}

	log.Printf("connecting to %s, answer = %t, session_id = %s (using ICE servers: %+v)", *connectAddr, *answer, *sessionID, iceServers)

	signorClient := signorclient.New(*connectAddr)
	ctx := context.Background()

	b, err := bundle.Load(ctx)
	if err != nil {
		log.Fatalf("failed to load bundle: %s", err)
	}

	peerConn, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: iceServers,
	})
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

	isAnswerer := *answer
	if !isAnswerer {
		if err := signorClient.Offer(ctx, []byte(*sessionID), peerConn); err != nil {
			log.Fatalf("failed to offer: %s", err)
		}
	} else {
		if err := signorClient.Answer(ctx, []byte(*sessionID), peerConn); err != nil {
			log.Fatalf("failed to answer: %s", err)
		}
	}

	log.Printf("connected!")

	randSource, seed, err := netsyncrand.Negotiate(ctx, dc)
	if err != nil {
		log.Fatalf("failed to negotiate randSource: %s", err)
	}

	log.Printf("negotiated rng, seed: %s", hex.EncodeToString(seed))

	g := game.New(b, dc, randSource, isAnswerer, *delaysWindowSize)
	go func() {
		if err := g.RunBackgroundTasks(ctx); err != nil {
			log.Fatalf("error running background tasks: %s", err)
		}
	}()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("failed to run game: %s", err)
	}
}
