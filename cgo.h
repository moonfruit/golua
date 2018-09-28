#ifndef _cgo_h
#define _cgo_h

#include <lua.h>

int load(lua_State *L, char *ctx, const char *name, const char *mode);

const char *goReader(lua_State *L, void *ud, size_t *sz);

#endif
