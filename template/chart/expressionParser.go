package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/postfix"
)

// ExpressionParser reads either a single call or a series of operations.
type ExpressionParser struct {
	out        postfix.Expression
	err        error
	argFactory ExpressionStateFactory // for testing
}

func (p ExpressionParser) GetExpression() (postfix.Expression, error) {
	return p.out, p.err
}

// NewRune starts with the first character of a directive's content.
func (p *ExpressionParser) NewRune(r rune) State {
	call := MakeCallParser(0, p.argFactory)
	series := SeriesParser{}
	para := MakeParallel(
		MakeChain(&series, StateExit(func() {
			if x, e := series.GetExpression(); e != nil {
				p.err = errutil.Append(p.err, e)
			} else if len(x) > 0 {
				p.out, p.err = x, nil // longest match wins; noting that operators can read too far ahead
				// particularly an issue with ! which provisionally matches !=
			}
		})),
		MakeChain(&call, StateExit(func() {
			if x, e := call.GetExpression(); e != nil {
				p.err = errutil.Append(p.err, e)
			} else if len(x) > 0 {
				p.out, p.err = x, nil // for equal matches, call wins.
			}
		})),
	)
	return para.NewRune(r)
}
