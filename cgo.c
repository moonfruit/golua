#include "cgo.h"

#include <lauxlib.h>
#include <lua.h>

#include "_cgo_export.h"

static int luaGo_callGoFunction(lua_State *L);

lua_State* luaGo_main(lua_State *L) {
    if (!lua_checkstack(L, 1)) {
        return NULL;
    }
    lua_geti(L, LUA_REGISTRYINDEX, LUA_RIDX_MAINTHREAD);
    lua_State *thread = lua_tothread(L, -1);
    lua_pop(L, 1);
    return thread;
}

int luaGo_load(lua_State *L, void *data, const char *chunkname, const char *mode) {
    return lua_load(L, (lua_Reader)goReader, data, chunkname, mode);
}

int luaGo_gc(lua_State *L) {
    GoUintptr *ud = (GoUintptr *)lua_touserdata(L, 1);
    if (ud != NULL) {
        goFree(L, *ud);
    }
    return 0;
}

void luaGo_pushGoFunction(lua_State *L, unsigned long ud) {
    *(GoUintptr *)lua_newuserdata(L, sizeof(GoUintptr)) = ud;

    if (luaL_newmetatable(L, META_GO_FUNCTION)) {
        lua_pushliteral(L, "__call");
        lua_pushcfunction(L, luaGo_callGoFunction);
        lua_settable(L, -3);

        lua_pushliteral(L, "__gc");
        lua_pushcfunction(L, luaGo_gc);
        lua_settable(L, -3);
    }

    lua_setmetatable(L, -2);
}

static int luaGo_callGoFunction(lua_State *L) {
    GoUintptr *ud = (GoUintptr *)luaL_checkudata(L, 1, META_GO_FUNCTION);
    lua_remove(L, 1);
    return goCall(L, *ud);
}
