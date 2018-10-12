package golua

// #include "cgo.h"
import "C"

type GoFunction func(state *State) int

func (s *State) PushGoFunction(fun GoFunction) {
	C.luaGo_pushGoFunction(s.L, C.ulong(s.RefGoValue(fun)))
}

func (s *State) ToGoFunction(idx int) (GoFunction, bool) {
	id := s.testUserData(idx, MetaGoFunction)
	if id == 0 {
		return nil, false
	}
	return s.GetGoValue(id).(GoFunction), true
}

func (s *State) CheckGoFunction(idx int) GoFunction {
	id := s.checkUserData(idx, MetaGoFunction)
	return s.GetGoValue(id).(GoFunction)
}
