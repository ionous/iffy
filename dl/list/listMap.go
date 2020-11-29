package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type Map struct {
	FromList, ToList, UsingPattern string // variable names
	activityCache
}

func (*Map) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_map",
		Group: "list",
		Desc: `Map List: Transform the values from one list and place the results in another list.
		The named pattern is called with two records 'in' and 'out' from the source and output lists respectively.`,
	}
}
func (op *Map) Execute(run rt.Runtime) (err error) {
	if e := op.remap(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Map) remap(run rt.Runtime) (err error) {
	if fromList, e := safe.List(run, op.FromList); e != nil {
		err = e
	} else if toList, e := safe.List(run, op.ToList); e != nil {
		err = e
	} else if e := op.cacheKinds(run, op.UsingPattern); e != nil {
		err = e
	} else {
		for it := g.ListIt(fromList); it.HasNext(); {
			ps := op.pk.NewRecord() // create a new set of parameters each loop
			if inVal, e := it.GetNext(); e != nil {
				err = e
				break
			} else if e := ps.SetFieldByIndex(op.in, inVal); e != nil {
				err = e
				break
			} else if e := op.call(run, ps); e != nil {
				err = e
				break
			} else if newVal, e := ps.GetFieldByIndex(op.out); e != nil {
				err = e
				break
			} else if src, dst := newVal.Affinity(), toList.Affinity(); src != affine.Element(dst) ||
				((src == affine.Record) && newVal.Type() != toList.Type()) {
				err = errutil.New("elements dont match")
				break
			} else {
				toList.Append(newVal)
			}
		}
		if err == nil {
			err = run.SetField(object.Variables, op.ToList, toList)
		}
	}
	return
}
