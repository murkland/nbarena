package packets

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"flag"
	"io"
	"log"

	"github.com/yumland/ctxwebrtc"
	"github.com/yumland/yumbattle/input"
)

var (
	debugLogPackets = flag.Bool("debug_log_packets", false, "log all packets (noisy!)")
)

var ErrUnknownPacket = errors.New("unknown packet")

type packetType uint8

const (
	packetTypePing   packetType = 0
	packetTypePong   packetType = 1
	packetTypeCommit packetType = 2
	packetTypeReveal packetType = 3
	packetTypeIntent packetType = 4
)

type Packet interface {
	packetType() packetType
}

type Ping struct {
	ID uint64
}

func (Ping) packetType() packetType { return packetTypePing }

type Pong struct {
	ID uint64
}

func (Pong) packetType() packetType { return packetTypePong }

type Commit struct {
	Commitment [32]uint8
}

func (Commit) packetType() packetType { return packetTypeCommit }

type Reveal struct {
	Nonce [16]uint8
}

func (Reveal) packetType() packetType { return packetTypeReveal }

type Intent struct {
	ForTick uint32
	Intent  input.Intent
}

func (Intent) packetType() packetType { return packetTypeIntent }

func Marshal(packet Packet) []byte {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, packet.packetType()); err != nil {
		panic(err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, packet); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func unmarshal[T Packet](r io.Reader) (T, error) {
	var packet T
	if err := binary.Read(r, binary.LittleEndian, &packet); err != nil {
		return packet, err
	}
	return packet, nil
}

func Unmarshal(raw []byte) (Packet, error) {
	r := bytes.NewReader(raw)
	var typ packetType
	if err := binary.Read(r, binary.LittleEndian, &typ); err != nil {
		return nil, err
	}

	switch typ {
	case packetTypePing:
		return unmarshal[Ping](r)
	case packetTypePong:
		return unmarshal[Pong](r)
	case packetTypeCommit:
		return unmarshal[Commit](r)
	case packetTypeReveal:
		return unmarshal[Reveal](r)
	case packetTypeIntent:
		return unmarshal[Intent](r)
	default:
		return nil, ErrUnknownPacket
	}
}

func Send(ctx context.Context, dc *ctxwebrtc.DataChannel, packet Packet) error {
	if *debugLogPackets {
		log.Printf("--> %d: %+v", packet.packetType(), packet)
	}
	return dc.Send(ctx, Marshal(packet))
}

func Recv(ctx context.Context, dc *ctxwebrtc.DataChannel) (Packet, error) {
	raw, err := dc.Recv(ctx)
	if err != nil {
		return nil, err
	}
	packet, err := Unmarshal(raw)
	if err != nil {
		return nil, err
	}
	if *debugLogPackets {
		log.Printf("<-- %d: %+v", packet.packetType(), packet)
	}
	return packet, nil
}
