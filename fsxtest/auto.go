package fsxtest

import (
	"fmt"
	"strings"
	"testing"
	"testing/fstest"
)

func AutoFile(t testing.TB, name string) *fstest.MapFile {
	t.Helper()

	return &fstest.MapFile{
		Data: []byte(strings.ToUpper(name)),
	}
}

func AutoFS(t testing.TB, files ...string) fstest.MapFS {
	t.Helper()

	m := fstest.MapFS{}

	for _, f := range files {
		m[f] = AutoFile(t, f)
	}

	return m
}

func AutoMod(t testing.TB, name string) *fstest.MapFile {
	t.Helper()

	return &fstest.MapFile{
		Data: []byte(fmt.Sprintf("module %s\n\ngo 1.19", name)),
	}
}
