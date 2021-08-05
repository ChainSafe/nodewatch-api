package reqresp

import (
	"fmt"
	"io"

	"github.com/golang/snappy"
)

type Compression interface {
	// Wraps a reader to decompress data as reads happen.
	Decompress(r io.Reader) io.Reader
	// Wraps a writer to compress data as writes happen.
	Compress(w io.WriteCloser) io.WriteCloser
	// Returns an error when the input size is too large to encode.
	MaxEncodedLen(msgLen uint64) (uint64, error)
	// The name of the compression that is suffixed to the actual encoding. E.g. "snappy", w.r.t. "ssz_snappy".
	Name() string
}

type SnappyCompression struct{}

func (c SnappyCompression) Decompress(reader io.Reader) io.Reader {
	return snappy.NewReader(reader)
}

func (c SnappyCompression) Compress(w io.WriteCloser) io.WriteCloser {
	return snappy.NewBufferedWriter(w)
}

func (c SnappyCompression) MaxEncodedLen(msgLen uint64) (uint64, error) {
	if msgLen&(1<<63) != 0 {
		return 0, fmt.Errorf("message length %d is too large to compress with snappy", msgLen)
	}
	m := snappy.MaxEncodedLen(int(msgLen))
	if m < 0 {
		return 0, fmt.Errorf("message length %d is too large to compress with snappy", msgLen)
	}
	return uint64(m), nil
}

func (c SnappyCompression) Name() string {
	return "snappy"
}
