package golua

import (
	"bytes"
	"testing"

	"github.com/moonfruit/golua/re"
	"github.com/stretchr/testify/require"
)

func TestNewState(t *testing.T) {
	state, err := NewState()
	require.NoError(t, err)

	code, err := re.Asset("re.lua")
	require.NoError(t, err)

	err = state.Load(bytes.NewReader(code), "re.lua", LoadText)
	require.NoError(t, err)

	state.PrintStack()

	//state.Call(0, 1)
	//state.PrintStack()

	state.Close()
}
