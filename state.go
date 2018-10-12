package golua

// #include <lua.h>
// #include <lauxlib.h>
// #include "cgo.h"
import "C"
import (
	"sync"
	"unsafe"
)

var (
	statePool      map[uintptr]*State
	statePoolMutex sync.RWMutex
)

func ptr(L *C.lua_State) uintptr {
	return uintptr(unsafe.Pointer(L))
}

func register(state *State) *State {
	statePoolMutex.Lock()
	defer statePoolMutex.Unlock()

	if statePool == nil {
		statePool = make(map[uintptr]*State)
	}

	statePool[ptr(state.L)] = state
	return state
}

func mainStateFor(L *C.lua_State) *State {
	statePoolMutex.RLock()
	defer statePoolMutex.RUnlock()

	state, ok := statePool[ptr(L)]
	if !ok && L != nil {
		L = C.luaGo_main(L)
		state = statePool[ptr(L)]
	}

	return state
}

func unregister(state *State) {
	statePoolMutex.Lock()
	defer statePoolMutex.Unlock()

	delete(statePool, ptr(state.L))
}

type State struct {
	L        *C.lua_State
	registry sliceRegistry
}

func NewState() *State {
	L := C.luaL_newstate()
	if L == nil {
		panic(ErrMem)
	}
	return register(&State{L: L})
}

func (s *State) Close() {
	if s.L != nil {
		C.lua_close(s.L)
		unregister(s)
		s.L = nil
	}
}

func (s *State) RefGoValue(val interface{}) uintptr {
	return s.registry.Ref(val)
}

func (s *State) UnRefGoValue(id uintptr) {
	s.registry.UnRef(id)
}

func (s *State) GetGoValue(id uintptr) interface{} {
	return s.registry.Get(id)
}

func (s *State) ToGoValue(idx int) interface{} {
	id := *(*uintptr)(C.lua_touserdata(s.L, C.int(idx)))
	return s.GetGoValue(id)
}
