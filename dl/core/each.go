package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/scope"
)

// DoNothing implements Execute, but .... does nothing.
type DoNothing struct{}

// ForEacNum visits values in a list of numbers.
// For each value visited it executes a block of statements, pushing a NumberCounter object into the scope as @.
// If the list is empty, it executes an alternative block of statements.
type ForEachNum struct {
	In       rt.NumListEval
	Go, Else *Activity
}

// ForEachText visits values in a text list.
// For each value visited it executes a block of statements, pushing a TextCounter object into the scope as @.
// If the list is empty, it executes an alternative block of statements.
type ForEachText struct {
	In       rt.TextListEval
	Go, Else *Activity
}

func (*DoNothing) Compose() composer.Spec {
	return composer.Spec{
		Name:  "do_nothing",
		Group: "exec",
		Desc:  "Do Nothing: Statement which does nothing.",
	}
}

func (DoNothing) Execute(rt.Runtime) error { return nil }

func (*ForEachNum) Compose() composer.Spec {
	return composer.Spec{
		Name:   "for_each_num",
		Group:  "exec",
		Desc:   "For Each Number: Loops over the passed list of numbers, or runs the 'else' activity if empty.",
		Locals: []string{"index", "first", "last", "num"},
	}
}

func (f *ForEachNum) Execute(run rt.Runtime) (err error) {
	if vs, e := rt.GetNumList(run, f.In); e != nil {
		err = e
	} else {
		it := g.SliceFloats(vs)
		err = scope.LoopOver(run, "num", it, f.Go, f.Else)
	}
	return
}

func (*ForEachText) Compose() composer.Spec {
	return composer.Spec{
		Name:   "for_each_text",
		Group:  "exec",
		Desc:   "For Each Text: Loops over the passed list of text, or runs the 'else' activity if empty.",
		Locals: []string{"index", "first", "last", "text"},
	}
}

func (f *ForEachText) Execute(run rt.Runtime) (err error) {
	if vs, e := rt.GetTextList(run, f.In); e != nil {
		err = e
	} else {
		it := g.SliceStrings(vs)
		err = scope.LoopOver(run, "text", it, f.Go, f.Else)
	}
	return
}
