package golua

// #include <lauxlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

func (s *State) GetMetaField(idx int, name string) Type {
	cstr := C.CString(name)
	defer cfree(cstr)
	return Type(C.luaL_getmetafield(s.L, C.int(idx), cstr))
}

func (s *State) CallMeta(idx int, name string) bool {
	cstr := C.CString(name)
	defer cfree(cstr)
	return C.luaL_callmeta(s.L, C.int(idx), cstr) != 0
}

func (s *State) ToString(idx int) string {
	var length C.size_t
	str := C.luaL_tolstring(s.L, C.int(idx), &length)
	return C.GoStringN(str, C.int(length))
}

func (s *State) ArgError(arg int, args ...interface{}) int {
	return s.argError(arg, fmt.Sprint(args...))
}

func (s *State) ArgErrorf(arg int, format string, args ...interface{}) int {
	return s.argError(arg, fmt.Sprintf(format, args...))
}

func (s *State) argError(arg int, msg string) int {
	cstr := C.CString(msg)
	defer cfree(cstr)
	return int(C.luaL_argerror(s.L, C.int(arg), cstr))
}

func (s *State) ArgCheck(cond bool, arg int, args ...interface{}) {
	if !cond {
		s.ArgError(arg, args...)
	}
}

func (s *State) ArgCheckf(cond bool, arg int, format string, args ...interface{}) {
	if !cond {
		s.ArgErrorf(arg, format, args...)
	}
}

func (s *State) CheckString(idx int) string {
	var length C.size_t
	str := C.luaL_checklstring(s.L, C.int(idx), &length)
	return C.GoStringN(str, C.int(length))
}

func (s *State) OptString(idx int, def string) string {
	if s.IsNoneOrNil(idx) {
		return def
	}
	return s.CheckString(idx)
}

func (s *State) CheckNumber(idx int) Number {
	return Number(C.luaL_checknumber(s.L, C.int(idx)))
}

func (s *State) OptNumber(idx int, def Number) Number {
	if s.IsNoneOrNil(idx) {
		return def
	}
	return s.CheckNumber(idx)
}

func (s *State) CheckInteger(idx int) Integer {
	return Integer(C.luaL_checkinteger(s.L, C.int(idx)))
}

func (s *State) OptInteger(idx int, def Integer) Integer {
	if s.IsNoneOrNil(idx) {
		return def
	}
	return s.CheckInteger(idx)
}

func (s *State) CheckStack(size int, args ...interface{}) {
	var cstr *C.char
	if len(args) != 0 {
		cstr = C.CString(fmt.Sprint(args...))
		defer cfree(cstr)
	}
	C.luaL_checkstack(s.L, C.int(size), cstr)
}

func (s *State) CheckStackf(size int, format string, args ...interface{}) {
	cstr := C.CString(fmt.Sprintf(format, args...))
	defer cfree(cstr)
	C.luaL_checkstack(s.L, C.int(size), cstr)
}

func (s *State) CheckType(idx int, ty Type) {
	C.luaL_checktype(s.L, C.int(idx), C.int(ty))
}

func (s *State) CheckAny(idx int) {
	C.luaL_checkany(s.L, C.int(idx))
}

func (s *State) NewMetaTable(name string) bool {
	cstr := C.CString(name)
	defer cfree(cstr)
	return C.luaL_newmetatable(s.L, cstr) != 0
}

func (s *State) SetMetaTable(name string) {
	cstr := C.CString(name)
	defer cfree(cstr)
	C.luaL_setmetatable(s.L, cstr)
}

func (s *State) GetMetaTable(name string) Type {
	return s.GetField(RegistryIndex, name)
}

func (s *State) TestUserData(idx int, name string) bool {
	cstr := C.CString(name)
	defer cfree(cstr)
	return C.luaL_testudata(s.L, C.int(idx), cstr) != nil
}

func (s *State) CheckUserData(idx int, name string) {
	cstr := C.CString(name)
	defer cfree(cstr)
	C.luaL_checkudata(s.L, C.int(idx), cstr)
}

func (s *State) Where(lvl int) {
	C.luaL_where(s.L, C.int(lvl))
}

func (s *State) Error(args ...interface{}) int {
	return s.error(fmt.Sprint(args...))
}

func (s *State) Errorf(format string, args ...interface{}) int {
	return s.error(fmt.Sprintf(format, args...))
}

func (s *State) error(msg string) int {
	s.Where(1)
	s.PushString(msg)
	s.Concat(2)
	return s.RawError()
}

func (s *State) CheckOption(idx int, def string, opts []string) int {
	str := s.OptString(idx, def)

	for i, opt := range opts {
		if str == opt {
			return i
		}
	}

	return s.ArgError(idx)
}

func (s *State) Ref(idx int) int {
	return int(C.luaL_ref(s.L, C.int(idx)))
}

func (s *State) UnRef(idx, ref int) {
	C.luaL_unref(s.L, C.int(idx), C.int(ref))
}

func (s *State) LoadFile(filename string) error {
	return s.LoadFileX(filename, LoadAll)
}

func (s *State) LoadFileX(filename string, mode LoadMode) error {
	cstr := C.CString(filename)
	defer cfree(cstr)
	return fromLua(C.luaL_loadfilex(s.L, cstr, mode.mode()))
}

func (s *State) LoadBytes(name string, buf []byte) error {
	return s.LoadBytesX(name, buf, LoadAll)
}

func (s *State) LoadBytesX(name string, buf []byte, mode LoadMode) error {
	cstr := C.CString(name)
	defer cfree(cstr)
	return fromLua(C.luaL_loadbufferx(s.L, (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)), cstr, mode.mode()))
}

func (s *State) LoadString(str string) error {
	cstr := C.CString(str)
	defer cfree(cstr)
	return fromLua(C.luaL_loadbufferx(s.L, cstr, C.size_t(len(str)), cstr, nil))
}

func (s *State) DoFile(filename string) error {
	if err := s.LoadFile(filename); err != nil {
		return err
	}
	return s.PCall(0, MultiRet, 0)
}

func (s *State) DoString(str string) error {
	if err := s.LoadString(str); err != nil {
		return err
	}
	return s.PCall(0, MultiRet, 0)
}

func (s *State) Length(idx int) Integer {
	return Integer(C.luaL_len(s.L, C.int(idx)))
}

func (s *State) EnsureTable(idx int, name string) bool {
	cstr := C.CString(name)
	defer cfree(cstr)
	return C.luaL_getsubtable(s.L, C.int(idx), cstr) != 0
}

// TODO: void (luaL_traceback) (lua_State *L, lua_State *L1, const char *msg, int level);

// TODO: int luaL_fileresult (lua_State *L, int stat, const char *fname);
// TODO: int luaL_execresult (lua_State *L, int stat);

// SKIP: void luaL_checkversion (lua_State *L);
// SKIP: const char *luaL_gsub (lua_State *L, const char *s, const char *p, const char *r);
