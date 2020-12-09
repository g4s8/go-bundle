package bundle

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Inp struct {
	bo  binary.ByteOrder
	src *bytes.Reader
	err error
}

func NewInput(order binary.ByteOrder) *Inp {
	return &Inp{bo: order}
}

func NewBEInput() *Inp {
	return NewInput(binary.BigEndian)
}

func BEInputFromBytes(bts []byte) *Inp {
	return &Inp{bo: binary.BigEndian, src: bytes.NewReader(bts)}
}

func NewLEInput() *Inp {
	return NewInput(binary.LittleEndian)
}

var errAlreadyUnmarshalled = errors.New("bunlde input already unmarshalled")

func (i *Inp) UnmarshalBinary(data []byte) error {
	if i.src != nil {
		return errAlreadyUnmarshalled
	}
	i.src = bytes.NewReader(data)
	return nil
}

func (i *Inp) Error() error {
	return i.err
}

func (i *Inp) GetInt32(out *int32) {
	i.readUnsafe(out)
}

func (i *Inp) GetUInt32(out *uint32) {
	i.readUnsafe(out)
}

func (i *Inp) GetInt64(out *int64) {
	i.readUnsafe(out)
}

func (i *Inp) GetUInt64(out *uint64) {
	i.readUnsafe(out)
}

func (i *Inp) GetBytes(out *[]byte) {
	i.GetUInts8(out)
}

func (i *Inp) GetUInts8(out *[]uint8) {
	var size uint32
	i.readUnsafe(&size)
	if i.err != nil {
		return
	}
	arr := make([]uint8, size)
	i.readUnsafe(arr)
	*out = arr
}

func (i *Inp) readUnsafe(out interface{}) {
	if i.err != nil {
		return
	}
	i.err = binary.Read(i.src, i.bo, out)
}
