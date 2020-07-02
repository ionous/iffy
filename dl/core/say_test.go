package core

import (
	"fmt"

	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/writer"
)

var helloThereWorld = []rt.Execute{
	&Say{&Text{"hello"}},
	&Say{&Text{"there"}},
	&Say{&Text{"world"}},
}

func ExampleSpan() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := rt.WriteText(&run, &Span{helloThereWorld}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// hello there world
}

func ExampleBracket() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := rt.WriteText(&run, &Bracket{helloThereWorld}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// ( hello there world )
}

func ExampleSlash() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := rt.WriteText(&run, &Slash{helloThereWorld}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// hello / there / world
}

func ExampleCommas() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := rt.WriteText(&run, &Commas{helloThereWorld}); e != nil {
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
	writer.Sink
}
