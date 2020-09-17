package core

import (
	"fmt"
	"strings"

	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
	"github.com/ionous/iffy/rt/writer"
)

// ExampleIndex verifies the loop index property.
func ExampleIndex() {
	var run forTester
	run.SetWriter(writer.NewStdout())
	if e := rt.RunOne(&run,
		&ForEachText{
			In: &Texts{oneTwoThree},
			Go: NewActivity(
				&Say{&PrintNum{&GetVar{T("index")}}},
			),
		},
	); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// 123
}

func ExampleChooseText() {
	var run forTester
	run.SetWriter(writer.NewStdout())
	if e := rt.RunOne(&run,
		&Say{&Commas{NewActivity(
			&ForEachText{
				In: &Texts{oneTwoThree},
				Go: NewActivity(
					&Say{&ChooseText{
						If:   &GetVar{T("last")},
						True: T("last"),
						False: &ChooseText{
							If:    &GetVar{T("first")},
							True:  T("first"),
							False: &GetVar{T("text")},
						},
					}}),
			}),
		}}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// first, two, and last

}

var oneTwoThree = []string{"one", "two", "three"}

type forTester struct {
	baseRuntime
	strings.Builder
	scope.ScopeStack
	writer.Sink
}

func T(t string) *Text {
	return &Text{t}
}
