package pattern

import "github.com/ionous/iffy/rt"

// BoolPattern finds the first matched rule and returns the result of that evaluation.
type BoolPattern struct {
	Name  string
	Rules []*BoolRule
}

// NumberPattern finds the first matched rule and returns the result of that evaluation.
// It implements rt.NumberEval.
type NumberPattern struct {
	Name  string
	Rules []*NumberRule
}

// TextPattern finds the first matched rule and returns the result of that evaluation.
// It implements rt.TextEval.
type TextPattern struct {
	Name  string
	Rules []*TextRule
}

// GetBool returns the first matching bool evaluation.
func (ps BoolPattern) GetBool(run rt.Runtime) (ret bool, err error) {
	for i, cnt := 0, len(ps.Rules); i < cnt; i++ {
		p := ps.Rules[cnt-i-1]
		if matched, e := rt.GetOptionalBool(run, p.Filter, true); e != nil {
			err = e
			break
		} else if matched {
			ret, err = rt.GetBool(run, p.BoolEval)
			break
		}
	}
	return
}

// GetNumber returns the first matching num evaluation.
func (ps NumberPattern) GetNumber(run rt.Runtime) (ret float64, err error) {
	for i, cnt := 0, len(ps.Rules); i < cnt; i++ {
		p := ps.Rules[cnt-i-1]
		if matched, e := rt.GetOptionalBool(run, p.Filter, true); e != nil {
			err = e
			break
		} else if matched {
			ret, err = rt.GetNumber(run, p.NumberEval)
			break
		}
	}
	return
}

// GetText returns the first matching text evaluation.
func (ps *TextPattern) GetText(run rt.Runtime) (ret string, err error) {
	for i, cnt := 0, len(ps.Rules); i < cnt; i++ {
		p := ps.Rules[cnt-i-1]
		if matched, e := rt.GetOptionalBool(run, p.Filter, true); e != nil {
			err = e
			break
		} else if matched {
			ret, err = rt.GetText(run, p.TextEval)
			break
		}
	}
	return
}
