package gopherfs

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// A Dir implements FileSystem using the native file system restricted to a
// specific directory tree.
//
// Dir is based on http.Dir and has the same caveats.
type Dir string

// Open implements FileSystem using os.Open, opening files for reading rooted
// and relative to the directory d.
func (d Dir) Open(name string) (File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("gopher: invalid character in file path")
	}
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	f, err := os.Open(fullName)
	if err != nil {
		return nil, os.ErrNotExist
	}
	return f, nil
}
