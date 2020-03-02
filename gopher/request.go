package gopher

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
)

type Request struct {
	url    URL
	body   io.ReadCloser
	format string

	// Server only. When a server accepts an actual connection, this will be set to the
	// remote address.  This field is ignored by the Gopher client.
	RemoteAddr *net.TCPAddr

	// Server only. Params is free to be set by your Server's Mux implementation. If you
	// have requirements that this can't satisfy, use the dreaded context.WithValue() to
	// add what you need.
	Params Params

	// Server only. SelectorPrefix is added to all generated selectors which are internal
	// to this server. Comes from the Server struct. Useful for things like configurable
	// sites deployable at different selector bases.
	SelectorPrefix string
}

func NewRequest(url URL, body io.Reader) *Request {
	var ok bool
	var rc io.ReadCloser
	if body == nil {
		rc = nilReadCloserVal
	} else {
		rc, ok = body.(io.ReadCloser)
		if !ok {
			rc = ioutil.NopCloser(body)
		}
	}
	return &Request{
		url:  url,
		body: rc,
	}
}

func NewFormatRequest(url URL, format string, body io.Reader) (*Request, error) {
	if url.Search != "" {
		return nil, fmt.Errorf("gopher: format request URL must not contain selector")
	}
	rq := NewRequest(url, body)
	rq.format = format
	return rq, nil
}

func (r *Request) URL() URL            { return r.url }
func (r *Request) Body() io.ReadCloser { return r.body }

// In addition to a selector string, a GopherIIbis-compliant request contains a *format*
// string, a data flag indicating the presence or absence of a data block, and an OPTIONAL
// data block.
//
// The reason for the inclusion of the format string is because GopherIIbis allows one
// selector to point to multiple versions of the same file, in multiple languages.
func (r *Request) Format() string {
	return r.format
}

func (r *Request) buildSelector(buf *bytes.Buffer) error {
	buf.WriteString(r.url.Selector)

	if r.url.Search == "" && r.format == "" && r.body == nil {
		goto done
	}

	buf.WriteByte('\t')
	buf.WriteString(r.url.Search)

	if r.format == "" && r.body == nil {
		goto done
	}

	buf.WriteByte('\t')
	buf.WriteString(r.format)

	if r.body != nil {
		buf.WriteByte('1')
	} else {
		buf.WriteByte('0')
	}

done:
	buf.WriteString("\r\n")
	return nil
}

type Params []Param

type Param struct {
	Key, Value string
}

func (params Params) Get(name string) string {
	for _, param := range params {
		if param.Key == name {
			return param.Value
		}
	}
	return ""
}
