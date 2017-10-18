package chart

import (
	"github.com/ionous/errutil"
)

type Functions map[string]bool

// head parser reads the "prelude" of a directive
type headParser struct {
	head     Argument
	err      error
	newArg   argFactory   // for arguments of head functions
	newBlock blockFactory // for sub directives
}

func newHeadParser(blocks blockFactory, args argFactory) *headParser {
	return &headParser{newBlock: blocks, newArg: args}
}

//
var headFactory argFactory = func() argParser {
	return &headParser{newBlock: subDirectiveFactory}
}

// attempt to build a quote, number, function, reference, or sub-block.
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

func (p headParser) GetArg() (Argument, error) {
	return p.head, p.err
}

func (p *headParser) setArg(s Argument) {
	if p.head != nil {
		panic("arg already set")
	} else {
		p.head = s
	}
}

// the passed rune is a bracket
func (p *headParser) parseDirective(r rune) State {
	dir := p.newBlock()
	return makeChain(dir, Statement(func(r rune) State {
		if block, e := dir.GetBlock(); e != nil {
			p.err = e
		} else if arg, ok := block.(Argument); !ok {
			p.err = errutil.Fmt("unknown block %T", block)
		} else {
			p.head = arg
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
			arg := &QuotedArg{v}
			p.setArg(arg)
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
			args := newCallParser(p.newArg)
			ret = makeChain(args, Statement(func(r rune) State {
				if args, e := args.GetArgs(); e != nil {
					p.err = e
				} else {
					arg := &FunctionArg{name, args}
					p.setArg(arg)
				}
				return nil // state exit action
			}))
		} else if r == '.' {
			fields := newFieldParser(name)
			ret = makeChain(fields, Statement(func(r rune) State {
				if fields, e := fields.GetFields(); e != nil {
					p.err = e
				} else {
					arg := &ReferenceArg{fields}
					p.setArg(arg)
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
			arg := &NumberArg{v}
			p.setArg(arg)
		}
		return nil // state exit action
	}))
}
