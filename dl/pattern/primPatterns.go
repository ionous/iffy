package pattern

import (
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// BoolPattern finds the first matched rule and returns the result of that evaluation.
type xBoolPattern struct {
	CommonPattern
	Rules []*xBoolRule
}

// NumberPattern finds the first matched rule and returns the result of that evaluation.
// It implements rt.NumberEval.
type xNumberPattern struct {
	CommonPattern
	Rules []*xNumberRule
}

// TextPattern finds the first matched rule and returns the result of that evaluation.
// It implements rt.TextEval.
type xTextPattern struct {
	CommonPattern
	Rules []*xTextRule
}

// GetBool returns the first matching bool evaluation.
func (ps *xBoolPattern) GetBool(run rt.Runtime) (ret g.Value, err error) {
	ret = g.False // provisionally
	for i, cnt := 0, len(ps.Rules); i < cnt; i++ {
		p := ps.Rules[cnt-i-1]
		if matched, e := safe.GetOptionalBool(run, p.Filter, true); e != nil {
			err = e
			break
		} else if matched.Bool() {
			ret, err = safe.GetBool(run, p.BoolEval)
			break
		}
	}
	return
}

// GetNumber returns the first matching num evaluation.
func (ps *xNumberPattern) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	ret = g.Zero // provisionally
	for i, cnt := 0, len(ps.Rules); i < cnt; i++ {
		p := ps.Rules[cnt-i-1]
		if matched, e := safe.GetOptionalBool(run, p.Filter, true); e != nil {
			err = e
			break
		} else if matched.Bool() {
			ret, err = safe.GetNumber(run, p.NumberEval)
			break
		}
	}
	return
}

// GetText returns the first matching text evaluation.
func (ps *xTextPattern) GetText(run rt.Runtime) (ret g.Value, err error) {
	ret = g.Empty // provisionally
	for i, cnt := 0, len(ps.Rules); i < cnt; i++ {
		p := ps.Rules[cnt-i-1]
		if matched, e := safe.GetOptionalBool(run, p.Filter, true); e != nil {
			err = e
			break
		} else if matched.Bool() {
			ret, err = safe.GetText(run, p.TextEval)
			break
		}
	}
	return
}
