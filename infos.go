package fsx

import (
	"encoding/json"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

type FileInfoJSON struct {
	IsDir   bool        `json:"is_dir,omitempty"`
	ModTime time.Time   `json:"mod_time,omitempty"`
	Mode    fs.FileMode `json:"mode,omitempty"`
	Name    string      `json:"name,omitempty"`
	Size    int64       `json:"size,omitempty"`
}

func NewFileInfoJSON(info fs.FileInfo) FileInfoJSON {
	return FileInfoJSON{
		IsDir:   info.IsDir(),
		ModTime: info.ModTime(),
		Mode:    info.Mode(),
		Name:    info.Name(),
		Size:    info.Size(),
	}
}

type InfoMap map[string]fs.FileInfo

func (im InfoMap) MarshalJSON() ([]byte, error) {
	m := map[string]FileInfoJSON{}

	for k, v := range im {
		m[k] = NewFileInfoJSON(v)
	}

	return json.Marshal(m)
}

func Infos(cab fs.FS) (InfoMap, error) {
	infos := InfoMap{}

	err := fs.WalkDir(cab, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		base := filepath.Base(path)
		if strings.HasPrefix(base, ".") {
			return nil
		}

		info, err := fs.Stat(cab, path)
		if err != nil {
			return err
		}

		infos[path] = info
		return nil
	})

	return infos, err
}
