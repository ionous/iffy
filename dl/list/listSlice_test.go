package list_test

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
)

func TestSlices(t *testing.T) {
	fruit := []string{"Banana", "Orange", "Lemon", "Apple", "Mango"}
	//
	if got := slice(0, 0, fruit); got != "Banana, Orange, Lemon, Apple, Mango" {
		t.Fatal(got)
	}
	// start uses a one based index
	if got := slice(3, 0, fruit); got != "Lemon, Apple, Mango" {
		t.Fatal(got)
	}
	// not including the fifth element
	if got := slice(3, 5, fruit); got != "Lemon, Apple" {
		t.Fatal(got)
	}
	// further than the end
	if got := slice(2, 8, fruit); got != "Orange, Lemon, Apple, Mango" {
		t.Fatal(got)
	}
	// starting from, and including, the second to last element
	if got := slice(-2, 0, fruit); got != "Apple, Mango" {
		t.Fatal(got)
	}
	// the third el, up to, though not including, the last
	if got := slice(3, -1, fruit); got != "Lemon, Apple" {
		t.Fatal(got)
	}
	// start beyond the end
	if got := slice(30, 0, fruit); got != "-" {
		t.Fatal(got)
	}
	// ending before the front
	if got := slice(0, -9, fruit); got != "-" {
		t.Fatal(got)
	}
	// reverse order
	if got := slice(3, 1, fruit); got != "-" {
		t.Fatal(got)
	}
	// math reverse order
	if got := slice(4, -4, fruit); got != "-" {
		t.Fatal(got)
	}
	// just the last element
	if got := slice(5, 0, fruit); got != "Mango" {
		t.Fatal(got)
	}
}

func slice(start, end int, src []string) (ret string) {
	if run, _, e := newListTime(src, nil); e != nil {
		ret = e.Error()
	} else {
		ret = joinText(run, &list.Slice{&core.Var{Name: "Source"}, I(start), I(end)})
	}
	return
}
