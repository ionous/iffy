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
				&Say{&PrintNum{&GetVar{"index"}}},
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
						If:   &GetVar{"last"},
						True: &Text{"last"},
						False: &ChooseText{
							If:    &GetVar{"first"},
							True:  &Text{"first"},
							False: &GetVar{"text"},
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
