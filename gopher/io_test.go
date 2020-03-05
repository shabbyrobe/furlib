package gopher

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

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
