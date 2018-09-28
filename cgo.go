package golua

// #cgo CFLAGS: -I${SRCDIR}/lua -I${SRCDIR}/lpeg
// #cgo CFLAGS: -DLUA_COMPAT_5_2 -DLUA_INT_TYPE=3 -DLUA_FLOAT_TYPE=2
// #cgo linux darwin CFLAGS: -DLUA_USE_POSIX
//
// #include <stdlib.h>
// #include <lua.h>
import "C"
import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"unsafe"
)

const (
	defaultBufSize = 4096
)

func cfree(p *C.char) {
	C.free(unsafe.Pointer(p))
}

//export goReader
func goReader(L *C.lua_State, ud unsafe.Pointer, sz *C.size_t) *C.char {
	val, ok := pool.Get(*(*uintptr)(ud))
	if !ok {
		panic("can not get reader")
	}

	ctx := val.(*goReaderCtx)
	if len(ctx.Bytes) == 0 {
		panic("no buffer")
	}

	var err error
	for err == nil {
		var n int
		n, err = ctx.Read(ctx.Bytes[:])
		if n > 0 {
			*sz = C.size_t(n)
			return (*C.char)(ctx.Pointer)
		}
	}
	if err == io.EOF {
		*sz = 0
		return nil
	}

	// FIXME: error handling
	errString := []byte(fmt.Sprintf("read error for `%v`", err))
	C.lua_pushlstring(L, (*C.char)(unsafe.Pointer(&errString[0])), C.size_t(len(errString)))
	C.lua_error(L)
	return nil
}

type goReaderCtx struct {
	io.Reader
	*hackedSlice
}

func newReaderCtxSize(r io.Reader, size int) *goReaderCtx {
	return &goReaderCtx{r, newSlice(size)}
}

func newReaderCtx(r io.Reader) *goReaderCtx {
	return newReaderCtxSize(r, defaultBufSize)
}

func (ctx *goReaderCtx) Close() {
	ctx.free()
}

type hackedSlice struct {
	Bytes    []byte
	Pointer  unsafe.Pointer
	needFree bool
}

func sliceAt(pointer unsafe.Pointer, size int) *hackedSlice {
	s := &hackedSlice{Pointer: pointer}

	h := (*reflect.SliceHeader)(unsafe.Pointer(&s.Bytes))
	h.Len = size
	h.Cap = size
	h.Data = uintptr(pointer)

	return s
}

func newSlice(size int) *hackedSlice {
	ptr := C.malloc(C.size_t(size))
	if ptr == nil {
		panic("memory allocation error")
	}

	s := sliceAt(ptr, size)
	s.needFree = true

	runtime.SetFinalizer(s, func(s *hackedSlice) {
		s.free()
	})

	return s
}

func (s *hackedSlice) free() {
	if s.Bytes != nil {
		s.Bytes = nil
		if s.needFree {
			C.free(s.Pointer)
			s.Pointer = nil
		}
	}
}
