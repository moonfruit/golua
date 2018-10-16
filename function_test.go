package golua

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoFunction(t *testing.T) {
	fun := func(state *State) int {
		state.CheckInteger(1)
		return 1
	}

	state := NewState()
	defer state.Close()

	expected := int64(10)
	state.PushGoFunction(fun)
	state.PushInteger(expected)
	state.Call(1, 1)

	actual, ok := state.ToIntegerX(-1)
	require.True(t, ok)
	require.Equal(t, expected, actual)
	state.Pop(1)

	state.PushGoFunction(fun)
	fun, ok = state.ToGoFunction(-1)
	require.True(t, ok)
	state.Pop(1)

	state.PushInteger(expected)
	fun(state)
	actual, ok = state.ToIntegerX(-1)
	require.True(t, ok)
	require.Equal(t, expected, actual)
	state.Pop(1)
}
