package pattern

import "github.com/ionous/iffy/rt"

// BoolRules finds the first matched rule and returns the result of that evaluation.
type BoolRules []*BoolRule

// BoolRule responds with a true/false result when its filters are satisfied.
// It implements rt.BoolEval.
type BoolRule struct {
	Filter rt.BoolEval
	rt.BoolEval
}

// NumberRules finds the first matched rule and returns the result of that evaluation.
// It implements rt.NumberEval.
type NumberRules []*NumberRule

// NumberRule responds with a single number when its filters are satisfied.
type NumberRule struct {
	Filter rt.BoolEval
	rt.NumberEval
}

// TextRules finds the first matched rule and returns the result of that evaluation.
// It implements rt.TextEval.
type TextRules []*TextRule

// TextRule responds with a bit of text when its filters are satisfied.
type TextRule struct {
	Filter rt.BoolEval
	rt.TextEval
}

func (*BoolRule) RuleDesc() RuleDesc {
	return RuleDesc{
		"bool_rule",
		(*BoolRules)(nil),
	}
}

func (*NumberRule) RuleDesc() RuleDesc {
	return RuleDesc{
		"number_rule",
		(*NumberRules)(nil),
	}
}

func (*TextRule) RuleDesc() RuleDesc {
	return RuleDesc{
		"text_rule",
		(*TextRules)(nil),
	}
}

// GetBool returns the first matching bool evaluation.
func (ps BoolRules) GetBool(run rt.Runtime) (ret bool, err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		p := ps[cnt-i-1]
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
func (ps NumberRules) GetNumber(run rt.Runtime) (ret float64, err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		p := ps[cnt-i-1]
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
func (ps TextRules) GetText(run rt.Runtime) (ret string, err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		p := ps[cnt-i-1]
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
