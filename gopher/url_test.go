package gopher

import (
	"encoding/json"
	"testing"
)

func TestURLIPv6(t *testing.T) {
	u, err := ParseURL("gopher://[::]:70/1/yep")
	if err != nil {
		t.Fatal(err)
	}

	if u.Hostname != "::" {
		t.Fatal(err)
	}
	if u.Port != "70" {
		t.Fatal(err)
	}
	if u.String() != "gopher://[::]/1/yep" {
		t.Fatal(u.String())
	}
}

func TestURLStringWithEmptyPort(t *testing.T) {
	u := URL{Hostname: "invalid", Selector: "foo"}
	if u.String() != "gopher://invalid/0foo" {
		t.Fatal(u.String())
	}
}

func TestURLJSON(t *testing.T) {
	url := "gopher://foo:7070/bar"
	u, err := ParseURL(url)
	if err != nil {
		t.Fatal(err)
	}

	bts, err := json.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}

	var r URL
	if err := json.Unmarshal(bts, &r); err != nil {
		t.Fatal(err)
	}

	if u != r {
		t.Fatal(u, "!=", r)
	}
}

func TestURLSelectorOnly(t *testing.T) {
	url := "0foo"
	u, err := ParseURL(url)
	if err != nil {
		t.Fatal(err)
	}
	if u.Selector != "foo" {
		t.Fatal(u)
	}
	if u.IsAbs() {
		t.Fatal()
	}
}
