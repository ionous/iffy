package chart

//  x: arg arg arg
type argParser struct {
	args          []Spec
	err           error
	newSpecParser specFactory
}

// creates new specParser;
// the primary implementation is headFactory.
type specFactory func() specParser

// the primary implementation is headParser.
type specParser interface {
	NewRune(rune) State
	GetSpec() (Spec, error)
}

func (p argParser) GetSpecs() ([]Spec, error) {
	return p.args, p.err
}

// first character past a function separator;
// each arg is read by a spec parser created by specFactory;
// args are separated by spaces
func (p *argParser) NewRune(r rune) State {
	return parseChain(r, spaces, Statement(p.readArg))
}

// r is the start of an arg
func (p *argParser) readArg(r rune) State {
	head := p.newSpecParser()
	return parseChain(r, head, Statement(func(r rune) (ret State) {
		if arg, e := head.GetSpec(); e != nil {
			p.err = e
		} else if arg != nil {
			p.args = append(p.args, arg)
			if isSpace(r) {
				ret = p // loop...
			}
		}
		return
	}))
}
