package core

import (
	"github.com/ionous/errutil"

	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

// Looper creates LoopScopes
type Looper struct {
	run  rt.Runtime
	obj  rt.Object
	loop []rt.Execute
}

func NewLooper(run rt.Runtime, cls string, loop []rt.Execute) (ret *Looper, err error) {
	if temp, e := run.NewObject(cls); e != nil {
		err = e
	} else {
		run.PushScope(scope.MultiFinder(
			scope.AtFinder(temp),
			scope.ModelFinder(run),
		))
		ret = &Looper{run, temp, loop}
	}
	return
}
func (l *Looper) End() {
	l.run.PopScope()
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
		if e := rt.ExecuteList(l.run, l.loop); e != nil {
			err = errutil.New("each", name, e)
		}
	}
	return
}
