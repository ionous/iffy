package list_test

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/kr/pretty"
)

func TestPop(t *testing.T) {
	// pop from the front of a list
	front := popTest(true, 5, "Orange", "Lemon", "Mango")
	if d := pretty.Diff(front, []string{"Orange", "Lemon", "Mango", "x", "x"}); len(d) > 0 {
		t.Fatal(d)
	}
	// pop from the back of a list
	back := popTest(false, 5, "Orange", "Lemon", "Mango")
	if d := pretty.Diff(back, []string{"Mango", "Lemon", "Orange", "x", "x"}); len(d) > 0 {
		t.Fatal(d)
	}
}

func popTest(front bool, amt int, strs ...string) []string {
	var out []string
	pop := &list.Pop{List: "strings",
		Front: list.FrontOrBack(front),
		Go:    core.NewActivity(&Write{&out, &core.GetVar{Name: T("text")}}),
		Else:  core.NewActivity(&Write{&out, T("x")}),
	}
	run := listTime{strings: strs}
	for i := 0; i < amt; i++ {
		if e := pop.Execute(&run); e != nil {
			panic(e)
		}
	}
	return out
}
