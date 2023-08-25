package fsx

import (
	"fmt"
	"io/fs"
	"sync"
)

type ArrayFS struct {
	cabs []fs.FS
	mu   sync.RWMutex
}

func (cab *ArrayFS) Append(c fs.FS) {
	if cab == nil && c == nil {
		return
	}

	cab.mu.Lock()
	defer cab.mu.Unlock()

	cab.cabs = append(cab.cabs, c)
}

func (cab *ArrayFS) Open(name string) (fs.File, error) {
	cab.mu.RLock()
	defer cab.mu.RUnlock()

	for _, c := range cab.cabs {
		f, err := c.Open(name)
		if err == nil {
			return f, nil
		}
	}

	return nil, fs.ErrNotExist
}

func (cab *ArrayFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if cab == nil {
		return nil, fmt.Errorf("fsx: ArrayFS is nil")
	}

	cab.mu.RLock()
	defer cab.mu.RUnlock()

	for _, c := range cab.cabs {
		if rd, ok := c.(fs.ReadDirFS); ok {
			f, err := rd.ReadDir(name)
			if err == nil {
				return f, nil
			}
		}
	}

	return nil, fs.ErrNotExist
}
