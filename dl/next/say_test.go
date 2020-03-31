package next

import (
	"fmt"
	"os"

	"github.com/ionous/iffy/qna"
	"github.com/ionous/iffy/rt"
)

var helloThereWorld = rt.Block{
	&Say{&Text{"hello"}},
	&Say{&Text{"there"}},
	&Say{&Text{"world"}},
}

func ExampleSpan() {
	var run sayTester
	span := &Span{helloThereWorld}
	if e := span.WriteText(&run, os.Stdout); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// hello there world
}

func ExampleBracket() {
	var run sayTester
	span := &Bracket{helloThereWorld}
	if e := span.WriteText(&run, os.Stdout); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// ( hello there world )
}
func ExampleSlash() {
	var run sayTester
	slash := &Slash{helloThereWorld}
	if e := slash.WriteText(&run, os.Stdout); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// hello / there / world
}

func ExampleCommas() {
	var run sayTester
	commas := &Commas{helloThereWorld}
	if e := commas.WriteText(&run, os.Stdout); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// hello, there, and world
}

type sayBase struct {
	rt.Panic
}
type sayTester struct {
	sayBase
	qna.WriterStack
}
