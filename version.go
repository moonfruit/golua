package golua

import "fmt"

var (
	Version     = version{5, 3, 5}
	LPegVersion = version{1, 0, 1}
)

type version struct {
	Major, Minor, Release int
}

func (v version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Release)
}
