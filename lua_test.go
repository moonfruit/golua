package golua

import (
	"bytes"
	"testing"

	"github.com/moonfruit/golua/re"
	"github.com/stretchr/testify/require"
)

func TestNewState(t *testing.T) {
	state := NewState()
	defer state.Close()

	state.PushString("abc")
	t.Log("----")
	state.PrintStackl(t)

	code, err := re.Asset("re.lua")
	require.NoError(t, err)

	err = state.Load(bytes.NewReader(code), "re.lua", LoadText)
	require.NoError(t, err)

	t.Log("----")
	state.PrintStackl(t)

	err = state.pcall(0, MultiRet, 0)
	require.Error(t, err)

	t.Log("----")
	state.PrintStackl(t)
}
