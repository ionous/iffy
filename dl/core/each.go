package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

type DoNothing struct{}

func (DoNothing) Execute(rt.Runtime) error { return nil }

type Len struct {
	List rt.ObjListEval
}

func (op *Len) GetNumber(run rt.Runtime) (ret float64, err error) {
	if os, e := op.List.GetObjectStream(run); e != nil {
		err = e
	} else if l, ok := os.(stream.Len); !ok {
		err = errutil.Fmt("unknown list type %T", os)
	} else {
		ret = float64(l.Len())
	}
	return
}

type EachEndStatus int

// EachCounter is a dl class used during ForEach* loops.
type EachCounter struct {
	Index       float64 // loop counter
	First, Last bool    // true if the loop is at the initial or final element respectively
}

// NumberCounter is a dl class used during ForEachNum loops.
type NumberCounter struct {
	EachCounter
	Num float64 // current value of ForEachNum
}

// TextCounter is a dl class used during ForEachText loops.
type TextCounter struct {
	EachCounter
	Text string // current value of ForEachText
}

// ForEacNum visits values in a list of numbers.
// For each value visited it executes a block of statements, pushing a NumberCounter object into the scope as @.
// If the list is empty, it executes an alternative block of statements.
type ForEachNum struct {
	In       rt.NumListEval
	Go, Else []rt.Execute
}

func (f *ForEachNum) Execute(run rt.Runtime) (err error) {
	if it, e := f.In.GetNumberStream(run); e != nil {
		err = e
	} else if !it.HasNext() {
		if e := rt.ExecuteList(f.Else).Execute(run); e != nil {
			err = errutil.New("failed each num else", e)
		}
	} else if l, e := NewLooper(run, "NumberCounter", f.Go); e != nil {
		err = e
	} else {
		for it.HasNext() {
			if v, e := it.GetNext(); e != nil {
				err = errutil.New("failed each num get", e)
				break
			} else if e := l.RunNext("Num", v, it.HasNext()); e != nil {
				err = e
				break
			}
		}
		l.End()
	}
	return
}

// ForEachText visits values in a text list.
// For each value visited it executes a block of statements, pushing a TextCounter object into the scope as @.
// If the list is empty, it executes an alternative block of statements.
type ForEachText struct {
	In       rt.TextListEval
	Go, Else []rt.Execute
}

func (f *ForEachText) Execute(run rt.Runtime) (err error) {
	if it, e := f.In.GetTextStream(run); e != nil {
		err = e
	} else if !it.HasNext() {
		if e := rt.ExecuteList(f.Else).Execute(run); e != nil {
			err = errutil.New("failed each num else", e)
		}
	} else if l, e := NewLooper(run, "TextCounter", f.Go); e != nil {
		err = e
	} else {
		for it.HasNext() {
			if v, e := it.GetNext(); e != nil {
				err = errutil.New("failed each num get", e)
				break
			} else if e := l.RunNext("Text", v, it.HasNext()); e != nil {
				err = e
				break
			}
		}
		l.End()
	}
	return
}
