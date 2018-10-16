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

type NamedGoFunction interface {
	Name() string
	Call(state *State) int
}

func (s *State) NewLib(funcs []NamedGoFunction) {
	s.CreateTable(0, len(funcs))
	s.SetFuncs(funcs, 0)
}

func (s *State) SetFuncs(funcs []NamedGoFunction, nup int) {
	s.CheckStack(nup, "too many upvalues")
	for _, item := range funcs {
		for i := 0; i < nup; i++ {
			s.PushValue(-nup)
		}
		s.pushClosure(item, nup)
		s.SetField(-(nup + 2), item.Name())
	}
	s.Pop(nup)
}

func (s *State) LoadLib(opener NamedGoFunction) {
	s.Require(opener, true)
	s.Pop(1)
}

func (s *State) PreloadLib(opener NamedGoFunction) {
	s.EnsureTable(RegistryIndex, KeyPreloadTable)
	s.pushClosure(opener, 0)
	s.SetField(-2, opener.Name())
	s.Pop(1)
}

func (s *State) Require(opener NamedGoFunction, global bool) {
	name := opener.Name()
	s.EnsureTable(RegistryIndex, KeyLoadedTable)
	s.GetField(-1, name) // _LOADED[name]
	if !s.ToBoolean(-1) {
		s.Pop(1)
		s.pushClosure(opener, 0)
		s.PushString(name) // argument to open function
		s.Call(1, 1)       // call 'opener' to open module
		s.PushValue(-1)
		s.SetField(-3, name) // _LOADED[name] = module
	}
	s.Remove(-2)
	if global {
		s.PushValue(-1)
		s.SetGlobal(name) // _G[name] = module
	}
}

type WithGoFunction interface {
	Func() GoFunction
}

func (s *State) pushClosure(fun NamedGoFunction, nup int) {
	if lib, ok := fun.(Lib); ok {
		opener := lib.opener()
		if opener != nil {
			C.lua_pushcclosure(s.L, opener, C.int(nup))
			return
		}
	}

	goFun := fun.Call
	if w, ok := fun.(WithGoFunction); ok {
		goFun = w.Func()
	}

	s.PushGoClosure(goFun, nup)
}

type GoFunctionHolder struct {
	FunName string
	Fun     GoFunction
}

func (h GoFunctionHolder) Name() string {
	return h.FunName
}

func (h GoFunctionHolder) Call(state *State) int {
	return h.Fun(state)
}

func (h GoFunctionHolder) Func() GoFunction {
	return h.Fun
}
