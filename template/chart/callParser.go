package chart

//  x: arg arg arg
type callParser struct {
	args         []Argument
	err          error
	newArgParser argFactory
}

func newCallParser(f argFactory) *callParser {
	return &callParser{newArgParser: f}
}

// creates new argParser;
// the primary implementation is newDefaultPrelude.
type argFactory func() argParser

// the primary implementation is preludeParser.
type argParser interface {
	NewRune(rune) State
	GetArg() (Argument, error)
}

// GetArgs returns the arguments for the called function.
func (p callParser) GetArgs() ([]Argument, error) {
	return p.args, p.err
}

// first character past a function separator;
// each arg is read by a arg parser created by argFactory;
// args are separated by spaces
func (p *callParser) NewRune(r rune) State {
	return parseChain(r, spaces, Statement(p.readArg))
}

// r is the start of an arg
func (p *callParser) readArg(r rune) (ret State) {
	if r != eof {
		argParser := p.newArgParser()
		ret = parseChain(r, argParser, Statement(func(r rune) (ret State) {
			if arg, e := argParser.GetArg(); e != nil {
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
	return
}
