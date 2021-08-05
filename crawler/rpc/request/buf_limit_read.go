package reqresp

import (
	"errors"
	"fmt"
	"io"
)

// Reader implements buffering for an io.Reader object.
type BufLimitReader struct {
	buf     []byte
	rd      io.Reader // reader provided by the client
	r, w    int       // buf read and write positions
	N       int       // max bytes remaining
	PerRead bool      // Limit applies per read, i.e. it is not affected at the end of the read.
}

// NewBufLimitReader returns a new Reader whose buffer has the specified size.
// The reader will return an error if Read crosses the limit.
func NewBufLimitReader(rd io.Reader, size int, limit int) *BufLimitReader {
	r := &BufLimitReader{
		buf: make([]byte, size, size),
		rd:  rd,
		N:   limit,
	}
	return r
}

var errNegativeRead = errors.New("reader returned negative count from Read")

// Read reads data into p.
// It returns the number of bytes read into p.
// The bytes are taken from at most one Read on the underlying Reader,
// hence N may be less than len(p).
// At EOF, the count will be zero and err will be io.EOF.
func (b *BufLimitReader) Read(p []byte) (n int, err error) {
	defer func() {
		if !b.PerRead {
			b.N -= n
		}
	}()
	if b.N <= 0 {
		return 0, io.EOF
	}
	if len(p) > b.N {
		if b.N == 0 {
			return 0, fmt.Errorf("reader BufLimitReader tried to read %d bytes, but limit was reached", len(p))
		}
		p = p[:b.N]
	}
	n = len(p)
	if n == 0 {
		return 0, nil
	}
	// if all buffered bytes have been written
	if b.r == b.w {
		if len(p) >= len(b.buf) {
			// Large read, empty buffer.
			// Read directly into p to avoid copy.
			n, err = b.rd.Read(p)
			if n < 0 {
				panic(errNegativeRead)
			}
			return n, err
		}
		b.r = 0
		b.w = 0
		to := b.N // read no more than allowed.
		if to > len(b.buf) {
			to = len(b.buf)
		}
		n, err = b.rd.Read(b.buf[:to])
		if n < 0 {
			panic(errNegativeRead)
		}
		b.w += n
		if err != nil {
			return n, err
		}
	}

	// copy as much as we can
	n = copy(p, b.buf[b.r:b.w])
	b.r += n
	return n, nil
}

func (b *BufLimitReader) ReadByte() (byte, error) {
	out := [1]byte{}
	n, err := b.Read(out[:])
	if n == 0 && err == nil {
		return 0, errors.New("failed to read single byte, but no error")
	}
	return out[0], err
}
