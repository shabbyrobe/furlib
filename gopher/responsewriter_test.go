package gopher

import (
	"bufio"
	"errors"
	"testing"
)

type errorWriter struct {
	err error
	n   int
}

func (ew errorWriter) Write(b []byte) (n int, err error) {
	return ew.n, ew.err
}

func TestDirWriterErrors(t *testing.T) {
	bork := errors.New("bork")
	rq := NewRequest(URL{Hostname: "yep", Port: "70"}, nil)

	{ // Dir() with buffering
		bw := bufio.NewWriterSize(errorWriter{err: bork}, 2048) // big enough to not flush on Dir()
		dw := NewDirWriterBuffer(bw, rq)
		if dw.Dir("yep", "yep") != nil {
			t.Fatal()
		}
		if dw.Flush() != bork {
			t.Fatal()
		}
	}

	{ // Dir() with implicit flush
		bw := bufio.NewWriterSize(errorWriter{err: bork}, 2) // small buffer should cause Dir() to flush
		dw := NewDirWriterBuffer(bw, rq)
		if dw.Dir("yep", "yep") != bork {
			t.Fatal()
		}
		if dw.Flush() != bork {
			t.Fatal()
		}
	}

	{ // Selector() with buffering
		bw := bufio.NewWriterSize(errorWriter{err: bork}, 2048) // big enough to not flush on Selector()
		dw := NewDirWriterBuffer(bw, rq)
		if dw.Selector(Text, "yep", "yep") != nil {
			t.Fatal()
		}
		if dw.Flush() != bork {
			t.Fatal()
		}
	}

	{ // Selector() with implicit flush
		bw := bufio.NewWriterSize(errorWriter{err: bork}, 2) // small buffer should cause Selector() to flush
		dw := NewDirWriterBuffer(bw, rq)
		if dw.Selector(Text, "yep", "yep") != bork {
			t.Fatal()
		}
		if dw.Flush() != bork {
			t.Fatal()
		}
	}

	{ // Error() with buffering
		bw := bufio.NewWriterSize(errorWriter{err: bork}, 2048) // big enough to not flush on Error()
		dw := NewDirWriterBuffer(bw, rq)
		if dw.Error("yep") != nil {
			t.Fatal()
		}
		if dw.Flush() != bork {
			t.Fatal()
		}
	}

	{ // Error() with implicit flush
		bw := bufio.NewWriterSize(errorWriter{err: bork}, 2) // small buffer should cause Error() to flush
		dw := NewDirWriterBuffer(bw, rq)
		if dw.Error("yep") != bork {
			t.Fatal()
		}
		if dw.Flush() != bork {
			t.Fatal()
		}
	}

	{ // Info() with buffering
		bw := bufio.NewWriterSize(errorWriter{err: bork}, 2048) // big enough to not flush on Info()
		dw := NewDirWriterBuffer(bw, rq)
		if dw.Info("yep") != nil {
			t.Fatal()
		}
		if dw.Flush() != bork {
			t.Fatal()
		}
	}

	{ // Info() with implicit flush
		bw := bufio.NewWriterSize(errorWriter{err: bork}, 2) // small buffer should cause Info() to flush
		dw := NewDirWriterBuffer(bw, rq)
		if dw.Info("yep") != bork {
			t.Fatal()
		}
		if dw.Flush() != bork {
			t.Fatal()
		}
	}
}
