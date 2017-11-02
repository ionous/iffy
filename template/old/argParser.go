package old

import (
	"github.com/ionous/iffy/template/postfix"
)

// ArgParser reads a single operand or function call.
type ArgParser struct {
	out   postfix.Shunt
	err   error
	nofun bool
}

// attempt to build a quote, number, function, reference, or sub-block.
func (p *ArgParser) NewRune(r rune) (ret State) {
	switch {
	case isQuote(r):
		ret = p.parseQuote(r)
	case isLetter(r):
		ret = p.parseIdent(r, p.seed.Reset())
	case isNumber(r) || r == '+' || r == '-':
		ret = p.parseNumber(r)
	}
	return
}

// GetArg can return nil and no error
func (p ArgParser) GetExpression() (ret postfix.Expression, err error) {
	if e := p.err; e != nil {
		err = e
	} else {
		ret = p.b.Expression()
	}
	return
}

// the passed rune starts some quoted text
func (p *ArgParser) parseQuote(r rune) State {
	var quote QuoteParser
	return parseChain(r, &quote, Statement(func(r rune) State {
		// hey look: its a state-transition action using a transient.
		if v, e := quote.GetString(); e != nil {
			p.err = e
		} else {
			arg := Quote(v)
			if e := p.out.AddFunction(arg); e != nil {
				p.err = e
			}
		}
		return nil // state exit action
	}))
}

// the passed rune starts a reference or a function
func (p *ArgParser) parseIdent(r rune, seed []rune) State {
	name := IdentParser{runes: Runes{seed}}
	return parseChain(r, &name, Statement(func(r rune) (ret State) {
		if name, e := name.GetName(); e != nil {
			p.err = e
		} else if !p.nofun && isSeparator(r) {
			arg := Command{name, len(args)}
			if e := p.out.AddFunction(arg); e != nil {
				p.err = e
			} else {
				ret = makeChain(spaces, MakeCallParser(&p.out))
			}
		} else if r == '.' {
			fields := newFieldParser(name)
			ret = makeChain(fields, Statement(func(r rune) State {
				if fields, e := fields.GetFields(); e != nil {
					p.err = e
				} else {
					arg := Reference{fields}
					if _, e := p.b.Write([]Function{arg}); e != nil {
						p.err = e
					}
				}
				return nil // state exit action
			}))
		} else {
			arg := Reference{[]string{name}}
			if _, e := p.b.Write([]Function{arg}); e != nil {
				p.err = e
			}
			// done.
		}
		return
	}))
}

// the passed rune starts a number
func (p *ArgParser) parseNumber(r rune) State {
	var num NumParser
	return parseChain(r, &num, Statement(func(r rune) State {
		if v, e := num.GetValue(); e != nil {
			p.err = e
		} else {
			arg := Number(v)
			if _, e := p.b.Write([]Function{arg}); e != nil {
				p.err = e
			}
		}
		return nil // state exit action
	}))
}
