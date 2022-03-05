package netsyncrand

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/yumland/ctxwebrtc"
	"github.com/yumland/nbarena/packets"
	"github.com/yumland/syncrand"
)

func Negotiate(ctx context.Context, dc *ctxwebrtc.DataChannel) (*syncrand.Source, []byte, error) {
	var nonce [16]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, nil, fmt.Errorf("failed to generate rng seed part: %w", err)
	}

	commitment := syncrand.Commit(nonce[:])
	var commitPacket packets.Commit
	copy(commitPacket.Commitment[:], commitment)
	if err := packets.Send(ctx, dc, commitPacket); err != nil {
		return nil, nil, fmt.Errorf("failed to send commit: %w", err)
	}

	theirCommitReply, err := packets.Recv(ctx, dc)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to receive commit: %w", err)
	}
	theirCommitment := theirCommitReply.(packets.Commit).Commitment

	if err := packets.Send(ctx, dc, packets.Reveal{Nonce: nonce}); err != nil {
		return nil, nil, fmt.Errorf("failed to send reveal: %w", err)
	}

	theirRevealReply, err := packets.Recv(ctx, dc)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to receive reveal: %w", err)
	}
	theirNonce := theirRevealReply.(packets.Reveal).Nonce

	if !syncrand.Verify(commitment, theirCommitment[:], theirNonce[:]) {
		return nil, nil, errors.New("failed to verify rng commitment")
	}

	seed := syncrand.MakeSeed(nonce[:], theirNonce[:])
	rng := syncrand.NewSource(seed)

	return rng, seed, nil
}
