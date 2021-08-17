package fsx

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
)

var (
	_ fmt.Stringer   = &FS{}
	_ fs.FS          = &FS{}
	_ fs.StatFS      = &FS{}
	_ fs.ReadFileFS  = &FS{}
	_ json.Marshaler = &FS{}
)

type FS struct {
	fs.FS
	Root string
}

func (cab FS) Exists(name string) bool {
	if _, err := cab.Stat(name); err != nil {
		return false
	}
	return true
}

func (cab FS) ReadFile(name string) ([]byte, error) {
	if rfs, ok := cab.FS.(fs.ReadFileFS); ok {
		return rfs.ReadFile(name)
	}
	return fs.ReadFile(cab.FS, name)
}

func (cab FS) Sub(path string) (*FS, error) {
	kid, err := fs.Sub(cab.FS, path)
	if err != nil {
		return nil, err
	}
	return &FS{
		FS:   kid,
		Root: filepath.Join(cab.Root, path),
	}, nil
}

func (cab FS) Abs(name string) (string, error) {
	if _, err := cab.Stat(name); err != nil {
		return name, err
	}
	return filepath.Join(cab.Root, name), nil
}

func (cab FS) MarshalJSON() ([]byte, error) {
	infos, err := Infos(cab)
	if err != nil {
		return nil, err
	}

	return json.Marshal(infos)
}

func (cab FS) Stat(path string) (fs.FileInfo, error) {
	return fs.Stat(cab.FS, path)
}

func (cab FS) String() string {
	b, err := json.Marshal(cab)
	if err != nil {
		return fmt.Sprintf("can not marshal cab %v %s", cab.FS, err)
	}
	return string(b)
}

func NewFS(cab fs.FS) *FS {
	return &FS{
		FS: cab,
	}
}
