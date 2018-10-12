package golua

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testRegistry(t *testing.T, r Registry) {
	require.Nil(t, r.Get(99))
	r.UnRef(99)

	expected := []byte{1, 2, 3}

	var ids []uintptr
	for _, i := range expected {
		ids = append(ids, r.Ref(i))
	}
	t.Log(r)

	for i, id := range ids {
		require.Equal(t, expected[i], r.Get(id))
	}

	require.Nil(t, r.Get(99))
	r.UnRef(99)

	r.UnRef(ids[1])
	t.Log(r)

	require.Nil(t, r.Get(ids[1]))

	id := r.Ref(4)
	t.Log(r)

	require.Equal(t, 4, r.Get(id))
}

func TestRegistry(t *testing.T) {
	t.Run("Map", func(t *testing.T) {
		testRegistry(t, &mapRegistry{})
	})
	t.Run("Slice", func(t *testing.T) {
		testRegistry(t, &sliceRegistry{})
	})
}

func benchmarkRegistryRef(b *testing.B, r Registry) {
	for i := 0; i < b.N; i++ {
		r.Ref(i)
	}
}

func benchmarkRegistryGet(b *testing.B, r Registry) {
	for i := 0; i < b.N; i++ {
		r.Get(uintptr(i))
	}
}

func benchmarkRegistryUnRef(b *testing.B, r Registry) {
	for i := 0; i < b.N; i++ {
		r.UnRef(uintptr(i))
	}
}

func benchmarkRegistryRef2(b *testing.B, r Registry) {
	for i := 0; i < b.N; i++ {
		if i%3 == 0 {
			r.UnRef(uintptr(i / 2))

		} else {
			r.Ref(i)
		}
	}
}

func BenchmarkRegistry(b *testing.B) {
	mr := &mapRegistry{}
	b.Run("Map/Ref", func(b *testing.B) {
		benchmarkRegistryRef(b, mr)
	})
	b.Run("Map/Get", func(b *testing.B) {
		benchmarkRegistryGet(b, mr)
	})
	b.Run("Map/UnRef", func(b *testing.B) {
		benchmarkRegistryUnRef(b, mr)
	})
	b.Run("Map/Ref2", func(b *testing.B) {
		benchmarkRegistryRef2(b, mr)
	})
	sr := &sliceRegistry{}
	b.Run("Slice/Ref", func(b *testing.B) {
		benchmarkRegistryRef(b, sr)
	})
	b.Run("Slice/Get", func(b *testing.B) {
		benchmarkRegistryGet(b, sr)
	})
	b.Run("Slice/UnRef", func(b *testing.B) {
		benchmarkRegistryUnRef(b, sr)
	})
	b.Run("Slice/Ref2", func(b *testing.B) {
		benchmarkRegistryRef2(b, sr)
	})
}
