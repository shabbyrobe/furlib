package gopherfs

import (
	"io"
	"os"
)

// A FileSystem implements access to a collection of named files. The elements in a file
// path are separated by forward slash ('/') characters, regardless of the host OS.
type FileSystem interface {
	// Open would probably be better named OpenURLPath, but is 'Open' for compatibility
	// with http.FileServer
	Open(urlPath string) (File, error)
}

// A File is returned by a FileSystem's Open method and can be served by the FileServer
// implementation.
//
// The methods should behave the same as those on an *os.File.
type File interface {
	io.Closer
	io.Reader
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}
