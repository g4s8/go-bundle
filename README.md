Binary encoding and decoding tool, type-safe wrapper for `encoding/binary` package.

Implements `BinaryMarshaler` and `BinaryUnmarshaler` from `encoding`, supports different
`encoding.BiteOrder`s. Errors accumulator to avoid each line error checks.
It could be useful to implement custom `encoding` binary protocols.

# Install

Add to imports: `"github.com/g4s8/go-bundle"`

# Usage

## Create new bundle output

There are different factory functions to create bundle outputs:
```go
out := NewOut(order) // creates bundle and specify byte order
outBE := NewBEOut() // creates bundle with big-endian byte order
outLE := NewLEOut() // creates bundle with little-endian byte order
```

## Put data into bundle

Each out bundle method has type-name suffix and accepts only
fixed-length and strongly typed parameters:
```go
out := NewBEOut()
out.PutInt32(int32(42))
out.PutUIint64(uint64(1111111111111))
out.PutFloat32(float32(1.4))
out.PutBytes([]byte{0x00, 0x01, 0x02})
```

## Errors

Output `Put*` methods doesn't return errors, but in case of error it will be
accumulated. It can be checked with `out.Error()` method. If error happens,
put method won't try to put next values to internal bytes buffer.

## Marshaling

Output implements `encoding.BinaryMarshaler` protocol. It can be used to
get byte array of the bundle:
```go
out := NewBEOut()
out.PutInt32(int32(42))
bytes, err := out.MarshalBinary()
```

It returns buffer bytes or first error.
