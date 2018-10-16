package golua

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLpeg(t *testing.T) {
	state := NewState()
	defer state.Close()

	state.OpenLibs()

	filename := "lpeg/test.lua"
	err := state.LoadFile(filename)
	require.NoError(t, err)

	state.Call(0, 0)
}
