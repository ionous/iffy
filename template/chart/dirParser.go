package chart

type directiveParser struct {
	arg     Argument
	exp     string
	filters []FunctionArg
	err     error
	canTrim bool
}

//
var topDirectiveFactory blockFactory = func() subBlockParser {
	dir := directiveParser{canTrim: true}
	return &dir
}

//
var subDirectiveFactory blockFactory = func() subBlockParser {
	dir := directiveParser{}
	return &dir
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
// FIX: need to support implicit directives ( which dont allow trailing trim )
func (p *directiveParser) NewRune(r rune) State {
	prelude := newHeadParser(subDirectiveFactory, headFactory)
	//
	return parseChain(r, prelude, Statement(func(r rune) (ret State) {
		if arg, e := prelude.GetArg(); e != nil {
			p.err = e
		} else if arg != nil {
			tail := tailParser{canTrim: p.canTrim} // expression
			ret = parseChain(r, &tail, Statement(func(r rune) (ret State) {
				if exp, ctrl, e := tail.GetTail(); e != nil {
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
			}))
		}
		return
	}))
}

// if the rune is a filter character, we add a new function
func (p *directiveParser) filter(r rune) (ret State) {
	if isFilter(r) {
		filter := newFilterParser(headFactory)
		ret = makeChain(filter, Statement(func(r rune) (ret State) {
			if f, e := filter.GetFunction(); e != nil {
				p.err = e
			} else {
				p.filters = append(p.filters, *f)
				// ************* does this make sense!????
				ret = p.filter(r)
			}
			return
		}))
	}
	return
}
