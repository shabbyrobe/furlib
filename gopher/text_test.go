package gopher

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

const textEnd = ".\n"
const textLineEnd = "\n" + textEnd

type textCase struct {
	dec string
	enc string
}

func (tc textCase) Enc() string {
	return strings.Replace(tc.enc, "\n", "\r\n", -1)
}

func (tc textCase) Dec(nl string) string {
	if nl == "" {
		return tc.dec
	}
	return strings.Replace(tc.dec, "\n", nl, -1)
}

var textCasesShared = []textCase{
	{"", textEnd},

	{"\n", "\n" + textEnd},
	{"-\n-", "-\n-" + textLineEnd},

	// FIXME: Leading-period escapes are not properly handled by textproto.DotReader:
	// {".", ".." + textLineEnd},
	// {"..", "..." + textLineEnd},
	// {"..\n..", "..\n.." + textLineEnd},
	// {"..\n..\n", "..\n..\n" + textEnd},
	// {"-\n-\n.\n", "-\n-\n..\n" + textEnd},
	// {"-\n-\n", "-\n-\n" + textEnd},
	// {"..\n-\n.", "...\n-\n.." + textLineEnd},
}

func TestTextWriter(t *testing.T) {
	cases := append(append(make([]textCase, 0), textCasesShared...), []textCase{}...)

	nlNames := []string{"unix", "win"}

	for idx, tc := range cases {
		for nlIdx, decNl := range []string{"\n", "\r\n"} {
			t.Run(fmt.Sprintf("%d/%q", idx, nlNames[nlIdx]), func(t *testing.T) {
				enc, dec := tc.Enc(), tc.Dec(decNl)

				var buf bytes.Buffer
				tw := NewTextWriter(&buf)
				tw.WriteString(dec)
				if err := tw.Flush(); err != nil {
					t.Fatal(err)
				}
				if enc != buf.String() {
					t.Fatalf("enc %q != wrt %q", enc, buf.String())
				}

				tr := NewTextReader(&buf)
				back, err := ioutil.ReadAll(tr)
				if err != nil {
					t.Fatal(err)
				}

				// When we go back, the string will have a final newline added
				// if it didn't have one before. We also have to use the original
				// decoded version, as TextWriter will always strip CR:
				decBack := string(tc.dec)
				if len(decBack) > 0 && !strings.HasSuffix(decBack, "\n") {
					decBack += "\n"
				}
				if string(back) != decBack {
					t.Fatalf("dec %q != bck %q", decBack, string(back))
				}
			})
		}
	}
}

func TestTextReader(t *testing.T) {
	for idx, tc := range []struct {
		in  string
		out string
	}{
		{"", ""},
		{"\n", "\n"},
		{"-\n-", "-\n-"},
		{"-\n-\n.", "-\n-\n"},
		{"-\n-\n.\n", "-\n-\n"},
		{"-\n-\n", "-\n-\n"},
		{"..\n-\n.", ".\n-\n"},
		{"..\n-\n.\n", ".\n-\n"},
		{"..\n..\n.", ".\n.\n"},
		{"..\n..\n.\n", ".\n.\n"},
		{},
	} {
		for _, nl := range []string{"\n", "\r\n"} {
			t.Run(fmt.Sprintf("%d/%q", idx, nl), func(t *testing.T) {
				in := strings.Replace(tc.in, "\n", nl, -1)
				tr := NewTextReader(strings.NewReader(in))
				result, _ := ioutil.ReadAll(tr)
				if string(result) != tc.out {
					t.Fatalf("%q != %q", string(result), tc.out)
				}
			})
		}
	}
}
