// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package reqresp

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// ResponseChunkHandler is a function that processes a response chunk. The index, size and result-code are already parsed.
// The contents (decompressed if previously compressed) can be read from r. Optionally an answer can be written back to w.
// If the response chunk could not be processed, an error may be returned.
type ResponseChunkHandler func(ctx context.Context, chunkIndex uint64, chunkSize uint64, result ResponseCode, r io.Reader, w io.Writer) error

// ResponseHandler processes a response by internally processing chunks, any error is propagated up.
type ResponseHandler func(ctx context.Context, r io.Reader, w io.WriteCloser) error

type OnRequested func()

// MakeResponseHandler builds a ResponseHandler, which won't take more than maxChunkCount chunks, or chunk contents larger than maxChunkContentSize.
// Compression is optional and may be nil. Chunks are processed by the given ResponseChunkHandler.
func (handleChunk ResponseChunkHandler) MakeResponseHandler(maxChunkCount uint64, maxChunkContentSize uint64, comp Compression) ResponseHandler {
	//		response  ::= <response_chunk>*
	//		response_chunk  ::= <result> | <encoding-dependent-header> | <encoded-payload>
	//		result    ::= “0” | “1” | “2” | [“128” ... ”255”]
	return func(ctx context.Context, r io.Reader, w io.WriteCloser) error {
		defer func() {
			_ = w.Close()
		}()
		if maxChunkCount == 0 {
			return nil
		}
		blr := NewBufLimitReader(r, 1024, 0)
		for chunkIndex := uint64(0); chunkIndex < maxChunkCount; chunkIndex++ {
			blr.N = 1
			resByte, err := blr.ReadByte()
			if errors.Is(err, io.EOF) { // no more chunks left.
				return nil
			}
			if err != nil {
				return fmt.Errorf("failed to read chunk %d result byte: %w", chunkIndex, err)
			}
			// varints need to be read byte by byte.
			blr.N = 1
			blr.PerRead = true
			chunkSize, err := binary.ReadUvarint(blr)
			blr.PerRead = false
			// TODO when input is incorrect, return a different type of error.
			if err != nil {
				return err
			}
			if resByte == byte(InvalidReqCode) || resByte == byte(ServerErrCode) {
				if chunkSize > MaxErrSize {
					return fmt.Errorf("chunk size %d of chunk %d exceeds error size limit %d", chunkSize, chunkIndex, MaxErrSize)
				}
				blr.N = MaxErrSize
			} else {
				if chunkSize > maxChunkContentSize {
					return fmt.Errorf("chunk size %d of chunk %d exceeds chunk limit %d", chunkSize, chunkIndex, maxChunkContentSize)
				}
				blr.N = int(maxChunkContentSize)
			}
			cr := io.Reader(blr)
			cw := w
			if comp != nil {
				cr = comp.Decompress(cr)
				cw = comp.Compress(cw)
			}
			if err := handleChunk(ctx, chunkIndex, chunkSize, ResponseCode(resByte), cr, cw); err != nil {
				_ = cw.Close()
				return err
			}
			if comp != nil {
				if err := cw.Close(); err != nil {
					return fmt.Errorf("failed to close response writer for chunk")
				}
			}
		}
		return nil
	}
}
