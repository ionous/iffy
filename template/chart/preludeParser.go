package chart

import (
	"github.com/ionous/errutil"
)

type Functions map[string]bool

// prelude parser reads the header of a directive or an argument in a function call
type preludeParser struct {
	prelude  Argument
	err      error
	newArg   argFactory   // for arguments of prelude functions
	newBlock blockFactory // for sub directives
}

func newCustomPrelude(blocks blockFactory, args argFactory) *preludeParser {
	return &preludeParser{newBlock: blocks, newArg: args}
}

//
func newDefaultPrelude() argParser {
	return &preludeParser{newBlock: newSubParser, newArg: newDefaultPrelude}
}

// attempt to build a quote, number, function, reference, or sub-block.
func (p *preludeParser) NewRune(r rune) (ret State) {
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

func (p preludeParser) GetArg() (Argument, error) {
	return p.prelude, p.err
}

func (p *preludeParser) setArg(s Argument) {
	if p.prelude != nil {
		panic("arg already set")
	} else {
		p.prelude = s
	}
}

// the passed rune is a bracket
func (p *preludeParser) parseDirective(r rune) State {
	dir := p.newBlock()
	return makeChain(dir, Statement(func(r rune) State {
		if block, e := dir.GetBlock(); e != nil {
			p.err = e
		} else if arg, ok := block.(Argument); !ok {
			p.err = errutil.Fmt("unknown block %T", block)
		} else {
			p.prelude = arg
		}
		return nil // state exit action
	}))
}

// the passed rune starts some quoted text
func (p *preludeParser) parseQuote(r rune) State {
	var quote quoteParser
	return parseChain(r, &quote, Statement(func(r rune) State {
		// hey look: its a state-transition action using a transient.
		if v, e := quote.GetString(); e != nil {
			p.err = e
		} else {
			arg := &Quote{v}
			p.setArg(arg)
		}
		return nil // state exit action
	}))
}

// the passed rune starts a reference or a function
func (p *preludeParser) parseIdent(r rune) State {
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
					arg := &Function{name, args}
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
					arg := &Reference{fields}
					p.setArg(arg)
				}
				return nil // state exit action
			}))
		} else {
			arg := &Reference{[]string{name}}
			p.setArg(arg)
			// done.
		}
		return
	}))
}

// the passed rune starts a number
func (p *preludeParser) parseNumber(r rune) State {
	var num numParser
	return parseChain(r, &num, Statement(func(r rune) State {
		if v, e := num.GetValue(); e != nil {
			p.err = e
		} else {
			arg := &Number{v}
			p.setArg(arg)
		}
		return nil // state exit action
	}))
}
