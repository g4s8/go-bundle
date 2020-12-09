package bundle

import (
	"bytes"
	"encoding"
	"encoding/binary"
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

func NewBEOut() *Out {
	return NewOut(binary.BigEndian)
}

func NewLEOut() *Out {
	return NewOut(binary.LittleEndian)
}

func (o *Out) Error() error {
	return o.err
}

func (o *Out) PutByte(x byte) {
	o.PutUInt8(x)
}

func (o *Out) PutUInt8(x uint8) {
	o.writeUnsafe(x)
}

func (o *Out) PutInt8(x int8) {
	o.writeUnsafe(x)
}

func (o *Out) PutUInt16(x uint16) {
	o.writeUnsafe(x)
}

func (o *Out) PutInt16(x int16) {
	o.writeUnsafe(x)
}

func (o *Out) PutUInt32(x uint32) {
	o.writeUnsafe(x)
}

func (o *Out) PutInt32(x int32) {
	o.writeUnsafe(x)
}

func (o *Out) PutInt64(x uint64) {
	o.writeUnsafe(x)
}

func (o *Out) PutUInt64(x uint64) {
	o.writeUnsafe(x)
}

func (o *Out) PutFloat32(x float32) {
	o.writeUnsafe(x)
}

func (o *Out) PutFloat64(x float64) {
	o.writeUnsafe(x)
}

func (o *Out) PutBytes(x []byte) {
	o.PutUints8(x)
}

func (o *Out) PutUints8(x []uint8) {
	o.PutUInt32(uint32(len(x)))
	o.writeUnsafe(x)
}

func (o *Out) PutBinary(bin encoding.BinaryMarshaler) {
	bytes, err := bin.MarshalBinary()
	if err != nil {
		o.err = err
	} else {
		o.PutBytes(bytes)
	}
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
