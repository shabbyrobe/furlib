package gopher

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseDirent(t *testing.T) {
	for idx, tc := range []struct {
		in   string
		out  Dirent
		flag DirentFlag
	}{
		{ // All fields set.
			in:  "0foo\tbar\tbaz\t70",
			out: Dirent{ItemType: Text, Display: "foo", Selector: "bar", Hostname: "baz", Port: "70"},
		},
		{ // Different item type
			in:  "1a\tb\tc\t70",
			out: Dirent{ItemType: Dir, Display: "a", Selector: "b", Hostname: "c", Port: "70"},
		},
		{ // Port '0' is ok.
			in:  "0foo\ta\tb\t0",
			out: Dirent{ItemType: Text, Display: "foo", Selector: "a", Hostname: "b", Port: "0"},
		},
		{ // Empty fields are valid here; apps may choose their own strictness around this
			in:  "0\t\t\t0",
			out: Dirent{ItemType: Text, Display: "", Selector: "", Hostname: "", Port: "0"},
		},
		{ // Port optional with flag
			in:   "0a\tb\tc",
			out:  Dirent{ItemType: Text, Display: "a", Selector: "b", Hostname: "c"},
			flag: DirentHostOptional,
		},
		{ // Host+port optional with flag
			in:   "0a\tb",
			out:  Dirent{ItemType: Text, Display: "a", Selector: "b"},
			flag: DirentHostOptional,
		},
		{ // Dodgy port is OK with the correct flag
			in:   "0a\tb\tc\td",
			out:  Dirent{ItemType: Text, Display: "a", Selector: "b", Hostname: "c", Port: "d"},
			flag: direntNoValidatePort,
		},
	} {

		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			var d Dirent
			if err := parseDirent(tc.in, 1, &d, tc.flag); err != nil {
				t.Fatal(err)
			}
			if tc.out.Raw == "" {
				tc.out.Raw = tc.in
			}
			if !reflect.DeepEqual(d, tc.out) {
				t.Fatalf("%#v != %#v", d, tc.out)
			}
		})
	}
}
