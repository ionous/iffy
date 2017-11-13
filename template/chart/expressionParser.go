package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/postfix"
)

// ExpressionParser reads either a call or a sequence.
type ExpressionParser struct {
	argFactory ExpressionStateFactory
	out        postfix.Expression
	err        error
}

func MakeExpressionParser(f ExpressionStateFactory) ExpressionParser {
	return ExpressionParser{argFactory: f}
}

func (p ExpressionParser) GetExpression() (postfix.Expression, error) {
	return p.out, p.err
}

// NewRune starts with the first character past the bar
func (p *ExpressionParser) NewRune(r rune) State {
	call := MakeCallParser(0, p.argFactory)
	seqp := SeriesParser{}
	para := MakeParallel(
		MakeChain(&call, StateExit(func() {
			if x, e := call.GetExpression(); e != nil {
				p.err = errutil.Append(p.err, e)
			} else if len(x) > 0 {
				p.out = x // longest match wins
			}
		})),
		MakeChain(&seqp, StateExit(func() {
			if x, e := seqp.GetExpression(); e != nil {
				p.err = errutil.Append(p.err, e)
			} else if len(x) > 0 {
				p.out = x // longest match wins
			}
		})),
	)
	return para.NewRune(r)
}
