package golua

// #include <lua.h>
// #include "cgo.h"
import "C"
import (
	"fmt"
	"io"
	"unsafe"
)

/*
 * State manipulation
 */

// TODO: lua_State *(lua_newthread) (lua_State *L);
// TODO: lua_CFunction (lua_atpanic) (lua_State *L, lua_CFunction panicf);

/*
 * Basic stack manipulation
 */

// AbsIndex converts the acceptable index idx into an equivalent absolute index
// (that is, one that does not depend on the stack top).
func (s *State) AbsIndex(idx int) int {
	return int(C.lua_absindex(s.L, C.int(idx)))
}

func (s *State) GetTop() int {
	return int(C.lua_gettop(s.L))
}

func (s *State) Pop(n int) {
	s.SetTop(-n - 1)
}

func (s *State) SetTop(idx int) {
	C.lua_settop(s.L, C.int(idx))
}

func (s *State) PushValue(idx int) {
	C.lua_pushvalue(s.L, C.int(idx))
}

func (s *State) Insert(idx int) {
	s.Rotate(idx, 1)
}

func (s *State) Rotate(idx, n int) {
	C.lua_rotate(s.L, C.int(idx), C.int(n))
}

func (s *State) Remove(idx int) {
	s.Rotate(idx, -1)
	s.Pop(1)
}

func (s *State) Replace(idx int) {
	s.Copy(-1, idx)
	s.Pop(1)
}

func (s *State) Copy(fromIdx, toIdx int) {
	C.lua_copy(s.L, C.int(fromIdx), C.int(toIdx))
}

// FIXME: CheckStack
func (s *State) RawCheckStack(n int) bool {
	return C.lua_checkstack(s.L, C.int(n)) != 0
}

// TODO: void (lua_xmove) (lua_State *from, lua_State *to, int n);

/*
 * Access functions  (stack -> Go)
 */

func (s *State) IsNumber(idx int) bool {
	return C.lua_isnumber(s.L, C.int(idx)) != 0
}

func (s *State) IsString(idx int) bool {
	return C.lua_isstring(s.L, C.int(idx)) != 0
}

// FIXME: IsGoFunction ??
func (s *State) IsCFunction(idx int) bool {
	return C.lua_iscfunction(s.L, C.int(idx)) != 0
}

func (s *State) IsInteger(idx int) bool {
	return C.lua_isinteger(s.L, C.int(idx)) != 0
}

func (s *State) IsUserData(idx int) bool {
	return C.lua_isuserdata(s.L, C.int(idx)) != 0
}

func (s *State) IsFunction(idx int) bool {
	return s.Type(idx) == TypeFunction
}

func (s *State) IsTable(idx int) bool {
	return s.Type(idx) == TypeTable
}

func (s *State) IsLightUserData(idx int) bool {
	return s.Type(idx) == TypeLightUserData
}

func (s *State) IsNil(idx int) bool {
	return s.Type(idx) == TypeNil
}

func (s *State) IsBoolean(idx int) bool {
	return s.Type(idx) == TypeBoolean
}

func (s *State) IsThread(idx int) bool {
	return s.Type(idx) == TypeThread
}

func (s *State) IsNone(idx int) bool {
	return s.Type(idx) == TypeNone
}

func (s *State) IsNoneOrNil(idx int) bool {
	ty := s.Type(idx)
	return ty == TypeNone || ty == TypeNil
}

func (s *State) Type(idx int) Type {
	return Type(C.lua_type(s.L, C.int(idx)))
}

func (s *State) ToNumber(idx int) Number {
	return Number(C.lua_tonumberx(s.L, C.int(idx), nil))
}

func (s *State) ToNumberX(idx int) (number Number, ok bool) {
	var flag C.int
	n := C.lua_tonumberx(s.L, C.int(idx), &flag)
	return Number(n), flag != 0
}

func (s *State) ToInteger(idx int) Integer {
	return Integer(C.lua_tointegerx(s.L, C.int(idx), nil))
}

func (s *State) ToIntegerX(idx int) (number Integer, ok bool) {
	var flag C.int
	n := C.lua_tointegerx(s.L, C.int(idx), &flag)
	return Integer(n), flag != 0
}

func (s *State) ToBoolean(idx int) bool {
	return C.lua_toboolean(s.L, C.int(idx)) != 0
}

func (s *State) RawToString(idx int) string {
	var length C.size_t
	str := C.lua_tolstring(s.L, C.int(idx), &length)
	return C.GoStringN(str, C.int(length))
}

func (s *State) RawLen(idx int) uint {
	return uint(C.lua_rawlen(s.L, C.int(idx)))
}

// TODO: lua_CFunction (lua_tocfunction) (lua_State *L, int idx);
// TODO: void *(lua_touserdata) (lua_State *L, int idx);
// TODO: lua_State *(lua_tothread) (lua_State *L, int idx);
// TODO: const void *(lua_topointer) (lua_State *L, int idx);

/*
 * Comparison and arithmetic functions
 */

func (s *State) Arith(op ArithOp) {
	C.lua_arith(s.L, C.int(op))
}

func (s *State) RawEqual(idx1, idx2 int) bool {
	return C.lua_rawequal(s.L, C.int(idx1), C.int(idx2)) != 0
}

func (s *State) Compare(idx1, idx2 int, op CompareOp) bool {
	return C.lua_compare(s.L, C.int(idx1), C.int(idx2), C.int(op)) != 0
}

/*
 * Push functions (C -> stack)
 */

func (s *State) PushNil() {
	C.lua_pushnil(s.L)
}

func (s *State) PushNumber(n Number) {
	C.lua_pushnumber(s.L, C.lua_Number(n))
}

func (s *State) PushInteger(n Integer) {
	C.lua_pushinteger(s.L, C.lua_Integer(n))
}

func (s *State) PushStringf(format string, args ...interface{}) {
	s.PushString(fmt.Sprintf(format, args...))
}

func (s *State) PushString(str string) {
	s.pushChars(C._GoStringPtr(str), len(str))
}

func (s *State) PushBytes(b []byte) {
	s.pushChars((*C.char)(unsafe.Pointer(&b[0])), len(b))
}

func (s *State) pushChars(chars *C.char, len int) {
	C.lua_pushlstring(s.L, chars, C.size_t(len))
}

// TODO: void (lua_pushcclosure) (lua_State *L, lua_CFunction fn, int n);

// TODO: void lua_pushcfunction (lua_State *L, lua_CFunction f);

func (s *State) PushBoolean(b bool) {
	var val C.int
	if b {
		val = 1
	}
	C.lua_pushboolean(s.L, val)
}

// TODO: void (lua_pushlightuserdata) (lua_State *L, void *p);

func (s *State) PushThread() (isMain bool) {
	main := C.lua_pushthread(s.L)
	return main != 0
}

func (s *State) PushGlobalTable() {
	s.RawGetI(RegistryIndex, RegistryGlobals)
}

/*
 * Get functions (Lua -> stack)
 */

func (s *State) GetGlobal(name string) Type {
	cstr := C.CString(name)
	defer cfree(cstr)
	return Type(C.lua_getglobal(s.L, cstr))
}

func (s *State) GetTable(idx int) Type {
	return Type(C.lua_gettable(s.L, C.int(idx)))
}

func (s *State) GetField(idx int, name string) Type {
	cstr := C.CString(name)
	defer cfree(cstr)
	return Type(C.lua_getfield(s.L, C.int(idx), cstr))
}

func (s *State) GetI(idx int, n Integer) Type {
	return Type(C.lua_geti(s.L, C.int(idx), C.lua_Integer(n)))
}

func (s *State) RawGet(idx int) Type {
	return Type(C.lua_rawget(s.L, C.int(idx)))
}

func (s *State) RawGetI(idx int, n Integer) Type {
	return Type(C.lua_rawgeti(s.L, C.int(idx), C.lua_Integer(n)))
}

// TODO: int (lua_rawgetp) (lua_State *L, int idx, const void *p);

func (s *State) NewTable() {
	s.CreateTable(0, 0)
}

func (s *State) CreateTable(nArr, nRec int) {
	C.lua_createtable(s.L, C.int(nArr), C.int(nRec))
}

// TODO: void *(lua_newuserdata) (lua_State *L, size_t sz);

func (s *State) RawGetMetaTable(idx int) (ok bool) {
	flag := C.lua_getmetatable(s.L, C.int(idx))
	return flag != 0
}

func (s *State) GetUserValue(idx int) Type {
	return Type(C.lua_getuservalue(s.L, C.int(idx)))
}

/*
 * Set functions (stack -> Lua)
 */

func (s *State) SetGlobal(name string) {
	cstr := C.CString(name)
	defer cfree(cstr)
	C.lua_setglobal(s.L, cstr)
}

func (s *State) SetTable(idx int) {
	C.lua_settable(s.L, C.int(idx))
}

func (s *State) SetField(idx int, name string) {
	cstr := C.CString(name)
	defer cfree(cstr)
	C.lua_setfield(s.L, C.int(idx), cstr)
}

func (s *State) SetI(idx int, n Integer) {
	C.lua_seti(s.L, C.int(idx), C.lua_Integer(n))
}

func (s *State) RawSet(idx int) {
	C.lua_rawset(s.L, C.int(idx))
}

func (s *State) RawSetI(idx int, n Integer) {
	C.lua_rawseti(s.L, C.int(idx), C.lua_Integer(n))
}

// TODO: void (lua_rawsetp) (lua_State *L, int idx, const void *p);

func (s *State) RawSetMetaTable(idx int) {
	C.lua_setmetatable(s.L, C.int(idx))
}

// FIXME: What does for ??
func (s *State) SetUserValue(idx int) {
	C.lua_setuservalue(s.L, C.int(idx))
}

/*
 * Load and run Lua code: 'load' and 'call' functions
 */

// TODO: void (lua_callk) (lua_State *L, int nargs, int nresults, lua_KContext ctx, lua_KFunction k);

func (s *State) Call(nArgs, nResults int) {
	C.lua_callk(s.L, C.int(nArgs), C.int(nResults), 0, nil)
}

// TODO: int (lua_pcallk) (lua_State *L, int nargs, int nresults, int errfunc, lua_KContext ctx, lua_KFunction k);

// FIXME: PCall()
func (s *State) pcall(nArgs, nResults, msgHandler int) error {
	return fromLua(C.lua_pcallk(s.L, C.int(nArgs), C.int(nResults), C.int(msgHandler), 0, nil))
}

func (s *State) Load(reader io.Reader, name string, mode LoadMode) error {
	cstr := C.CString(name)
	defer cfree(cstr)

	ctx := newReaderCtx(reader)
	id := s.RefGoValue(ctx)
	defer s.UnRefGoValue(id)

	err := fromLua(C.luaGo_load(s.L, unsafe.Pointer(id), cstr, mode.mode()))
	if ctx.err != nil && ctx.err != io.EOF {
		return ctx.err
	}

	return err
}

// TODO: int (lua_dump) (lua_State *L, lua_Writer writer, void *data, int strip);

/*
 * Coroutine functions
 */

// TODO: int (lua_yieldk) (lua_State *L, int nresults, lua_KContext ctx, lua_KFunction k);

func (s *State) Yield(nResults int) int {
	return int(C.lua_yieldk(s.L, C.int(nResults), 0, nil))
}

// TODO: int (lua_resume) (lua_State *L, lua_State *from, int narg);

func (s *State) Status() error {
	return fromLua(C.lua_status(s.L))
}

func (s *State) IsYieldable() bool {
	return C.lua_isyieldable(s.L) != 0
}

/*
 * Garbage-collection functionx
 */

func (s *State) GC(what GcOption, data int) int {
	return int(C.lua_gc(s.L, C.int(what), C.int(data)))
}

/*
 * Miscellaneous functions
 */

func (s *State) RawError() int {
	return int(C.lua_error(s.L))
}

func (s *State) Next(idx int) int {
	return int(C.lua_next(s.L, C.int(idx)))
}

func (s *State) Concat(n int) {
	C.lua_concat(s.L, C.int(n))
}

func (s *State) Len(idx int) {
	C.lua_len(s.L, C.int(idx))
}

func (s *State) StringToNumber(str string) uint {
	cstr := C.CString(str)
	defer cfree(cstr)
	return uint(C.lua_stringtonumber(s.L, cstr))
}

/*
 * Other
 */

// TODO: void *lua_getextraspace (lua_State *L);

// TODO: void lua_register (lua_State *L, const char *name, lua_CFunction f);

// SKIP: lua_State *lua_newstate (lua_Alloc f, void *ud);
// SKIP: lua_Alloc lua_getallocf (lua_State *L, void **ud);
// SKIP: void lua_setallocf (lua_State *L, lua_Alloc f, void *ud);
// SKIP: int lua_numbertointeger (lua_Number n, lua_Integer *p);
// SKIP: const char *lua_pushliteral (lua_State *L, const char *s);
