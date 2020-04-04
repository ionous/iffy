package next

import (
	"github.com/ionous/iffy/assign"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
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
	if it, e := rt.GetNumberStream(run, f.In); e != nil {
		err = e
	} else {
		err = loop(run, it, f.Go, f.Else, func() (ret scope.ReadOnly, err error) {
			var num float64
			if e := it.GetNext(&num); e != nil {
				err = e
			} else {
				ret = &readOnlyValue{"num", func(pv interface{}) error {
					return assign.ToFloat(pv, num)
				}}
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
	if it, e := rt.GetTextStream(run, f.In); e != nil {
		err = e
	} else {
		err = loop(run, it, f.Go, f.Else, func() (ret scope.ReadOnly, err error) {
			var txt string
			if e := it.GetNext(&txt); e != nil {
				err = e
			} else {
				ret = &readOnlyValue{"text", func(pv interface{}) error {
					return assign.ToString(pv, txt)
				}}
			}
			return
		})
	}
	return
}

type readOnlyValue struct {
	name    string
	convert func(pv interface{}) error
}

func (h *readOnlyValue) GetVariable(n string, pv interface{}) (err error) {
	if n == h.name {
		err = h.convert(pv)
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
		var lf scope.LoopFactory
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
