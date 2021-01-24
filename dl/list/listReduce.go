package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// A normal reduce would return a value, instead we accumulate into a variable
type Reduce struct {
	IntoValue    string
	FromList     core.Assignment
	UsingPattern pattern.PatternName
	activityCache
}

func (*Reduce) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_reduce",
		Group: "list",
		Desc: `Reduce List: Transform the values from one list by combining them into a single value.
		The named pattern is called with two parameters: 'in' ( each element of the list ) and 'out' ( ex. a record ).`,
	}
}

func (op *Reduce) Execute(run rt.Runtime) (err error) {
	if e := op.reduce(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Reduce) reduce(run rt.Runtime) (err error) {
	if fromList, e := core.GetAssignedValue(run, op.FromList); e != nil {
		err = e
	} else if outVal, e := safe.CheckVariable(run, op.IntoValue, ""); e != nil {
		err = e
	} else if e := op.cacheKinds(run, op.UsingPattern); e != nil {
		err = e
	} else {
		ps := op.pk.NewRecord()
		for it := g.ListIt(fromList); it.HasNext(); {
			in, out := 0, 1
			if inVal, e := it.GetNext(); e != nil {
				err = e
				break
			} else if e := ps.SetFieldByIndex(in, inVal); e != nil {
				err = e
				break
			} else if e := ps.SetFieldByIndex(out, outVal); e != nil {
				err = e
				break
			} else if e := op.call(run, ps); e != nil {
				err = e
				break
			} else if newVal, e := ps.GetFieldByIndex(out); e != nil {
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
