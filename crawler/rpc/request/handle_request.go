package reqresp

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

const requestBufferSize = 2048

// RequestPayloadHandler processes a request (decompressed if previously compressed), read from r.
// The handler can respond by writing to w. After returning the writer will automatically be closed.
// If the input is already known to be invalid, e.g. the request size is invalid, then `invalidInputErr != nil`, and r will not read anything more.
type RequestPayloadHandler func(ctx context.Context, peerId peer.ID, requestLen uint64, r io.Reader, w io.Writer, comp Compression, invalidInputErr error)

type StreamCtxFn func() context.Context

// startReqRPC registers a request handler for the given protocol. Compression is optional and may be nil.
func (handle RequestPayloadHandler) MakeStreamHandler(newCtx StreamCtxFn, comp Compression, minRequestContentSize, maxRequestContentSize uint64) network.StreamHandler {
	return func(stream network.Stream) {
		peerId := stream.Conn().RemotePeer()
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
			handle(ctx, peerId, 0, nil, w, comp, nil)
			return
		}

		var invalidInputErr error

		// TODO: pool this
		blr := NewBufLimitReader(stream, requestBufferSize, 0)
		blr.N = 1 // var ints need to be read byte by byte
		blr.PerRead = true
		reqLen, err := binary.ReadUvarint(blr)
		blr.PerRead = false
		if err != nil {
			invalidInputErr = err
		} else if reqLen < minRequestContentSize {
			// Check against raw content size minimum (without compression applied)
			invalidInputErr = fmt.Errorf("request length %d is unexpectedly small, request size minimum is %d", reqLen, minRequestContentSize)
		} else if reqLen > maxRequestContentSize {
			// Check against raw content size limit (without compression applied)
			invalidInputErr = fmt.Errorf("request length %d exceeds request size limit %d", reqLen, maxRequestContentSize)
		} else if comp != nil {
			// Now apply compression adjustment for size limit, and use that as the limit for the buffered-limited-reader.
			s, err := comp.MaxEncodedLen(maxRequestContentSize)
			if err != nil {
				invalidInputErr = err
			} else {
				maxRequestContentSize = s
			}
		}
		// If the input is invalid, never read it.
		if invalidInputErr != nil {
			maxRequestContentSize = 0
		}
		if comp == nil {
			blr.N = int(maxRequestContentSize)
		} else {
			v, err := comp.MaxEncodedLen(maxRequestContentSize)
			if err != nil {
				blr.N = int(maxRequestContentSize)
			} else {
				blr.N = int(v)
			}
		}
		r := io.Reader(blr)
		handle(ctx, peerId, reqLen, r, w, comp, invalidInputErr)
	}
}
