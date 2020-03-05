package gopher

import (
	"bufio"
	"io"
)

type DirReader struct {
	scn  *bufio.Scanner
	rdr  io.Reader
	err  error
	line int
}

func NewDirReader(rdr io.Reader) *DirReader {
	dot := NewTextReader(rdr)
	scn := bufio.NewScanner(dot)
	return &DirReader{
		scn: scn,
		rdr: dot,
	}
}

func (br *DirReader) ReadErr() error {
	return br.err
}

func (br *DirReader) Read(dir *Dirent) bool {
	if br.err != nil {
		return false
	}

retry:
	if !br.scn.Scan() {
		br.err = br.scn.Err()
		return false
	}
	br.line++

	txt := br.scn.Text()
	if len(txt) == 0 {
		goto retry
	}

	if err := parseDirent(txt, br.line, dir, 0); err != nil {
		br.err = err
		return false
	}

	return true
}
