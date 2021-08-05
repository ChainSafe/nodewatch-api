package reqresp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type payloadBuffer bytes.Buffer

func (p *payloadBuffer) Close() error {
	return nil
}

func (p *payloadBuffer) Write(b []byte) (n int, err error) {
	return (*bytes.Buffer)(p).Write(b)
}

func (p *payloadBuffer) OutputSizeVarint(w io.Writer) error {
	size := (*bytes.Buffer)(p).Len()
	sizeBytes := [binary.MaxVarintLen64]byte{}
	sizeByteLen := binary.PutUvarint(sizeBytes[:], uint64(size))
	_, err := w.Write(sizeBytes[:sizeByteLen])
	return err
}

func (p *payloadBuffer) WriteTo(w io.Writer) (n int64, err error) {
	return (*bytes.Buffer)(p).WriteTo(w)
}

type noCloseWriter struct {
	w io.Writer
}

func (nw *noCloseWriter) Write(p []byte) (n int, err error) {
	return nw.w.Write(p)
}

func (nw *noCloseWriter) Close() error {
	return nil
}

// EncodeHeaderAndPayload reads a payload, buffers (and optionally compresses) the payload,
// then computes the header-data (varint of byte size). And then writes header and payload.
func EncodeHeaderAndPayload(r io.Reader, w io.Writer, comp Compression) error {
	var buf payloadBuffer
	if _, err := io.Copy(&buf, r); err != nil {
		return err
	}
	if err := buf.OutputSizeVarint(w); err != nil {
		return err
	}
	if comp != nil {
		compressedWriter := comp.Compress(&noCloseWriter{w: w})
		defer compressedWriter.Close()
		if _, err := buf.WriteTo(compressedWriter); err != nil {
			return err
		}
	} else {
		if _, err := buf.WriteTo(w); err != nil {
			return err
		}
	}
	return nil
}

// StreamHeaderAndPayload reads a payload and streams (and optionally compresses) it to the writer.
// To do so, it requires the (uncompressed) payload length to be known in advance.
func StreamHeaderAndPayload(size uint64, r io.Reader, w io.Writer, comp Compression) error {
	sizeBytes := [binary.MaxVarintLen64]byte{}
	sizeByteLen := binary.PutUvarint(sizeBytes[:], size)
	n, err := w.Write(sizeBytes[:sizeByteLen])
	if err != nil {
		return fmt.Errorf("failed to write size bytes: %v", err)
	}
	if n != sizeByteLen {
		return fmt.Errorf("failed to write size bytes fully: %d/%d", n, sizeByteLen)
	}
	if comp != nil {
		compressedWriter := comp.Compress(&noCloseWriter{w: w})
		defer compressedWriter.Close()
		if _, err := io.Copy(compressedWriter, r); err != nil {
			return fmt.Errorf("failed to write payload through compressed writer: %v", err)
		}
		return nil
	} else {
		if _, err := io.Copy(w, r); err != nil {
			return fmt.Errorf("failed to write payload: %v", err)
		}
		return nil
	}
}

// EncodeResult writes the result code to the output writer.
func EncodeResult(result ResponseCode, w io.Writer) error {
	_, err := w.Write([]byte{uint8(result)})
	return err
}

// EncodeChunk reads (decompressed) response message from the msg io.Reader,
// and writes it as a chunk with given result code to the output writer. The compression is optional and may be nil.
func EncodeChunk(result ResponseCode, r io.Reader, w io.Writer, comp Compression) error {
	if err := EncodeResult(result, w); err != nil {
		return err
	}
	return EncodeHeaderAndPayload(r, w, comp)
}

// EncodeChunk reads (decompressed) response message from the msg io.Reader,
// and writes it as a chunk with given result code to the output writer. The compression is optional and may be nil.
func StreamChunk(result ResponseCode, size uint64, r io.Reader, w io.Writer, comp Compression) error {
	if err := EncodeResult(result, w); err != nil {
		return err
	}
	return StreamHeaderAndPayload(size, r, w, comp)
}
