package golua

// #include "cgo.h"
import "C"

type GoFunction func(state *State) int

func (s *State) PushGoClosure(fun GoFunction, n int) {
	s.PushGoValue(fun)
	C.lua_pushcclosure(s.L, (C.lua_CFunction)(C.luaGo_callGoFunction), C.int(n+1))
}

func (s *State) PushGoFunction(fun GoFunction) {
	s.PushGoClosure(fun, 0)
}

func (s *State) ToGoFunction(idx int) (GoFunction, bool) {
	if ret := C.luaGo_getGoFunction(s.L, C.int(idx)); ret == 0 {
		return nil, false
	}
	defer s.Pop(1)

	fun, ok := s.ToGoValue(-1).(GoFunction)
	if !ok {
		return nil, false
	}

	return fun, true
}
