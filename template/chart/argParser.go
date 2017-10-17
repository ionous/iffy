package chart

//  x: parameters |
type argParser struct {
	args []Spec
	err  error
}

func (p argParser) GetArgs() ([]Spec, error) {
	return p.args, p.err
}

// first character past the separator
func (p *argParser) NewRune(r rune) State {
	return parseChain(r, spaces, Statement(p.readArg))
}

// r is the start of an arg
func (p *argParser) readArg(r rune) State {
	var head headParser
	return parseChain(r, &head, Statement(func(r rune) (ret State) {
		if arg, e := head.GetSpec(); e != nil {
			p.err = e
		} else if arg != nil {
			p.args = append(p.args, arg)
			if isSpace(r) {
				ret = p.NewRune(r)
			}
		}
		return
	}))
}
