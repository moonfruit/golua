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
	"unsafe"
)

const (
	defaultBufSize = 4096
)

func cfree(p *C.char) {
	C.free(unsafe.Pointer(p))
}

//export goReader
func goReader(L *C.lua_State, ud uintptr, sz *C.size_t) uintptr {
	val, _ := pool.Get(ud)
	ctx := val.(*goReaderCtx)

	for ctx.err == nil {
		var n int
		n, ctx.err = ctx.Read(ctx.buf[:])
		if n > 0 {
			*sz = C.size_t(n)
			return (uintptr)(unsafe.Pointer(&ctx.buf[0]))
		}
	}

	*sz = 0
	return 0
}

type goReaderCtx struct {
	io.Reader
	buf []byte
	err error
}

func newReaderCtxSize(r io.Reader, size int) *goReaderCtx {
	if size <= 0 {
		size = defaultBufSize
	}
	return &goReaderCtx{Reader: r, buf: make([]byte, size)}
}

func newReaderCtx(r io.Reader) *goReaderCtx {
	return newReaderCtxSize(r, defaultBufSize)
}
