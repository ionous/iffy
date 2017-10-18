package chart

import (
	"github.com/ionous/errutil"
)

// starts first past the bar, it reads a single function and its parameters.
type filterParser struct {
	name         string
	args         []Argument
	err          error
	newArgParser argFactory
}

func newFilterParser(f argFactory) *filterParser {
	return &filterParser{newArgParser: f}
}

// GetFunction returns
func (p *filterParser) GetFunction() (ret *Function, err error) {
	if e := p.err; e != nil {
		err = e
	} else if len(p.name) == 0 {
		err = errutil.New("missing function call after filter")
	} else {
		ret = &Function{p.name, p.args}
	}
	return
}

// first character past the bar
func (p *filterParser) NewRune(r rune) State {

	var id identParser
	return parseChain(r, spaces, makeChain(&id, Statement(func(r rune) (ret State) {
		// read an identifier, which ends with any unknown character.
		if n, e := id.GetName(); e != nil {
			p.err = e
		} else {

			// if that character was a separator: start parsing args
			if isSeparator(r) {
				args := newCallParser(p.newArgParser)
				// use makeChain to skip the separator itself
				ret = makeChain(args, Statement(func(r rune) State {
					if args, e := args.GetArgs(); e != nil {
						p.err = e
					} else {
						p.name = n
						p.args = args
					}
					return nil // state exit action
				}))
			}
		}
		return
	})))
}
