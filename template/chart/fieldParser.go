package chart

import (
	"github.com/ionous/errutil"
)

// FieldParser reads identifiers separated by dots.
type FieldParser struct {
	err     error
	fields  []string
	pending bool
}

// NewRune starts on the first letter of a new field.
func (p *FieldParser) NewRune(r rune) State {
	var id IdentParser
	return parseChain(r, &id, Statement(func(r rune) (ret State) {
		if n := id.GetName(); len(n) > 0 {
			p.fields = append(p.fields, n)
			if r == '.' {
				p.pending = true
				ret = p // loop...
			} else {
				p.pending = false
			}
		}
		return
	}))
}

// GetFields returns an array of parsed identifiers.
func (p FieldParser) GetFields() (ret []string, err error) {
	if e := p.err; e != nil {
		err = e
	} else if p.pending {
		err = errutil.New("incomplete fields")
	} else {
		ret = p.fields
	}
	return
}
