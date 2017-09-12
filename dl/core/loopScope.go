package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

// Looper creates LoopScopes
type Looper struct {
	run  rt.Runtime
	obj  rt.Object
	loop rt.ExecuteList
}

func MakeLooper(run rt.Runtime, temp interface{}, loop rt.ExecuteList) Looper {
	obj := run.Emplace(temp)
	runat := rt.AtFinder(run, obj)
	return Looper{runat, obj, loop}
}

func (l *Looper) RunNext(name string, value interface{}, hasNext bool) (err error) {
	var index int
	if e := l.obj.GetValue("index", &index); e != nil {
		err = e
	} else if e := l.obj.SetValue("index", index+1); e != nil {
		err = e
	} else if e := l.obj.SetValue("first", index == 0); e != nil {
		err = e
	} else if e := l.obj.SetValue("last", !hasNext); e != nil {
		err = e
	} else if e := l.obj.SetValue(name, value); e != nil {
		err = e
	} else if l.loop != nil {
		if e := l.loop.Execute(l.run); e != nil {
			err = errutil.New("each", name, e)
		}
	}
	return
}
