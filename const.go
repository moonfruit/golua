package golua

/*
#include <lua.h>
#include <lauxlib.h>

const char *empty = "";
const char *modeB = "b";
const char *modeT = "t";
const char *modeBT = "bt";
*/
import "C"
import (
	"fmt"
)

type Integer = int64
type Number = float64

var empty = C.empty

// MultiRet is the option for multiple returns in `Call()` and `PCall()`.
const MultiRet = C.LUA_MULTRET

const (
	RegistryIndex      = C.LUA_REGISTRYINDEX
	RegistryMainThread = C.LUA_RIDX_MAINTHREAD
	RegistryGlobals    = C.LUA_RIDX_GLOBALS
)

const (
	KeyLoadedTable  = "_LOADED"
	KeyPreLoadTable = "_PRELOAD"
)

func UpValue(idx int) int {
	return RegistryIndex - idx
}

type Status int

//noinspection GoVarAndConstTypeMayBeOmitted
const (
	statusOK    Status = C.LUA_OK
	StatusYield Status = C.LUA_YIELD
	ErrRun      Status = C.LUA_ERRRUN
	ErrSyntax   Status = C.LUA_ERRSYNTAX
	ErrMem      Status = C.LUA_ERRMEM
	ErrGcMM     Status = C.LUA_ERRGCMM
	ErrErr      Status = C.LUA_ERRERR
	ErrFile     Status = C.LUA_ERRFILE

	ErrInvalidBuffer Status = 100
)

func fromLua(st C.int) error {
	s := Status(st)
	if s == statusOK {
		return nil
	}
	return s
}

func (s Status) Error() string {
	switch s {
	case statusOK:
		return "success"
	case StatusYield:
		return "yield"
	case ErrRun:
		return "runtime error"
	case ErrSyntax:
		return "syntax error"
	case ErrMem:
		return "memory allocation error"
	case ErrGcMM:
		return "error while running a __gc metamethod"
	case ErrErr:
		return "error while running the message handler"
	case ErrFile:
		return "file error"
	case ErrInvalidBuffer:
		return "invalid buffer"
	default:
		return fmt.Sprintf("unknown error `%d`", s)
	}
}

type Type int

//noinspection GoVarAndConstTypeMayBeOmitted
const (
	TypeNone          Type = C.LUA_TNONE
	TypeNil           Type = C.LUA_TNIL
	TypeBoolean       Type = C.LUA_TBOOLEAN
	TypeLightUserData Type = C.LUA_TLIGHTUSERDATA
	TypeNumber        Type = C.LUA_TNUMBER
	TypeString        Type = C.LUA_TSTRING
	TypeTable         Type = C.LUA_TTABLE
	TypeFunction      Type = C.LUA_TFUNCTION
	TypeUserData      Type = C.LUA_TUSERDATA
	TypeThread        Type = C.LUA_TTHREAD
)

func (ty Type) String() string {
	return C.GoString(C.lua_typename(nil, C.int(ty)))
}

type ArithOp int

//noinspection GoVarAndConstTypeMayBeOmitted
const (
	OPADD  ArithOp = C.LUA_OPADD
	OPSUB  ArithOp = C.LUA_OPSUB
	OPMUL  ArithOp = C.LUA_OPMUL
	OPMOD  ArithOp = C.LUA_OPMOD
	OPPOW  ArithOp = C.LUA_OPPOW
	OPDIV  ArithOp = C.LUA_OPDIV
	OPIDIV ArithOp = C.LUA_OPIDIV
	OPBAND ArithOp = C.LUA_OPBAND
	OPBOR  ArithOp = C.LUA_OPBOR
	OPBXOR ArithOp = C.LUA_OPBXOR
	OPSHL  ArithOp = C.LUA_OPSHL
	OPSHR  ArithOp = C.LUA_OPSHR
	OPUNM  ArithOp = C.LUA_OPUNM
	OPBNOT ArithOp = C.LUA_OPBNOT
)

type CompareOp int

//noinspection GoVarAndConstTypeMayBeOmitted
const (
	OPEQ CompareOp = C.LUA_OPEQ
	OPLT CompareOp = C.LUA_OPLT
	OPLE CompareOp = C.LUA_OPLE
)

type GcOption int

//noinspection GoVarAndConstTypeMayBeOmitted
const (
	GcStop       GcOption = C.LUA_GCSTOP
	GcRestart    GcOption = C.LUA_GCRESTART
	GcCollect    GcOption = C.LUA_GCCOLLECT
	GcCount      GcOption = C.LUA_GCCOUNT
	GcCountB     GcOption = C.LUA_GCCOUNTB
	GcStep       GcOption = C.LUA_GCSTEP
	GcSetPause   GcOption = C.LUA_GCSETPAUSE
	GcSetStepMul GcOption = C.LUA_GCSETSTEPMUL
	GcIsRunning  GcOption = C.LUA_GCISRUNNING
)

type LoadMode int

const (
	LoadBinary LoadMode = 1
	LoadText   LoadMode = 1 << 1
	LoadAll             = LoadBinary | LoadText
)

func (m LoadMode) mode() *C.char {
	if m&LoadBinary != 0 {
		if m&LoadText != 0 {
			return nil
		} else {
			return C.modeB
		}

	} else if m&LoadText != 0 {
		return C.modeT
	}

	return empty
}
