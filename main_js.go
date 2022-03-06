//go:build js

package main

import (
	"github.com/pion/webrtc/v3"
)

func webRTCAPI() (*webrtc.API, error) {
	return webrtc.NewAPI(), nil
}

func loaderCallback(path string, i int, n int) {
}
