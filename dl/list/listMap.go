package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/scope"
)

type Map struct {
	FromList, ToList, Pattern string // variable names
	pk, lk                    *g.Kind
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
	var pat pattern.ActivityPattern
	if fromList, e := run.GetField(object.Variables, op.FromList); e != nil {
		err = e
	} else if toList, e := run.GetField(object.Variables, op.ToList); e != nil {
		err = e
	} else if e := run.GetEvalByName(op.Pattern, &pat); e != nil {
		err = e
	} else if pk, e := op.newParams(run, &pat); e != nil {
		err = e
	} else if lk, e := op.newLocals(run, &pat); e != nil {
		err = e
	} else if in := pk.FieldIndex("in"); in < 0 {
		err = errutil.New("pattern expected an 'in' parameter")
	} else if out := pk.FieldIndex("out"); out < 0 {
		err = errutil.New("pattern expected an 'out' parameter")
	} else {
		for it := g.ListIt(fromList); it.HasNext(); {
			ps, ls := pk.NewRecord(), lk.NewRecord()
			if inVal, e := it.GetNext(); e != nil {
				err = e
				break
			} else if e := ps.SetFieldByIndex(in, inVal); e != nil {
				err = e
				break
			} else if e := op.mapOne(run, &pat, ps, ls); e != nil {
				err = e
				break
			} else if outVal, e := ps.GetFieldByIndex(out); e != nil {
				err = e
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

func (op *Map) newParams(run rt.Runtime, pat *pattern.ActivityPattern) (ret *g.Kind, err error) {
	// create variables for all the known parameters
	if op.pk != nil && !op.pk.IsStaleKind(run) {
		ret = op.pk
	} else {
		var parms term.Terms
		if e := pat.Prepare(run, &parms); e != nil {
			err = e
		} else {
			pk := parms.NewKind(run)
			ret, op.pk = pk, pk
		}
	}
	return
}

func (op *Map) newLocals(run rt.Runtime, pat *pattern.ActivityPattern) (ret *g.Kind, err error) {
	// create variables for all the known parameters
	if op.lk != nil && !op.lk.IsStaleKind(run) {
		ret = op.lk
	} else {
		var locals term.Terms
		if e := pat.ComputeLocals(run, &locals); e != nil {
			err = e
		} else {
			lk := locals.NewKind(run)
			ret, op.lk = lk, lk
		}
	}
	return
}

// see also pattern.Stitch
func (op *Map) mapOne(run rt.Runtime, pat *pattern.ActivityPattern, ps, ls *g.Record) (err error) {
	run.PushScope(&scope.TargetRecord{object.Variables, ps})
	run.PushScope(&scope.TargetRecord{object.Variables, ls})
	err = pat.Execute(run)
	run.PopScope()
	run.PopScope()
	return
}
