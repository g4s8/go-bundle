package bundle

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
)

// Out bundle of binary data to write
type Out struct {
	bo  binary.ByteOrder
	buf *bytes.Buffer
	err error
}

var bytesEmpty = []byte{}

// NewOut creates empty output to write data
func NewOut(order binary.ByteOrder) *Out {
	return &Out{bo: order, buf: bytes.NewBuffer(bytesEmpty)}
}

// NewBEOut creates new output with big-endian byte order
func NewBEOut() *Out {
	return NewOut(binary.BigEndian)
}

// NewLEOut creates new output with little-endian byte order
func NewLEOut() *Out {
	return NewOut(binary.LittleEndian)
}

func (o *Out) Error() error {
	return o.err
}

// PutBool buts a bool to bundle
func (o *Out) PutBool(x bool) {
	o.writeUnsafe(x)
}

// PutByte puts a byte to bundle
func (o *Out) PutByte(x byte) {
	o.PutUInt8(x)
}

// PutUInt8 puts uint8 to bundle
func (o *Out) PutUInt8(x uint8) {
	o.writeUnsafe(x)
}

// PutInt8 puts int8 to bundle
func (o *Out) PutInt8(x int8) {
	o.writeUnsafe(x)
}

// PutUInt16 puts uint16 to bundle
func (o *Out) PutUInt16(x uint16) {
	o.writeUnsafe(x)
}

// PutInt16 puts int16 to bundle
func (o *Out) PutInt16(x int16) {
	o.writeUnsafe(x)
}

// PutUInt32 puts uint32 to bundle
func (o *Out) PutUInt32(x uint32) {
	o.writeUnsafe(x)
}

// PutInt32 puts int32 to bundle
func (o *Out) PutInt32(x int32) {
	o.writeUnsafe(x)
}

// PutRune puts rune to bundle
func (o *Out) PutRune(x rune) {
	o.writeUnsafe(int32(x))
}

// PutInt64 puts int64 to bundle
func (o *Out) PutInt64(x int64) {
	o.writeUnsafe(x)
}

// PutUInt64 puts uint64 to bundle
func (o *Out) PutUInt64(x uint64) {
	o.writeUnsafe(x)
}

// PutFloat32 puts float32 to bundle
func (o *Out) PutFloat32(x float32) {
	o.writeUnsafe(x)
}

// PutFloat64 puts float64 to bundle
func (o *Out) PutFloat64(x float64) {
	o.writeUnsafe(x)
}

// PutBytes puts byte array to bundle
func (o *Out) PutBytes(x []byte) {
	o.PutUInts8(x)
}

// PutUInts8 puts uint8 array to bundle
func (o *Out) PutUInts8(x []uint8) {
	o.PutUInt32(uint32(len(x)))
	o.writeUnsafe(x)
}

// PutString puts string to bundle
func (o *Out) PutString(s string) {
	o.PutUInts8([]uint8(s))
}

// PutBinary puts binary marshaler to bundle
func (o *Out) PutBinary(bin encoding.BinaryMarshaler) {
	bytes, err := bin.MarshalBinary()
	if err != nil {
		o.err = err
	} else {
		o.PutBytes(bytes)
	}
}

var errUnsupportedByteOred = errors.New("unsupported byte order")

const (
	encBoBE = 1
	encBoLE = 2
)

// PutBundle puts nested bundle
func (o *Out) PutBundle(b *Out) {
	if o.err != nil {
		return
	}
	if b.bo == binary.BigEndian {
		o.PutUInt8(encBoBE)
	} else if b.bo == binary.LittleEndian {
		o.PutUInt8(encBoLE)
	} else {
		o.err = errUnsupportedByteOred
		return
	}
	o.PutBinary(b)
}

// Flip bundle onput to input
func (o *Out) Flip() *Inp {
	i := new(Inp)
	if o.err != nil {
		i.err = o.err
	} else {
		i.src = bytes.NewReader(o.buf.Bytes())
		i.bo = o.bo
	}
	return i
}

func (o *Out) MarshalBinary() ([]byte, error) {
	if o.err != nil {
		return nil, o.err
	}
	return o.buf.Bytes(), nil
}

func (o *Out) writeUnsafe(x interface{}) {
	if o.err != nil {
		return
	}
	o.err = binary.Write(o.buf, o.bo, x)
}
