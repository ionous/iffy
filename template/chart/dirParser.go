package chart

type directiveParser struct {
	spec    Spec
	exp     string
	filters []FunctionSpec
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
	} else if p.spec != nil {
		ret = &Directive{p.spec, p.exp, p.filters}
	}
	return
}

// NewRune starts just after the opening of a directive or its trim.
// FIX: need to support implicit directives ( which dont allow trailing trim )
func (p *directiveParser) NewRune(r rune) State {
	head := headParser{newBlock: subDirectiveFactory}
	return parseChain(r, &head, Statement(func(r rune) (ret State) {
		if spec, e := head.GetSpec(); e != nil {
			p.err = e
		} else if spec != nil {
			tail := tailParser{canTrim: p.canTrim} // expression
			ret = parseChain(r, &tail, Statement(func(r rune) (ret State) {
				if exp, ctrl, e := tail.GetTail(); e != nil {
					p.err = e
				} else {
					p.spec = spec
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
