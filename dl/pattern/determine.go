package pattern

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"

	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/scope"
)

// FromPattern helps runs a pattern
type FromPattern struct {
	Pattern   PatternName     // pattern name, i guess a text eval here would be like a function pointer.
	Arguments *core.Arguments // arguments passed to the pattern. kept as a pointer for composer formatting...
	// each is a name targeting some parameter, and an "assignment"
	ps, ls *g.Kind
}

type DetermineAct FromPattern
type DetermineNum FromPattern
type DetermineText FromPattern
type DetermineBool FromPattern
type DetermineNumList FromPattern
type DetermineTextList FromPattern

// Stitch finds the pattern, builds the scope, and executes the passed callback to generate a result.
// It's an adapter from the the specific DetermineActivity, DetermineNumber, etc. statements.
func (op *FromPattern) Stitch(run rt.Runtime, pat Pattern, fn func() error) (err error) {
	// find the pattern (p), qna's implementation assembles the rules by querying the db.
	if e := run.GetEvalByName(op.Pattern.String(), pat); e != nil {
		err = e
	} else if pk, e := op.newParams(run, pat); e != nil {
		err = e
	} else if lk, e := op.newLocals(run, pat); e != nil {
		err = e
	} else {
		ps, ls := pk.NewRecord(), lk.NewRecord()
		if op.Arguments != nil {
			// read from each argument and store into the parameters
			err = op.Arguments.Distill(run, ps)
		}
		if err == nil {
			// fix: dont love the double map creation, double scope....
			run.PushScope(&scope.TargetRecord{object.Variables, ps})
			run.PushScope(&scope.TargetRecord{object.Variables, ls})
			err = fn()
			run.PopScope()
			run.PopScope()
		}
	}
	return
}

func (op *FromPattern) newParams(run rt.Runtime, pat Pattern) (ret *g.Kind, err error) {
	// create variables for all the known parameters
	if op.ps != nil {
		ret = op.ps
	} else {
		var parms term.Terms
		if e := pat.ComputeParams(run, &parms); e != nil {
			err = e
		} else {
			ps := parms.NewKind(run)
			ret, op.ps = ps, ps
		}
	}
	return
}
func (op *FromPattern) newLocals(run rt.Runtime, pat Pattern) (ret *g.Kind, err error) {
	// create variables for all the known parameters
	if op.ls != nil {
		ret = op.ls
	} else {
		var locals term.Terms
		if e := pat.ComputeLocals(run, &locals); e != nil {
			err = e
		} else {
			ls := locals.NewKind(run)
			ret, op.ls = ls, ls
		}
	}
	return
}

func (*DetermineAct) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_act",
		Spec:  "determine {activity%name:pattern_name}{?arguments}",
		Group: "patterns",
		Desc:  "Determine an activity",
		Stub:  true,
	}
}

func (op *DetermineAct) Execute(run rt.Runtime) (err error) {
	var pat ActivityPattern
	if e := (*FromPattern)(op).Stitch(run, &pat, func() (err error) {
		err = pat.Execute(run)
		return
	}); e != nil {
		err = cmdErrorCtx(op, op.Pattern.String(), e)
	}
	return
}

func (*DetermineNum) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_num",
		Spec:  "{number pattern%name:pattern_name}{?arguments}",
		Group: "patterns",
		Desc:  "Determine a number",
		Stub:  true,
	}
}

// GetNumber returns the first matching num evaluation.
func (op *DetermineNum) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	var pat NumberPattern
	if e := (*FromPattern)(op).Stitch(run, &pat, func() (err error) {
		ret, err = pat.GetNumber(run)
		return
	}); e != nil {
		err = cmdErrorCtx(op, op.Pattern.String(), e)
	}
	return
}

func (*DetermineText) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_text",
		Spec:  "{text pattern%name:pattern_name}{?arguments}",
		Group: "patterns",
		Desc:  "Determine some text",
		Stub:  true,
	}
}

// GetText returns the first matching text evaluation.
func (op *DetermineText) GetText(run rt.Runtime) (ret g.Value, err error) {
	var pat TextPattern
	if e := (*FromPattern)(op).Stitch(run, &pat, func() (err error) {
		ret, err = pat.GetText(run)
		return
	}); e != nil {
		err = cmdErrorCtx(op, op.Pattern.String(), e)
	}
	return
}

func (*DetermineBool) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_bool",
		Spec:  "{true/false pattern%name:pattern_name}{?arguments}",
		Group: "patterns",
		Desc:  "Determine a true/false value",
		Stub:  true,
	}
}

// GetBool returns the first matching bool evaluation.
func (op *DetermineBool) GetBool(run rt.Runtime) (ret g.Value, err error) {
	var pat BoolPattern
	if e := (*FromPattern)(op).Stitch(run, &pat, func() (err error) {
		ret, err = pat.GetBool(run)
		return
	}); e != nil {
		err = cmdErrorCtx(op, op.Pattern.String(), e)
	}
	return
}

func (*DetermineNumList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_num_list",
		Spec:  "{number list pattern%name:pattern_name}{?arguments}",
		Group: "patterns",
		Desc:  "Determine a list of numbers",
		Stub:  true,
	}
}

func (op *DetermineNumList) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	var pat NumListPattern
	if e := (*FromPattern)(op).Stitch(run, &pat, func() (err error) {
		ret, err = pat.GetNumList(run)
		return
	}); e != nil {
		err = cmdErrorCtx(op, op.Pattern.String(), e)
	}
	return
}

func (*DetermineTextList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_text_list",
		Spec:  "{text list pattern%name:pattern_name}{?arguments}",
		Group: "patterns",
		Desc:  "Determine a list of text",
		Stub:  true,
	}
}

func (op *DetermineTextList) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	var pat TextListPattern
	if e := (*FromPattern)(op).Stitch(run, &pat, func() (err error) {
		ret, err = pat.GetTextList(run)
		return
	}); e != nil {
		err = cmdErrorCtx(op, op.Pattern.String(), e)
	}
	return
}
