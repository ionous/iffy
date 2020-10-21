package list_test

import (
	"testing"

	"github.com/ionous/iffy/dl/list"
)

func TestPush(t *testing.T) {
	fruit := []string{"Lemon"}
	if got := push(fruit, true, "Banana", "Orange"); got != "3; Banana, Orange, Lemon" {
		t.Fatal(got)
	}
	if got := push(fruit, true, "Apple"); got != "2; Apple, Lemon" {
		t.Fatal(got)
	}
	if got := push(fruit, true); got != "1; Lemon" {
		t.Fatal(got)
	}
	if got := push(fruit, false); got != "1; Lemon" {
		t.Fatal(got)
	}
	if got := push(fruit, false, "Mango"); got != "2; Lemon, Mango" {
		t.Fatal(got)
	}
	if got := push(fruit, false, "Mango", "Grape"); got != "3; Lemon, Mango, Grape" {
		t.Fatal(got)
	}
}

func push(src []string, front bool, add ...string) (ret string) {
	run := listTime{strings: append([]string{}, src...)}
	num := getNum(&run, &list.Push{"strings", list.FrontOrBack(front), FromTs(add)})
	next := joinStrings(run.strings) // get the variable set by splice
	return num + "; " + next
}
