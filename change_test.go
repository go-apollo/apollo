//Copyright (c) 2017 Phil
package apollo

import "testing"

func TestChangeType(t *testing.T) {
	var tps = []ChangeType{ADD, MODIFY, DELETE, ChangeType(-1)}
	var strs = []string{"ADD", "MODIFY", "DELETE", "UNKNOW"}
	for i, tp := range tps {
		if tp.String() != strs[i] {
			t.FailNow()
		}
	}
}

func TestMakeDeleteChange(t *testing.T) {
	change := makeDeleteChange("key", []byte("val"))
	if change.ChangeType != DELETE || string(change.OldValue) != "val" {
		t.FailNow()
	}
}

func TestMakeModifyChange(t *testing.T) {
	change := makeModifyChange("key", []byte("old"), []byte("new"))
	if change.ChangeType != MODIFY || string(change.OldValue) != "old" || string(change.NewValue) != "new" {
		t.FailNow()
	}
}

func TestMakeAddChange(t *testing.T) {
	change := makeAddChange("key", []byte("value"))
	if change.ChangeType != ADD || string(change.NewValue) != "value" {
		t.FailNow()
	}
}
