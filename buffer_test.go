package golua

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuffer(t *testing.T) {
	state, err := NewState()
	require.NoError(t, err)

	expected := "abc"

	buf := state.NewBuffer()
	buf.AddChar(expected[0])
	buf.AddStringf("%c%c", expected[1], expected[2])
	buf.PushResult()

	state.PrintStackl(t)

	actual := state.ToString(-1)
	require.Equal(t, expected, actual)

	state.Pop(2)
}
