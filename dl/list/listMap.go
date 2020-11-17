package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
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
	if e := op.execute(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Map) execute(run rt.Runtime) (err error) {
	if fromList, e := run.GetField(object.Variables, op.FromList); e != nil {
		err = e
	} else if toList, e := run.GetField(object.Variables, op.ToList); e != nil {
		err = e
	} else if e := op.cacheKinds(run, op.UsingPattern); e != nil {
		err = e
	} else {
		ps := op.pk.NewRecord()
		for it := g.ListIt(fromList); it.HasNext(); {
			if inVal, e := it.GetNext(); e != nil {
				err = e
				break
			} else if e := ps.SetFieldByIndex(op.in, inVal); e != nil {
				err = e
			} else if e := op.call(run, ps); e != nil {
				err = e
				break
			} else if newVal, e := ps.GetFieldByIndex(op.out); e != nil {
				err = e
				break
			} else if vs, e := toList.Append(newVal); e != nil {
				err = e
				break
			} else {
				toList = vs
			}
		}
		if err == nil {
			err = run.SetField(object.Variables, op.ToList, toList)
		}
	}
	return
}
