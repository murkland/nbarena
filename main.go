package main

import (
	"context"
	"encoding/hex"
	"flag"
	"log"

	"github.com/pion/webrtc/v3"
	"github.com/yumland/clone"
	"github.com/yumland/ctxwebrtc"
	signorclient "github.com/yumland/signor/client"
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

	var selfEntityID uint32
	var opponentEntityID uint32

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

	if !*answer {
		if err := signorClient.Offer(ctx, []byte(*sessionID), peerConn); err != nil {
			log.Fatalf("failed to offer: %s", err)
		}

		selfEntityID = 1
		opponentEntityID = 2
	} else {
		if err := signorClient.Answer(ctx, []byte(*sessionID), peerConn); err != nil {
			log.Fatalf("failed to offer: %s", err)
		}

		selfEntityID = 2
		opponentEntityID = 1
	}

	_ = selfEntityID
	_ = opponentEntityID

	rng, seed, err := netsyncrand.Negotiate(ctx, dc)
	if err != nil {
		log.Fatalf("failed to negotiate rng: %s", err)
	}
	_ = rng

	log.Printf("negotiated rng, seed: %s", hex.EncodeToString(seed))
}
