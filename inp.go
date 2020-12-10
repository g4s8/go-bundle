package bundle

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
)

// Inp is bundle input type, it wraps source byte array
// and read data with specified bytes order
type Inp struct {
	bo  binary.ByteOrder
	src *bytes.Reader
	err error
}

// NewInput creates new empty input with specified byte order,
// the data could be added with `UnmarshalBinary` method.
func NewInput(order binary.ByteOrder) *Inp {
	return &Inp{bo: order}
}

// NewBEInput creates new empty input with big-endian byte order,
// the data could be added with `UnmarshalBinary` method.
func NewBEInput() *Inp {
	return NewInput(binary.BigEndian)
}

// NewBEInput creates new empty input with little-endian byte order,
// the data could be added with `UnmarshalBinary` method.
func NewLEInput() *Inp {
	return NewInput(binary.LittleEndian)
}

// InpWrapBytes creates new input from bytes source with
// specified byte order.
func InpWrapBytes(bts []byte, order binary.ByteOrder) *Inp {
	return &Inp{bo: order, src: bytes.NewReader(bts)}
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

// GetBool reads bool type to `out` param
func (i *Inp) GetBool(out *bool) {
	i.readUnsafe(out)
}

// GetUInt8 reads int8 type to `out` param
func (i *Inp) GetUInt8(out *uint8) {
	i.readUnsafe(out)
}

// GetInt8 reads int8 type to `out` param
func (i *Inp) GetInt8(out *int8) {
	i.readUnsafe(out)
}

// GetUInt16 reads uint16 type to `out` param
func (i *Inp) GetUInt16(out *uint16) {
	i.readUnsafe(out)
}

// GetInt16 reads int16 type to `out` param
func (i *Inp) GetInt16(out *int16) {
	i.readUnsafe(out)
}

// GetRune reads rune type to `out` param
func (i *Inp) GetRune(out *rune) {
	var tmp int32
	i.readUnsafe(&i)
	if i.Error() == nil {
		*out = tmp
	}
}

// GetInt32 reads int32 type to `out` param
func (i *Inp) GetInt32(out *int32) {
	i.readUnsafe(out)
}

// GetUInt32 reads uint32 type to `out` param
func (i *Inp) GetUInt32(out *uint32) {
	i.readUnsafe(out)
}

// GetInt64 reads int64 type to `out` param
func (i *Inp) GetInt64(out *int64) {
	i.readUnsafe(out)
}

// GetUInt64 reads uint64 type to `out` param
func (i *Inp) GetUInt64(out *uint64) {
	i.readUnsafe(out)
}

// GetBytes reads byte array to `out` param
func (i *Inp) GetBytes(out *[]byte) {
	i.GetUInts8(out)
}

// GetFloat32 reads float32 to `out` param
func (i *Inp) GetFloat32(out *float32) {
	i.readUnsafe(out)
}

// GetFloat64 reads float64 to `out` param
func (i *Inp) GetFloat64(out *float64) {
	i.readUnsafe(out)
}

// GetUInts8 reads unit8 array to `out` param
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

// GetString reads strig from bundle
func (i *Inp) GetString(out *string) {
	var b []uint8
	i.GetUInts8(&b)
	if i.err == nil {
		*out = string(b)
	}
}

// GetBinary reads marshaled type to binary unmarshaler
func (i *Inp) GetBinary(bin encoding.BinaryUnmarshaler) {
	var bytes []byte
	i.GetBytes(&bytes)
	if i.err == nil {
		i.err = bin.UnmarshalBinary(bytes)
	}
}

func (i *Inp) GetBundle(b *Inp) {
	if i.err != nil {
		return
	}
	var bo uint8
	i.GetUInt8(&bo)
	switch bo {
	case encBoBE:
		b.bo = binary.BigEndian
	case encBoLE:
		b.bo = binary.LittleEndian
	default:
		i.err = errUnsupportedByteOred
		return
	}
	i.GetBinary(b)
}

func (i *Inp) String() string {
	return fmt.Sprintf("Input unread: %d, order: %s",
		i.src.Len(), i.bo)
}

func (i *Inp) GoString() string {
	return fmt.Sprintf("bundle.Input(Buffer(%#v), Order(%#v))",
		i.src, i.bo)
}

var errNoInputSource = errors.New("input was not unmarshaled: source is nil")

func (i *Inp) readUnsafe(out interface{}) {
	if i.err != nil {
		return
	}
	if i.src == nil {
		i.err = errNoInputSource
		return
	}
	i.err = binary.Read(i.src, i.bo, out)
}
