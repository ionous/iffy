package next

import (
	"fmt"

	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/print"
	"github.com/ionous/iffy/rt/scope"
)

// ExampleIndex verifies the loop index property.
func ExampleIndex() {
	var run forTester
	if e := rt.Run(&run, &ForEachText{
		In: &Texts{oneTwoThree},
		Go: &Say{&PrintNum{&GetVar{"index"}}},
	},
	); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// 123
}

func ExampleEndings() {
	var run forTester
	if e := rt.Run(&run,
		&Say{&Commas{
			&ForEachText{
				In: &Texts{oneTwoThree},
				Go: &Say{&ChooseText{
					If:   &GetVar{"last"},
					True: &Text{"last"},
					False: &ChooseText{
						If:    &GetVar{"first"},
						True:  &Text{"first"},
						False: &GetVar{"text"},
					},
				}},
			},
		}}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// first, two, and last

}

var oneTwoThree = []string{"one", "two", "three"}

type forTester struct {
	baseRuntime
	print.Stack
	scope.ScopeStack
}
