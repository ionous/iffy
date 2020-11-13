package list_test

import (
	"testing"

	"github.com/ionous/iffy/dl/list"
	g "github.com/ionous/iffy/rt/generic"
)

func TestPush(t *testing.T) {
	fruit := []string{"Lemon"}
	if got, e := push(fruit, true, "Banana", "Orange"); e != nil {
		t.Fatal(e)
	} else if got != "3; Banana, Orange, Lemon" {
		t.Fatal(got)
	}
	if got, e := push(fruit, true, "Apple"); e != nil {
		t.Fatal(e)
	} else if got != "2; Apple, Lemon" {
		t.Fatal(got)
	}
	if got, e := push(fruit, true); e != nil {
		t.Fatal(e)
	} else if got != "1; Lemon" {
		t.Fatal(got)
	}
	if got, e := push(fruit, false); e != nil {
		t.Fatal(e)
	} else if got != "1; Lemon" {
		t.Fatal(got)
	}
	if got, e := push(fruit, false, "Mango"); e != nil {
		t.Fatal(e)
	} else if got != "2; Lemon, Mango" {
		t.Fatal(got)
	}
	if got, e := push(fruit, false, "Mango", "Grape"); e != nil {
		t.Fatal(e)
	} else if got != "3; Lemon, Mango, Grape" {
		t.Fatal(got)
	}
}

func push(src []string, front bool, add ...string) (ret string, err error) {
	lt := listTime{vals: values{"src": g.StringsOf(append([]string{}, src...))}}
	num := getNum(&lt, &list.Push{"src", list.FrontOrBack(front), FromTs(add)})
	if strs, e := lt.vals["src"].GetTextList(); e != nil {
		err = e
	} else {
		next := joinStrings(strs) // get the variable set by splice
		ret = num + "; " + next
	}
	return
}
