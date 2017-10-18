package chart

type directiveParser struct {
	arg     Argument
	exp     string
	filters []Function
	err     error
	canTrim bool
}

// newTopParser is used in blocks
func newTopParser() subBlockParser {
	return &directiveParser{canTrim: true}
}

// newSubParser is used inside of other directives.
func newSubParser() subBlockParser {
	return &directiveParser{canTrim: false}
}

// GetBlock
func (p directiveParser) GetBlock() (ret Block, err error) {
	if e := p.err; e != nil {
		err = e
	} else if p.arg != nil {
		ret = &Directive{p.arg, p.exp, p.filters}
	}
	return
}

// NewRune starts just after the opening of a directive or its trim.
func (p *directiveParser) NewRune(r rune) State {
	prelude := newCustomPrelude(newSubParser, newDefaultPrelude)
	return parseChain(r, prelude, Statement(func(r rune) (ret State) {
		if arg, e := prelude.GetArg(); e != nil {
			p.err = e
		} else if arg != nil {
			epilouge := newEpilogueParser(p.canTrim) // expression
			ret = parseChain(r, spaces, makeChain(epilouge, Statement(func(r rune) (ret State) {
				if exp, ctrl, e := epilouge.GetResult(); e != nil {
					p.err = e
				} else {
					p.arg = arg
					p.exp = exp
					switch {
					case isTrim(ctrl):
						ret = spaces
					case isFilter(ctrl):
						ret = p.filter(r)
					}
				}
				return
			})))
		}
		return
	}))
}

// r is the rune just after a filter.
func (p *directiveParser) filter(r rune) (ret State) {
	filter := newFilterParser(newDefaultPrelude)
	ret = parseChain(r, filter, Statement(func(r rune) (ret State) {
		if f, e := filter.GetFunction(); e != nil {
			p.err = e
		} else {
			p.filters = append(p.filters, *f)
			ret = parseChain(r, spaces, Statement(func(r rune) (ret State) {
				if isFilter(r) {
					ret = Statement(p.filter)
				}
				return
			}))
		}
		return
	}))
	return
}
