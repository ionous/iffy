package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
	"github.com/ionous/iffy/rt/scope"
)

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

func (*ForEachNum) Compose() composer.Spec {
	return composer.Spec{
		Name:   "for_each_num",
		Group:  "exec",
		Desc:   "For Each Number: Loops over the passed list of numbers, or runs the 'else' activity if empty.",
		Locals: []string{"index", "first", "last", "num"},
	}
}

func (op *ForEachNum) Execute(run rt.Runtime) (err error) {
	if e := op.forEach(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ForEachNum) forEach(run rt.Runtime) (err error) {
	if vs, e := safe.GetNumList(run, op.In); e != nil {
		err = e
	} else {
		err = scope.LoopOver(run, "num", g.ListIt(vs), op.Go, op.Else)
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

func (op *ForEachText) Execute(run rt.Runtime) (err error) {
	if e := op.forEach(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ForEachText) forEach(run rt.Runtime) (err error) {
	if vs, e := safe.GetTextList(run, op.In); e != nil {
		err = e
	} else {
		err = scope.LoopOver(run, "text", g.ListIt(vs), op.Go, op.Else)
	}
	return
}
