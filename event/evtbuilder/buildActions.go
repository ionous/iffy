package evtbuilder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ops"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
)

func NewActions(classes *ref.Classes, ops *ops.Ops) *Actions {
	return &Actions{classes, ops, ref.NewClasses(), make(event.ActionMap)}
}

type Actions struct {
	objectClasses *ref.Classes
	ops           *ops.Ops
	dataClasses   *ref.Classes
	event.ActionMap
}

// Add registers a new action builder class.
func (a *Actions) Add(name, targetClass string, dataClass interface{}) (err error) {
	id := id.MakeId(name)
	if act, exists := a.ActionMap[id]; exists {
		err = errutil.New("duplicate action registered", name, act.Name)
	} else if targetCls, ok := a.objectClasses.GetClass(targetClass); !ok {
		err = errutil.New("target class not found", targetClass)
	} else if rtype, e := unique.TypePtr(dataClass); e != nil {
		err = e
	} else if dataCls, e := a.objectClasses.RegisterClass(rtype); e != nil {
		err = e
	} else {
		a.ActionMap[id] = &event.Action{
			id,
			name,
			targetCls,
			dataCls,
			nil,
		}
	}
	return
}

// On adds a default action.
func (a *Actions) On(name string, fn BuildOps) (err error) {
	id := id.MakeId(name)
	if act, exists := a.ActionMap[id]; !exists {
		err = errutil.New("unknown action", name)
	} else if exec, e := fn.Build(a.ops); e != nil {
		err = e
	} else {
		act.DefaultActions = append(act.DefaultActions, exec)
	}
	return
}
