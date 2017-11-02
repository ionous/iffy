package old

import (
	"github.com/ionous/errutil"
)

// starts first past the bar, it reads a single function and its arguments.
type filterParser struct {
	name         string
	args         []Argument
	err          error
	newArgParser ArgFactory
}

func newFilterParser(f ArgFactory) *filterParser {
	return &filterParser{newArgParser: f}
}

// GetFunc returns one function and its arguments.
func (p *filterParser) GetFunc() (ret *Func, err error) {
	if e := p.err; e != nil {
		err = e
	} else if len(p.name) == 0 {
		err = errutil.New("missing function call after filter")
	} else {
		ret = &Func{p.name, p.args}
	}
	return
}

// NewRune starts with the first character past the bar
func (p *filterParser) NewRune(r rune) State {
	var id IdentParser
	return parseChain(r, spaces, makeChain(&id, Statement(func(r rune) (ret State) {
		// read an identifier, which ends with any unknown character.
		if n, e := id.GetName(); e != nil {
			p.err = e
		} else {
			// if that character was a separator: start parsing args
			if isSeparator(r) {
				args := newCallParser(p.newArgParser)
				// use makeChain to skip the separator itself
				ret = makeChain(spaces, makeChain(args, Statement(func(r rune) State {
					if args, e := args.GetArgs(); e != nil {
						p.err = e
					} else {
						p.name = n
						p.args = args
					}
					return nil // state exit action
				})))
			}
		}
		return
	})))
}
