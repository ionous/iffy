package pattern

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/rt"

	g "github.com/ionous/iffy/rt/generic"
)

// Determine helps run a pattern.
// It implements every core evaluation,
// erroring if the value requested doesnt support the error returned.
type Determine struct {
	Pattern   PatternName     // a text eval here would be like a function pointer maybe...
	Arguments *core.Arguments // pattern args kept as a pointer for composer formatting...
}

// Stitch finds the pattern, builds the scope, and executes the passed callback to generate a result.
// It's an adapter from the the specific DetermineActivity, DetermineNumber, etc. statements.
func (op *Determine) newScope(run rt.Runtime, pat *Pattern) (ret *patScope, err error) {
	// find the pattern (p), qna's implementation assembles the rules by querying the db.
	if e := run.GetEvalByName(op.Pattern.String(), pat); e != nil {
		err = e
	} else {
		// pack down parameters and locals into a single set of "terms"
		var ts term.Terms
		if _, e := pat.ComputeParams(run, &ts); e != nil {
			err = e
		} else if _, e := pat.ComputeLocals(run, &ts); e != nil {
			err = e
		} else if retVar, e := pat.ComputeReturn(run, &ts); e != nil {
			err = e
		} else {
			// create a kind and a record from that.
			// ( fix: should be happening in the assembler )
			vs := ts.NewKind(run).NewRecord()
			if op.Arguments != nil {
				// read from each argument and store into the parameters
				err = op.Arguments.Distill(run, ts, vs)
			}
			if err == nil {
				ret = newScope(retVar, vs)
			}
		}
	}
	return
}

func (op *Determine) runPattern(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	var pat Pattern
	if x, e := op.newScope(run, &pat); e != nil {
		err = e
	} else {
		run.PushScope(x)
		if e := pat.Execute(run); e != nil {
			err = e
		} else if res, e := x.GetValue(aff); e != nil {
			err = e
		} else {
			ret = res
		}
		run.PopScope()
	}
	return
}

func (op *Determine) determine(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if v, e := op.runPattern(run, aff); e != nil {
		err = cmdErrorCtx(op, op.Pattern.String(), e)
	} else {
		ret = v
	}
	return
}

func (*Determine) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine",
		Spec:  "{pattern%name:pattern_name}{?arguments}",
		Group: "patterns",
		Desc:  "Runs a pattern, and potentially returns a value.",
		Stub:  true,
	}
}

func (op *Determine) Execute(run rt.Runtime) error {
	_, err := op.determine(run, "")
	return err
}

// GetNumber returns the first matching num evaluation.
func (op *Determine) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Number)
}

// GetText returns the first matching text evaluation.
func (op *Determine) GetText(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Text)
}

// GetBool returns the first matching bool evaluation.
func (op *Determine) GetBool(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Bool)
}

func (op *Determine) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Record)
}

func (op *Determine) GetNumList(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.NumList)
}

func (op *Determine) GetTextList(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.TextList)
}

func (op *Determine) GetRecordList(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.RecordList)
}
