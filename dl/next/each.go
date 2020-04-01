package next

import (
	"github.com/ionous/iffy/qna"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/scope"
)

// DoNothing implements Execute, but .... does nothing.
type DoNothing struct{}

func (DoNothing) Execute(rt.Runtime) error { return nil }

// ForEacNum visits values in a list of numbers.
// For each value visited it executes a block of statements, pushing a NumberCounter object into the scope as @.
// If the list is empty, it executes an alternative block of statements.
type ForEachNum struct {
	In       rt.NumListEval
	Go, Else rt.Execute
}

func (f *ForEachNum) Execute(run rt.Runtime) (err error) {
	if it, e := f.In.GetNumberStream(run); e != nil {
		err = e
	} else {
		err = loop(run, it, f.Go, f.Else, func() (ret scope.ReadOnly, err error) {
			if v, e := it.GetNumber(); e != nil {
				err = e
			} else {
				ret = eachNumber(v)
			}
			return
		})
	}
	return
}

// ForEachText visits values in a text list.
// For each value visited it executes a block of statements, pushing a TextCounter object into the scope as @.
// If the list is empty, it executes an alternative block of statements.
type ForEachText struct {
	In       rt.TextListEval
	Go, Else rt.Execute
}

func (f *ForEachText) Execute(run rt.Runtime) (err error) {
	if it, e := f.In.GetTextStream(run); e != nil {
		err = e
	} else {
		err = loop(run, it, f.Go, f.Else, func() (ret scope.ReadOnly, err error) {
			if v, e := it.GetText(); e != nil {
				err = e
			} else {
				ret = eachText(v)
			}
			return
		})
	}
	return
}

type eachNumber float64

func (h eachNumber) GetVariable(n string, pv interface{}) (err error) {
	if n == "num" {
		err = qna.Assign(pv, float64(h))
	} else {
		err = scope.UnknownVariable(n)
	}
	return
}

type eachText string

func (h eachText) GetVariable(n string, pv interface{}) (err error) {
	if n == "text" {
		err = qna.Assign(pv, string(h))
	} else {
		err = scope.UnknownVariable(n)
	}
	return
}

func loop(run rt.Runtime, it interface{ HasNext() bool }, Go, Else rt.Execute, next func() (scope.ReadOnly, error)) (err error) {
	if !it.HasNext() {
		if e := Else.Execute(run); e != nil {
			err = e
		}
	} else {
		var lf qna.LoopFactory
		for it.HasNext() {
			if varscope, e := next(); e != nil {
				err = e
			} else if e := rt.ScopeBlock(run, lf.NextScope(varscope, it.HasNext()), Go); e != nil {
				err = e
				break
			}
		}
	}
	return
}
