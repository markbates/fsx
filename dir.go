package fsx

import (
	"os"
	"path/filepath"
)

func DirFS(path string) (*FS, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	cab := os.DirFS(path)
	return &FS{
		FS: cab,
	}, nil
}
