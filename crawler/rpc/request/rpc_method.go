package reqresp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/protolambda/ztyp/codec"
)

type Request interface {
	fmt.Stringer
}

type Codec interface {
	MinByteLen() uint64
	MaxByteLen() uint64
	Encode(w io.Writer, input codec.Serializable) error
	Decode(r io.Reader, bytesLen uint64, dest codec.Deserializable) error
	Alloc() SerDes
}

type SerDes interface {
	codec.Serializable
	codec.Deserializable
}

type SSZCodec struct {
	alloc      func() SerDes
	minByteLen uint64
	maxByteLen uint64
}

func NewSSZCodec(alloc func() SerDes, minByteLen uint64, maxByteLen uint64) *SSZCodec {
	return &SSZCodec{
		alloc:      alloc,
		minByteLen: minByteLen,
		maxByteLen: maxByteLen,
	}
}

func (c *SSZCodec) MinByteLen() uint64 {
	if c == nil {
		return 0
	}
	return c.minByteLen
}

func (c *SSZCodec) MaxByteLen() uint64 {
	if c == nil {
		return 0
	}
	return c.maxByteLen
}

func (c *SSZCodec) Encode(w io.Writer, input codec.Serializable) error {
	if c == nil {
		if input != nil {
			return errors.New("expected empty data, nil input. This codec is no-op")
		}
		return nil
	}
	return input.Serialize(codec.NewEncodingWriter(w))
}

func (c *SSZCodec) Decode(r io.Reader, bytesLen uint64, dest codec.Deserializable) error {
	if c == nil {
		if bytesLen != 0 {
			return errors.New("expected 0 bytes, no definition")
		}
		return nil
	}
	return dest.Deserialize(codec.NewDecodingReader(r, bytesLen))
}

func (c *SSZCodec) Alloc() SerDes {
	if c == nil {
		return nil
	}
	return c.alloc()
}

type RPCMethod struct {
	Protocol                  protocol.ID
	RequestCodec              Codec
	ResponseChunkCodec        Codec
	DefaultResponseChunkCount uint64
}

type ResponseCode uint8

const (
	SuccessCode    ResponseCode = 0
	InvalidReqCode              = 1
	ServerErrCode               = 2
)

// 256 bytes max error size
const MAX_ERR_SIZE = 256

type OnResponseListener func(chunk ChunkedResponseHandler) error

type RequestInput interface {
	Reader(c Codec) (io.Reader, error)
}

type RequestSSZInput struct {
	Obj codec.Serializable
}

func (v RequestSSZInput) Reader(c Codec) (io.Reader, error) {
	var buf bytes.Buffer
	if err := c.Encode(&buf, v.Obj); err != nil {
		return nil, err
	}
	return &buf, nil
}

type RequestBytesInput []byte

func (v RequestBytesInput) Reader(_ Codec) (io.Reader, error) {
	return bytes.NewReader(v), nil
}

type ChunkedResponseHandler interface {
	ChunkSize() uint64
	ChunkIndex() uint64
	ResultCode() ResponseCode
	ReadRaw() ([]byte, error)
	ReadErrMsg() (string, error)
	ReadObj(dest codec.Deserializable) error
}

type chRespHandler struct {
	m          *RPCMethod
	r          io.Reader
	result     ResponseCode
	chunkSize  uint64
	chunkIndex uint64
}

func (c *chRespHandler) ChunkSize() uint64 {
	return c.chunkSize
}

func (c *chRespHandler) ChunkIndex() uint64 {
	return c.chunkIndex
}

func (c *chRespHandler) ResultCode() ResponseCode {
	return c.result
}

func (c *chRespHandler) ReadRaw() ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(io.LimitReader(c.r, int64(c.chunkSize)))
	return buf.Bytes(), err
}

func (c *chRespHandler) ReadErrMsg() (string, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(io.LimitReader(c.r, int64(c.chunkSize)))
	return string(buf.Bytes()), err
}

func (c *chRespHandler) ReadObj(dest codec.Deserializable) error {
	return c.m.ResponseChunkCodec.Decode(c.r, c.chunkSize, dest)
}

func (m *RPCMethod) RunRequest(ctx context.Context, newStreamFn NewStreamFn,
	peerId peer.ID, comp Compression, req RequestInput, maxRespChunks uint64, madeRequest func() error,
	onResponse OnResponseListener) error {

	handleChunks := ResponseChunkHandler(func(ctx context.Context, chunkIndex uint64, chunkSize uint64, result ResponseCode, r io.Reader, w io.Writer) error {
		return onResponse(&chRespHandler{
			m:          m,
			r:          r,
			result:     result,
			chunkSize:  chunkSize,
			chunkIndex: chunkIndex,
		})
	})

	reqR, err := req.Reader(m.RequestCodec)
	if err != nil {
		return err
	}

	protocolId := m.Protocol
	maxChunkContentSize := m.ResponseChunkCodec.MaxByteLen()
	if comp != nil {
		protocolId += protocol.ID("_" + comp.Name())
		if s, err := comp.MaxEncodedLen(maxChunkContentSize); err != nil {
			return err
		} else {
			maxChunkContentSize = s
		}
	}

	respHandler := handleChunks.MakeResponseHandler(maxRespChunks, maxChunkContentSize, comp)

	handler := ResponseHandler(func(ctx context.Context, r io.Reader, w io.WriteCloser) error {
		if err := madeRequest(); err != nil {
			return fmt.Errorf("made request, but could not continue: %v", err)
		}
		return respHandler(ctx, r, w)
	})

	// Runs the request in sync, which processes responses,
	// and then finally closes the channel through the earlier deferred close.
	return newStreamFn.Request(ctx, peerId, protocolId, reqR, comp, handler)
}

type ReadRequestFn func(dest interface{}) error
type WriteSuccessChunkFn func(data interface{}) error
type WriteMsgFn func(msg string) error

type RequestReader interface {
	// nil if not an invalid input
	InvalidInput() error
	ReadRequest(dest codec.Deserializable) error
	RawRequest() ([]byte, error)
}

type RequestResponder interface {
	WriteResponseChunk(code ResponseCode, data codec.Serializable) error
	WriteRawResponseChunk(code ResponseCode, chunk []byte) error
	StreamResponseChunk(code ResponseCode, size uint64, r io.Reader) error
	WriteErrorChunk(code ResponseCode, msg string) error
}

type ChunkedRequestHandler interface {
	RequestReader
	RequestResponder
}

type chReqHandler struct {
	m               *RPCMethod
	comp            Compression
	respBuf         bytes.Buffer
	reqLen          uint64
	r               io.Reader
	w               io.Writer
	invalidInputErr error
}

func (h *chReqHandler) InvalidInput() error {
	return h.invalidInputErr
}

func (h *chReqHandler) ReadRequest(dest codec.Deserializable) error {
	if h.invalidInputErr != nil {
		return h.invalidInputErr
	}
	r := h.r
	if h.comp != nil {
		r = h.comp.Decompress(r)
	}
	return h.m.RequestCodec.Decode(r, h.reqLen, dest)
}

func (h *chReqHandler) RawRequest() ([]byte, error) {
	if h.invalidInputErr != nil {
		return nil, h.invalidInputErr
	}
	var buf bytes.Buffer
	r := h.r
	if h.comp != nil {
		r = h.comp.Decompress(r)
	}
	if _, err := buf.ReadFrom(io.LimitReader(r, int64(h.reqLen))); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *chReqHandler) WriteResponseChunk(code ResponseCode, data codec.Serializable) error {
	h.respBuf.Reset() // re-use buffer for each response chunk
	if err := h.m.ResponseChunkCodec.Encode(&h.respBuf, data); err != nil {
		return err
	}
	b := h.respBuf.Bytes()
	return StreamChunk(code, uint64(len(b)), bytes.NewReader(b), h.w, h.comp)
}

func (h *chReqHandler) WriteRawResponseChunk(code ResponseCode, chunk []byte) error {
	return StreamChunk(code, uint64(len(chunk)), bytes.NewReader(chunk), h.w, h.comp)
}

func (h *chReqHandler) StreamResponseChunk(code ResponseCode, size uint64, r io.Reader) error {
	return StreamChunk(code, size, r, h.w, h.comp)
}

func (h *chReqHandler) WriteErrorChunk(code ResponseCode, msg string) error {
	if len(msg) > MAX_ERR_SIZE {
		msg = msg[:MAX_ERR_SIZE-3]
		msg += "..."
	}
	b := []byte(msg)
	return StreamChunk(code, uint64(len(b)), bytes.NewReader(b), h.w, h.comp)
}

type OnRequestListener func(ctx context.Context, peerId peer.ID, handler ChunkedRequestHandler)

func (m *RPCMethod) MakeStreamHandler(newCtx StreamCtxFn, comp Compression, listener OnRequestListener) network.StreamHandler {
	return RequestPayloadHandler(func(ctx context.Context, peerId peer.ID, requestLen uint64, r io.Reader, w io.Writer, comp Compression, invalidInputErr error) {
		listener(ctx, peerId, &chReqHandler{
			m: m, comp: comp, reqLen: requestLen, r: r, w: w, invalidInputErr: invalidInputErr,
		})
	}).MakeStreamHandler(newCtx, comp, m.RequestCodec.MinByteLen(), m.RequestCodec.MaxByteLen())
}
