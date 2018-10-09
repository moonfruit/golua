package golua

/*
#include <lauxlib.h>

const size_t bufSize = sizeof(luaL_Buffer);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Buffer struct {
	s    State
	buf  *C.luaL_Buffer
	data []byte
}

func (s State) NewBuffer() *Buffer {
	data := make([]byte, C.bufSize)
	buf := (*C.luaL_Buffer)(unsafe.Pointer(&data[0]))
	C.luaL_buffinit(s.L, buf)
	return &Buffer{s, buf, data}
}

func (b *Buffer) State() State {
	return b.s
}

func (b *Buffer) AddChar(c byte) error {
	if b.buf == nil {
		return ErrInvalidBuffer
	}
	if b.buf.n >= b.buf.size {
		C.luaL_prepbuffsize(b.buf, 1)
	}
	*((*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(b.buf.b)) + uintptr(b.buf.n)))) = C.char(c)
	b.buf.n++
	return nil
}

func (b *Buffer) AddStringf(format string, args ...interface{}) error {
	return b.AddString(fmt.Sprintf(format, args...))
}

func (b *Buffer) AddString(str string) error {
	return b.AddBytes([]byte(str))
}

func (b *Buffer) AddBytes(bytes []byte) error {
	if b.buf == nil {
		return ErrInvalidBuffer
	}
	C.luaL_addlstring(b.buf, (*C.char)(unsafe.Pointer(&bytes[0])), C.size_t(len(bytes)))
	return nil
}

func (b *Buffer) AddValue() error {
	if b.buf == nil {
		return ErrInvalidBuffer
	}
	C.luaL_addvalue(b.buf)
	return nil
}

func (b *Buffer) PushResult() {
	C.luaL_pushresult(b.buf)
}
