package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/scope"
)

// temp: to call activities with values that we can manipulate semi-efficiently
// we should be able to build the kinds in the assembler
type activityCache struct {
	pat    pattern.Pattern
	pk, lk *g.Kind
	run    rt.Runtime
}

// see also pattern.CreateScope
func (op *activityCache) call(run rt.Runtime, ps *g.Record) (err error) {
	ls := op.lk.NewRecord()
	run.PushScope(&scope.TargetRecord{object.Variables, ps})
	run.PushScope(&scope.TargetRecord{object.Variables, ls})
	if e := op.pat.Execute(run); e != nil {
		err = errutil.Append(errutil.Fmt("error in pattern %q", op.pat.Name), e)
	}
	run.PopScope()
	run.PopScope()
	return
}

func (op *activityCache) cacheKinds(run rt.Runtime, pat pattern.PatternName) (err error) {
	if run != op.run {
		if e := run.GetEvalByName(pat.String(), &op.pat); e != nil {
			err = e
		} else if pk, e := op.newParams(run); e != nil {
			err = e
		} else if lk, e := op.newLocals(run); e != nil {
			err = e
		} else if pk.NumField() < 2 {
			err = errutil.New("pattern expected at least two parameters, an input and an output")
		} else {
			op.pk, op.lk = pk, lk
			op.run = run
		}
	}
	return
}

func (op *activityCache) newParams(run rt.Runtime) (ret *g.Kind, err error) {
	// create variables for all the known parameters
	if op.pk != nil && !op.pk.IsStaleKind(run) {
		ret = op.pk
	} else {
		var parms term.Terms
		if _, e := op.pat.ComputeParams(run, &parms); e != nil {
			err = e
		} else {
			pk := parms.NewKind(run)
			ret, op.pk = pk, pk
		}
	}
	return
}

func (op *activityCache) newLocals(run rt.Runtime) (ret *g.Kind, err error) {
	// create variables for all the known parameters
	if op.lk != nil && !op.lk.IsStaleKind(run) {
		ret = op.lk
	} else {
		var locals term.Terms
		if _, e := op.pat.ComputeLocals(run, &locals); e != nil {
			err = e
		} else {
			lk := locals.NewKind(run)
			ret, op.lk = lk, lk
		}
	}
	return
}
