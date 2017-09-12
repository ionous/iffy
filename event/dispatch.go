package event

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/class"
	"github.com/ionous/iffy/ref/kindOf"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

func Trigger(run rt.Runtime, events EventMap, data rt.Object) (err error) {
	id := class.Id(data.Type())
	if els, ok := events[id]; !ok {
		err = errutil.New("no such event", id)
	} else if target, e := TargetOf(run, data); e != nil {
		err = e
	} else if path, e := els.CollectAncestors(run, target); e != nil {
		err = e
	} else {
		//
		at := els.CollectTargets(target, nil)
		if len(path) == 0 && len(at) == 0 {
			err = run.ExecuteMatching(run, data)
		} else {
			evt := &EventObject{
				Id:            id,
				Data:          data,
				Bubbles:       true,
				Cancelable:    true,
				CurrentTarget: target,
			}
			// FIX: class finder and hints
			ac, e := NewFrame(run, evt)
			if e != nil {
				err = e
			} else if e := ac.DispatchFrame(at, path); e != nil {
				err = e
			} else if !evt.PreventDefault {
				if e := run.ExecuteMatching(run, data); e != nil {
					err = e
				} else {
					evt.Phase = AfterPhase
					evt.CurrentTarget = target
					err = ac.queue.Flush(run, evt)
				}
			}
		}
	}
	return
}

// TargetOf returns the target object for the passed event data as described by Field().
func TargetOf(run rt.Runtime, data rt.Object) (ret rt.Object, err error) {
	if field, ok := Field(data.Type()); !ok {
		err = errutil.New("no target found", data)
	} else if e := data.GetValue(field.Name, &ret); e != nil {
		err = e
	}
	return
}

// Field returns the struct tagged "if:target" in rtype, or the first struct field in rtype if there is no target tag.
func Field(rtype r.Type) (ret *r.StructField, okay bool) {
	foundDefault := false
	pathOf := func(f *r.StructField, path []int) (done bool) {
		if kindOf.IdentId(f.Type) {
			t := unique.Tag(f.Tag)
			if _, ok := t.Find("target"); ok {
				ret, done, okay = f, true, true
			} else if !foundDefault {
				ret, okay = f, true // keep going
			}
		}
		return
	}
	unique.WalkProperties(rtype, pathOf)
	return
}
