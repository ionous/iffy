package internal

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/pattern"
	"github.com/ionous/iffy/rt"
)

type ruleBuilder struct {
	name      string // ex. NumberRule
	rule      interface{}
	filterPtr *rt.BoolEval
	filters   []rt.BoolEval
}

func (rb *ruleBuilder) addFilter(eval rt.BoolEval) {
	rb.filters = append(rb.filters, eval)
}

func (rb *ruleBuilder) distill() interface{} {
	// build the filters
	var w rt.BoolEval
	switch cnt := len(rb.filters); cnt {
	case 0:
		// no filters? then always true
		w = &core.Bool{true}
	case 1:
		w = rb.filters[0]
	default:
		// if there are multiple elements in the slice --
		// return a AllTrue
		w = &core.AllTrue{rb.filters}
	}
	(*rb.filterPtr) = w
	// return the rule itself
	return rb.rule
}

func newBoolRule(i rt.BoolEval) *ruleBuilder {
	ptr := &pattern.BoolRule{BoolEval: i}
	return &ruleBuilder{"bool_rule", ptr, &ptr.Filter, nil}
}
func newNumberRule(i rt.NumberEval) *ruleBuilder {
	ptr := &pattern.NumberRule{NumberEval: i}
	return &ruleBuilder{"number_rule", ptr, &ptr.Filter, nil}
}
func newTextRule(i rt.TextEval) *ruleBuilder {
	ptr := &pattern.TextRule{TextEval: i}
	return &ruleBuilder{"text_rule", ptr, &ptr.Filter, nil}
}
func newExecuteRule(exes []rt.Execute) *ruleBuilder {
	ptr := &pattern.ExecuteRule{Go: exes}
	return &ruleBuilder{"exec_rule", ptr, &ptr.Filter, nil}
}
func newTextListRule(i rt.TextListEval) *ruleBuilder {
	ptr := &pattern.TextListRule{TextListEval: i}
	return &ruleBuilder{"text_list_rule", ptr, &ptr.Filter, nil}
}
func newNumListRule(i rt.NumListEval) *ruleBuilder {
	ptr := &pattern.NumListRule{NumListEval: i}
	return &ruleBuilder{"num_list_rule", ptr, &ptr.Filter, nil}
}
