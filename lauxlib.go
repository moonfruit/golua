package golua

// #include <lauxlib.h>
import "C"

func NewState() (*State, error) {
	L := C.luaL_newstate()
	if L == nil {
		return nil, ErrMem
	}
	return &State{L}, nil
}
