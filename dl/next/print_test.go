package next

import (
	"fmt"
	"os"

	"github.com/ionous/iffy/rt"
)

func ExamplePrintNum() {
	var run sayTester
	if e := rt.WriteText(&run, os.Stdout, &PrintNum{&Number{213}}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// 213
}

func ExamplePrintNumWord() {
	var run sayTester
	if e := rt.WriteText(&run, os.Stdout, &PrintNumWord{&Number{213}}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// two hundred thirteen
}
