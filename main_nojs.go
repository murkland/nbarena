//go:build !js

package main

import (
	"flag"
	"log"
	"net"

	"github.com/pion/webrtc/v3"
)

var (
	webRTCListenAddr = flag.String("webrtc_listen_addr", "", "address to listen on for WebRTC")
)

func webRTCAPI() (*webrtc.API, error) {
	var opts []func(*webrtc.API)

	if *webRTCListenAddr != "" {
		udpLis, err := net.ListenPacket("udp", *webRTCListenAddr)
		if err != nil {
			return nil, err
		}
		log.Printf("listening for webrtc on %s", udpLis.LocalAddr())

		settingEngine := webrtc.SettingEngine{}
		settingEngine.SetICEUDPMux(webrtc.NewICEUDPMux(nil, udpLis))
		opts = append(opts, webrtc.WithSettingEngine(settingEngine))
	}

	return webrtc.NewAPI(opts...), nil
}

func loaderCallback(path string, i int, n int) {
	log.Printf("loaded %d/%d: %s", i, n, path)
}
