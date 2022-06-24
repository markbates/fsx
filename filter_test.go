package fsx

import (
	"io/fs"
	"sort"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func Test_Filter(t *testing.T) {
	t.Parallel()

	files := []string{
		"module.md",
		"assets/foo.png",
		"src/foo.go",
		"sub/module.md",
		"sub/assets/bar.png",
		"sub/src/bar.go",
		"js/foo.js",
		"mysrcfolder/foo.go",
		"thesrc/foo.go",
		"src-files/foo.go",
	}

	cab := fstest.MapFS{}

	for _, file := range files {
		cab[file] = &fstest.MapFile{
			Data: []byte(file),
		}
	}

	tcs := []struct {
		in  string
		out []string
		err bool
	}{
		{in: "", out: files},
		{in: `(^|/)src/`, out: []string{"src/foo.go", "sub/src/bar.go"}},
	}

	for _, tc := range tcs {
		name := tc.in
		if len(name) == 0 {
			name = "empty"
		}

		t.Run(name, func(t *testing.T) {
			r := require.New(t)

			var act []string

			err := Filter(cab, tc.in, func(path string, d fs.DirEntry, err error) error {

				if err != nil {
					return err
				}

				if d.IsDir() {
					return nil
				}

				act = append(act, path)

				return nil
			})

			if tc.err {
				r.Error(err)
			}

			r.NoError(err)

			sort.Strings(act)
			sort.Strings(tc.out)

			r.Equal(tc.out, act)
		})
	}
}
