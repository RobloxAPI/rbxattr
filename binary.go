package rbxattr

import (
	"encoding/binary"
	"io"
	"math"
)

// Returns the size of an integer.
func numberDataSize(data interface{}) int {
	switch data.(type) {
	case int8, *int8, uint8, *uint8:
		return 1
	case int16, *int16, uint16, *uint16:
		return 2
	case int32, *int32, uint32, *uint32, float32, *float32:
		return 4
	case int64, *int64, uint64, *uint64, float64, *float64:
		return 8
	}
	return 0
}

// Reader wrapper that keeps track of the number of bytes written.
type binaryReader struct {
	r   io.Reader
	n   int64
	err error
}

func newBinaryReader(r io.Reader) *binaryReader {
	return &binaryReader{r: r}
}

func (br *binaryReader) N() (n int64) {
	return br.n
}

func (br *binaryReader) Err() (err error) {
	return br.err
}

func (br *binaryReader) End() (n int64, err error) {
	return br.n, br.err
}

// Add receives the results of a ReadFrom and adds them to br.
func (br *binaryReader) Add(n int64, err error) (failed bool) {
	if br.err != nil {
		return true
	}

	br.n += n
	br.err = err

	if br.err != nil {
		return true
	}
	return false
}

func (br *binaryReader) Bytes(p []byte) (failed bool) {
	if br.err != nil {
		return true
	}

	var n int
	n, br.err = io.ReadFull(br.r, p)
	br.n += int64(n)

	if br.err != nil {
		return true
	}
	return false
}

func (br *binaryReader) Number(data interface{}) (failed bool) {
	if br.err != nil {
		return true
	}

	if m := numberDataSize(data); m != 0 {
		var b [8]byte
		bs := b[:m]
		if br.Bytes(bs) {
			return true
		}
		switch data := data.(type) {
		case *int8:
			*data = int8(b[0])
		case *uint8:
			*data = b[0]
		case *int16:
			*data = int16(binary.LittleEndian.Uint16(bs))
		case *uint16:
			*data = binary.LittleEndian.Uint16(bs)
		case *int32:
			*data = int32(binary.LittleEndian.Uint32(bs))
		case *uint32:
			*data = binary.LittleEndian.Uint32(bs)
		case *int64:
			*data = int64(binary.LittleEndian.Uint64(bs))
		case *uint64:
			*data = binary.LittleEndian.Uint64(bs)
		case *float32:
			*data = math.Float32frombits(binary.LittleEndian.Uint32(bs))
		case *float64:
			*data = math.Float64frombits(binary.LittleEndian.Uint64(bs))
		default:
			goto invalid
		}
		return false
	}

invalid:
	panic("invalid type")
}

func (br *binaryReader) String(data *string) (failed bool) {
	if br.err != nil {
		return true
	}

	var length uint32
	if br.Number(&length) {
		return true
	}
	s := make([]byte, length)
	if br.Bytes(s) {
		return true
	}
	*data = string(s)

	return false
}

// Writer wrapper that keeps track of the number of bytes written.
type binaryWriter struct {
	w   io.Writer
	n   int64
	err error
}

func newBinaryWriter(w io.Writer) *binaryWriter {
	return &binaryWriter{w: w}
}

func (bw *binaryWriter) N() (n int64) {
	return bw.n
}

func (bw *binaryWriter) Err() (err error) {
	return bw.err
}

func (bw *binaryWriter) End() (n int64, err error) {
	return bw.n, bw.err
}

// Add receives the results of a WriteTo and adds them to bw.
func (bw *binaryWriter) Add(n int64, err error) (failed bool) {
	if bw.err != nil {
		return true
	}

	bw.n += n
	bw.err = err

	if bw.err != nil {
		return true
	}
	return false
}

func (bw *binaryWriter) Bytes(p []byte) (failed bool) {
	if bw.err != nil {
		return true
	}

	var n int
	n, bw.err = bw.w.Write(p)
	bw.n += int64(n)
	if n < len(p) {
		return true
	}

	return false
}

func (bw *binaryWriter) Number(data interface{}) (failed bool) {
	if bw.err != nil {
		return true
	}

	if m := numberDataSize(data); m != 0 {
		b := make([]byte, 8)
		switch data := data.(type) {
		case int8:
			b[0] = uint8(data)
		case uint8:
			b[0] = data
		case int16:
			binary.LittleEndian.PutUint16(b, uint16(data))
		case uint16:
			binary.LittleEndian.PutUint16(b, data)
		case int32:
			binary.LittleEndian.PutUint32(b, uint32(data))
		case uint32:
			binary.LittleEndian.PutUint32(b, data)
		case int64:
			binary.LittleEndian.PutUint64(b, uint64(data))
		case uint64:
			binary.LittleEndian.PutUint64(b, data)
		case float32:
			binary.LittleEndian.PutUint32(b, math.Float32bits(data))
		case float64:
			binary.LittleEndian.PutUint64(b, math.Float64bits(data))
		default:
			goto invalid
		}
		return bw.Bytes(b[:m])
	}

invalid:
	panic("invalid type")
}

func (bw *binaryWriter) String(data string) (failed bool) {
	if bw.err != nil {
		return true
	}

	if bw.Number(uint32(len(data))) {
		return true
	}

	return bw.Bytes([]byte(data))
}
