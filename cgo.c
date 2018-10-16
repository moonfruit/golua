#include "cgo.h"

#include <lauxlib.h>
#include <lua.h>
#include <lualib.h>

#include "_cgo_export.h"

static int luaGo_GoValueEqual(lua_State *L);
static int luaGo_GoValueToString(lua_State *L);
static int luaGo_GoValueClose(lua_State *L);

static int luaGo_upValueCount(lua_State *L);
static int luaGo_getUpValueCount(lua_State *L, int idx);

static const luaL_Reg loadedlibs[] = {{"_G", luaopen_base},
                                      {LUA_LOADLIBNAME, luaopen_package},
                                      {LUA_COLIBNAME, luaopen_coroutine},
                                      {LUA_TABLIBNAME, luaopen_table},
                                      {LUA_STRLIBNAME, luaopen_string},
                                      {LUA_MATHLIBNAME, luaopen_math},
                                      {LUA_UTF8LIBNAME, luaopen_utf8},
                                      {LUA_DBLIBNAME, luaopen_debug},
                                      {LUA_BITLIBNAME, luaopen_bit32},
                                      {NULL, NULL}};

void luaGo_openBasicLibs(lua_State *L) {
    const luaL_Reg *lib;
    /* "require" functions from 'loadedlibs' and set results to global table */
    for (lib = loadedlibs; lib->func; lib++) {
        luaL_requiref(L, lib->name, lib->func, 1);
        lua_pop(L, 1); /* remove lib */
    }
}

void luaGo_preload(lua_State *L, const char *modname, lua_CFunction f) {
    luaL_getsubtable(L, LUA_REGISTRYINDEX, LUA_PRELOAD_TABLE);
    lua_pushcfunction(L, f);
    lua_setfield(L, -2, modname);
    lua_pop(L, 1); // remove PRELOAD table
}

lua_State *luaGo_main(lua_State *L) {
    if (!lua_checkstack(L, 1)) {
        return NULL;
    }
    lua_geti(L, LUA_REGISTRYINDEX, LUA_RIDX_MAINTHREAD);
    lua_State *thread = lua_tothread(L, -1);
    lua_pop(L, 1);
    return thread;
}

int luaGo_call(lua_State *L, lua_CFunction f) { return f(L); }

int luaGo_load(lua_State *L, void *data, const char *chunkname, const char *mode) {
    return lua_load(L, (lua_Reader)goReader, data, chunkname, mode);
}

void luaGo_pushGoValue(lua_State *L, unsigned long ud) {
    *(GoUintptr *)lua_newuserdata(L, sizeof(GoUintptr)) = ud;

    if (luaL_newmetatable(L, META_GO_VALUE)) {
        lua_pushliteral(L, "__gc");
        lua_pushcfunction(L, luaGo_GoValueClose);
        lua_settable(L, -3);

        lua_pushliteral(L, "__eq");
        lua_pushcfunction(L, luaGo_GoValueEqual);
        lua_settable(L, -3);

        lua_pushliteral(L, "__tostring");
        lua_pushcfunction(L, luaGo_GoValueToString);
        lua_settable(L, -3);
    }

    lua_setmetatable(L, -2);
}

static int luaGo_GoValueEqual(lua_State *L) {
    GoUintptr *ud1 = (GoUintptr *)luaL_testudata(L, 1, META_GO_VALUE);
    if (ud1 == NULL) {
        lua_pushboolean(L, 0);
        return 1;
    }

    GoUintptr *ud2 = (GoUintptr *)luaL_testudata(L, 2, META_GO_VALUE);
    if (ud2 == NULL) {
        lua_pushboolean(L, 0);
        return 1;
    }

    lua_pushboolean(L, *ud1 == *ud2);
    return 1;
}

static int luaGo_GoValueToString(lua_State *L) {
    GoUintptr *ud = (GoUintptr *)luaL_checkudata(L, 1, META_GO_VALUE);
    lua_pushfstring(L, META_GO_VALUE ": %p", *ud);
    return 1;
}

static int luaGo_GoValueClose(lua_State *L) {
    GoUintptr *ud = (GoUintptr *)luaL_testudata(L, 1, META_GO_VALUE);
    if (ud != NULL) {
        goFree(L, *ud);
    }
    return 0;
}

int luaGo_callGoFunction(lua_State *L) {
    int idx = lua_upvalueindex(luaGo_upValueCount(L));
    GoUintptr *ud = (GoUintptr *)luaL_checkudata(L, idx, META_GO_VALUE);
    return goCall(L, *ud);
}

int luaGo_getGoFunction(lua_State *L, int idx) {
    if (lua_tocfunction(L, idx) != luaGo_callGoFunction) {
        return 0;
    }

    int value = luaGo_getUpValueCount(L, idx);
    if (value == 0) {
        return 0;
    }

    if (lua_getupvalue(L, idx, value) == NULL) {
        return 0;
    }

    return 1;
}

static int luaGo_upValueCount(lua_State *L) {
    lua_Debug ar;

    if (!lua_getstack(L, 0, &ar)) {
        return 0;
    }
    if (!lua_getinfo(L, "u", &ar)) {
        return 0;
    }

    return ar.nups;
}

static int luaGo_getUpValueCount(lua_State *L, int idx) {
    lua_Debug ar;

    lua_pushvalue(L, idx);
    if (!lua_getinfo(L, ">u", &ar)) {
        return 0;
    }

    return ar.nups;
}
