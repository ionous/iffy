package chart

import (
	"github.com/ionous/errutil"
)

type fieldParser struct {
	err     error
	fields  []string
	pending bool
}

func (p fieldParser) GetFields() (ret []string, err error) {
	if e := p.err; e != nil {
		err = e
	} else if p.pending {
		err = errutil.New("incomplete fields")
	} else if len(p.fields) == 0 {
		err = errutil.New("empty fields")
	} else {
		ret = p.fields
	}
	return
}

// NewRune where r is the first letter of the field
func (p *fieldParser) NewRune(r rune) State {
	var id identParser
	return parseChain(r, &id, Statement(func(r rune) (ret State) {
		if n, e := id.GetName(); e != nil {
			p.err = e
		} else {
			p.fields = append(p.fields, n)
			if r == '.' {
				p.pending = true
				ret = p // loop... more idents.
			} else {
				p.pending = false
			}
		}
		return
	}))
}
