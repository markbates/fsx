package fsx

import (
	"fmt"
	"io/fs"
	"regexp"
)

func Filter(cab fs.FS, filter string, fn fs.WalkDirFunc) error {
	rx, err := regexp.Compile(filter)
	if err != nil {
		return fmt.Errorf("invalid filter: %q: %w", filter, err)
	}

	err = fs.WalkDir(cab, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if rx.MatchString(path) {
			return fn(path, d, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking dir: %q: %w", filter, err)
	}

	return nil
}
