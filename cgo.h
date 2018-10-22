#ifndef _cgo_h
#define _cgo_h

#include <lua.h>

#define META_GO_VALUE "golua.GoValue"
#define META_GO_FUNCTION "golua.GoFunction"

// c helper functions
void luaGo_openBasicLibs(lua_State *L);
void luaGo_preload(lua_State *L, const char *modname, lua_CFunction f);

lua_State *luaGo_main(lua_State *L);

int luaGo_call(lua_State *L, lua_CFunction f);

int luaGo_load(lua_State *L, void *data, const char *chunkname, const char *mode);

void luaGo_pushGoValue(lua_State *L, uintptr_t ud);

int luaGo_callGoFunction(lua_State *L);
int luaGo_getGoFunction(lua_State *L, int idx);

#endif
