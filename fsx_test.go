package fsx

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func Test_FS_Open(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := &FS{
		Backing: fstest.MapFS{
			"module.md": &fstest.MapFile{
				Data: []byte("# Module"),
			},
		},
		extras: fstest.MapFS{
			"assets/foo.png": &fstest.MapFile{
				Data: []byte("foo"),
			},
		},
	}

	b, err := fs.ReadFile(cab, "module.md")
	r.NoError(err)

	r.Equal([]byte("# Module"), b)

	b, err = fs.ReadFile(cab, "assets/foo.png")
	r.NoError(err)

	r.Equal([]byte("foo"), b)

	_, err = fs.ReadFile(cab, "src/foo.go")
	r.Error(err)
}

func Test_FS_Write(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := &FS{
		Backing: fstest.MapFS{},
	}

	err := cab.WriteFile("module.md", []byte("# Module"))
	r.NoError(err)

	b, err := fs.ReadFile(cab, "module.md")
	r.NoError(err)

	r.Equal([]byte("# Module"), b)

}

func Test_FS_ReadDir_Walk(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := &FS{
		Backing: fstest.MapFS{
			"module.md": &fstest.MapFile{
				Data: []byte("# Module"),
			},
			"assets/foo.png": &fstest.MapFile{
				Data: []byte("foo"),
			},
		},
		extras: fstest.MapFS{
			"src/foo.go": &fstest.MapFile{
				Data: []byte("package foo"),
			},
		},
	}

	var walked []string

	err := fs.WalkDir(cab, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		walked = append(walked, path)

		return nil
	})

	r.NoError(err)

	exp := []string{"assets/foo.png", "module.md", "src/foo.go"}
	r.Len(walked, len(exp))

	r.Equal(exp, walked)
}

func Test_FS_Stat(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := &FS{
		Backing: fstest.MapFS{
			"module.md": &fstest.MapFile{
				Data: []byte("# Module"),
			},
		},
		extras: fstest.MapFS{
			"src/foo.go": &fstest.MapFile{
				Data: []byte("package foo"),
			},
		},
	}

	fi, err := fs.Stat(cab, "module.md")
	r.NoError(err)
	r.NotNil(fi)

	fi, err = fs.Stat(cab, "src/foo.go")
	r.NoError(err)
	r.NotNil(fi)

	_, err = fs.Stat(cab, "assets/foo.png")
	r.Error(err)
}
