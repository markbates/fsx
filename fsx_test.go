package fsx

import (
	"testing"
	"testing/fstest"

	"github.com/markbates/fsx/tfs"
	"github.com/stretchr/testify/require"
)

func ModuleFS(t *testing.T) fstest.MapFS {
	t.Helper()
	cab := fstest.MapFS{
		"module.md": tfs.File(t, "module.md"),
	}
	return cab
}

func Test_FS_Exists(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	mfs := ModuleFS(t)

	cab := NewFS(mfs)
	r.NotNil(cab)

	r.True(cab.Exists("module.md"))
	r.False(cab.Exists("404"))
}

func Test_FS_ReadFile(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	mfs := ModuleFS(t)
	cab := NewFS(mfs)
	r.NotNil(cab)

	b, err := cab.ReadFile("module.md")
	r.NoError(err)
	r.Equal(`MODULE.MD`, string(b))
}

func Test_FS_Sub(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	mfs := ModuleFS(t)
	tfs.AppendFile(t, mfs, "assets/foo.png")

	cab := NewFS(mfs)
	r.NotNil(cab)
	r.False(cab.Exists("foo.png"))

	kid, err := cab.Sub("assets")
	r.NoError(err)
	r.True(kid.Exists("foo.png"))
}

func Test_FS_MarshalJSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	mfs := ModuleFS(t)
	tfs.AppendFile(t, mfs, ".DS_Store")
	tfs.AppendFile(t, mfs, "assets/foo.png")

	cab := NewFS(mfs)
	r.NotNil(cab)

	b, err := cab.MarshalJSON()
	r.NoError(err)
	act := string(b)
	r.Contains(act, `"name":"foo.png","size":14`)
}

func Test_FS_Stat(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	mfs := ModuleFS(t)
	cab := NewFS(mfs)
	r.NotNil(cab)

	info, err := cab.Stat("module.md")
	r.NoError(err)
	r.NotNil(info)
}

func Test_Paths(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	mfs := ModuleFS(t)
	tfs.AppendFile(t, mfs, ".DS_Store")
	tfs.AppendFile(t, mfs, "assets/foo.png")

	act, err := Paths(mfs)
	r.NoError(err)

	exp := []string{
		"assets",
		"assets/foo.png",
		"module.md",
	}
	r.Equal(exp, act)
}

func Test_DirFS(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab, err := DirFS(".")
	r.NoError(err)
	r.NotNil(cab)
	r.True(cab.Exists("fsx_test.go"))
}
