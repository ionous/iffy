package list_test

import (
	"strconv"
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/rt/safe"
)

func TestPush(t *testing.T) {
	fruit := []string{"Lemon"}
	if got, e := pushFront(fruit, "Banana", "Orange"); e != nil {
		t.Fatal(e)
	} else if got != "3; Banana, Orange, Lemon" {
		t.Fatal(got)
	}
	if got, e := pushFront(fruit, "Apple"); e != nil {
		t.Fatal(e)
	} else if got != "2; Apple, Lemon" {
		t.Fatal(got)
	}
	if got, e := pushFront(fruit); e != nil {
		t.Fatal(e)
	} else if got != "1; Lemon" {
		t.Fatal(got)
	}
	if got, e := pushBack(fruit); e != nil {
		t.Fatal(e)
	} else if got != "1; Lemon" {
		t.Fatal(got)
	}
	if got, e := pushBack(fruit, "Mango"); e != nil {
		t.Fatal(e)
	} else if got != "2; Lemon, Mango" {
		t.Fatal(got)
	}
	if got, e := pushBack(fruit, "Mango", "Grape"); e != nil {
		t.Fatal(e)
	} else if got != "3; Lemon, Mango, Grape" {
		t.Fatal(got)
	}
}

func pushFront(src []string, ins ...string) (ret string, err error) {
	return push(src, true, ins)
}
func pushBack(src []string, ins ...string) (ret string, err error) {
	return push(src, false, ins)
}
func push(src []string, front bool, ins []string) (ret string, err error) {
	if run, vals, e := newListTime(src, nil); e != nil {
		err = e
	} else {
		front := list.Edge(front)
		if e := safe.Run(run, &list.PutEdge{Into: &list.IntoTxtList{core.Variable{Str: "Source"}}, From: FromTs(ins), AtEdge: front}); e != nil {
			err = e
		} else if strs, e := vals.GetNamedField("Source"); e != nil {
			err = e
		} else {
			strs := strs.Strings()
			next := joinStrings(strs) // get the variable set by splice
			ret = strconv.Itoa(len(strs)) + "; " + next
		}
	}
	return
}
