package next

import (
	"fmt"
	"os"
)

func ExamplePrintNum() {
	var run sayTester
	prn := &PrintNum{&Number{213}}
	if e := prn.WriteText(&run, os.Stdout); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// 213
}

func ExamplePrintNumWord() {
	var run sayTester
	prn := &PrintNumWord{&Number{213}}
	if e := prn.WriteText(&run, os.Stdout); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// two hundred thirteen
}
