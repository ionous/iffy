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
	if e := rt.WriteText(&run, os.Stdout, &Span{helloThereWorld}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// hello there world
}

func ExampleBracket() {
	var run sayTester
	if e := rt.WriteText(&run, os.Stdout, &Bracket{helloThereWorld}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// ( hello there world )
}

func ExampleSlash() {
	var run sayTester
	if e := rt.WriteText(&run, os.Stdout, &Slash{helloThereWorld}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// hello / there / world
}

func ExampleCommas() {
	var run sayTester
	if e := rt.WriteText(&run, os.Stdout, &Commas{helloThereWorld}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// hello, there, and world
}

type baseRuntime struct {
	rt.Panic
}
type sayTester struct {
	baseRuntime
	qna.WriterStack
}
