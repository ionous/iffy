package pattern

import (
	"strconv"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
	"github.com/ionous/iffy/rt/stream"
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

// Stitch finds the pattern, builds the scope, and executes the passed callback to generate a result.
// It's an adapter from the the specific DetermineActivity, DetermineNumber, etc. statements.
func (op *FromPattern) Stitch(run rt.Runtime, patType string, fn func(p interface{}) error) (err error) {
	// find the pattern (p), qna's implementation assembles the rules by querying the db.
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
		if rules, ok := p.([]*ExecuteRule); !ok {
			err = errutil.New("Pattern", op.Pattern, "not an activity")
		} else if inds, e := splitExe(run, rules); e != nil {
			err = e
		} else {
			for _, i := range inds {
				if e := rt.RunOne(run, rules[i].Execute); e != nil {
					err = e
					break
				}
				// NOTE: if we need to differentiate between "ran" and "not found",
				// "didnt run" should probably become an error code.
			}
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

// GetNumber returns the first matching num evaluation.
func (op *DetermineNum) GetNumber(run rt.Runtime) (ret float64, err error) {
	err = (*FromPattern)(op).Stitch(run, object.NumberRule, func(p interface{}) (err error) {
		if rules, ok := p.([]*NumberRule); !ok {
			err = errutil.New("Pattern", op.Pattern, "not a number")
		} else {
			for i, cnt := 0, len(rules); i < cnt; i++ {
				p := rules[cnt-i-1]
				if matched, e := rt.GetOptionalBool(run, p.Filter, true); e != nil {
					err = e
					break
				} else if matched {
					ret, err = rt.GetNumber(run, p.NumberEval)
					break
				}
			}
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

// GetText returns the first matching text evaluation.
func (op *DetermineText) GetText(run rt.Runtime) (ret string, err error) {
	err = (*FromPattern)(op).Stitch(run, object.TextRule, func(p interface{}) (err error) {
		if rules, ok := p.([]*TextRule); !ok {
			err = errutil.New("Pattern", op.Pattern, "not text")
		} else {
			for i, cnt := 0, len(rules); i < cnt; i++ {
				p := rules[cnt-i-1]
				if matched, e := rt.GetOptionalBool(run, p.Filter, true); e != nil {
					err = e
					break
				} else if matched {
					ret, err = rt.GetText(run, p.TextEval)
					break
				}
			}
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

// GetBool returns the first matching bool evaluation.
func (op *DetermineBool) GetBool(run rt.Runtime) (ret bool, err error) {
	err = (*FromPattern)(op).Stitch(run, object.BoolRule, func(p interface{}) (err error) {
		if rules, ok := p.([]*BoolRule); !ok {
			err = errutil.New("Pattern", op.Pattern, "not a boolean")
		} else {
			for i, cnt := 0, len(rules); i < cnt; i++ {
				p := rules[cnt-i-1]
				if matched, e := rt.GetOptionalBool(run, p.Filter, true); e != nil {
					err = e
					break
				} else if matched {
					ret, err = rt.GetBool(run, p.BoolEval)
					break
				}
			}
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
		if rules, ok := p.([]*NumListRule); !ok {
			err = errutil.New("Pattern", op.Pattern, "not a boolean")
		} else if inds, e := splitNumbers(run, rules); e != nil {
			err = e
		} else {
			it := numIterator{run, rules, inds, 0}
			ret = stream.NewNumberChain(&it)
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
		if rules, ok := p.([]*TextListRule); !ok {
			err = errutil.New("Pattern", op.Pattern, "not a boolean")
		} else if inds, e := splitText(run, rules); e != nil {
			err = e
		} else {
			it := textIterator{run, rules, inds, 0}
			ret = stream.NewTextChain(&it)
		}
		return
	})
	return
}
