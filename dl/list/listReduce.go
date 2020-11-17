package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

// A normal reduce would return a value, instead we accumulate into a variable
type Reduce struct {
	FromList, IntoValue, UsingPattern string // variable names
	activityCache
}

func (*Reduce) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_map",
		Group: "list",
		Desc: `Reduce List: Transform the values from one list by combining them into a single value.
		The named pattern is called with two parameters: 'in' ( each element of the list ) and 'out' ( ex. a record ).`,
	}
}

func (op *Reduce) Execute(run rt.Runtime) (err error) {
	if e := op.execute(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Reduce) execute(run rt.Runtime) (err error) {
	if fromList, e := run.GetField(object.Variables, op.FromList); e != nil {
		err = e
	} else if outVal, e := run.GetField(object.Variables, op.IntoValue); e != nil {
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
				break
			} else if e := ps.SetFieldByIndex(op.out, outVal); e != nil {
				err = e
				break
			} else if e := op.call(run, ps); e != nil {
				err = e
				break
			} else if newVal, e := ps.GetFieldByIndex(op.out); e != nil {
				err = e
				break
			} else {
				// send it back in for the next time.
				outVal = newVal
			}
		}
		if err == nil {
			err = run.SetField(object.Variables, op.IntoValue, outVal)
		}
	}
	return
}
