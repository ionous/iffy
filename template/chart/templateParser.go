package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/postfix"
)

type TemplateParser struct {
	out, pending Section
	err          error
	delegate     Delegate
}

func MakeTemplateParser() TemplateParser {
	return TemplateParser{delegate: defaultParser}
}

func MakeSubParser(d Delegate, xs postfix.Expression) TemplateParser {
	t := TemplateParser{delegate: d}
	t.out.Append(xs)
	return t
}

type Delegate func(*TemplateParser, Directive) (State, error)

func (p *TemplateParser) GetExpression() (ret postfix.Expression, err error) {
	return p.reduce(template.Span)
}

// words { directive } words { directive }
func (p *TemplateParser) NewRune(r rune) State {
	var left LeftParser
	return ParseChain(r, &left, Statement(func(r rune) State {
		if text := left.GetText(); len(text) > 0 {
			// println("got text", text)
			p.pending.Append(quote(text))
		}
		return ParseChain(r, spaces, Statement(func(r rune) (ret State) {
			if r != eof {
				var right RightParser
				ret = ParseChain(r, &right, Statement(func(r rune) (ret State) {
					if v, e := right.GetDirective(); e != nil {
						p.err = e
					} else {
						if len(v.Key) == 0 {
							p.pending.Append(v.Expression)
							ret = p.NewRune(r) // loop back to left half.
						} else if next, e := p.delegate(p, v); e != nil {
							p.err = e
						} else if next != nil {
							ret = next.NewRune(r) // delegate this rune to the next handler.
						}
						// if a rune is unhandled; we "rewind" to the parent state.
					}
					return
				}))
			}
			return
		}))
	}))
}

// reduce returns the output of the template as a command of the passed type:
// a span of elements, an if statement with branches, a sequence with cycling text, etc.
func (p *TemplateParser) reduce(kind template.BuiltinType) (ret postfix.Expression, err error) {
	if p.err != nil {
		err = p.err
	} else {
		p.endSection(false)
		ret = p.out.Reduce(kind)
	}
	return
}

// forceSpan is true if SPAN/0 is written for empty sections.
func (p *TemplateParser) endSection(forceSpan bool) {
	if forceSpan || len(p.pending.list) > 0 {
		p.out.Append(p.pending.Reduce(template.Span))
	}
	p.pending = Section{}
}

// baseParser handles the common functionality of all keyword directives.
func baseParser(p *TemplateParser, v Directive) (ret State, err error) {
	switch k := builtin(v.Key); v.Key {
	case "once", "cycle", "shuffle":
		if e := UnexpectedExpression(v); e != nil {
			err = e
		} else {
			t := MakeSubParser(sequenceParser, nil)
			ret = MakeChain(&t, Statement(func(r rune) (ret State) {

				if res, e := t.reduce(k); e != nil {
					p.err = e
				} else {
					p.pending.Append(res)
					ret = p.NewRune(r)
				}
				return
			}))
		}
	case "if", "unless":
		if e := ExpectedExpression(v); e != nil {
			err = e
		} else {
			t := MakeSubParser(conditionParser, v.Expression)
			ret = MakeChain(&t, Statement(func(r rune) (ret State) {
				// 	if len(t.pending.list) > 0 {
				// 		t.out.Append(t.pending.Reduce(template.Span))
				// 	}
				// 	p.pending.Append(t.out.Reduce(k))
				// 	ret = p.NewRune(r)

				if res, e := t.reduce(k); e != nil {
					p.err = e
				} else {
					p.pending.Append(res)
					ret = p.NewRune(r)
				}
				return
			}))
		}
	}
	return
}

func defaultParser(p *TemplateParser, v Directive) (ret State, err error) {
	// println("defaultParser", v.Key)
	if next, e := baseParser(p, v); e != nil {
		err = e
	} else if next != nil {
		ret = next
	} else {
		err = errutil.New("default parser", UnknownDirective(v))
	}
	return
}

func sequenceParser(p *TemplateParser, v Directive) (ret State, err error) {
	// println("sequenceParser", v.Key)
	if next, e := baseParser(p, v); e != nil {
		err = e
	} else if next != nil {
		ret = next
	} else if e := UnexpectedExpression(v); e != nil {
		err = e
	} else if v.Key == "or" {
		p.endSection(true) // put everything accumulated during this block into our final output;
		ret = p            // continue with the existing sequence parser
	} else if v.Key != "end" {
		err = errutil.New("sequence parser", UnknownDirective(v))
	}
	return
}

func conditionParser(p *TemplateParser, v Directive) (ret State, err error) {
	// println("condition parser", v.Key)
	if next, e := branchParser(p, v); e != nil {
		err = e
	} else if next != nil {
		ret = next
	} else if v.Key != "end" {
		err = errutil.New("condition parser", UnknownDirective(v))
	}
	return
}

func endingParser(p *TemplateParser, v Directive) (ret State, err error) {
	// println("ending parser", v.Key)
	if next, e := baseParser(p, v); e != nil {
		err = e
	} else if next != nil {
		ret = next
	} else if v.Key != "end" {
		err = errutil.New("ending parser", UnknownDirective(v))
	}
	return
}

func branchParser(p *TemplateParser, v Directive) (ret State, err error) {
	// println("branch parser", v.Key)
	switch v.Key {
	case "else", "otherwise":
		if e := UnexpectedExpression(v); e != nil {
			err = e
		} else {
			p.endSection(true)
			// we expect to see "end" next (plus/or minus subexpressions)
			// which is the end of our branch; the template.IfStatement handler
			// which is our parent, will get the next crack at the rune stream.
			t := MakeSubParser(endingParser, nil)
			ret = MakeChain(&t, StateExit(func() {
				if res, e := t.reduce(template.Span); e != nil {
					p.err = e
				} else {
					p.pending.Append(res)
				}
				return
			}))
		}
	case "elseIf", "otherwiseIf", "elseUnless", "otherwiseUnless":
		if e := ExpectedExpression(v); e != nil {
			err = e
		} else {
			k := builtin(v.Key)
			p.endSection(true)
			t := MakeSubParser(branchParser, v.Expression)
			ret = MakeChain(&t, StateExit(func() {
				if res, e := t.reduce(k); e != nil {
					p.err = e
				} else {
					p.pending.Append(res)
				}
				return
			}))
		}
	}
	return
}

func builtin(key string) template.BuiltinType {
	which := map[string]template.BuiltinType{
		"once":    template.Stopping,
		"cycle":   template.Cycle,
		"shuffle": template.Shuffle,
		//
		"if":          template.IfStatement,
		"elseIf":      template.IfStatement,
		"otherwiseIf": template.IfStatement,
		//
		"unless":          template.UnlessStatement,
		"elseUnless":      template.UnlessStatement,
		"otherwiseUnless": template.UnlessStatement,
	}
	return which[key]
}
