package core

import (
	"fmt"

	"github.com/ionous/iffy/rt"
)

func ExamplePrintNum() {
	var run sayTester
	if e := rt.WriteText(&run, &PrintNum{&Number{213}}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// 213
}

func ExamplePrintNumWord() {
	var run sayTester
	if e := rt.WriteText(&run, &PrintNumWord{&Number{213}}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// two hundred thirteen
}
