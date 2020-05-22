package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/postfix"
)

// ExpressionStateFactory is used for testing with mocks.
// when no factory is specified, the pipe or subdirective parser is used depending on context.
type ExpressionStateFactory interface {
	NewExpressionState() ExpressionState
}

type ExpressionState interface {
	State
	GetExpression() (postfix.Expression, error)
}

// DirectiveParser reads a key-expression pair where both elements are optional.
// ( compare to KeyParser which always has a key, followed by an optional expression. )
type DirectiveParser struct {
	factory ExpressionStateFactory
	out     Directive
	err     error
}

func (p *DirectiveParser) StateName() string {
	return "directive parser"
}

func (p *DirectiveParser) GetDirective() (Directive, error) {
	return p.out, p.err
}

// rune at the start of a directive's content.
func (p *DirectiveParser) NewRune(r rune) State {
	keyp := KeyParser{exp: p.newExpressionParser()}
	expp := p.newExpressionParser()
	//
	para := MakeParallel("key or expression",
		MakeChain(expp, StateExit("expression", func() {
			if exp, e := expp.GetExpression(); e != nil {
				p.err = errutil.Append(p.err, e)
			} else if len(exp) > 0 {
				p.out = Directive{Expression: exp} // last match wins
			}
		})),
		MakeChain(&keyp, StateExit("key", func() {
			if d, e := keyp.GetDirective(); e != nil {
				p.err = errutil.Append(p.err, e)
			} else {
				p.out = d // last match wins
			}
		})),
	)
	return para.NewRune(r)
}

func (p *DirectiveParser) newExpressionParser() (ret ExpressionState) {
	if p.factory != nil {
		ret = p.factory.NewExpressionState()
	} else {
		ret = new(PipeParser)
	}
	return
}
