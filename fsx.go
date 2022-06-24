package fsx

import (
	"fmt"
	"io/fs"
	"sort"
	"sync"
	"testing/fstest"
	"time"
)

type FS struct {
	Backing fs.FS

	extras fstest.MapFS

	sync.RWMutex
}

func (cab *FS) WriteFile(name string, data []byte) error {
	if err := cab.validate(); err != nil {
		return err
	}

	cab.Lock()
	defer cab.Unlock()

	cab.extras[name] = &fstest.MapFile{
		Data:    data,
		ModTime: time.Now(),
	}

	return nil
}

func (cab *FS) Open(name string) (fs.File, error) {
	if err := cab.validate(); err != nil {
		return nil, err
	}

	cab.RLock()
	defer cab.RUnlock()

	if f, err := cab.extras.Open(name); err == nil {
		return f, nil
	}

	return cab.Backing.Open(name)
}

func (cab *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	if err := cab.validate(); err != nil {
		return nil, err
	}

	cab.RLock()
	defer cab.RUnlock()

	dir := map[string][]fs.DirEntry{}

	if d, err := cab.extras.ReadDir(name); err == nil {
		dir[name] = d
	}

	if d, err := fs.ReadDir(cab.Backing, name); err == nil {
		dir[name] = append(dir[name], d...)
	}

	if len(dir) == 0 {
		return nil, fs.ErrNotExist
	}

	entries := dir[name]
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	return entries, nil
}

func (cab *FS) Stat(name string) (fs.FileInfo, error) {
	if err := cab.validate(); err != nil {
		return nil, err
	}

	f, err := cab.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.Stat()
}

func (cab *FS) validate() error {
	if cab == nil {
		return fmt.Errorf("FS is nil")
	}

	cab.Lock()
	defer cab.Unlock()

	if cab.Backing == nil {
		cab.Backing = fstest.MapFS{}
	}

	if cab.extras == nil {
		cab.extras = fstest.MapFS{}
	}

	return nil
}
