package core

import (
	"strconv"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

// FromPattern helps runs a pattern
type FromPattern struct {
	Pattern    string      // pattern name, i guess a text eval here would be like a function pointer.
	Parameters *Parameters // for optional parameters ( and formatting )
}

type DetermineAct FromPattern
type DetermineNum FromPattern
type DetermineText FromPattern
type DetermineBool FromPattern
type DetermineNumList FromPattern
type DetermineTextList FromPattern

type Parameters struct {
	Params []*Parameter
}

type Parameter struct {
	Name string // parameter name
	From Assignment
}

func (*Parameters) Compose() composer.Spec {
	return composer.Spec{
		Name:  "parameters",
		Spec:  " when {parameters%params+parameter}",
		Group: "patterns",
	}
}

func (*Parameter) Compose() composer.Spec {
	return composer.Spec{
		Name:  "parameter",
		Spec:  "its {name:variable_name} is {from:assignment}",
		Group: "patterns",
	}
}

// Stitch find the pattern, builds the scope, and executes the passed callback to generate a result.
// Its an adapter from the the specific DetermineActivity, DetermineNumber, etc. statements.
func (op *FromPattern) Stitch(run rt.Runtime, patType string, fn func(p interface{}) error) (err error) {
	// find the pattern (p):
	patName := op.Pattern
	if pat, e := run.GetField(patName, patType); e != nil {
		err = e
	} else {
		// bake the parameters down
		parms := make(scope.Parameters)
		if op.Parameters != nil {
			// read from each argument
			for _, param := range op.Parameters.Params {
				if e := param.From.Assign(run, func(i interface{}) (err error) {
					if n, e := unpack(run, patName, param.Name); e != nil {
						err = e
					} else {
						err = parms.SetVariable(n, i)
					}
					return
				}); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
		if err == nil {
			run.PushScope(&parms)
			err = fn(pat)
			run.PopScope()
		}
	}
	return
}

// change a param name ( which could be an index ) into a valid param name
func unpack(run rt.Runtime, pattern, param string) (ret string, err error) {
	if usesIndex := len(param) > 1 && param[:1] == "$"; !usesIndex {
		ret = param
	} else if idx, e := strconv.Atoi(param[1:]); e != nil {
		err = e
	} else {
		ret, err = run.GetFieldByIndex(pattern, idx)
	}
	return
}

func (*DetermineAct) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_act",
		Group: "patterns",
		Desc:  "Determine an activity",
	}
}

func (op *DetermineAct) Execute(run rt.Runtime) (err error) {
	err = (*FromPattern)(op).Stitch(run, object.ExecuteRule, func(p interface{}) (err error) {
		// cast the pattern to an execute
		// fix: this may require a []cast instead.
		if exe, ok := p.(rt.Execute); !ok {
			err = errutil.New("Pattern", op.Pattern, "not an activity")
		} else {
			err = exe.Execute(run)
		}
		return
	})
	return
}

func (*DetermineNum) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_num",
		Spec:  "the {number pattern%name:pattern_name}{?parameters}",
		Group: "patterns",
		Desc:  "Determine a number",
	}
}

func (op *DetermineNum) GetNumber(run rt.Runtime) (ret float64, err error) {
	err = (*FromPattern)(op).Stitch(run, object.NumberRule, func(p interface{}) (err error) {
		if eval, ok := p.(rt.NumberEval); !ok {
			err = errutil.New("Pattern", op.Pattern, "not a number")
		} else {
			ret, err = rt.GetNumber(run, eval)
		}
		return
	})
	return
}

func (*DetermineText) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_text",
		Spec:  "the {text pattern%name:pattern_name}{?parameters}",
		Group: "patterns",
		Desc:  "Determine some text",
	}
}

func (op *DetermineText) GetText(run rt.Runtime) (ret string, err error) {
	err = (*FromPattern)(op).Stitch(run, object.TextRule, func(p interface{}) (err error) {
		if eval, ok := p.(rt.TextEval); !ok {
			err = errutil.New("Pattern", op.Pattern, "not text")
		} else {
			ret, err = rt.GetText(run, eval)
		}
		return
	})
	return
}

func (*DetermineBool) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_bool",
		Spec:  "the {true/false pattern%name:pattern_name}{?parameters}",
		Group: "patterns",
		Desc:  "Determine a true/false value",
	}
}

func (op *DetermineBool) GetBool(run rt.Runtime) (ret bool, err error) {
	err = (*FromPattern)(op).Stitch(run, object.BoolRule, func(p interface{}) (err error) {
		if eval, ok := p.(rt.BoolEval); !ok {
			err = errutil.New("Pattern", op.Pattern, "not a boolean")
		} else {
			ret, err = rt.GetBool(run, eval)
		}
		return
	})
	return
}

func (*DetermineNumList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_num_list",
		Spec:  "the {number list pattern%name:pattern_name}{?parameters}",
		Group: "patterns",
		Desc:  "Determine a list of numbers",
	}
}

func (op *DetermineNumList) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	err = (*FromPattern)(op).Stitch(run, object.NumListRule, func(p interface{}) (err error) {
		if eval, ok := p.(rt.NumListEval); !ok {
			err = errutil.New("Pattern", op.Pattern, "not a boolean")
		} else {
			ret, err = rt.GetNumberStream(run, eval)
		}
		return
	})
	return
}

func (*DetermineTextList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_text_list",
		Spec:  "the {text list pattern%name:pattern_name}{?parameters}",
		Group: "patterns",
		Desc:  "Determine a list of text",
	}
}

func (op *DetermineTextList) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	err = (*FromPattern)(op).Stitch(run, object.TextListRule, func(p interface{}) (err error) {
		if eval, ok := p.(rt.TextListEval); !ok {
			err = errutil.New("Pattern", op.Pattern, "not a boolean")
		} else {
			ret, err = rt.GetTextStream(run, eval)
		}
		return
	})
	return
}
