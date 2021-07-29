package tfs

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/require"
)

func File(t *testing.T, name string) *fstest.MapFile {
	t.Helper()
	f := &fstest.MapFile{
		Data:    []byte(strings.ToUpper(name)),
		ModTime: time.Now(),
		Mode:    0755,
	}
	return f
}

func Dir(t *testing.T, name string) *fstest.MapFile {
	t.Helper()

	r := require.New(t)
	r.Len(filepath.Ext(name), 0)
	return &fstest.MapFile{
		ModTime: time.Now(),
		Mode:    fs.ModeDir,
	}
}

func AppendFile(t *testing.T, cab fstest.MapFS, name string) {
	t.Helper()
	cab[name] = File(t, name)
}
