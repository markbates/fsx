package fsx

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func Test_ArrayFS_Open(t *testing.T) {
	t.Parallel()

	cab := &ArrayFS{}
	cab.Append(fstest.MapFS{
		"a/a.txt":  &fstest.MapFile{},
		"a/aa.txt": &fstest.MapFile{},
	})
	cab.Append(fstest.MapFS{
		"b/b.txt":  &fstest.MapFile{},
		"b/bb.txt": &fstest.MapFile{},
	})
	cab.Append(fstest.MapFS{
		"c/c.txt":  &fstest.MapFile{},
		"c/cc.txt": &fstest.MapFile{},
	})

	tcs := []struct {
		name string
		fn   string
		err  error
	}{
		{
			name: "known file",
			fn:   "b/b.txt",
			err:  nil,
		},
		{
			name: "unknown file",
			fn:   "d/d.txt",
			err:  fs.ErrNotExist,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			f, err := cab.Open(tc.fn)
			if tc.err != nil {
				r.Error(err)
				r.True(errors.Is(err, tc.err))
				return
			}

			r.NoError(err)
			defer f.Close()

		})
	}
}

func Test_ArrayFS_ReadDir(t *testing.T) {
	t.Parallel()

	cab := &ArrayFS{}
	cab.Append(fstest.MapFS{
		"a/a.txt":  &fstest.MapFile{},
		"a/aa.txt": &fstest.MapFile{},
	})
	cab.Append(fstest.MapFS{
		"b/b.txt":  &fstest.MapFile{},
		"b/bb.txt": &fstest.MapFile{},
	})
	cab.Append(fstest.MapFS{
		"c/c.txt":  &fstest.MapFile{},
		"c/cc.txt": &fstest.MapFile{},
	})

	tcs := []struct {
		name string
		fn   string
		err  error
	}{
		{
			name: "known file",
			fn:   "b",
			err:  nil,
		},
		{
			name: "unknown file",
			fn:   "d",
			err:  fs.ErrNotExist,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			files, err := cab.ReadDir(tc.fn)
			if tc.err != nil {
				r.Error(err)
				r.True(errors.Is(err, tc.err))
				return
			}

			r.NoError(err)
			r.Len(files, 2)
		})
	}
}
