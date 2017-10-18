package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/sliceOf"
)

type Functions map[string]bool

type headParser struct {
	spec     Spec
	err      error
	newBlock blockFactory
}

//
var headFactory specFactory = func() specParser {
	return &headParser{newBlock: subDirectiveFactory}
}

// attempt to build one of:
// text, number, function, reference, or sub-directive.
func (p *headParser) NewRune(r rune) (ret State) {
	switch {
	case isQuote(r):
		ret = p.parseQuote(r)
	case isLetter(r):
		ret = p.parseIdent(r)
	case isNumber(r) || r == '+' || r == '-':
		ret = p.parseNumber(r)
	case isOpenBracket(r):
		ret = p.parseDirective(r)
	}
	return
}

func (p headParser) GetSpec() (Spec, error) {
	return p.spec, p.err
}

func (p *headParser) setSpec(s Spec) {
	if p.spec != nil {
		panic("spec already set")
	} else {
		p.spec = s
	}
}

// the passed rune is a bracket
func (p *headParser) parseDirective(r rune) State {
	dir := p.newBlock()
	return makeChain(dir, Statement(func(r rune) State {
		if block, e := dir.GetBlock(); e != nil {
			p.err = e
		} else if spec, ok := block.(Spec); !ok {
			p.err = errutil.Fmt("unknown block %T", block)
		} else {
			p.spec = spec
		}
		return nil // state exit action
	}))
}

// the passed rune starts some quoted text
func (p *headParser) parseQuote(r rune) State {
	var quote quoteParser
	return parseChain(r, &quote, Statement(func(r rune) State {
		// hey look: its a state-transition action using a transient.
		if v, e := quote.GetString(); e != nil {
			p.err = e
		} else {
			spec := &TextSpec{v}
			p.setSpec(spec)
		}
		return nil // state exit action
	}))
}

// the passed rune starts a reference or a function
func (p *headParser) parseIdent(r rune) State {
	var name identParser
	return parseChain(r, &name, Statement(func(r rune) (ret State) {
		if name, e := name.GetName(); e != nil {
			p.err = e
		} else if isSeparator(r) {
			args := newArgParser(headFactory)
			ret = makeChain(args, Statement(func(r rune) State {
				if args, e := args.GetSpecs(); e != nil {
					p.err = e
				} else {
					spec := &FunctionSpec{name, args}
					p.setSpec(spec)
				}
				return nil // state exit action
			}))
		} else if r == '.' {
			fields := fieldParser{fields: sliceOf.String()}
			ret = makeChain(&fields, Statement(func(r rune) State {
				if fields, e := fields.GetFields(); e != nil {
					p.err = e
				} else {
					spec := &ReferenceSpec{fields}
					p.setSpec(spec)
				}
				return nil // state exit action
			}))
		}
		return
	}))
}

// the passed rune starts a number
func (p *headParser) parseNumber(r rune) State {
	var num numParser
	return parseChain(r, &num, Statement(func(r rune) State {
		if v, e := num.GetValue(); e != nil {
			p.err = e
		} else {
			spec := &NumberSpec{v}
			p.setSpec(spec)
		}
		return nil // state exit action
	}))
}
