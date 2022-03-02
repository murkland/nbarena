//go:build js

package main

import (
	"github.com/pion/webrtc/v3"
)

func WebRTCAPI() (*webrtc.API, error) {
	return webrtc.NewAPI(), nil
}
