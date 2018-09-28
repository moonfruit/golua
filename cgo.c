#include "cgo.h"

int load(lua_State *L, char *ctx, const char *name, const char *mode) {
    return lua_load(L, goReader, ctx, name, mode);
}
