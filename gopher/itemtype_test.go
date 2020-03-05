package gopher

import (
	"encoding/json"
	"testing"
)

func TestItemTypeMarshal(t *testing.T) {
	v := Text
	var r ItemType
	bts, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(bts, &r); err != nil {
		t.Fatal(err)
	}
	if r != v {
		t.Fatal(r, "!=", v)
	}
}
