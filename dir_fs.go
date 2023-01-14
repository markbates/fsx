package fsx

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var _ fs.FS = DirFS("")
var _ WriteableFS = DirFS("")

type DirFS string

func (cab DirFS) Open(name string) (fs.File, error) {
	return os.DirFS(string(cab)).Open(name)
}

func (cab DirFS) Create(name string) (io.WriteCloser, error) {
	return os.Create(filepath.Join(string(cab), name))
}
