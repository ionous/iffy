package list

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type Map struct {
	FromList, ToList, Pattern string // variable names
}

func (*Map) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_map",
		Group: "list",
		Desc: `Map List: Transform the values from one list and place the results in another list.
		The named pattern is called with two parameters: 'in' and 'out'`,
	}
}
func (op *Map) Execute(run rt.Runtime) (err error) {
	if e := op.execute(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Map) execute(run rt.Runtime) (err error) {
	var pat pattern.ActivityPattern
	if fromList, e := run.GetField(object.Variables, op.FromList); e != nil {
		err = e
	} else if toList, e := run.GetField(object.Variables, op.ToList); e != nil {
		err = e
	} else if e := run.GetEvalByName(op.Pattern, &pat); e != nil {
		err = e
	} else {
		a, t := affine.Element(toList.Affinity()), toList.Type()
		for it := g.ListIt(fromList); it.HasNext(); {
			if inVal, e := it.GetNext(); e != nil {
				err = e
				break
			} else if defVal, e := g.DefaultFor(run, a, t); e != nil {
				err = e
				break
			} else if outVal, e := op.mapOne(run, pat, inVal, defVal); e != nil {
				err = e
				break
			} else if vs, e := toList.Append(outVal); e != nil {
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

// see also pattern.Stitch
func (op *Map) mapOne(run rt.Runtime, pat pattern.ActivityPattern, inVal, outVal g.Value) (ret g.Value, err error) {
	var parms term.Terms
	if e := pat.Prepare(run, &parms); e != nil {
		err = e
	} else if e := parms.SetField(object.Variables, "in", inVal); e != nil {
		err = e
	} else if e := parms.SetField(object.Variables, "out", outVal); e != nil {
		err = e
	} else {
		run.PushScope(&parms)
		var locals term.Terms
		if e := pat.ComputeLocals(run, &locals); e != nil {
			err = e
		} else {
			run.PushScope(&locals)
			if e := pat.Execute(run); e != nil {
				err = e
			} else if v, e := parms.GetField(object.Variables, "out"); e != nil {
				err = e
			} else {
				ret = v
			}
			run.PopScope()
		}
		run.PopScope()

	}
	return
}
