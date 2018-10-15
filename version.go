package golua

import "fmt"

var (
	luaVersion  = Version{5, 3, 5}
	lpegVersion = Version{1, 0, 1}
)

type Version struct {
	Major, Minor, Release int
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Release)
}

func LuaVersion() Version {
	return luaVersion
}

func LPegVersion() Version {
	return lpegVersion
}
