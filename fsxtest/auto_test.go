package fsxtest

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AutoFile(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	f := AutoFile(t, "test.txt")
	r.NotNil(f)

	r.Equal("TEST.TXT", string(f.Data))
}

func Test_AutoFS(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	files := []string{
		"a.txt",
		"b.txt",
		"c.txt",
	}

	fs := AutoFS(t, files...)
	r.NotNil(fs)

	for _, n := range files {
		f, ok := fs[n]
		r.True(ok)
		r.Equal(strings.ToUpper(n), string(f.Data))
	}
}

func Test_AutoMod(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	f := AutoMod(t, "demo")
	r.NotNil(f)

	r.Contains(string(f.Data), "module demo\n\ngo 1.19")
}
