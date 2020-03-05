package gopher

import (
	"bufio"
	"crypto/tls"
	"io"
	"net"
	"os"

	"github.com/shabbyrobe/furlib/internal/uuencode"
)

type ResponseClass int

const (
	UnknownClass ResponseClass = iota
	BinaryClass
	DirClass
	TextClass
)

var lineEnding = []byte{'\r', '\n'}

type ResponseInfo struct {
	Request *Request

	// TLS contains information about the TLS connection on which the
	// response was received. It is nil for unencrypted responses.
	// The pointer is shared between responses and should not be
	// modified.
	TLS *tls.ConnectionState
}

func (ri *ResponseInfo) URL() URL { return ri.Request.url }

func newResponseInfo(conn net.Conn, rq *Request) *ResponseInfo {
	ri := &ResponseInfo{
		Request: rq,
	}
	if tls, ok := conn.(interface{ ConnectionState() tls.ConnectionState }); ok {
		cs := tls.ConnectionState()
		ri.TLS = &cs
	}
	return ri
}

type Response interface {
	Reader() io.ReadCloser
	Info() *ResponseInfo
	Class() ResponseClass
	Close() error
}

type BinaryResponse struct {
	info  *ResponseInfo
	inner io.ReadCloser
}

var _ Response = &BinaryResponse{}

func NewBinaryResponse(info *ResponseInfo, rdr io.ReadCloser) *BinaryResponse {
	return &BinaryResponse{info: info, inner: rdr}
}

func (br *BinaryResponse) Class() ResponseClass  { return BinaryClass }
func (br *BinaryResponse) Info() *ResponseInfo   { return br.info }
func (br *BinaryResponse) Reader() io.ReadCloser { return br }
func (br *BinaryResponse) Close() error          { return br.inner.Close() }

func (br *BinaryResponse) Read(b []byte) (n int, err error) {
	return br.inner.Read(b)
}

type UUEncodedResponse struct {
	info *ResponseInfo
	uu   *uuencode.Reader
	cls  io.Closer
}

var _ Response = &UUEncodedResponse{}

func NewUUEncodedResponse(info *ResponseInfo, rdr io.ReadCloser) *UUEncodedResponse {
	uu := uuencode.NewReader(NewTextReader(rdr), nil)
	return &UUEncodedResponse{info: info, uu: uu, cls: rdr}
}

func (br *UUEncodedResponse) File() (string, bool)      { return br.uu.File() }
func (br *UUEncodedResponse) Mode() (os.FileMode, bool) { return br.uu.Mode() }

func (br *UUEncodedResponse) Class() ResponseClass  { return BinaryClass }
func (br *UUEncodedResponse) Info() *ResponseInfo   { return br.info }
func (br *UUEncodedResponse) Reader() io.ReadCloser { return br }
func (br *UUEncodedResponse) Close() error          { return br.cls.Close() }

func (br *UUEncodedResponse) Read(b []byte) (n int, err error) {
	return br.uu.Read(b)
}

type TextResponse struct {
	info *ResponseInfo
	rdr  io.Reader
	cls  io.ReadCloser
}

var _ Response = &TextResponse{}

func NewTextResponse(info *ResponseInfo, rdr io.ReadCloser) *TextResponse {
	return &TextResponse{info: info, rdr: NewTextReader(rdr), cls: rdr}
}

func (br *TextResponse) Class() ResponseClass  { return TextClass }
func (br *TextResponse) Info() *ResponseInfo   { return br.info }
func (br *TextResponse) Reader() io.ReadCloser { return br }
func (br *TextResponse) Close() error          { return br.cls.Close() }

func (br *TextResponse) Read(b []byte) (n int, err error) {
	return br.rdr.Read(b)
}

type DirResponse struct {
	info *ResponseInfo
	cls  io.Closer
	scn  *bufio.Scanner
	rdr  io.Reader
	err  error
	line int
}

var _ Response = &DirResponse{}

func NewDirResponse(info *ResponseInfo, rdr io.ReadCloser) *DirResponse {
	dot := NewTextReader(rdr)
	scn := bufio.NewScanner(dot)
	return &DirResponse{
		info: info,
		cls:  rdr,
		scn:  scn,
		rdr:  dot,
	}
}

func (br *DirResponse) Class() ResponseClass { return DirClass }
func (br *DirResponse) Info() *ResponseInfo  { return br.info }

func (br *DirResponse) Reader() io.ReadCloser {
	return &readCloser{
		readFn:  br.rdr.Read,
		closeFn: br.Close,
	}
}

func (br *DirResponse) Close() error {
	err := br.err
	if err == io.EOF {
		err = nil
	}
	if cerr := br.cls.Close(); err == nil && cerr != nil {
		err = cerr
	}
	return err
}

func (br *DirResponse) Next(dir *Dirent) bool {
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
