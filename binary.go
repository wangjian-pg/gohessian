package hessian

import (
	"bytes"
	"errors"
	"io"
)

const (
	BC_BIN_FIN   = byte('B') // final chunk
	BC_BIN_CHUNK = byte('A') // non-final chunk

	BC_COMPACT_BIN     = byte(0x20)
	BC_COMPACT_BIN_MAX = byte(0x2f)

	BC_SHORT_BIN     = byte(0x34)
	BC_SHORT_BIN_MAX = byte(0x37)

	maxShortBinSize = 0x3ff
	maxBinChunkSize = 0xffff
)

func ReadBytes(reader io.Reader) ([]byte, error) {
	fin := false
	buffer := new(bytes.Buffer)
	bs := make([]byte, 2)
	var err error

	for !fin {
		_, err = io.ReadFull(reader, bs[:1])
		if err != nil {
			return nil, err
		}
		flag := bs[0]
		len := 0
		switch {
		case flag == BC_NIL:
			return nil, nil
		case flag >= BC_COMPACT_BIN && flag <= BC_COMPACT_BIN_MAX:
			len = int(flag - BC_COMPACT_BIN)
			fin = true
		case flag >= BC_SHORT_BIN && flag <= BC_SHORT_BIN_MAX:
			if _, err := io.ReadFull(reader, bs[:1]); err != nil {
				return nil, err
			} else {
				len = int(flag-BC_SHORT_BIN)<<8 + int(bs[0])
			}
			fin = true
		case flag == BC_BIN_FIN:
			fin = true
			fallthrough
		case flag == BC_BIN_CHUNK:
			if _, err := io.ReadFull(reader, bs); err != nil {
				return nil, err
			} else {
				len = int(bs[0])<<8 + int(bs[1])
			}
		default:
			return nil, errors.New("Not encoded as hessian binary")
		}
		chunk := make([]byte, len)
		if _, err := io.ReadFull(reader, chunk); err != nil {
			return nil, err
		} else {
			buffer.Write(chunk)
		}
	}

	return buffer.Bytes(), nil
}

func WriteBytes(w io.Writer, buf []byte) (int, error) {
	if buf == nil {
		return w.Write([]byte{BC_NIL})
	}

	left := len(buf)
	if left <= int(BC_COMPACT_BIN_MAX-BC_COMPACT_BIN) {
		p := make([]byte, left+1)
		flag := byte(left) + BC_COMPACT_BIN
		p[0] = flag
		copy(p[1:], buf)
		return w.Write(p)
	} else if left <= maxShortBinSize {
		p := make([]byte, 2)
		p[0] = byte((left >> 8)) + BC_SHORT_BIN
		p[1] = byte(left)
		if _, err := w.Write(p); err != nil {
			return -1, err
		}
		if _, err := w.Write(buf); err != nil {
			return -1, err
		} else {
			return left + 2, nil
		}
	}

	var flag byte
	var chunkSize, n int
	for left > 0 {
		start := len(buf) - left
		if left > maxBinChunkSize {
			chunkSize = maxBinChunkSize
			flag = BC_BIN_CHUNK
		} else {
			chunkSize = left
			flag = BC_BIN_FIN
		}

		lenBitLow := byte(chunkSize)
		lenBitHigh := byte(chunkSize >> 8)

		if _, err := w.Write([]byte{flag, lenBitHigh, lenBitLow}); err != nil {
			return -1, err
		} else {
			n += 3
		}
		end := start + chunkSize
		if v, err := w.Write(buf[start:end]); err != nil {
			return -1, err
		} else {
			n += v
		}
		left -= chunkSize
	}
	return n, nil
}
