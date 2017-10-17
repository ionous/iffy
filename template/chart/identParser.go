package chart

import (
	"github.com/ionous/errutil"
)

type identParser struct {
	err   error
	runes []rune
}

func (p identParser) GetName() (ret string, err error) {
	if e := p.err; e != nil {
		err = e
	} else if s := string(p.runes); len(s) == 0 {
		err = errutil.New("empty identifier")
	} else {
		ret = s
	}
	return
}

// first character of the identifier
func (p *identParser) NewRune(r rune) (ret State) {
	if isLetter(r) {
		p.runes = append(p.runes, r)
		ret = Statement(p.body)
	} else {
		p.err = errutil.New("identifier too short")
	}
	return
}

// subsequent characters can be letters or numbers
// noting that fields are separated by dots "."
func (p *identParser) body(r rune) (ret State) {
	if isLetter(r) || isNumber(r) {
		p.runes = append(p.runes, r)
		ret = Statement(p.body) // loop ...
	}
	return
}
