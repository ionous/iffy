package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

// RelationEmpty returns true if the requested object has no related objects.
type RelationEmpty struct {
	Relation string
	Object   rt.ObjectEval
}

// RelatedList returns a stream of objects related to the requested object.
type RelatedList struct {
	Relation string
	Object   rt.ObjectEval
}

func (a *RelationEmpty) GetBool(run rt.Runtime) (okay bool, err error) {
	if r, ok := run.GetRelation(a.Relation); !ok {
		err = errutil.New("unknown relation", a.Relation)
	} else if obj, e := a.Object.GetObject(run); e != nil {
		err = e
	} else {
		_, hasChild := r.GetTable().Primary.FindFirst(0, obj.Id().Name)
		okay = !hasChild
	}
	return
}

func (a *RelatedList) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
	if r, ok := run.GetRelation(a.Relation); !ok {
		err = errutil.New("unknown relation", a.Relation)
	} else if obj, e := a.Object.GetObject(run); e != nil {
		err = e
	} else {
		// for now, copy everything all at once.
		// in the far future, could have some sort of lock.
		var list []rt.Object
		r.GetTable().Primary.Walk(obj.Id().Name, func(other string) (done bool) {
			if obj, ok := run.GetObject(other); !ok {
				err = errutil.New("unknown related object", a.Object, a.Relation, other)
				done = true
			} else {
				list = append(list, obj)
			}
			return
		})
		if err == nil {
			ret = stream.NewObjectStream(stream.FromList(list))
		}
	}
	return
}
