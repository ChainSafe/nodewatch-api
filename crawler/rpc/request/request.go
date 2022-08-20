// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package reqresp

import (
	"bytes"
	"context"
	"io"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type NewStreamFn func(ctx context.Context, peerId peer.ID, protocolId ...protocol.ID) (network.Stream, error)

func (newStreamFn NewStreamFn) Request(ctx context.Context, peerID peer.ID, protocolID protocol.ID, r io.Reader, comp Compression, handle ResponseHandler) error {
	stream, err := newStreamFn(ctx, peerID, protocolID)
	if err != nil {
		return err
	}
	defer func() {
		_ = stream.Close()
	}()

	var buf bytes.Buffer
	if err := EncodeHeaderAndPayload(r, &buf, comp); err != nil {
		return err
	}
	if _, err := stream.Write(buf.Bytes()); err != nil {
		return err
	}
	return handle(ctx, stream, stream)
}
