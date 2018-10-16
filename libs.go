package golua

/*
#include <lualib.h>
#include "cgo.h"

extern int luaopen_lpeg(lua_State *L);
const char *lpeg_name = "lpeg";
*/
import "C"
import (
	"bytes"
	"fmt"

	"github.com/moonfruit/golua/re"
)

func (s *State) OpenLibs() {
	C.luaL_openlibs(s.L)
	C.luaGo_preload(s.L, C.lpeg_name, (C.lua_CFunction)(C.luaopen_lpeg))
	s.PreloadLib(RE)
}

func (s *State) OpenBasicLibs() {
	C.luaGo_openBasicLibs(s.L)
	C.luaGo_preload(s.L, C.lpeg_name, (C.lua_CFunction)(C.luaopen_lpeg))
	s.PreloadLib(RE)
}

type Lib byte

const (
	Base Lib = iota
	Package
	Coroutine
	Table
	IO
	OS
	String
	Math
	UTF8
	Debug
	Bit32
	LPeg
	RE
)

func (l Lib) Name() string {
	if l == Base {
		return "_G"
	} else if l <= RE {
		return l.String()
	} else {
		panic(l.String())
	}
}

func (l Lib) String() string {
	switch l {
	case Base:
		return "base"
	case Package:
		return "package"
	case Coroutine:
		return "coroutine"
	case Table:
		return "table"
	case IO:
		return "io"
	case OS:
		return "os"
	case String:
		return "string"
	case Math:
		return "math"
	case UTF8:
		return "utf8"
	case Debug:
		return "debug"
	case Bit32:
		return "bit32"
	case LPeg:
		return "lpeg"
	case RE:
		return "re"
	}
	return fmt.Sprintf("unknown lib `%d`", l)
}

func (l Lib) Call(state *State) int {
	switch l {
	case Base:
		return int(C.luaopen_base(state.L))
	case Package:
		return int(C.luaopen_package(state.L))
	case Coroutine:
		return int(C.luaopen_coroutine(state.L))
	case Table:
		return int(C.luaopen_table(state.L))
	case IO:
		return int(C.luaopen_io(state.L))
	case OS:
		return int(C.luaopen_os(state.L))
	case String:
		return int(C.luaopen_string(state.L))
	case Math:
		return int(C.luaopen_math(state.L))
	case UTF8:
		return int(C.luaopen_utf8(state.L))
	case Debug:
		return int(C.luaopen_debug(state.L))
	case Bit32:
		return int(C.luaopen_bit32(state.L))
	case LPeg:
		return int(C.luaopen_lpeg(state.L))
	case RE:
		return openAsset(state, "re.lua")
	}
	return state.Error(l)
}

func (l Lib) opener() C.lua_CFunction {
	switch l {
	case Base:
		return (C.lua_CFunction)(C.luaopen_base)
	case Package:
		return (C.lua_CFunction)(C.luaopen_package)
	case Coroutine:
		return (C.lua_CFunction)(C.luaopen_coroutine)
	case Table:
		return (C.lua_CFunction)(C.luaopen_table)
	case IO:
		return (C.lua_CFunction)(C.luaopen_io)
	case OS:
		return (C.lua_CFunction)(C.luaopen_os)
	case String:
		return (C.lua_CFunction)(C.luaopen_string)
	case Math:
		return (C.lua_CFunction)(C.luaopen_math)
	case UTF8:
		return (C.lua_CFunction)(C.luaopen_utf8)
	case Debug:
		return (C.lua_CFunction)(C.luaopen_debug)
	case Bit32:
		return (C.lua_CFunction)(C.luaopen_bit32)
	case LPeg:
		return (C.lua_CFunction)(C.luaopen_lpeg)
	}
	return nil
}

func openAsset(state *State, asset string) int {
	reader := bytes.NewReader(re.MustAsset(asset))
	if err := state.Load(reader, asset, LoadAll); err != nil {
		return state.Errorf("load asset `%v` error: %v", asset, err)
	}
	state.PushValue(1)
	state.Call(1, 1)
	return 1
}
