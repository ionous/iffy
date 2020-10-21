package list_test

import (
	"testing"

	"github.com/ionous/iffy/dl/list"
)

func TestSplices(t *testing.T) {
	months := []string{"Banana", "Orange", "Lemon", "Apple"}
	// insert by making the new element the second element
	if got := splice(months, 2, 0, "Mango"); got != "Banana, Mango, Orange, Lemon, Apple; -" {
		t.Fatal(got)
	}
	// replace one element
	if got := splice(months, 4, 1, "Mango"); got != "Banana, Orange, Lemon, Mango; Apple" {
		t.Fatal(got)
	}
	// remove two element
	if got := splice(months, 2, 2); got != "Banana, Apple; Orange, Lemon" {
		t.Fatal(got)
	}
	// do nothing
	if got := splice(months, 0, 0); got != "Banana, Orange, Lemon, Apple; -" {
		t.Fatal(got)
	}
	// remove them all
	if got := splice(months, 0, len(months)); got != "-; Banana, Orange, Lemon, Apple" {
		t.Fatal(got)
	}
	// negative start
	if got := splice(months, -2, 2); got != "Banana, Orange; Lemon, Apple" {
		t.Fatal(got)
	}
	// too negative is the same as starting at the front
	if got := splice(months, -20, 2); got != "Lemon, Apple; Banana, Orange" {
		t.Fatal(got)
	}
	// negative lengths do nothing
	if got := splice(months, 3, -20); got != "Banana, Orange, Lemon, Apple; -" {
		t.Fatal(got)
	}
}

func splice(src []string, start, cnt int, add ...string) (ret string) {
	run := listTime{strings: append([]string{}, src...)}
	rub := joinText(&run, &list.Splice{"strings", I(start), I(cnt), FromTs(add)})
	next := joinStrings(run.strings) // get the variable set by splice
	return next + "; " + rub
}
