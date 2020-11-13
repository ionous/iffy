package list_test

import (
	"testing"

	"github.com/ionous/iffy/dl/list"
	g "github.com/ionous/iffy/rt/generic"
)

func TestSplices(t *testing.T) {
	fruit := []string{"Banana", "Orange", "Lemon", "Apple"}
	// insert by making the new element the second element
	if got, e := splice(fruit, 2, 0, "Mango"); e != nil {
		t.Fatal(e)
	} else if got != "Banana, Mango, Orange, Lemon, Apple; -" {
		t.Fatal(got)
	}
	// replace one element
	if got, e := splice(fruit, 4, 1, "Mango"); e != nil {
		t.Fatal(e)
	} else if got != "Banana, Orange, Lemon, Mango; Apple" {
		t.Fatal(got)
	}
	// remove two element
	if got, e := splice(fruit, 2, 2); e != nil {
		t.Fatal(e)
	} else if got != "Banana, Apple; Orange, Lemon" {
		t.Fatal(got)
	}
	// do nothing
	if got, e := splice(fruit, 0, 0); e != nil {
		t.Fatal(e)
	} else if got != "Banana, Orange, Lemon, Apple; -" {
		t.Fatal(got)
	}
	// remove them all
	if got, e := splice(fruit, 0, len(fruit)); e != nil {
		t.Fatal(e)
	} else if got != "-; Banana, Orange, Lemon, Apple" {
		t.Fatal(got)
	}
	// negative start
	if got, e := splice(fruit, -2, 2); e != nil {
		t.Fatal(e)
	} else if got != "Banana, Orange; Lemon, Apple" {
		t.Fatal(got)
	}
	// too negative is the same as starting at the front
	if got, e := splice(fruit, -20, 2); e != nil {
		t.Fatal(e)
	} else if got != "Lemon, Apple; Banana, Orange" {
		t.Fatal(got)
	}
	// negative lengths do nothing
	if got, e := splice(fruit, 3, -20); e != nil {
		t.Fatal(e)
	} else if got != "Banana, Orange, Lemon, Apple; -" {
		t.Fatal(got)
	}
	if got, e := splice(nil, 1, -1); e != nil {
		t.Fatal(e)
	} else if got != "-; -" {
		t.Fatal(got)
	}
}

func splice(src []string, start, cnt int, add ...string) (ret string, err error) {
	lt := listTime{vals: values{"src": g.StringsOf(append([]string{}, src...))}}
	rub := joinText(&lt, &list.Splice{"src", I(start), I(cnt), FromTs(add)})
	if strs, e := lt.vals["src"].GetTextList(); e != nil {
		err = e
	} else {
		next := joinStrings(strs) // get the variable set by splice
		ret = next + "; " + rub
	}
	return
}
