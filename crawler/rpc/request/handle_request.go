// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package reqresp

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

const requestBufferSize = 2048

// RequestPayloadHandler processes a request (decompressed if previously compressed), read from r.
// The handler can respond by writing to w. After returning the writer will automatically be closed.
// If the input is already known to be invalid, e.g. the request size is invalid, then `invalidInputErr != nil`, and r will not read anything more.
type RequestPayloadHandler func(ctx context.Context, peerId peer.ID, requestLen uint64, r io.Reader, w io.Writer, comp Compression, invalidInputErr error)

type StreamCtxFn func() context.Context

// MakeStreamHandler startReqRPC registers a request handler for the given protocol. Compression is optional and may be nil.
func (handle RequestPayloadHandler) MakeStreamHandler(newCtx StreamCtxFn, comp Compression, minRequestContentSize, maxRequestContentSize uint64) network.StreamHandler {
	return func(stream network.Stream) {
		peerID := stream.Conn().RemotePeer()
		ctx, cancel := context.WithCancel(newCtx())
		defer cancel()

		go func() {
			<-ctx.Done()
			// TODO: should this be a stream reset?
			_ = stream.Close() // Close stream after ctx closes.
		}()

		w := io.WriteCloser(stream)
		// If no request data, then do not even read a length from the stream.
		if maxRequestContentSize == 0 {
			handle(ctx, peerID, 0, nil, w, comp, nil)
			return
		}

		var invalidInputErr error

		// TODO: pool this
		blr := NewBufLimitReader(stream, requestBufferSize, 0)
		blr.N = 1 // var ints need to be read byte by byte
		blr.PerRead = true
		reqLen, err := binary.ReadUvarint(blr)
		blr.PerRead = false
		switch {
		case err != nil:
			invalidInputErr = err
		case reqLen < minRequestContentSize:
			// Check against raw content size minimum (without compression applied)
			invalidInputErr = fmt.Errorf("request length %d is unexpectedly small, request size minimum is %d", reqLen, minRequestContentSize)
		case reqLen > maxRequestContentSize:
			// Check against raw content size limit (without compression applied)
			invalidInputErr = fmt.Errorf("request length %d exceeds request size limit %d", reqLen, maxRequestContentSize)
		case comp != nil:
			// Now apply compression adjustment for size limit, and use that as the limit for the buffered-limited-reader.
			s, err := comp.MaxEncodedLen(maxRequestContentSize)
			if err != nil {
				invalidInputErr = err
			} else {
				maxRequestContentSize = s
			}
		}
		switch {
		case invalidInputErr != nil: // If the input is invalid, never read it.
			maxRequestContentSize = 0
		case comp == nil:
			blr.N = int(maxRequestContentSize)
		default:
			v, err := comp.MaxEncodedLen(maxRequestContentSize)
			if err != nil {
				blr.N = int(maxRequestContentSize)
			} else {
				blr.N = int(v)
			}
		}
		r := io.Reader(blr)
		handle(ctx, peerID, reqLen, r, w, comp, invalidInputErr)
	}
}
