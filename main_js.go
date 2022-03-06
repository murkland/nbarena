//go:build js

package main

import (
	"syscall/js"

	"github.com/pion/webrtc/v3"
)

func webRTCAPI() (*webrtc.API, error) {
	return webrtc.NewAPI(), nil
}

func loaderCallback(path string, i int, n int) {
	global := js.Global()
	cb := global.Get("loaderCallback")
	cb.Invoke(path, i, n)
}
