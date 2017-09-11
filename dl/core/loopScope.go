package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/rt"
)

// Looper creates LoopScopes
type Looper struct {
	run  rt.Runtime
	obj  rt.Object
	loop rt.ExecuteList
}

func MakeLooper(run rt.Runtime, temp interface{}, loop rt.ExecuteList) Looper {
	obj := obj.Emplace(temp)
	runat := rt.AtFinder(run, obj)
	return Looper{runat, obj, loop}
}

func (l *Looper) RunNext(name string, value interface{}, hasNext bool) (err error) {
	var index int
	if e := l.run.GetValue(l.obj, "index", &index); e != nil {
		err = e
	} else if e := l.run.SetValue(l.obj, "index", index+1); e != nil {
		err = e
	} else if e := l.run.SetValue(l.obj, "first", index == 0); e != nil {
		err = e
	} else if e := l.run.SetValue(l.obj, "last", !hasNext); e != nil {
		err = e
	} else if e := l.run.SetValue(l.obj, name, value); e != nil {
		err = e
	} else if l.loop != nil {
		if e := l.loop.Execute(l.run); e != nil {
			err = errutil.New("each", name, e)
		}
	}
	return
}
