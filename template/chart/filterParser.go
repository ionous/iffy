package chart

//  x: parameteres |
// uses argParser
type filterParser struct {
	name          string
	args          []Spec
	err           error
	newSpecParser specFactory
}

func (p *filterParser) GetFunction() (ret *FunctionSpec, err error) {
	if e := p.err; e != nil {
		err = e
	} else {
		ret = &FunctionSpec{p.name, p.args}
	}
	return
}

// first character past the bar
func (p *filterParser) NewRune(r rune) State {
	var name identParser
	return parseChain(r, spaces, makeChain(&name, Statement(func(r rune) (ret State) {
		if n, e := name.GetName(); e != nil {
			p.err = e
		} else if isSeparator(r) {
			args := argParser{newSpecParser: p.newSpecParser}
			ret = parseChain(r, &args, Statement(func(r rune) State {
				if args, e := args.GetSpecs(); e != nil {
					p.err = e
				} else {
					p.name = n
					p.args = args
				}
				return nil // state exit action
			}))
		}
		return
	})))
}
