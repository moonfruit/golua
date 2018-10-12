#ifndef _cgo_h
#define _cgo_h

#include <lua.h>

#define META_GO_FUNCTION "golua.GoFunction"

// c helper functions
lua_State* luaGo_main(lua_State *L);
int luaGo_load(lua_State *L, void *data, const char *chunkname, const char *mode);
int luaGo_gc(lua_State *L);
void luaGo_pushGoFunction(lua_State *L, unsigned long ud);

#endif
