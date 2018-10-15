package golua

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuffer(t *testing.T) {
	state := NewState()
	defer state.Close()

	expected := "abc"

	buf := state.NewBuffer()
	buf.AddChar(expected[0])
	buf.AddStringf("%c%c", expected[1], expected[2])
	buf.PushResult()

	actual := state.ToString(-1)
	require.Equal(t, expected, actual)

	state.Pop(2)
}
