package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/murkland/clone"
	"github.com/murkland/ctxwebrtc"
	"github.com/murkland/moreflag"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/game"
	"github.com/murkland/nbarena/netsyncrand"
	signorclient "github.com/murkland/signor/client"
	"github.com/pion/webrtc/v3"
)

var defaultWebRTCConfig = (func() string {
	s, err := json.Marshal(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
			{
				URLs: []string{"stun:stun1.l.google.com:19302"},
			},
			{
				URLs: []string{"stun:stun2.l.google.com:19302"},
			},
			{
				URLs: []string{"stun:stun3.l.google.com:19302"},
			},
			{
				URLs: []string{"stun:stun4.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	return string(s)
})()

var (
	connectAddr      = flag.String("connect_addr", "http://localhost:12345", "address to connect to")
	answer           = flag.Bool("answer", false, "if true, answers a session instead of offers")
	sessionID        = flag.String("session_id", "test-session", "session to join to")
	webrtcConfig     = flag.String("webrtc_config", defaultWebRTCConfig, "webrtc configuration")
	delaysWindowSize = flag.Int("delays_window_size", 5, "size of window for calculating delay")
	inputFrameDelay  = flag.Int("input_frame_delay", 0, "additional input frame delay, if any")
)

func main() {
	moreflag.Parse()
	ctx := context.Background()

	b, err := bundle.Load(ctx, loaderCallback)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("you are missing assets from BN6. please see README.md in the assets directory for instructions on how to dump assets from a ROM of BN6.")
		}
		log.Fatalf("failed to load bundle: %s", err)
	}

	var peerConnConfig webrtc.Configuration
	if err := json.Unmarshal([]byte(*webrtcConfig), &peerConnConfig); err != nil {
		log.Fatalf("failed to parse webrtc config: %s", err)
	}

	log.Printf("connecting to %s, answer = %t, session_id = %s (using peer config: %+v)", *connectAddr, *answer, *sessionID, peerConnConfig)

	signorClient := signorclient.New(*connectAddr)

	api, err := webRTCAPI()
	if err != nil {
		log.Fatalf("failed to get WebRTC API: %s", err)
	}

	peerConn, err := api.NewPeerConnection(peerConnConfig)
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

	log.Printf("signaling complete!")
	log.Printf("local SDP: %s", peerConn.LocalDescription().SDP)
	log.Printf("remote SDP: %s", peerConn.RemoteDescription().SDP)

	randSource, seed, err := netsyncrand.Negotiate(ctx, dc)
	if err != nil {
		log.Fatalf("failed to negotiate randSource: %s", err)
	}

	log.Printf("negotiated rng, seed: %s", hex.EncodeToString(seed))

	g := game.New(b, dc, randSource, isAnswerer, *delaysWindowSize, *inputFrameDelay)
	go func() {
		if err := g.RunBackgroundTasks(ctx); err != nil {
			log.Fatalf("error running background tasks: %s", err)
		}
	}()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("failed to run game: %s", err)
	}
}
