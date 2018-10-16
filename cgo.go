package golua

/*
#cgo CFLAGS: -I${SRCDIR}/lua -I${SRCDIR}/lpeg
#cgo CFLAGS: -DLUA_COMPAT_5_2 -DLUA_INT_TYPE=3 -DLUA_FLOAT_TYPE=2
#cgo linux darwin CFLAGS: -DLUA_USE_POSIX

#include <stdlib.h>
#include <lua.h>
*/
import "C"
import (
	"io"
	"reflect"
	"runtime"
	"unsafe"
)

const (
	defaultBufSize           = 4096
	maxConsecutiveEmptyReads = 100
)

func cfree(p *C.char) {
	C.free(unsafe.Pointer(p))
}

//export goFree
func goFree(L *C.lua_State, ud uintptr) {
	mainStateFor(L).UnRefGoValue(ud)
}

//export goCall
func goCall(L *C.lua_State, ud uintptr) C.int {
	state := mainStateFor(L)
	fun := state.GetGoValue(ud).(GoFunction)
	return C.int(fun(state))
}

//export goReader
func goReader(L *C.lua_State, ud unsafe.Pointer, sz *C.size_t) *C.char {
	ctx := mainStateFor(L).GetGoValue(uintptr(ud)).(*goReaderCtx)

	for i := 0; ctx.err == nil && i < maxConsecutiveEmptyReads; i++ {
		var n int
		n, ctx.err = ctx.Read(ctx.Bytes)
		if n > 0 {
			*sz = C.size_t(n)
			return (*C.char)(ctx.Pointer)
		}
	}

	if ctx.err == nil {
		ctx.err = io.ErrNoProgress
	}
	*sz = 0
	return nil
}

type goReaderCtx struct {
	io.Reader
	*hackedSlice
	err error
}

func newReaderCtxSize(r io.Reader, size int) *goReaderCtx {
	if size <= 0 {
		size = defaultBufSize
	}
	return &goReaderCtx{Reader: r, hackedSlice: newSlice(size)}
}

func newReaderCtx(r io.Reader) *goReaderCtx {
	return newReaderCtxSize(r, defaultBufSize)
}

type hackedSlice struct {
	Bytes    []byte
	Pointer  unsafe.Pointer
	needFree bool
}

func newSlice(size int) *hackedSlice {
	return newSliceAt(C.malloc(C.size_t(size)), size, true)
}

func newSliceAt(pointer unsafe.Pointer, size int, needFree bool) *hackedSlice {
	s := &hackedSlice{Pointer: pointer, needFree: needFree}

	if needFree {
		runtime.SetFinalizer(s, func(s *hackedSlice) {
			s.free()
		})
	}

	h := (*reflect.SliceHeader)(unsafe.Pointer(&s.Bytes))
	h.Len = size
	h.Cap = size
	h.Data = uintptr(pointer)

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
