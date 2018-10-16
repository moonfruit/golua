package golua

import (
	"testing"
)

func TestNewState(t *testing.T) {
	state := NewState()
	defer state.Close()

	state.OpenLibs()
}
