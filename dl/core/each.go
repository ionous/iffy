package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
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
		Desc:   "For Each Number: Loops over the passed list of numbers, or runs the 'else' statement if empty.",
		Locals: []string{"index", "first", "last", "num"},
	}
}

func (f *ForEachNum) Execute(run rt.Runtime) (err error) {
	if it, e := rt.GetNumberStream(run, f.In); e != nil {
		err = e
	} else {
		err = loop(run, it, f.Go, f.Else, func() (retn string, retv rt.Value, err error) {
			var num float64
			if e := it.GetNext(&num); e != nil {
				err = e
			} else {
				retn, retv = "num", &rt.NumberValue{Value: num}
			}
			return
		})
	}
	return
}

func (*ForEachText) Compose() composer.Spec {
	return composer.Spec{
		Name:   "for_each_text",
		Group:  "exec",
		Desc:   "For Each Text: Loops over the passed list of text, or runs the 'else' statement if empty.",
		Locals: []string{"index", "first", "last", "text"},
	}
}

func (f *ForEachText) Execute(run rt.Runtime) (err error) {
	if it, e := rt.GetTextStream(run, f.In); e != nil {
		err = e
	} else {
		err = loop(run, it, f.Go, f.Else, func() (retn string, retv rt.Value, err error) {
			var text string
			if e := it.GetNext(&text); e != nil {
				err = e
			} else {
				retn, retv = "text", &rt.TextValue{Value: text}
			}
			return
		})
	}
	return
}

// iterate over the go and else statements; introduce the loop counter
func loop(
	run rt.Runtime,
	it interface{ HasNext() bool },
	Run, Else *Activity,
	next func() (string, rt.Value, error),
) (err error) {
	if !it.HasNext() {
		if e := rt.RunOne(run, Else); e != nil {
			err = e
		}
	} else {
		var lf scope.LoopFactory
		for it.HasNext() {
			if name, val, e := next(); e != nil {
				err = e
				break
			} else {
				// brings the names of an object's properties into scope for the duration of fn.
				run.PushScope(lf.NextScope(name, val, it.HasNext()))
				e := rt.RunOne(run, Run)
				run.PopScope()
				if e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}
