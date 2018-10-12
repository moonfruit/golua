package golua

//#include <lauxlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

type Buffer struct {
	s    *State
	buf  *C.luaL_Buffer
	data []byte
}

func (s *State) NewBuffer() *Buffer {
	data := make([]byte, C.sizeof_luaL_Buffer)
	buf := (*C.luaL_Buffer)(unsafe.Pointer(&data[0]))
	C.luaL_buffinit(s.L, buf)
	return &Buffer{s, buf, data}
}

func (b *Buffer) State() *State {
	return b.s
}

func (b *Buffer) AddChar(c byte) {
	if b.buf.n >= b.buf.size {
		C.luaL_prepbuffsize(b.buf, 1)
	}
	*((*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(b.buf.b)) + uintptr(b.buf.n)))) = C.char(c)
	b.buf.n++
}

func (b *Buffer) AddStringf(format string, args ...interface{}) {
	b.AddString(fmt.Sprintf(format, args...))
}

func (b *Buffer) AddString(str string) {
	b.addChars(C._GoStringPtr(str), len(str))
}

func (b *Buffer) AddBytes(bytes []byte) {
	b.addChars((*C.char)(unsafe.Pointer(&bytes[0])), len(bytes))
}

func (b *Buffer) addChars(chars *C.char, len int) {
	C.luaL_addlstring(b.buf, chars, C.size_t(len))
}

func (b *Buffer) AddValue() {
	C.luaL_addvalue(b.buf)
}

func (b *Buffer) PushResult() {
	C.luaL_pushresult(b.buf)
}
