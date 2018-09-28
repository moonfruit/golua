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

	var err error
	for err == nil {
		var n int
		n, err = ctx.Read(ctx.buf[:])
		if n > 0 {
			*sz = C.size_t(n)
			return (uintptr)(unsafe.Pointer(&ctx.buf[0]))
		}
	}
	if err == io.EOF {
		*sz = 0
		return 0
	}

	// FIXME: error handling
	errString := []byte(fmt.Sprintf("read error for `%v`", err))
	C.lua_pushlstring(L, (*C.char)(unsafe.Pointer(&errString[0])), C.size_t(len(errString)))
	C.lua_error(L)
	return 0
}

type goReaderCtx struct {
	io.Reader
	buf []byte
}

func newReaderCtxSize(r io.Reader, size int) *goReaderCtx {
	return &goReaderCtx{r, make([]byte, size)}
}

func newReaderCtx(r io.Reader) *goReaderCtx {
	return newReaderCtxSize(r, defaultBufSize)
}
