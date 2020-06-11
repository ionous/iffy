package core

import (
	"fmt"

	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/print"
	"github.com/ionous/iffy/rt/scope"
)

// ExampleIndex verifies the loop index property.
func ExampleIndex() {
	var run forTester
	if e := rt.RunAll(&run, []rt.Execute{
		&ForEachText{
			In: &Texts{oneTwoThree},
			Go: []rt.Execute{
				&Say{&PrintNum{&GetVar{"index"}}}},
		},
	}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// 123
}

func ExampleChooseText() {
	var run forTester
	if e := rt.RunAll(&run, []rt.Execute{
		&Say{&Commas{[]rt.Execute{
			&ForEachText{
				In: &Texts{oneTwoThree},
				Go: []rt.Execute{
					&Say{&ChooseText{
						If:   &GetVar{"last"},
						True: &Text{"last"},
						False: &ChooseText{
							If:    &GetVar{"first"},
							True:  &Text{"first"},
							False: &GetVar{"text"},
						},
					}}},
			}},
		}}}); e != nil {
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
