Binary encoding and decoding tool, type-safe wrapper for `encoding/binary` package.

Implements `BinaryMarshaler` and `BinaryUnmarshaler` from `encoding`, supports different
`encoding.BiteOrder`s. Errors accumulator to avoid each line error checks.
It could be useful to implement custom `encoding` binary protocols.

# Install

Add to imports: `"github.com/g4s8/go-bundle"`

# Usage

The bundle package provides two main types:
`Out` and `Inp`. `Out` (output) is used to put data into bundle
and produce byte array from it, `Inp` (input) is used to read
data from bytes array.

## Output

### Create new bundle output

There are different factory functions to create bundle outputs:
```go
out := NewOut(order) // creates bundle and specify byte order
outBE := NewBEOut() // creates bundle with big-endian byte order
outLE := NewLEOut() // creates bundle with little-endian byte order
```

### Put data into bundle

Each out bundle method has type-name suffix and accepts only
fixed-length and strongly typed parameters:
```go
out := NewBEOut()
out.PutInt32(int32(42))
out.PutUIint64(uint64(1111111111111))
out.PutFloat32(float32(1.4))
out.PutBytes([]byte{0x00, 0x01, 0x02})
out.PutString("hello bundle!")
```

Also, output accepts `encoding.BinaryMarshaler` types to put:
`out.PutBinary(binary)`.

### Errors

Output `Put*` methods doesn't return errors, but in case of error it will be
accumulated. It can be checked with `out.Error()` method. If error happens,
put method won't try to put next values to internal bytes buffer.

### Marshaling

Output implements `encoding.BinaryMarshaler` protocol. It can be used to
get byte array of the bundle:
```go
out := NewBEOut()
out.PutInt32(int32(42))
bytes, err := out.MarshalBinary()
```

It returns buffer bytes or first error.

### Flip

`Flip` methods creates input from output: `func (o *Out) Flip() *Inp`.
It can be useful for testing.

## Input

### Creating new input

These factory functions can be used to wrap byte array with
bundle input type:
```go

i1 := NewInput(order) // creates input with specified order
i1 := NewBEInput() // creates big-endian input
i3 := NewLEInput() // creates little-endian input
i4 := InpWrapBytes(order, bts) // creates input from bytes with order
```

### Filling input with data

Bundle input implements `encoding.BinaryUnmarshaler` protocol, it means
it can be read from byte array:
```go
inp := NewInput(binary.BigEndian)
err := inp.UnmarshalBinary(data)
```
Alternatively, it can be constructed directly from bytes using wrap factory method:
```go
inp := InpWrapBytes(order, bytes)
```

Input can't be read twise, the call to `UnmarshalBinary` returns an error
if input was already read.

### Reading data from input

Input bundle methods has type suffixes and accepts output reference
parameters to read the data:
```go
// Reads int32
var i32 int32
inp.ReadInt32(&i32)

// Reads byte-array
var bts []byte
inp.ReadBytes(&bts)

// Read binary unmarshaler
// subInput - implements `encoding.BinaryUnmarshaler`
inp.ReadBinary(subInput)
```

### Errors

Errors are accumulated and can be accessed with `Error()` method.
In case of error, next read methods won't be performed.

## Example

Assuming you need to implement `encoding.BinaryMarshaler` and
`encoding.BinaryUnmarshaler` for some struct. You can use `bundle`
to do that:
```go
import (
  "encoding/binary"
  "github.com/g4s8/go-bundle"
  ma "github.com/multiformats/go-multiaddr"
)

type Foo struct {
  ID int32
  Description string
  internal []byte
  Addr ma.Multiaddr
}

var byteOrder = binary.BigEndian

func (f *Foo) MarshalBinary() ([]byte, error) {
  out := bundle.NewOut(byteOrder)
  out.PutInt32(f.ID)
  out.PutString(f.Description)
  out.PutBytes(f.internal)
  out.PutBinary(f.Addr)
  return out.MarshalBinary()
}

func (f *Foo) UnmarshalBinary(data []byte) error {
  inp := bundle.NewInput(byteOrder)
  if err := inp.UnmarshalBinary(data); err != nil {
         return err
  }
  inp.GetInt32(&f.ID)
  inp.GetString(&f.Description)
  inp.GetBytes(&f.internal)
  var rawAddr []byte
  inp.GetBytes(&rawAddr)
  addr, err := ma.NewMultiaddrBytes(rawAddr)
  if err != nil {
    return err
  }
  f.Addr = addr
  return nil
}
```
