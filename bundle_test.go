package bundle

import (
	m "github.com/g4s8/go-matchers"
	"testing"
)

func Test_PutGetInt32(t *testing.T) {
	assert := m.Assert(t)
	out := NewBEOut()
	out.PutInt32(42)
	inp := out.Flip()
	var i int32
	inp.GetInt32(&i)
	assert.That("Read without errors", inp.Error(), m.Nil())
	assert.That("Read int32 correctly", i, m.Eq(int32(42)))
}

func Test_PutGetByteArray(t *testing.T) {
	assert := m.Assert(t)
	out := NewBEOut()
	out.PutBytes([]byte{0x00, 0x01, 0xab})
	inp := out.Flip()
	var i []byte
	inp.GetBytes(&i)
	assert.That("Read without errors", inp.Error(), m.Nil())
	assert.That("Read int32 correctly", i, m.EqBytes([]byte{0x00, 0x01, 0xab}))
}

func Test_ReadString(t *testing.T) {
	assert := m.Assert(t)
	out := NewLEOut()
	out.PutString("Hello bundle!")
	inp := out.Flip()
	var s string
	inp.GetString(&s)
	assert.That("Read without errors", inp.Error(), m.Nil())
	assert.That("Reads string correctly", s, m.Eq("Hello bundle!"))
}
