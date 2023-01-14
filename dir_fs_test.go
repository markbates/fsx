package fsx

import (
	"io"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DirFS_Open(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := DirFS(".")

	f, err := cab.Open("go.mod")
	r.NoError(err)
	defer f.Close()

	b, err := io.ReadAll(f)
	r.NoError(err)

	r.Contains(string(b), "module github.com/markbates/fsx")
}

func Test_DirFS_Create(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	temp := t.TempDir()
	cab := DirFS(temp)

	f, err := cab.Create("test.txt")
	r.NoError(err)

	_, err = f.Write([]byte("Hello World!"))
	r.NoError(err)

	err = f.Close()
	r.NoError(err)

	ff, err := cab.Open("test.txt")
	r.NoError(err)

	b, err := io.ReadAll(ff)
	r.NoError(err)

	r.Equal("Hello World!", string(b))

	panic(runtime.Version())
}
