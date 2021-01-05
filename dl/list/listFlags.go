package list

import "github.com/ionous/iffy/dl/composer"

type Edge bool

func (op *Edge) Front() bool { return *op != false }

func (*Edge) Compose() composer.Spec {
	return composer.Spec{
		Name:    "list_edge",
		Fluent:  &composer.Fluid{},
		Group:   "list",
		Strings: []string{"at_back", "at_front"},
		Desc:    "List Edge: Put elements at the front or back of a list.",
		Stub:    true, // the stub parse the flag
	}
}

type Order bool

func (op *Order) Descending() bool { return *op != false }

func (*Order) Compose() composer.Spec {
	return composer.Spec{
		Name:    "list_order",
		Group:   "list",
		Strings: []string{"ascending", "descending"},
		Desc:    "List Order: Sort larger values towards the end of a list.",
		Stub:    true, // the stub parse the flag
	}
}

type Case bool

func (op *Case) IgnoreCase() bool { return *op != false }

func (*Case) Compose() composer.Spec {
	return composer.Spec{
		Name:    "list_case",
		Group:   "list",
		Strings: []string{"include_case", "ignore_case"},
		Desc:    "List Case: When sorting, treat uppercase and lowercase versions of letters the same.",
		Stub:    true, // the stub parse the flag
	}
}
