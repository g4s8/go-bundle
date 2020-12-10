package bundle

import (
	"encoding/binary"
	"fmt"
	"testing"

	m "github.com/g4s8/go-matchers"
)

func TestNesting(t *testing.T) {
	assert := m.Assert(t)
	one := NewLEOut()
	one.PutInt32(4)
	two := NewLEOut()
	two.PutInt64(8)
	two.PutBundle(one)
	two.PutBool(true)
	inp := two.Flip()
	var i64 int64
	inp.GetInt64(&i64)
	nested := new(Inp)
	inp.GetBundle(nested)
	var i32 int32
	nested.GetInt32(&i32)
	var b bool
	inp.GetBool(&b)
	assert.That("Reads parent without errors", inp.Error(), m.Nil())
	assert.That("Reads nested without errors", nested.Error(), m.Nil())
	assert.That("Reads i64", i64, m.Eq(int64(8)))
	assert.That("Reads i32", i32, m.Eq(int32(4)))
	assert.That("Reads bool", b, m.Eq(true))
}

func testerEncDec(bo binary.ByteOrder,
	init func(o *Out), verify func(m.Assertion, *Inp)) func(t *testing.T) {
	return func(t *testing.T) {
		assert := m.Assert(t)
		out := NewOut(bo)
		init(out)
		assert.That("Initialized successfully", out.Error(), m.Nil())
		inp := out.Flip()
		verify(assert, inp)
		assert.That("Verified successfully", inp.Error(), m.Nil())
	}
}

func TestEncodeDecodeTypes(t *testing.T) {
	for _, bo := range []binary.ByteOrder{binary.BigEndian, binary.LittleEndian} {
		t.Run(fmt.Sprintf("bools with %s", bo), func(t *testing.T) {
			for _, b := range []bool{true, false} {
				t.Run(fmt.Sprintf("bool: %v", b), testerEncDec(bo,
					func(out *Out) {
						out.PutBool(b)
					}, func(assert m.Assertion, inp *Inp) {
						var res bool
						inp.GetBool(&res)
						assert.That("reads same bool", res, m.Is(b))
					}))
			}
		})
		t.Run(fmt.Sprintf("int8s with %s", bo), func(t *testing.T) {
			for _, i := range []int8{-128, -42, 0, 42, 127} {
				t.Run(fmt.Sprintf("int8: %d", i), testerEncDec(bo,
					func(out *Out) {
						out.PutInt8(i)
					}, func(assert m.Assertion, inp *Inp) {
						var res int8
						inp.GetInt8(&res)
						assert.That("reads same int8", res, m.Is(i))
					}))
			}
		})
		t.Run(fmt.Sprintf("uint8s with %s", bo), func(t *testing.T) {
			for _, i := range []uint8{0, 32, 64, 127, 255} {
				t.Run(fmt.Sprintf("uint8: %d", i), testerEncDec(bo,
					func(out *Out) {
						out.PutUInt8(i)
					}, func(assert m.Assertion, inp *Inp) {
						var res uint8
						inp.GetUInt8(&res)
						assert.That("reads same uint8", res, m.Is(i))
					}))
			}
		})
		t.Run(fmt.Sprintf("int16s with %s", bo), func(t *testing.T) {

			for _, i := range []int16{-32768, -257, 0, 127, 258, 32767} {
				t.Run(fmt.Sprintf("int16: %d", i), testerEncDec(bo,
					func(out *Out) {
						out.PutInt16(i)
					}, func(assert m.Assertion, inp *Inp) {
						var res int16
						inp.GetInt16(&res)
						assert.That("reads same int16", res, m.Is(i))
					}))
			}
		})
		t.Run(fmt.Sprintf("uint16s with %s", bo), func(t *testing.T) {
			for _, i := range []uint16{0, 33, 66, 130, 241, 1025, 4099, 65535} {
				t.Run(fmt.Sprintf("uint16: %d", i), testerEncDec(bo,
					func(out *Out) {
						out.PutUInt16(i)
					}, func(assert m.Assertion, inp *Inp) {
						var res uint16
						inp.GetUInt16(&res)
						assert.That("reads same uint16", res, m.Is(i))
					}))
			}
		})
		t.Run(fmt.Sprintf("int32s with %s", bo), func(t *testing.T) {
			for _, i := range []int32{-2147483648, -65546, -129, 0, 130, 257, 65536, 2147483647} {
				t.Run(fmt.Sprintf("int32: %d", i), testerEncDec(bo,
					func(out *Out) {
						out.PutInt32(i)
					}, func(assert m.Assertion, inp *Inp) {
						var res int32
						inp.GetInt32(&res)
						assert.That("reads same int32", res, m.Is(i))
					}))
			}
		})
		t.Run(fmt.Sprintf("uint32s with %s", bo), func(t *testing.T) {
			for _, i := range []uint32{0, 130, 257, 65536, 2147483647, 4294967295} {
				t.Run(fmt.Sprintf("int32: %d", i), testerEncDec(bo,
					func(out *Out) {
						out.PutUInt32(i)
					}, func(assert m.Assertion, inp *Inp) {
						var res uint32
						inp.GetUInt32(&res)
						assert.That("reads same uint32", res, m.Is(i))
					}))
			}
		})
		t.Run(fmt.Sprintf("int64s with %s", bo), func(t *testing.T) {
			for _, i := range []int64{-9223372036854775808, -4294967296, 0, 32,
				4294967295, 9223372036854775807} {
				t.Run(fmt.Sprintf("int64: %d", i), testerEncDec(bo,
					func(out *Out) {
						out.PutInt64(i)
					}, func(assert m.Assertion, inp *Inp) {
						var res int64
						inp.GetInt64(&res)
						assert.That("reads same int64", res, m.Is(i))
					}))
			}
		})
		t.Run(fmt.Sprintf("int64s with %s", bo), func(t *testing.T) {
			for _, i := range []uint64{0, 32, 4294967295,
				9223372036854775807, 18446744073709551615} {
				t.Run(fmt.Sprintf("uint64: %d", i), testerEncDec(bo,
					func(out *Out) {
						out.PutUInt64(i)
					}, func(assert m.Assertion, inp *Inp) {
						var res uint64
						inp.GetUInt64(&res)
						assert.That("reads same uint64", res, m.Is(i))
					}))
			}
		})
		t.Run(fmt.Sprintf("float32 with %s", bo), func(t *testing.T) {
			for _, i := range []float32{-1.2, -3.000003, 0, 127, 9312312451} {
				t.Run(fmt.Sprintf("float32: %f", i), testerEncDec(bo,
					func(out *Out) {
						out.PutFloat32(i)
					}, func(assert m.Assertion, inp *Inp) {
						var res float32
						inp.GetFloat32(&res)
						assert.That("reads same float32", res, m.Is(i))
					}))
			}
		})
		t.Run(fmt.Sprintf("float64 with %s", bo), func(t *testing.T) {
			for _, i := range []float64{-103200000.0000002, -3.000003, 0, 127, 9312312451} {
				t.Run(fmt.Sprintf("float32: %f", i), testerEncDec(bo,
					func(out *Out) {
						out.PutFloat64(i)
					}, func(assert m.Assertion, inp *Inp) {
						var res float64
						inp.GetFloat64(&res)
						assert.That("reads same float64", res, m.Is(i))
					}))
			}
		})
		t.Run(fmt.Sprintf("strings with %s", bo), func(t *testing.T) {
			for _, i := range []string{"hello", "юникод", "\n\ba"} {
				t.Run(fmt.Sprintf("string: %s", i), testerEncDec(bo,
					func(out *Out) {
						out.PutString(i)
					}, func(assert m.Assertion, inp *Inp) {
						var res string
						inp.GetString(&res)
						assert.That("reads same string", res, m.Eq(i))
					}))
			}
		})
		t.Run(fmt.Sprintf("byte arrays with %s", bo), func(t *testing.T) {
			for _, i := range [][]byte{{}, {0x00}, {0xca, 0xfe, 0xba, 0xbe}} {
				t.Run(fmt.Sprintf("[]byte: %v", i), testerEncDec(bo,
					func(out *Out) {
						out.PutBytes(i)
					}, func(assert m.Assertion, inp *Inp) {
						var res []byte
						inp.GetBytes(&res)
						assert.That("reads same bytes", res, m.EqBytes(i))
					}))
			}
		})
	}
}

func sizeKb(kb int, item int) int {
	return kb * 1024 * 1024 / item
}

const (
	sizeInt64 = 8
	sizeInt32 = 4
	sizeInt16 = 2
	sizeInt8  = 1
)

func BenchmarkEncDec(b *testing.B) {
	for _, bo := range []binary.ByteOrder{binary.BigEndian, binary.LittleEndian} {
		b.Run(fmt.Sprintf("encode int64 with %s", bo), func(b *testing.B) {
			for _, size := range []int{sizeKb(1, sizeInt64), sizeKb(4, sizeInt64),
				sizeKb(16, sizeInt64), sizeKb(32, sizeInt64), sizeKb(64, sizeInt64)} {
				b.Run(fmt.Sprintf("with size %d kb", size), func(b *testing.B) {
					arr := make([]int64, size)
					for i := 0; i < size; i++ {
						arr[i] = int64(i*2 - 1)
					}
					out := NewOut(bo)
					b.ResetTimer()
					for i := range arr {
						out.PutInt64(arr[i])
					}
				})
			}
		})
		b.Run(fmt.Sprintf("encode int32 with %s", bo), func(b *testing.B) {
			for _, size := range []int{sizeKb(1, sizeInt32), sizeKb(4, sizeInt32),
				sizeKb(16, sizeInt32), sizeKb(32, sizeInt32), sizeKb(64, sizeInt32)} {
				b.Run(fmt.Sprintf("with size %d kb", size), func(b *testing.B) {
					arr := make([]int32, size)
					for i := 0; i < size; i++ {
						arr[i] = int32(i*2 - 1)
					}
					out := NewOut(bo)
					b.ResetTimer()
					for i := range arr {
						out.PutInt32(arr[i])
					}
				})
			}
		})
		b.Run(fmt.Sprintf("encode int16 with %s", bo), func(b *testing.B) {
			for _, size := range []int{sizeKb(1, sizeInt16), sizeKb(4, sizeInt16),
				sizeKb(16, sizeInt16), sizeKb(32, sizeInt16), sizeKb(64, sizeInt16)} {
				b.Run(fmt.Sprintf("with size %d kb", size), func(b *testing.B) {
					arr := make([]int16, size)
					for i := 0; i < size; i++ {
						arr[i] = int16(i*2 - 1)
					}
					out := NewOut(bo)
					b.ResetTimer()
					for i := range arr {
						out.PutInt16(arr[i])
					}
				})
			}
		})
		b.Run(fmt.Sprintf("encode int8 with %s", bo), func(b *testing.B) {
			for _, size := range []int{sizeKb(1, sizeInt8), sizeKb(4, sizeInt8),
				sizeKb(16, sizeInt8), sizeKb(32, sizeInt8), sizeKb(64, sizeInt8)} {
				b.Run(fmt.Sprintf("with size %d kb", size), func(b *testing.B) {
					arr := make([]int8, size)
					for i := 0; i < size; i++ {
						arr[i] = int8(i*2 - 1)
					}
					out := NewOut(bo)
					b.ResetTimer()
					for i := range arr {
						out.PutInt8(arr[i])
					}
				})
			}
		})
	}
}
